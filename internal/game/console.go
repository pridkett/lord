package game

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pridkett/doorgames-common/door"
)

// Console is the part of dorkit's dk.console that lord.js uses: an ANSI
// screen with a current CGA attribute, key input with named extended keys
// ("KEY_UP", ...), and getstr line editing.
type Console struct {
	sess    *door.Session
	Attr    int // dk.console.attr.value: CGA attribute (blink<<7 | bg<<4 | fg)
	sent    int // attribute last sent to the terminal (-1: unknown)
	Rows    int
	Cols    int
	pending []string
	hangup  bool
}

// NewConsole wraps a door session.
func NewConsole(sess *door.Session) *Console {
	c := &Console{sess: sess, Attr: 7, sent: -1, Rows: 24, Cols: 80}
	if sess.Ctx.Terminal.Rows > 0 {
		c.Rows = sess.Ctx.Terminal.Rows
	}
	if sess.Ctx.Terminal.Cols > 0 {
		c.Cols = sess.Ctx.Terminal.Cols
	}
	return c
}

// CGA colour order to ANSI colour numbers.
var cgaToANSI = [8]int{0, 4, 2, 6, 1, 5, 3, 7}

func (c *Console) syncAttr() {
	if c.Attr == c.sent {
		return
	}
	a := c.Attr
	seq := "\x1b[0"
	if a&0x08 != 0 {
		seq += ";1"
	}
	if a&0x80 != 0 {
		seq += ";5"
	}
	seq += ";" + strconv.Itoa(30+cgaToANSI[a&7]) + ";" + strconv.Itoa(40+cgaToANSI[(a>>4)&7]) + "m"
	c.sess.Print(seq)
	c.sent = a
}

// Print writes text (CP437 bytes) in the current attribute.
func (c *Console) Print(s string) {
	if s == "" {
		return
	}
	c.syncAttr()
	c.sess.Print(s)
}

// Raw writes bytes without touching the attribute.
func (c *Console) Raw(s string) { c.sess.Print(s) }

// Clear clears the screen in the current attribute and homes the cursor.
func (c *Console) Clear() {
	c.syncAttr()
	c.sess.Print("\x1b[2J\x1b[H")
}

// ClearEOL clears to the end of the line.
func (c *Console) ClearEOL() {
	c.syncAttr()
	c.sess.Print("\x1b[K")
}

// GotoXY moves the cursor to 0-based column x, row y.
func (c *Console) GotoXY(x, y int) {
	c.sess.Print(fmt.Sprintf("\x1b[%d;%dH", y+1, x+1))
}

// MoveX moves the cursor n columns right (left when negative).
func (c *Console) MoveX(n int) {
	switch {
	case n > 0:
		c.sess.Print(fmt.Sprintf("\x1b[%dC", n))
	case n < 0:
		c.sess.Print(fmt.Sprintf("\x1b[%dD", -n))
	}
}

// Flush sends buffered output.
func (c *Console) Flush() { _ = c.sess.Flush() }

// keyName turns a decoded key into dorkit's key string.
func keyName(k door.Key) string {
	switch k {
	case door.KeyUp:
		return "KEY_UP"
	case door.KeyDown:
		return "KEY_DOWN"
	case door.KeyLeft:
		return "KEY_LEFT"
	case door.KeyRight:
		return "KEY_RIGHT"
	case door.KeyHome:
		return "KEY_HOME"
	case door.KeyEnd:
		return "KEY_END"
	case door.KeyInsert:
		return "KEY_INS"
	case door.KeyDelete:
		return "KEY_DEL"
	case door.KeyPgUp:
		return "KEY_PGUP"
	case door.KeyPgDn:
		return "KEY_PGDN"
	case door.KeyF1, door.KeyF2, door.KeyF3, door.KeyF4, door.KeyF5,
		door.KeyF6, door.KeyF7, door.KeyF8, door.KeyF9, door.KeyF10:
		return "KEY_F" + strconv.Itoa(int(k-door.KeyF1)+1)
	case door.KeyUnknown:
		return "KEY_UNKNOWN"
	}
	if k.IsByte() {
		return string([]byte{byte(k)})
	}
	return "KEY_UNKNOWN"
}

