// Package game is the Go port of LORD 5.02 JS, Synchronet's JavaScript
// rewrite of Legend of the Red Dragon (xtrn/lord/lord.js, by Deuce).
//
// lord.js is a script of global variables and functions; here the globals
// are fields of Game and each function is a method with the camel-cased
// name of the original (forest() is (*Game).forest, check_mail() is
// (*Game).checkMail).  ref/synchronet holds the source that was ported.
//
// Left out on purpose: RIPscrip graphics and the SyncTERM icon cache (the
// `rip` flag is always false), the remote game server (lordsrv.js / psock),
// and JavaScript in-game modules (IGMs), which can't run in Go.
package game

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"math/rand/v2"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/pridkett/doorgames-common/door"
	"github.com/pridkett/lord/internal/assets"
)

// Version is lord.js's `ver`.
const Version = "5.02 JS"

// Game holds lord.js's global state for one caller.
type Game struct {
	Dir  string // the world (exec) directory
	Ctx  door.Context
	con  *Console
	rng  *rand.Rand
	now  func() time.Time
	sess *door.Session

	Settings   Settings
	Tournament TournamentSettings

	pfile     *RecordFile[Player]
	statefile *RecordFile[State]
	player    *Player
	state     *State

	txt      []string       // lordtxt.lrd
	txtindex map[string]int // @#NAME -> index of the line after it

	whichCastle  int64
	curlinenum   int
	curcolnum    int
	morechk      bool
	cleanupFiles []string
	onExitArmed  bool
	online       *os.File // liveness lock (see liveness.go)
	mutexes      map[string]*os.File

	// rip is always false: RIPscrip is not supported.  It's kept so the
	// translated code reads like the original.
	rip bool
}

// exitSignal unwinds the stack for lord.js's exit(code).
type exitSignal struct{ code int }

// New prepares a game in dir for the caller in sess.
func New(dir string, sess *door.Session) *Game {
	g := &Game{
		Dir:        dir,
		sess:       sess,
		rng:        rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(os.Getpid()))),
		now:        time.Now,
		Settings:   DefaultSettings(),
		Tournament: TournamentSettings{},
		txtindex:   map[string]int{},
		curlinenum: 1,
		curcolnum:  1,
		morechk:    true,
	}
	if sess != nil {
		g.Ctx = sess.Ctx
		g.con = NewConsole(sess)
	}
	g.whichCastle = g.random(5) + 1
	return g
}

// exit is lord.js's exit(): it ends the door with code.
func (g *Game) exit(code int) { panic(exitSignal{code}) }

// fatalf is `throw new Error(...)`.
func (g *Game) fatalf(format string, args ...any) {
	panic(fmt.Errorf(format, args...))
}

// Run is the script body: cmdline() then start().  It returns the exit
// code.  Errors thrown by the game are returned with code 1.
func (g *Game) Run() (code int, err error) {
	defer func() {
		r := recover()
		if r == nil {
			r = exitSignal{0}
		}
		if ex, ok := r.(exitSignal); ok {
			code = ex.code
		} else {
			code = 1
			err = fmt.Errorf("%v\n%s", r, debug.Stack())
		}
		g.onExit()
	}()
	g.start()
	return 0, nil
}

// onExit is the js.on_exit handler start() registers.
func (g *Game) onExit() {
	defer func() { _ = recover() }()
	if g.onExitArmed && g.player != nil {
		g.player.OnNow = false
		g.player.Put()
	}
	for len(g.cleanupFiles) > 0 {
		_ = os.Remove(g.cleanupFiles[0])
		g.releaseMutex(g.cleanupFiles[0])
		g.cleanupFiles = g.cleanupFiles[1:]
	}
	if g.con != nil {
		g.con.Flush()
	}
}

// random is Synchronet's random(n): 0..n-1, and 0 when n < 2.
func (g *Game) random(n int64) int64 {
	if n < 2 {
		return 0
	}
	return g.rng.Int64N(n)
}

// randomf is random() applied to a JavaScript number that may have a
// fraction (it is truncated like ToInt32).
func (g *Game) randomf(n float64) int64 { return g.random(int64(n)) }

// time is Synchronet's time(): Unix seconds.
func (g *Game) time() int64 { return g.now().Unix() }

