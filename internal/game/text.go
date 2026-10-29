package game

import (
	"strconv"
	"strings"
	"time"
)

// ---- output ---------------------------------------------------------------

// sw outputs text without ` parsing and without a CRLF.
func (g *Game) sw(str string) {
	if str == "" {
		return
	}
	if str == "\r\n" {
		g.curcolnum = 1
	}
	g.con.Print(str)
	if str != "\r\n" {
		g.curcolnum += len(str)
	}
}

// sln outputs a line and its CRLF, pausing at a full screen when morechk
// is on.  *All* CRLFs are sent from here.
func (g *Game) sln(str string) {
	if str != "" {
		g.sw(str)
	}
	g.sw("\r\n")
	g.curlinenum++
	if g.morechk {
		if g.curlinenum > g.con.Rows-1 {
			g.moreNomail()
			g.curlinenum = 1
		}
	} else {
		g.curlinenum = 1
	}
}

// sclrscr clears the screen keeping the attribute.
func (g *Game) sclrscr() {
	oa := g.con.Attr
	g.con.Clear()
	g.con.Attr = oa
	g.curlinenum = 1
}

// foreground sets the foreground colour; values over 15 blink.
func (g *Game) foreground(col int) {
	if col > 15 {
		col = 0x80 | (col & 0x0f)
	}
	g.con.Attr = (g.con.Attr & 0x70) | col
}

// background sets the background colour (0-7).
func (g *Game) background(col int) {
	if col > 7 || col < 0 {
		return
	}
	g.con.Attr = (g.con.Attr & 0x8f) | (col << 4)
}

// lln outputs a line with ` codes; ext enables the mail action codes.
func (g *Game) lln(str string, ext ...bool) {
	g.lw(str, ext...)
	g.sln("")
}

