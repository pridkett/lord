package game

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// RecordFile is recordfile.js: fixed-length binary records described by a
// field list, with one byte-range lock per record.  The field list comes
// from the `rec` struct tags of T, in declaration order, using the type names
// of recorddefs.js ("String:20", "SignedInteger16", "Boolean", ...).
type RecordFile[T any] struct {
	f      *os.File
	fields []recField
	reclen int
	locks  map[int]bool
}

type recField struct {
	index int    // struct field index
	typ   string // recordfile.js type name
	size  int
	def   string // default value (tag "def=")
}

func parseRecFields(t reflect.Type) ([]recField, int) {
	var out []recField
	total := 0
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("rec")
		if tag == "" {
			continue
		}
		typ, def, _ := strings.Cut(tag, ",def=")
		size := 0
		switch {
		case typ == "SignedInteger8" || typ == "Integer8" || typ == "Boolean":
			size = 1
		case typ == "SignedInteger16" || typ == "Integer16":
			size = 2
		case typ == "SignedInteger" || typ == "Integer":
			size = 4
		case strings.HasPrefix(typ, "String:"):
			size, _ = strconv.Atoi(typ[len("String:"):])
		default:
			panic("unknown record type " + typ)
		}
		out = append(out, recField{index: i, typ: typ, size: size, def: def})
		total += size
	}
	return out, total
}

// OpenRecordFile opens (creating) a record file.
func OpenRecordFile[T any](path string) (*RecordFile[T], error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, err
	}
	fields, n := parseRecFields(reflect.TypeFor[T]())
	return &RecordFile[T]{f: f, fields: fields, reclen: n, locks: map[int]bool{}}, nil
}

// Close releases every lock and closes the file.
func (r *RecordFile[T]) Close() {
	for rec := range r.locks {
		r.UnLock(rec)
	}
	_ = r.f.Close()
}

// Name is the file's path.
func (r *RecordFile[T]) Name() string { return r.f.Name() }

// Length is the number of whole records in the file.
func (r *RecordFile[T]) Length() int {
	st, err := r.f.Stat()
	if err != nil {
		return 0
	}
	return int(st.Size() / int64(r.reclen))
}