// mswait sleeps ms milliseconds, flushing output first.
func (g *Game) mswait(ms int) {
	if g.con != nil {
		g.con.Flush()
	}
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

// ---- files -------------------------------------------------------------

// gamedir is gamedir(): the path of a game data file.  GamePrefix may name
// a directory or a file-name prefix, as in lord.js.
func (g *Game) gamedir(fname string) string {
	pre := g.Settings.GamePrefix
	var base string
	switch {
	case pre == "":
		base = g.Dir + string(filepath.Separator)
	case filepath.IsAbs(pre):
		base = pre
	default:
		base = filepath.Join(g.Dir, pre)
	}
	if st, err := os.Stat(base); err == nil && st.IsDir() && !strings.HasSuffix(base, string(filepath.Separator)) {
		base += string(filepath.Separator)
	}
	return base + fname
}

// gameOrExec is game_or_exec(): a sysop's copy of a distribution file in
// MenuDir or the game directory, or "" when the built-in copy (the
// original's exec directory) should be used.
func (g *Game) gameOrExec(fname string) string {
	if md := g.Settings.MenuDir; md != "" {
		if !filepath.IsAbs(md) {
			md = filepath.Join(g.Dir, md)
		}
		if st, err := os.Stat(md); err == nil && st.IsDir() {
			if p := filepath.Join(md, fname); fileExists(p) {
				return p
			}
			return g.execFile(fname)
		}
	}
	if p := g.gamedir(fname); fileExists(p) {
		return p
	}
	return g.execFile(fname)
}

// execFile is js.exec_dir + fname: the world directory's copy if there is
// one (found case-insensitively), else "" for the built-in asset.
func (g *Game) execFile(fname string) string {
	if p := findFold(g.Dir, fname); p != "" {
		return p
	}
	return ""
}

// findFold finds fname in dir ignoring case.
func findFold(dir, fname string) string {
	p := filepath.Join(dir, fname)
	if fileExists(p) {
		return p
	}
	ents, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	for _, e := range ents {
		if strings.EqualFold(e.Name(), fname) && !e.IsDir() {
			return filepath.Join(dir, e.Name())
		}
	}
	return ""
}

// assetBytes reads a distribution file through gameOrExec, falling back to
// the copy built into the binary.
func (g *Game) assetBytes(fname string) ([]byte, error) {
	if p := g.gameOrExec(fname); p != "" {
		return os.ReadFile(p)
	}
	return assets.Read(fname)
}

// assetLines is assetBytes split into lines (File.readln without the
// line ends).
func (g *Game) assetLines(fname string) ([]string, error) {
	b, err := g.assetBytes(fname)
	if err != nil {
		return nil, err
	}
	return splitLines(b), nil
}

// assetExists reports whether a distribution file is available.
func (g *Game) assetExists(fname string) bool {
	if g.gameOrExec(fname) != "" {
		return true
	}
	_, err := assets.Read(fname)
	return err == nil
}

func splitLines(b []byte) []string {
	var out []string
	sc := bufio.NewScanner(bytes.NewReader(b))
	sc.Buffer(make([]byte, 64*1024), 1024*1024)
	for sc.Scan() {
		out = append(out, strings.TrimRight(sc.Text(), "\r"))
	}
	return out
}

func fileExists(p string) bool {
	st, err := os.Stat(p)
	return err == nil && !st.IsDir()
}

func fileRemove(p string) bool { return os.Remove(p) == nil }

func fileRename(from, to string) bool { return os.Rename(from, to) == nil }

// fileDate is file_date(): modification time in Unix seconds (-1 if none).
func fileDate(p string) int64 {
	st, err := os.Stat(p)
	if err != nil {
		return -1
	}
	return st.ModTime().Unix()
}

// readLines reads a text file's lines; ok is false when it can't be opened.
func readLines(p string) (lines []string, ok bool) {
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, false
	}
	return splitLines(b), true
}

// appendLines appends lines (File.open('a') + writeln) to p.
func appendLines(p string, lines ...string) error {
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, l := range lines {
		if _, err := f.WriteString(l + "\r\n"); err != nil {
			return err
		}
	}
	return nil
}

// writeLines replaces p with lines (File.open('w') + writeln).
func writeLines(p string, lines ...string) error {
	var b strings.Builder
	for _, l := range lines {
		b.WriteString(l)
		b.WriteString("\r\n")
	}
	return os.WriteFile(p, []byte(b.String()), 0o644)
}

// fmutex is fmutex(): create fname+".lock" exclusively, retrying for 15
// seconds.  str is the lock file's content.
//
// DIFF: the creating session also holds a flock on the lock file until
// rfmutex.  Renegade SIGKILLs a door on disconnect, which leaves its lock
// files behind; a lock file that nobody holds is stale and is removed.
func (g *Game) fmutex(fname string, str ...string) bool {
	end := time.Now().Add(15 * time.Second)
	fname += ".lock"
	content := "lord " + Version + "\n"
	if len(str) > 0 {
		content = str[0]
	}
	for {
		f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o644)
		if err == nil {
			_ = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
			_, _ = f.WriteString(content)
			if g.mutexes == nil {
				g.mutexes = map[string]*os.File{}
			}
			g.mutexes[fname] = f
			break
		}
		if !errors.Is(err, fs.ErrExist) {
			g.fatalf("Unable to create %s please notify Sysop", fname)
		}
		if staleMutex(fname) {
			_ = os.Remove(fname)
			continue
		}
		if time.Now().After(end) {
			g.fatalf("Unable to create %s please notify Sysop", fname)
		}
		time.Sleep(100 * time.Millisecond)
	}
	g.cleanupFiles = append(g.cleanupFiles, fname)
	return true
}