// lw outputs text with LORD's ` colour and control codes.  With ext (mail
// lines), a code at the start of the line performs a mail action instead.
func (g *Game) lw(str string, extArg ...bool) {
	ext := len(extArg) > 0 && extArg[0]
	snip := ""
	for i := 0; i < len(str); i++ {
		if str[i] != '`' {
			snip += str[i : i+1]
			continue
		}
		g.sw(snip)
		snip = ""
		i++
		if i >= len(str) {
			break
		}
		switch str[i] {
		case '.':
		case 'n':
			if g.curcolnum <= g.con.Cols {
				g.sln("")
			}
		case '`':
			g.getkey()
		case 'l':
			g.foreground(2)
			g.sw("  -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			g.foreground(int(str[i] - '0'))
		case '0':
			g.foreground(10)
		case '!':
			g.foreground(11)
		case '@':
			g.foreground(12)
		case '#':
			g.foreground(13)
		case '$':
			g.foreground(14)
		case '%':
			g.foreground(15)
		case '^':
			g.foreground(0)
		case ')':
			g.foreground(20)
		case 'c':
			g.sclrscr()
			g.sln("")
			g.sln("")
		case 'C':
			g.sclrscr()
		case '|':
			fileRemove(g.quoteFile())
		case 'B':
			i++
			if i >= len(str) {
				break
			}
			switch c := str[i]; c {
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				g.foreground(16 + int(c-'0'))
			case '0':
				g.foreground(26)
			case '!':
				g.foreground(27)
			case '@':
				g.foreground(28)
			case '#':
				g.foreground(29)
			case '$':
				g.foreground(30)
			case '%':
				g.foreground(31)
			case '^':
				g.foreground(16)
			}
		case 'r':
			i++
			if i >= len(str) {
				break
			}
			if c := str[i]; c >= '0' && c <= '7' {
				g.background(int(c - '0'))
			}
		default:
			if i == 1 && ext {
				if g.lwMailCode(str) {
					return
				}
			}
		}
	}
	g.sw(snip)
}

// lwMailCode performs the mail action of a line starting with `X (the
// `default:` branch of lw).  It reports whether lw should stop there.
func (g *Game) lwMailCode(str string) bool {
	p := g.player
	arg, _ := jsParseInt(str[2:])
	switch str[1] {
	case '-':
		g.answerMail(arg)
	case 'b':
		p.Bank = min(p.Bank+arg, 2000000000)
	case 'E':
		p.Exp = min(p.Exp+arg, 2000000000)
		g.tournamentCheck()
	case '{':
		p.Laid++
		g.tournamentCheck()
	case 'G':
		p.Gold = min(p.Gold+arg, 2000000000)
	case ',':
		p.ForestFights = min(p.ForestFights+arg, 32000)
	case ':':
		p.PvPFights = min(p.PvPFights+arg, 32000)
	case ';':
		p.HP = min(p.HP+arg, 2000000000)
	case '+':
		p.Cha = arg
	case '}':
		p.Cha = min(p.Cha+1, 32000)
	case 'K', 'k':
		p.Kids = min(p.Kids+1, 32000)
	case 'V':
		p.Exp = max(p.Exp-arg, 1)
	case 'M':
		p.Str = min(p.Str+arg, 2000000000)
	case 'D':
		p.Def = min(p.Def+arg, 2000000000)
	case 'T':
		g.smileMail(arg)
	case 'Y':
		g.kissMail(arg)
	case 'U':
		g.dinnerMail(arg)
	case 'I':
		g.sleepMail(arg)
	case '?':
		to := arg
		if to > -1 {
			oop := g.playerGet(to)
			if oop.MarriedTo == int64(p.Record) {
				p.MarriedTo = to
				p.Put()
			}
		}
		fileRemove(g.quoteFile())
	case 'o':
		// TODO (lord.js): This is extra fragile since it's holding the mail lock...
		fileRemove(g.quoteFile())
		oop := g.playerGet(arg)
		if oop.OnNow {
			g.sclrscr()
			g.sln("")
			g.sln("")
			g.foreground(15)
			g.sln("                             ** ONLINE BATTLE **")
			g.foreground(10)
			g.sln("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-==-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
			p.Put()
			g.lln("  `0" + oop.Name + " `2has challenged you to an Online Duel!")
			g.sln("")
			g.lw("  Do you accept? [`%N`2] : ")
			ch := strings.ToUpper(g.getkey())
			if ch != "Y" {
				ch = "N"
			}
			g.sln(ch)
			if ch == "N" {
				g.obSendResp(oop, "R")
				g.obCleanup()
				return false
			}
			g.sln("")
			g.sln("  As is the tradition, you are allowed to completely heal yourself first.")
			g.sln("")
			g.fmutex(g.gamedir("fight"+strconv.Itoa(p.Record)+".rec"), "IN BATTLE\n")
			p.HP = p.HPMax
			g.foreground(15)
			if g.random(10) < 5 {
				g.onBattle(oop, true, "")
			} else {
				g.obSendResp(oop, "A")
				g.onBattle(oop, false, "")
			}
			g.obCleanup()
			if p.HP < 1 || p.Dead {
				p.HP = 0
				p.Dead = true
				p.Put()
				g.exit(0)
			}
		}
		return false
	case 'P':
		oop := g.playerGet(arg)
		g.sln("")
		g.getState(false)
		st := g.state
		if p.MarriedTo > -1 || st.MarriedToSeth == int64(p.Record) || st.MarriedToViolet == int64(p.Record) {
			g.foreground(15)
			g.sln("  BECAUSE YOU ARE ALREADY MARRIED, YOU ARE FORCED TO DECLINE!")
			g.sln("")
			g.mailTo(arg, "`%  ** Romantic Mail Return From "+p.Name+"`% **\n"+
				"`2-=-=-=-=-=-=-=-=-=-=--=-=-=-=-=-=-=-=-=-=-=-=-=--=-=-=-=-=-=-\n"+
				"`0  "+p.Name+"`2 is already married, and cannot say yes.")
			g.logLine("`%** `0" + p.Name + " `2 had to decline `0" + oop.Name + "`2's marriage proposal! `%**")
			g.moreNomail()
		}
		if oop.MarriedTo > -1 || st.MarriedToSeth == int64(oop.Record) || st.MarriedToViolet == int64(oop.Record) {
			g.lln("  `%" + oop.Name + " `0is already married!  No is the only answer.")
			g.sln("")
			// lord.js leaves out the "\n" after the rule here.
			g.mailTo(arg, "`%  ** Romantic Mail Return From "+p.Name+"`% **\n"+
				"`2-=-=-=-=-=-=-=-=-=-=--=-=-=-=-=-=-=-=-=-=-=-=-=--=-=-=-=-=-=-"+
				"`0  "+p.Name+" `2refuses to marry you!  You are already married!")
			g.logLine("`%** `0" + p.Name + " `2decides not to marry a married person. `%**")
			g.moreNomail()
		}
		g.marryMail(arg)
	case 'S':
		switch p.Clss {
		case 1:
			p.SkillW = min(p.SkillW+1, 40)
		case 2:
			p.SkillM = min(p.SkillM+1, 40)
		case 3:
			// BUG (lord.js): increments player.skillr, a property that
			// doesn't exist, so a thief gains nothing.  Kept.
		}
	default:
		return false
	}
	return true
}

// jsParseInt is parseInt(s, 10): the leading integer of s.  ok is false
// where JavaScript would produce NaN (the value is then 0).
func jsParseInt(s string) (int64, bool) {
	s = strings.TrimLeft(s, " \t\r\n")
	end := 0
	if end < len(s) && (s[end] == '-' || s[end] == '+') {
		end++
	}
	start := end
	for end < len(s) && s[end] >= '0' && s[end] <= '9' {
		end++
	}
	if end == start {
		return 0, false
	}
	n, err := strconv.ParseInt(s[:end], 10, 64)
	if err != nil {
		return 0, false
	}
	return n, true
}

// quoteFile is gamedir('quote'+player.Record+'.lrd').
func (g *Game) quoteFile() string {
	return g.gamedir("quote" + strconv.Itoa(g.player.Record) + ".lrd")
}

// ---- input ------------------------------------------------------------------

// getkeyw waits for a key, enforcing the BBS time limit (less 30 seconds)
// and a five-minute idle limit.
func (g *Game) getkeyw() string {
	timeout := g.time() + 60*5
	for {
		now := g.time()
		tl := g.secondsRemaining() - 30
		if tl < 1 {
			if g.Settings.NoTime != "" {
				g.displayFile(g.Settings.NoTime, false, true)
			} else {
				g.sln("(*** GOD STRIKES YOU UNCONSCIOUS (OUT OF BBS TIME!) ***)")
			}
			if g.player != nil {
				g.player.OnNow = false
				g.player.Put()
			}
			g.exit(0)
		}
		if now >= timeout {
			g.sln("")
			g.sln("")
			g.sln("  Ack!  Apperently you didn't find this door very exciting, because")
			g.sln("  you have obviously fallen asleep.  Please come back sometime when")
			g.sln("  you feel like actually playing.")
			g.sln("")
			if g.player != nil {
				g.player.OnNow = false
				g.player.Put()
			}
			g.exit(0)
		}
		if g.con.WaitKey(1000) {
			break
		}
	}
	return g.con.GetKey()
}

// secondsRemaining is the caller's remaining time in seconds
// (dk.user.seconds_remaining + seconds_remaining_from - time()).
func (g *Game) secondsRemaining() int64 {
	return int64(g.sess.TimeLeft() / time.Second)
}

// getkey reads a key, turning control and high characters into '.'.
func (g *Game) getkey() string {
	g.curlinenum = 1
	for {
		ch := g.getkeyw()
		if len(ch) != 1 {
			continue
		}
		a := ch[0]
		switch {
		case a == 0:
			continue
		case a >= 1 && a <= 7, a >= 9 && a <= 12, a >= 14 && a <= 26, a >= 127:
			ch = "."
		}
		return ch
	}
}

// getstr is lord.js's getstr(x, y, len, c, c1, str, mode): an input box at
// (x, y) (1-based; 0,0 = here) in colours c (background) and c1.
func (g *Game) getstr(x, y, length, c, c1 int, str string, mode GetStrOpts) string {
	oa := g.con.Attr
	mode.Len = length
	mode.Edit = str
	mode.Attr = c<<4 | c1
	mode.InputBox = true
	mode.NoCRLF = true
	if str != "" {
		mode.Select = true
	}
	g.curlinenum = 1
	if x != 0 && y != 0 {
		g.con.GotoXY(x-1, y-1)
	}
	g.background(c)
	g.foreground(c1)
	ret := g.con.GetStr(mode)
	g.con.Attr = oa
	return ret
}

// consoleGetstr is dk.console.getstr(opts) with dorkit's defaults (Attr
// -1 = current, CRLF after input).
func (g *Game) consoleGetstr(o GetStrOpts) string {
	g.curlinenum = 1
	if o.Attr == 0 {
		o.Attr = -1
	}
	return g.con.GetStr(o)
}

// ---- more prompts ----------------------------------------------------------

// moreNomail is more_nomail(): a <MORE> prompt that doesn't check mail.
func (g *Game) moreNomail() {
	oa := g.con.Attr
	g.curlinenum = 1
	g.lw("  `2<`0MORE`2>")
	g.getkey()
	if !g.Settings.UseFancyMore {
		g.sln("")
		g.curlinenum = 1
	} else {
		g.con.Print("\r")
		g.con.ClearEOL()
	}
	g.con.Attr = oa
}

// more saves the player, delivers mail, then shows a <MORE> prompt.
func (g *Game) more() {
	oa := g.con.Attr
	g.curlinenum = 1
	g.player.Put()
	g.checkMail()
	g.lw("  `2<`0MORE`2>")
	g.getkey()
	if !g.Settings.UseFancyMore {
		g.sln("")
		g.curlinenum = 1
	} else {
		g.con.Print("\r")
		g.con.ClearEOL()
	}
	g.con.Attr = oa
}

// ---- files on screen --------------------------------------------------------

// displayFile shows a text file with ` codes; canskip (default true) lets
// ^C or space stop it, more (default true) pauses each screen.
func (g *Game) displayFile(fname string, canskip, more bool) {
	lines, ok := readLines(fname)
	if !ok {
		g.fatalf("Unable to open %s", fname)
	}
	ms := g.morechk
	g.morechk = false
	lc := 0
loop:
	for _, ln := range lines {
		lc++
		if more && lc%(g.con.Rows-1) == 0 {
			g.moreNomail()
		}
		g.lln(ln, true)
		if canskip {
			for {
				ch := g.con.GetKey()
				if ch == "\x03" || ch == " " {
					g.sln("")
					break loop
				}
				if ch == "" {
					break
				}
			}
		}
	}
	g.morechk = ms
}

// showLordFile shows a mail file, optionally copying it to the quote file.
func (g *Game) showLordFile(fname string, quote, mail bool) {
	lines, ok := readLines(fname)
	if !ok {
		return
	}
	g.showLordLines(lines, quote, mail)
}

// showLordBuffer is show_lord_buffer(): showLordFile for text in memory.
func (g *Game) showLordBuffer(buf string, quote, mail bool) {
	g.showLordLines(strings.Split(strings.ReplaceAll(buf, "\r\n", "\n"), "\n"), quote, mail)
}

func (g *Game) showLordLines(lines []string, quote, mail bool) {
	ms := g.morechk
	g.morechk = false
	lc := 0
	ns := false
	for _, ln := range lines {
		if quote {
			if err := appendLines(g.quoteFile(), ln); err != nil {
				g.fatalf("Unable to open %s", g.quoteFile())
			}
		}
		g.lln(ln, mail)
		lc++
		if !ns && lc >= g.con.Rows-1 {
			g.lw(" `2(`5C`2)ontinue, (`5S`2)top, (`5N`2)onstop `0: `%")
			ch := strings.ToUpper(g.getkey())
			lc = 0
			if ch == "S" {
				g.con.Print("\r")
				g.con.ClearEOL()
				break
			}
			if ch == "N" {
				ns = true
			}
			g.con.Print("\r")
			g.con.ClearEOL()
		}
	}
	g.morechk = ms
}

// ---- string helpers ------------------------------------------------------------

// prettyInt formats n with thousands separators.
func prettyInt(n int64) string {
	ret := strconv.FormatInt(n, 10)
	for i := len(ret) - 3; i > 0; i -= 3 {
		if i == 1 && ret[0] == '-' {
			break
		}
		ret = ret[:i] + "," + ret[i:]
	}
	return ret
}

// looks is looks(): the line of mlooks.txt/flooks.txt for a charm value.
func (g *Game) looks(sex string, charm int64) string {
	name := "flooks.txt"
	if sex == "M" {
		name = "mlooks.txt"
	}
	lines, err := g.assetLines(name)
	if err != nil {
		g.fatalf("Unable to open %s", name)
	}
	if charm < 1 {
		charm = 1
	}
	if charm > 100 {
		charm = 101
	}
	if int(charm) > len(lines) {
		return ""
	}
	return lines[charm-1]
}

// dispLen is the displayed length of a string with ` codes.
func dispLen(str string) int {
	n := 0
	for i := 0; i < len(str); i++ {
		if str[i] == '`' {
			i++
			if i < len(str) && str[i] == 'r' {
				i++
			}
		} else {
			n++
		}
	}
	return n
}

// dispStr removes ` codes.
func dispStr(str string) string {
	var b strings.Builder
	for i := 0; i < len(str); i++ {
		if str[i] == '`' {
			i++
			if i < len(str) && str[i] == 'r' {
				i++
			}
		} else {
			b.WriteByte(str[i])
		}
	}
	return b.String()
}

// center centres a string with ` codes on an 80-column line.
func center(str string) string {
	n := (80 - dispLen(str)) / 2
	n = max(0, min(40, n))
	return strings.Repeat(" ", n) + str
}

// spacePad pads str to len displayed columns (on the left with left).
func spacePad(str string, length int, left ...bool) string {
	l := len(left) > 0 && left[0]
	for dispLen(str) < length {
		if l {
			str = " " + str
		} else {
			str += " "
		}
	}
	return str
}

// cleanStr removes every ` code except colours.
func cleanStr(str string) string {
	var b strings.Builder
	for i := 0; i < len(str); i++ {
		if str[i] != '`' {
			b.WriteByte(str[i])
			continue
		}
		var next byte
		if i+1 < len(str) {
			next = str[i+1]
		}
		if strings.IndexByte("0123456789!@#$%", next) >= 0 && next != 0 {
			b.WriteByte(str[i])
		} else {
			i++
		}
	}
	// TODO (lord.js): This is where badwords replacements should happen...
	return b.String()
}

// lordToAnsi converts ` codes to ANSI escape sequences.
func lordToAnsi(str string) string {
	var b strings.Builder
	codes := map[byte]string{
		'c': "\x1b[2J\x1b[H", '1': "\x1b[0;34m", '*': "\x1b[0;30;41m", '2': "\x1b[0;32m",
		'3': "\x1b[0;36m", '4': "\x1b[0;31m", '5': "\x1b[0;35m", '6': "\x1b[0;33m",
		'7': "\x1b[0;37m", '8': "\x1b[1;30m", '9': "\x1b[1;34m", '0': "\x1b[1;32m",
		'!': "\x1b[1;36m", '@': "\x1b[1;31m", '#': "\x1b[1;35m", '$': "\x1b[1;33m",
		'%': "\x1b[1;37m",
	}
	for i := 0; i < len(str); i++ {
		if str[i] != '`' {
			b.WriteByte(str[i])
			continue
		}
		i++
		if i >= len(str) {
			break
		}
		if s, ok := codes[str[i]]; ok {
			b.WriteString(s)
		} else if str[i] == 'r' {
			// lord.js compares the character with numbers here, which
			// never matches, so `rX produces nothing.
			i++
		}
	}
	return b.String()
}

// boolStr helps translate JavaScript ternaries on a player's sex.
func pick[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}