// Lock takes the byte-range lock of record rec, waiting up to timeout
// (forever when timeout < 0, a single try when zero).
func (r *RecordFile[T]) Lock(rec int, timeout time.Duration) bool {
	if r.locks[rec] {
		return true
	}
	end := time.Now().Add(timeout)
	for {
		lk := syscall.Flock_t{Type: syscall.F_WRLCK, Whence: io.SeekStart,
			Start: int64(rec * r.reclen), Len: int64(r.reclen)}
		if err := syscall.FcntlFlock(r.f.Fd(), syscall.F_SETLK, &lk); err == nil {
			r.locks[rec] = true
			return true
		}
		if timeout >= 0 && time.Now().After(end) {
			return false
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// UnLock releases record rec.
func (r *RecordFile[T]) UnLock(rec int) {
	if !r.locks[rec] {
		return
	}
	lk := syscall.Flock_t{Type: syscall.F_UNLCK, Whence: io.SeekStart,
		Start: int64(rec * r.reclen), Len: int64(r.reclen)}
	_ = syscall.FcntlFlock(r.f.Fd(), syscall.F_SETLK, &lk)
	delete(r.locks, rec)
}

// Locked reports whether this process holds rec's lock.
func (r *RecordFile[T]) Locked(rec int) bool { return r.locks[rec] }

// Read decodes record rec into v (taking the lock for the read unless it is
// already held, like RecordFile.get).
func (r *RecordFile[T]) Read(rec int, v *T) error {
	held := r.locks[rec]
	if !held {
		r.Lock(rec, -1)
		defer r.UnLock(rec)
	}
	buf := make([]byte, r.reclen)
	if _, err := r.f.ReadAt(buf, int64(rec*r.reclen)); err != nil {
		return err
	}
	decodeRecord(r.fields, buf, reflect.ValueOf(v).Elem())
	return nil
}

// Write encodes v as record rec.
func (r *RecordFile[T]) Write(rec int, v *T) error {
	held := r.locks[rec]
	if !held {
		r.Lock(rec, -1)
		defer r.UnLock(rec)
	}
	buf := encodeRecord(r.fields, r.reclen, reflect.ValueOf(v).Elem())
	_, err := r.f.WriteAt(buf, int64(rec*r.reclen))
	return err
}

// Defaults sets every record field of v to its recorddefs.js default.
func (r *RecordFile[T]) Defaults(v *T) { setDefaults(r.fields, reflect.ValueOf(v).Elem()) }

func setDefaults(fields []recField, v reflect.Value) {
	for _, f := range fields {
		fv := v.Field(f.index)
		switch fv.Kind() {
		case reflect.String:
			fv.SetString(f.def)
		case reflect.Bool:
			fv.SetBool(f.def == "true")
		case reflect.Int64, reflect.Int:
			n, _ := strconv.ParseInt(f.def, 10, 64)
			fv.SetInt(n)
		}
	}
}

func decodeRecord(fields []recField, buf []byte, v reflect.Value) {
	off := 0
	for _, f := range fields {
		b := buf[off : off+f.size]
		off += f.size
		fv := v.Field(f.index)
		switch f.typ {
		case "SignedInteger8":
			// DIFF: recordfile.js tests >= 127, reading a stored 127 as -129.
			fv.SetInt(int64(int8(b[0])))
		case "Integer8":
			fv.SetInt(int64(b[0]))
		case "SignedInteger16":
			fv.SetInt(int64(int16(binary.LittleEndian.Uint16(b))))
		case "Integer16":
			fv.SetInt(int64(binary.LittleEndian.Uint16(b)))
		case "SignedInteger":
			fv.SetInt(int64(int32(binary.LittleEndian.Uint32(b))))
		case "Integer":
			fv.SetInt(int64(binary.LittleEndian.Uint32(b)))
		case "Boolean":
			fv.SetBool(b[0] > 0)
		default: // String:n
			s := strings.ReplaceAll(string(b), "\x00", "")
			fv.SetString(s)
		}
	}
}

func clamp(n, lo, hi int64) int64 {
	return max(lo, min(hi, n))
}

func encodeRecord(fields []recField, n int, v reflect.Value) []byte {
	buf := make([]byte, n)
	off := 0
	for _, f := range fields {
		b := buf[off : off+f.size]
		off += f.size
		fv := v.Field(f.index)
		switch f.typ {
		case "SignedInteger8":
			b[0] = byte(int8(clamp(fv.Int(), -128, 127)))
		case "Integer8":
			b[0] = byte(clamp(fv.Int(), 0, 255))
		case "SignedInteger16":
			binary.LittleEndian.PutUint16(b, uint16(int16(clamp(fv.Int(), -32768, 32767))))
		case "Integer16":
			binary.LittleEndian.PutUint16(b, uint16(clamp(fv.Int(), 0, 65535)))
		case "SignedInteger":
			binary.LittleEndian.PutUint32(b, uint32(int32(clamp(fv.Int(), -2147483648, 2147483647))))
		case "Integer":
			binary.LittleEndian.PutUint32(b, uint32(clamp(fv.Int(), 0, 4294967295)))
		case "Boolean":
			if fv.Bool() {
				b[0] = 255
			}
		default:
			copy(b, fv.String())
		}
	}
	return buf
}

// readRecordAt reads record num of a read-only table file (dragon.bin,
// lenemy.bin, trainers.bin, weapon.bin, armour.bin).
func readRecordAt[T any](path string, num int) (*T, int, error) {
	rf, err := os.Open(path)
	if err != nil {
		return nil, 0, err
	}
	defer rf.Close()
	fields, n := parseRecFields(reflect.TypeFor[T]())
	st, err := rf.Stat()
	if err != nil {
		return nil, 0, err
	}
	count := int(st.Size() / int64(n))
	if num < 0 || num >= count {
		return nil, count, fmt.Errorf("%s: no record %d", path, num)
	}
	buf := make([]byte, n)
	if _, err := rf.ReadAt(buf, int64(num*n)); err != nil {
		return nil, count, err
	}
	var v T
	decodeRecord(fields, buf, reflect.ValueOf(&v).Elem())
	return &v, count, nil
}