// WaitKey reports whether a key arrives within ms milliseconds
// (dk.console.waitkey).
func (c *Console) WaitKey(ms int) bool {
	if len(c.pending) > 0 || c.hangup {
		return true
	}
	d := time.Duration(ms) * time.Millisecond
	if d <= 0 {
		d = time.Millisecond
	}
	k, ok, err := c.sess.KeyWithin(d)
	if err != nil {
		if errors.Is(err, door.ErrTimeUp) {
			// lord.js checks the time itself (getkeyw); report a key so
			// it notices, and let GetKey return nothing.
			return false
		}
		c.hangup = true
		return true
	}
	if ok {
		c.pending = append(c.pending, keyName(k))
	}
	return ok
}

// GetKey returns a waiting key, or "" when none is waiting
// (dk.console.getkey returning undefined).
func (c *Console) GetKey() string {
	if len(c.pending) == 0 && !c.WaitKey(0) {
		return ""
	}
	if len(c.pending) > 0 {
		k := c.pending[0]
		c.pending = c.pending[1:]
		return k
	}
	if c.hangup {
		// dorkit's connection-closed handling ends the script.
		panic(exitSignal{0})
	}
	return ""
}

// GetStrOpts are dorkit's getstr options that lord.js uses.
type GetStrOpts struct {
	Len      int    // max length (default 80)
	Edit     string // initial value
	Integer  bool   // digits only
	InputBox bool   // draw the input box
	NoCRLF   bool   // crlf:false
	Select   bool   // select the edit value (typing replaces it)
	Attr     int    // input attribute; -1 = current
	Hotkeys  string
	Timeout  time.Duration
}

// GetStr is dk.console.getstr: a line editor.
func (c *Console) GetStr(o GetStrOpts) string {
	if o.Len == 0 {
		o.Len = 80
	}
	attr := c.Attr
	if o.Attr >= 0 && o.InputBox {
		attr = o.Attr
	}
	selAttr := attr&0x88 | (attr&0x07)<<4 | (attr>>4)&0x07
	str := o.Edit
	sel := o.Select && str != ""
	c.Attr = attr
	if o.InputBox {
		c.Print(strings.Repeat(" ", o.Len))
		c.MoveX(-o.Len)
	}
	if sel {
		c.Attr = selAttr
	}
	c.Print(str)
	c.Attr = attr
	pos := len(str)
	selectKeep := func() {
		if sel {
			sel = false
			c.MoveX(-pos)
			c.Print(str)
			pos = len(str)
		}
	}
	selectErase := func() {
		if sel {
			sel = false
			c.MoveX(-pos)
			c.Print(strings.Repeat(" ", len(str)))
			c.MoveX(-len(str))
			str = ""
			pos = 0
		}
	}
	for {
		wait := 1000
		if o.Timeout > 0 {
			wait = int(o.Timeout / time.Millisecond)
		}
		if !c.WaitKey(wait) {
			if o.Timeout > 0 {
				return str
			}
			continue
		}
		key := c.GetKey()
		if o.Hotkeys != "" && str == "" && len(key) == 1 && strings.Contains(o.Hotkeys, key) {
			if key[0] >= 32 {
				c.Print(key)
			}
			if !o.NoCRLF {
				c.Print("\r\n")
			}
			return key
		}
		switch key {
		case "KEY_HOME":
			selectKeep()
			c.MoveX(-pos)
			pos = 0
		case "KEY_END":
			selectKeep()
			c.MoveX(len(str) - pos)
			pos = len(str)
		case "KEY_LEFT":
			selectKeep()
			if pos > 0 {
				pos--
				c.MoveX(-1)
			}
		case "KEY_RIGHT":
			selectKeep()
			if pos < len(str) {
				pos++
				c.MoveX(1)
			}
		case "\x7f", "\b":
			selectErase()
			if pos == 0 {
				break
			}
			str = str[:pos-1] + str[pos:]
			pos--
			c.MoveX(-1)
			c.Print(str[pos:] + " ")
			c.MoveX(-1 - (len(str) - pos))
		case "KEY_DEL":
			selectErase()
			if pos >= len(str) {
				break
			}
			str = str[:pos] + str[pos+1:]
			c.Print(str[pos:] + " ")
			c.MoveX(-1 - (len(str) - pos))
		case "\r":
			if !o.NoCRLF {
				c.Print("\r\n")
			}
			return str
		case "":
		default:
			if len(key) > 1 || key[0] < 32 || key[0] > 127 {
				break
			}
			if o.Integer && (key[0] < '0' || key[0] > '9') {
				break
			}
			selectErase()
			if len(str) >= o.Len {
				break
			}
			str = str[:pos] + key + str[pos:]
			c.Print(str[pos:])
			pos++
			c.MoveX(-(len(str) - pos))
		}
	}
}