// rfmutex is rfmutex(): remove a lock taken by fmutex.
func (g *Game) rfmutex(fname string) {
	fname += ".lock"
	for i, f := range g.cleanupFiles {
		if f == fname {
			_ = os.Remove(fname)
			g.releaseMutex(fname)
			g.cleanupFiles = append(g.cleanupFiles[:i], g.cleanupFiles[i+1:]...)
			return
		}
	}
}

func (g *Game) releaseMutex(fname string) {
	if f := g.mutexes[fname]; f != nil {
		_ = f.Close()
		delete(g.mutexes, fname)
	}
}

// staleMutex reports whether lock file fname was left by a session that
// is gone: it is over a second old and nobody holds its flock.
func staleMutex(fname string) bool {
	st, err := os.Stat(fname)
	if err != nil || time.Since(st.ModTime()) < time.Second {
		return false
	}
	f, err := os.Open(fname)
	if err != nil {
		return false
	}
	defer f.Close()
	if syscall.Flock(int(f.Fd()), syscall.LOCK_SH|syscall.LOCK_NB) != nil {
		return false
	}
	_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	return true
}

// mutexHeld reports whether fmutex's lock for fname (fname+".lock") is
// held by a live session, removing a stale one.
func mutexHeld(fname string) bool {
	fname += ".lock"
	if !fileExists(fname) {
		return false
	}
	if staleMutex(fname) {
		_ = os.Remove(fname)
		return false
	}
	return true
}

// fileMutexExists is mutexHeld (file_exists on a name fmutex created).
func fileMutexExists(fname string) bool { return mutexHeld(fname) }

// ---- state --------------------------------------------------------------

// getState is get_state(): read state.bin; with lock the state mutex stays
// held until putState.
func (g *Game) getState(lock bool) {
	if g.statefile == nil {
		sf, err := OpenRecordFile[State](g.gamedir("state.bin"))
		if err != nil {
			panic(err)
		}
		g.statefile = sf
	}
	g.fmutex(g.statefile.Name())
	st := &State{sf: g.statefile}
	if g.statefile.Length() < 1 {
		g.statefile.Defaults(st)
		st.Put()
	} else if err := g.statefile.Read(0, st); err != nil {
		panic(err)
	}
	g.state = st
	if !lock {
		g.rfmutex(g.statefile.Name())
	}
}

// putState is put_state(): write state.bin and drop its mutex.
func (g *Game) putState() {
	g.state.Put()
	g.rfmutex(g.statefile.Name())
}

// ---- lordtxt.lrd ----------------------------------------------------------

// buildTxtIndex is build_txt_index().
func (g *Game) buildTxtIndex() {
	lines, err := g.assetLines("lordtxt.lrd")
	if err != nil {
		g.sln("Unable to open lordtxt.lrd!")
		g.exit(1)
	}
	g.txt = lines
	for i, l := range lines {
		if strings.HasPrefix(l, "@#") {
			g.txtindex[l[2:]] = i + 1
		}
	}
}

// lrdfile is lrdfile(): show section fname of lordtxt.lrd.  more enables
// the <MORE> prompts (default false).
func (g *Game) lrdfile(fname string, more ...bool) {
	idx, ok := g.txtindex[fname]
	if !ok {
		return
	}
	mc := g.morechk
	if len(more) == 0 || !more[0] {
		g.morechk = false
	}
	for i := idx; i < len(g.txt); i++ {
		if strings.HasPrefix(g.txt[i], "@#") {
			break
		}
		g.lln(g.txt[i])
	}
	g.morechk = mc
	g.curlinenum = 1
}

// lrdrip is lrdrip(): RIP screens are not supported.
func (g *Game) lrdrip(string, ...bool) {}

// itoa formats an int (record numbers in file names).
func itoa[T ~int | ~int64](n T) string { return strconv.FormatInt(int64(n), 10) }

// castles is lord.js's `castles`.
var castles = []string{
	"Dumbass World",
	"Castle Coldrake",
	"Fortress Liddux",
	"Gannon Keep",
	"Penyon Manor",
	"Dema's Lair",
	"Why are you snooping? Go away!",
}
