package game

import (
	"regexp"
	"strings"
)

// warriorsOnNow is warriors_on_now(): who is playing right now.  inhello
// tells the others that the player has arrived.
func (g *Game) warriorsOnNow(inhello bool) {
	var on []*Player

	g.sln("")
	g.foreground(15)
	g.background(0)
	g.sln("                         Warriors in the Realm Now")
	g.lln("`2-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
	for _, op := range g.allPlayers() {
		if op.Name != "X" && op.OnNow {
			on = append(on, op)
		}
	}
	for _, op := range on {
		// DIFF: lord.js shows where a player is from out<n>.lrd, which only
		// IGMs write; IGMs are not supported, so it's always the arrival
		// time.
		g.lln("  `0" + spacePad(op.Name, 25) + "   `2Arrived At`%                   " + op.TimeOn)
		if inhello && op.Record != g.player.Record {
			g.mailTo(int64(op.Record), "`0  "+g.player.Name+" `2has entered the realm.")
		}
	}
}

// showLog is show_log(): today's (and yesterday's) happenings.
func (g *Game) showLog() {
	var ch string

	g.createLog(false)

	doLog := func(today bool) {
		if !today {
			g.showLordFile(g.gamedir("logold.lrd"), false, false)
		} else {
			g.showLordFile(g.gamedir("lognow.lrd"), false, false)
		}
	}

	for {
		g.sclrscr()
		g.sln("")
		g.sln("")
		doLog(true)
		for {
			g.sln("")
			g.foreground(2)
			g.lw("`2  (`5C`2)ontinue   (`5T`2)odays happenings again  (`5Y`2)esterdays `0[`5C`0] : ")
			ch = strings.ToUpper(g.getkey())
			if !strings.Contains("TY", ch) {
				ch = "C"
			}
			g.sw(ch)
			if ch == "T" {
				break
			}
			if ch == "Y" {
				if !fileExists(g.gamedir("logold.lrd")) {
					g.sln("")
					g.sln("")
					g.sln("  Apparently nothing of importance happened yesterday.")
				} else {
					g.sclrscr()
					g.sln("")
					g.sln("")
					doLog(false)
				}
			}
			if ch == "C" {
				break
			}
		}
		if ch == "C" {
			break
		}
	}
	g.sln("")
}

// showStats is show_stats(): the player's stats screen.
func (g *Game) showStats() {
	p := g.player

	g.sclrscr()
	g.sln("")
	g.lln("`%  " + p.Name + "`2's Stats...")
	g.lln("`0-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
	g.lln("  `2Experience   : `0" + prettyInt(p.Exp))
	g.lln("`2  Level        : `0" + spacePad(prettyInt(p.Level), 17) + " `2HitPoints          : `0" + prettyInt(p.HP) + "`2 of `0" + prettyInt(p.HPMax))
	g.lln("`2  Forest Fights: `0" + spacePad(prettyInt(p.ForestFights), 17) + " `2Player Fights Left : `0" + prettyInt(p.PvPFights))
	g.lln("`2  Gold In Hand : `0" + spacePad(prettyInt(p.Gold), 17) + " `2Gold In Bank       : `0" + prettyInt(p.Bank))
	g.lln("`2  Weapon       : `0" + spacePad(p.Weapon, 17) + " `2Attack Strength    : `0" + prettyInt(p.Str))
	g.lln("`2  Armour       : `0" + spacePad(p.Arm, 17) + " `2Defensive Strength : `0" + prettyInt(p.Def))
	g.lln("`2  Charm        : `0" + spacePad(prettyInt(p.Cha), 17) + " `2Gems               : `0" + prettyInt(p.Gem))
	g.sln("")
	g.getState(false)
	if g.state.MarriedToViolet == int64(p.Record) {
		g.lln("  `2You are married to `#Violet`2.")
		g.sln("")
	}
	if g.state.MarriedToSeth == int64(p.Record) {
		g.lln("  `2You are married to `%Seth Able`2.")
		g.sln("")
	}
	if p.MarriedTo != -1 {
		// DIFF: lord.js tests for undefined, but player_get returns null for
		// a missing record and the next line throws; skip it instead.
		if tp := g.playerGet(p.MarriedTo); tp != nil {
			g.lln("  `2You are married to `%" + tp.Name + "`2.")
			g.sln("")
		}
	}
	if p.Kids == 1 {
		g.lln("  `2You have `01`2 child.")
		g.sln("")
	}
	if p.Kids > 1 {
		g.lln("  `2You have `0" + prettyInt(p.Kids) + "`2 children.")
		g.sln("")
	}
	if p.Horse {
		g.lln("  `2You are on `%horseback`2.")
	}
	if p.HasFairy {
		g.lln("  `2You have a `#fairy`2 in your pocket.")
	}
	skill := func(label string, level, sk int64) {
		if level > 0 || sk > 0 {
			g.lw(label)
			if sk == 0 {
				g.sw("NONE       ")
			} else if sk >= 40 {
				g.sw("MASTERED   ")
			} else if sk > 0 {
				g.sw(spacePad(prettyInt(sk), 11))
			}
			g.lln("`2 Uses Today: (`0" + prettyInt(level) + "`2)")
		}
	}
	skill("`2  Death Knight Skills: `0", p.LevelW, p.SkillW)
	skill("`2  The Mystical Skills: `0", p.LevelM, p.SkillM)
	skill("`2  The Thieving Skills: `0", p.LevelT, p.SkillT)
	switch p.Clss {
	case 1:
		g.lln("  `0You are currently interested in `%Death Knight`0 skills.")
	case 2:
		g.lln("  `0You are currently interested in `%The Mystical`0 skills.")
	case 3:
		g.lln("  `0You are currently interested in `%The Thieving`0 skills.")
	}
	if p.Amulet {
		g.lln("  `2You are wearing an `%Amulet of Accuracy`2.")
		g.sln("")
	}
	g.more()
}

var digitsRE = regexp.MustCompile(`^[0-9]+$`)

// hello is hello(): the player enters the realm.
func (g *Game) hello() {
	p := g.player
	no := g.gamedir("nodeon." + itoa(g.Ctx.Node))

	fileRemove(g.gamedir("quote" + itoa(p.Record) + ".lrd"))
	fileRemove(g.gamedir("msg" + itoa(p.Record) + ".lrd"))
	fileRemove(g.gamedir("temp" + itoa(p.Record)))

	if g.Settings.SafeNode {
		if fileExists(no) {
			if lines, ok := readLines(no); ok {
				tmp := ""
				if len(lines) > 0 {
					tmp = lines[0]
				}
				if digitsRE.MatchString(tmp) {
					n, _ := jsParseInt(tmp)
					if n != int64(p.Record) {
						if op := g.playerGet(n); op != nil {
							op.OnNow = false
							op.Put()
						}
					}
				}
			}
		}
		if writeLines(no, itoa(p.Record)) == nil {
			g.cleanupFiles = append(g.cleanupFiles, no)
		}
	}

	g.sln("")
	// NOTE: Likely wrap checks... not needed.
	if p.Bank < 0 {
		p.Bank = 0
	}
	if p.Gold < 0 {
		p.Gold = 0
	}
	if p.Exp < 0 {
		p.Exp = 0
	}
	p.OnNow = true
	p.Put()

	p.TimeOn = g.now().Format("03:04") // strftime('%I:%M', time())
	p.Put()
	g.warriorsOnNow(true)
	g.sln("")
	g.createLog(false)
	if p.Time == g.state.Days {
		g.foreground(2)
		g.sln("")
		g.sln("  You have already been on today.")
		if p.Inn {
			g.sln("  You wake up early, sheathe your weapon, and head down to the bar.")
		}
		if p.Dead || p.HP < 1 {
			p.Dead = true
			g.lln("  `2You have been `4killed`2 today.  You will not be allowed")
			g.sln("  to play until TOMORROW.")
			p.OnNow = false
			p.Put()
		}
		g.sln("")
	} else {
		g.wakeUp()
	}
	if p.MarriedTo > -1 {
		op := g.playerGet(p.MarriedTo)
		// DIFF: lord.js throws when the spouse's record doesn't exist; treat
		// that like a deleted spouse.
		if op == nil || op.Name == "X" {
			g.sln("")
			g.lln("  `2Your lover has been declared MISSING!")
			g.sln("")
			g.lln("  `%YOU FEEL HEARTBROKEN!  (CHARM DROPS 50%)")
			g.sln("")
			p.MarriedTo = -1
			p.Cha = p.Cha / 2
			g.moreNomail()
		}
	}
	g.moreNomail()
	g.showStats()
	if g.mailCheck() {
		g.sln("")
		g.lln("  `%You have mail today!")
		g.sln("")
		g.moreNomail()
		g.checkMail()
	}
	g.showLog()
	if p.Clss == 0 {
		g.sln("")
		g.sln("")
		g.sln("  ** You feel listless, you need a goal in life **")
		p.Clss = g.chooseProfession(false)
	}
	p.Put()
}

// newPlayer is new_player(): create a character for a new caller.
func (g *Game) newPlayer() {
	var name, ch string
	var npName, npSex string
	var npClss int64

	// np_to_player copies the new player's choices into the record.
	npToPlayer := func() {
		g.player.Name = npName
		g.player.Sex = npSex
		g.player.Clss = npClss
	}

	g.sclrscr()
	if g.playerLength() > 147 {
		g.sln("")
		g.lln("`c  `%                      ** THE REALM IS FULL **")
		g.lln("`0-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
		g.lln("  `2An over-population problem is keeping you from moving in.")
		g.sln("")
		g.lln("  `0Take heart. `2 Although `%LORD " + Version + "`2 is NEVER reset, players who do not")
		g.sln("  visit for a certain number of days are deleted every morning.")
		g.sln("")
		g.sln("  (You may want to ask your Sysop to make the # of days less)")
		g.sln("")
		g.sln("  Please try again tomorrow.  Until then, you can only dream of the")
		g.sln("  wonderful world that awaits...")
		g.sln("")
		g.moreNomail()
		g.exit(0)
	}
	g.sln("")
	g.sln("")
	g.sln("                      ** Welcome to the realm new warrior! **")
	ch = ""
	for {
		g.sln("")
		g.sln("")
		g.lln("  `0What would you like as an alias? ")
		g.lln("")
		g.lw("  `2Name: `%")
		name = g.consoleGetstr(GetStrOpts{Len: 20})
		if name == "" {
			g.sln("")
			g.lln("  `)YOU'RE NOTHING BUT A WEENIE!")
			g.sln("")
			if g.player != nil {
				g.player.OnNow = false
				g.player.Put()
			}
			g.exit(0)
		}
		g.foreground(10)
		name = strings.TrimSpace(name)
		name = cleanStr(name)
		npName = name
		if len(name) > 18 {
			g.sln("  Try a shorter name!")
			g.sln("")
		} else if len(name) < 3 {
			g.sln("  Try a longer name!")
			g.sln("")
		} else if dispLen(name) < 3 {
			g.sln("  You are such a sneaky boy.  Yet so simplistic.")
			g.sln("")
		} else {
			plen := g.playerLength()
			i := 0
			for ; i < plen; i++ {
				tp := g.playerGet(int64(i))
				if dispStr(tp.Name) == dispStr(name) {
					g.sln("")
					g.lln("  `2You recall hearing that another warrior currently goes by that name.  Not")
					g.lln("  wanting to dishonor `0" + tp.Name + " `2you decide to pick another.")
					g.sln("")
					break
				}
			}
			if i == plen {
				if g.checkName(name) {
					g.lw("  `0" + name + "`2? `2[`0Y`2] : `%")
					ch = strings.ToUpper(g.getkey())
					if ch != "N" {
						ch = "Y"
					}
					g.foreground(15)
					g.sln(ch)
					g.foreground(10)
				}
			}
		}
		if ch == "Y" {
			break
		}
	}
	// TODO: Race condition in name... it's checked for dupes when you enter it, not at creation.
	g.lw("  `2And your gender?  (`0M`2/`0F`2) [`0M`2]: `%")
	ch = strings.ToUpper(g.getkey())
	if ch != "F" {
		ch = "M"
	}
	g.sln(ch)
	npSex = ch
	if ch == "M" {
		switch g.random(5) {
		case 0:
			g.lln("  With a name like \"" + npName + "`0\", no one is going to believe it.")
		case 1:
			g.sln("  ALL RIGHT!!  A member of the more ADVANCED sex.  You had better win.")
		case 2:
			g.sln("  Good.  Men rule this earth.  We own and run EVERYTHING.")
		case 3:
			g.sln("  Then don't be wearing any dresses, eh.")
		case 4:
			g.sln("  Very good.  If a woman ever beats you in battle, go into exile.")
		}
	} else {
		switch g.random(5) {
		case 0:
			g.sln("  Good.  Teach those men that they do NOT rule the world.")
		case 1:
			g.sln("  ALL RIGHT!!  A member of the more ADVANCED sex.  You had better win.")
		case 2:
			g.sln("  Excellent.  Taunt the men, tease them, and break their hearts!")
		case 3:
			g.sln("  Be warned, you are going to have to fight, kill and maim here.")
		case 4:
			g.sln("  Good.  There are way too many men in this land..")
		}
	}
	g.sln("")
	npClss = g.chooseProfession(false)
	plen := g.playerLength()
	i := 0
	for ; i < plen; i++ {
		g.player = g.playerGetLocked(int64(i))
		if g.player != nil && g.player.Name == "X" {
			g.player.ReInit()
			npToPlayer()
			// TODO: This should be user.number since real name isn't unique!
			g.player.RealName = g.Ctx.PlayerKey()
			g.player.Put()
			g.player.UnLock()
			break
		}
		if g.player != nil {
			g.player.UnLock()
		}
	}
	if i == plen {
		g.player = g.playerNew(true)
		if g.player == nil {
			g.fatalf("Unable to create a new player record")
		}
		npToPlayer()
		// TODO: This should be user.number since real name isn't unique!
		g.player.RealName = g.Ctx.PlayerKey()
		g.player.Put()
		g.player.UnLock()
	}
	if g.Settings.ResDays > 0 {
		g.player.Gone = 0
	} else {
		g.player.Gone = -1
	}
	g.killmail(g.player.Record)
	g.getState(true)
	// Can this even happen?  Isn't this done when the player is deleted?
	if g.state.MarriedToSeth == int64(g.player.Record) {
		g.state.MarriedToSeth = -1
	}
	if g.state.MarriedToViolet == int64(g.player.Record) {
		g.state.MarriedToViolet = -1
	}
	g.putState()
	g.hello()
}

// loadPlayer is load_player(): find the caller's character; with create,
// offer to make one.
func (g *Game) loadPlayer(create bool) {
	key := g.Ctx.PlayerKey()

	plen := g.playerLength()
	var found *Player
	for i := 0; i < plen; i++ {
		p := g.playerGet(int64(i))
		// TODO: This should be user.number since real name isn't unique!
		if p != nil && p.RealName == key {
			found = p
			break
		}
	}
	// DIFF: lord.js leaves `player` set to the last record it looked at
	// when there is no match, so later code (the WEENIE exit in
	// new_player, check_gameover, the exit handler) marks some other
	// warrior offline.  The player stays unset here instead.
	g.player = found
	if found == nil {
		if create {
			g.sclrscr()
			g.sln("")
			g.sln("")
			g.foreground(15)
			g.sln("  Joining The Game")
			g.foreground(10)
			g.sln("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
			g.foreground(2)
			g.sln("  You have never visited the realm before.  Do you want to join this place")
			g.sln("  of Dragons, Knights, Magic, Friends & Foes and Good & Evil?")
			g.sln("")
			g.lln("  `2(`0Y`2)es.")
			g.lln("  `2(`0N`2)o.")
			g.sln("")
			g.foreground(10)
			g.lw("  `2Your choice?  `2[`0Y`2] : `%")
			ch := strings.ToUpper(g.getkey())
			if ch != "N" {
				ch = "Y"
			}
			g.sln(ch)
			if ch != "Y" {
				g.quitMsg()
				g.exit(0)
			}
			g.newPlayer()
		}
		return
	}

	p := g.player
	p.OnNow = true
	p.Put()
	f := g.gamedir("battle" + itoa(p.Record) + ".rec.lock")
	if fileExists(f) {
		if lines, ok := readLines(f); ok && len(lines) > 0 {
			n, _ := jsParseInt(lines[0])
			op := g.playerGet(n)
			// DIFF: lord.js throws if the lock names no player.
			if op != nil {
				g.lln("  `%" + op.Name + "`2 is currently battling you")
				for {
					g.mswait(500)
					p.ReLoad()
					op.ReLoad()
					g.sw("\r")
					g.con.Print("\r")
					g.con.ClearEOL()
					g.lw("Your HP: " + prettyInt(p.HP) + "\tHis HP: " + prettyInt(op.HP))
					if !fileExists(f) {
						break
					}
				}
			}
		}
	}
	g.hello()
}

// showAlone is show_alone().
func (g *Game) showAlone() {
	alone := true

	for _, tp := range g.allPlayers() {
		if tp.OnNow {
			alone = false
		}
	}
	if !alone {
		g.lln(center("`0You are `%NOT ALONE`0 here.  `2(Someone else is playing also)"))
	}
}

// introMenu is intro_menu(): the title-screen menu.
func (g *Game) introMenu() string {
	var ch string

	for {
		g.sln("")
		g.lln("                         `2(`0E`2)nter the realm of the Dragon")
		g.lln("                         `2(`0I`2)nstructions")
		g.lln("                         `2(`0L`2)ist Warriors")
		g.lln("                         `2(`0Q`2)uit back to BBS")
		for {
			g.sln("")
			g.lw("                         `2Your choice, warrior? [`0E`2]: `%")
			ch = strings.ToUpper(g.getkey())
			if !strings.Contains("SE\rILQ1234567890", ch) {
				g.sln("")
				g.sln("")
				g.lln(center("`2Try entering a valid option, it might actually do something."))
				continue
			}
			break
		}
		g.sln(ch)
		switch ch {
		case "1":
			g.lrdfile("INTRO1")
		case "2":
			g.lrdfile("DRAG")
		case "3":
			g.lrdfile("INTRO2")
		case "4":
			g.lrdfile("DRAGON3")
		case "5":
			g.lrdfile("SM-LORD")
		case "6":
			g.lrdfile("LORD")
		case "7":
			g.lrdfile("FOOT")
		case "8":
			g.lrdfile("DEMON")
		case "9":
			g.lrdfile("LONGINTRO")
		case "0":
			g.lrdfile("ACCESSD")
		case "Q":
			g.quitMsg()
		case "S":
			g.lln("`c  `2You have pressed the `0S`2 key.  `0S`2 is for story.")
			g.sln("")
			g.moreNomail()
			g.lrdfile("STORY", true)
			g.moreNomail()
			g.sln("")
		case "L":
			if !fileExists(g.gamedir("player.bin")) {
				g.sln("")
				g.sln("")
				g.lln(center("The game has never been played before."))
			} else {
				tf := g.gamedir("tempn" + itoa(g.Ctx.Node))
				g.generateRankings(tf, true, true, false)
				g.sclrscr()
				g.displayFile(tf, true, true)
				fileRemove(tf)
				g.sln("")
				g.sln("")
			}
		case "I":
			g.instructions()
		}
		if strings.Contains("E\rQ", ch) {
			break
		}
	}
	return ch
}

// checkGameover is check_gameover(): the game has been won; show the
// winner and leave.
func (g *Game) checkGameover() {
	g.getState(false)
	if g.state.WonBy >= 0 {
		wb := g.playerGet(g.state.WonBy)
		g.loadPlayer(false)
		if g.player == nil {
			// DIFF: lord.js uses {Record:-1}, whose missing put() throws
			// below; this player's Put does nothing.
			g.player = &Player{Record: -1}
		}
		g.lrdfile("WIN")
		g.moreNomail()
		g.lln("`c                      `%PAY HOMAGE TO YOUR BETTER!     ")
		g.sln("")
		g.lln("  `0The incredible warrior whose deeds will grace every tongue of")
		g.sln("  every minstrel in every town every song of every day is the master")
		wbName := ""
		if wb != nil {
			wbName = wb.Name
		}
		g.lln("  warrior known as `%" + wbName + "`2.  ")
		g.sln("")
		g.moreNomail()
		if g.state.WonBy == int64(g.player.Record) {
			g.lln("  `0You smile modestly.  If only people knew, that incredible")
			g.sln("  warrior was you.")
		} else {
			g.lln("  `2You bow your head in reverence - vowing to follow the teachings")
			g.sln("  of this great person - to learn whatever this Godlike wonder can")
			g.sln("  show you.")
			g.sln("")
			g.lln("  `#(ASK YOUR SYSOP TO USE THE (R)ESET OPTION IN LORDCFG.EXE)`%")
		}
		g.sln("")
		g.moreNomail()
		g.player.OnNow = false
		g.player.Put()
		g.exit(0)
	}
}

// start is start(): the door's entry point after cmdline().
func (g *Game) start() {
	g.buildTxtIndex()

	// Should this allow a global file?
	if hf := g.gamedir("hello.lrd"); fileExists(hf) {
		g.displayFile(hf, false, false)
	}

	g.sclrscr()

	switch g.random(10) {
	case 0:
		g.lrdfile("INTRO1")
	case 1:
		g.lrdfile("DRAG")
	case 2:
		g.lrdfile("INTRO2")
	case 3:
		g.lrdfile("DRAGON3")
	case 4:
		g.lrdfile("SM-LORD")
	case 5:
		g.lrdfile("LORD")
	case 6:
		g.lrdfile("FOOT")
	case 7:
		g.lrdfile("DEMON")
	case 8:
		g.lrdfile("LONGINTRO")
	case 9:
		g.lrdfile("ACCESSD")
	}
	g.getkey()
	g.sclrscr()
	g.sln("")
	g.sln("")
	g.lln(center("`0L`2egend `0O`2f The `0R`2ed `0D`2ragon"))
	g.sln("")
	g.getState(false)
	g.lln(center("`2The current game has been running `0" + itoa(g.state.Days) + "`2 days."))
	g.lln(center("`2Players are deleted after `0" + itoa(g.Settings.DeleteDays) + "`2 days of inactivity."))
	g.showAlone()

	if g.introMenu() == "Q" {
		if g.player != nil {
			g.player.OnNow = false
			g.player.Put()
		}
		g.exit(0)
	}

	g.checkGameover()
	g.loadPlayer(true)
	// TODO: LOCKOUT.DAT support

	// js.on_exit(...): the exit handler saves the player as offline.
	g.onExitArmed = true

	// DIFF: lord.js would throw here if load_player(true) left no player;
	// it can't (declining exits), but guard anyway.
	if g.player != nil && g.player.OnNow {
		if !g.player.Dead && g.player.HP > 0 {
			// IGM resume (out<n>.lrd) is not supported.
			g.main()
		}

		g.player.OnNow = false
		g.player.Put()
		g.sclrscr()
		g.sln("")
		g.sln("")
		g.generateRankings(g.gamedir("scores.lrd"), true, false, false)
		g.sln("  Leaving the realm.")
	}
}
