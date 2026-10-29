package game

import "strings"

// announce is announce(): put a message in the daily happenings.
func (g *Game) announce() {
	p := g.player
	g.sln("")
	g.sln("")
	g.foreground(2)
	g.sln("  Are you sure you want to announce something?  It will appear to")
	g.sln("  EVERYONE in the daily happenings.")
	g.sln("")
	g.lw("  Make Announcement? [`0Y`2] :`% ")
	ch := strings.ToUpper(g.getkey())
	if ch != "N" {
		ch = "Y"
	}
	g.sln(ch)
	if ch == "N" {
		return
	}

	msg := " \n  `0" + p.Name + "`2 Announces:`%"
	g.sln("")
	g.sln("  Enter message now..Blank line quits!")
	for {
		g.lw(" `2>`% ")
		l := "  `%" + cleanStr(g.consoleGetstr(GetStrOpts{Len: 75}))
		if l == "  `%" {
			break
		}
		msg = msg + "\n" + l
	}
	g.logLine(msg)
	g.sln("")
	g.sln("  Announcement Made!")
}

// singleViolet is single_violet(): a man flirting with an unmarried
// Violet.
func (g *Game) singleViolet() {
	p := g.player
	g.sln("")
	g.lrdfile("VIOLET")
	for {
		g.sln("")
		g.lw("  `0`2Your choice?  (`0? for menu`2) : ")
		ch := strings.ToUpper(g.getkey())
		if ch == "\r" {
			ch = "R"
		}
		g.sln(ch)
		g.sln("")
		switch ch {
		case "?":
			g.sln("")
			g.lrdfile("VIOLET")
		case "N", "R", "Q":
			return
		case "K":
			g.singleVioletKiss()
		case "W":
			g.singleVioletWink()
		case "P":
			g.singleVioletPeck()
		case "S":
			g.singleVioletSit()
		case "G":
			g.singleVioletGrab()
		case "C":
			if g.Settings.CleanMode {
				g.sln("")
				g.sln("")
				g.sln("  Sorry, your sysop has disabled that function!")
				g.sln("")
			} else {
				g.singleVioletCarry()
			}
		case "M":
			g.getState(false)
			if p.MarriedTo > -1 {
				op := g.playerGet(p.MarriedTo)
				g.sln("")
				g.lln("  `0You kind of doubt you wife `)" + op.Name + " `0would")
				g.sln("  like that...")
				g.moreNomail()
				return
			} else if g.state.MarriedToSeth == int64(p.Record) {
				g.sln("")
				g.lln("  `0You kind of doubt you husband `%Seth Able `0would")
				g.sln("  like that...")
				g.moreNomail()
				return
			}
			g.singleVioletMarriage()
			return
		}
		if p.SeenViolet {
			break
		}
	}
}

// singleVioletNo is single_violet()'s no().
func (g *Game) singleVioletNo() {
	p := g.player
	g.sln("")
	g.sln("  Your hopeful smile slowly turns into a look of deep sadness.")
	g.sln("")
	g.lln("  `%THE ENTIRE BAR LAUGHS AT YOUR MISFORTUNE!")
	g.sln("")
	p.SeenViolet = true
	g.logLine("`2  `#Violet`2 has refused to marry `0" + p.Name + "`2!")
	g.more()
}

// singleVioletMarry is single_violet()'s marry().
func (g *Game) singleVioletMarry() {
	p := g.player
	g.getState(true)
	if g.state.MarriedToViolet != -1 && g.state.MarriedToViolet != int64(p.Record) {
		g.putState()
		mop := g.playerGet(g.state.MarriedToViolet)
		g.lln("`c                        `%** THE BLESSED DAY ARRIVES **`0")
		g.sln("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
		g.lln("  `2As you walk up to the chapel, you see Violet walking out...")
		g.lln("  `2on the arm of " + mop.Name + "!")
		g.sln("")
		return
	}
	g.state.MarriedToViolet = int64(p.Record)
	g.putState()
	g.lln("`c                        `%** THE BLESSED DAY ARRIVES **`0")
	g.sln("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
	g.lln("  `2Violet agrees to marry you!")
	g.sln("")
	g.sln("  After a short ceremony you are finally able to take her into your arms.")
	g.sln("  (Well, you've done that quite a few times, but not as your wife!)")
	g.sln("")
	g.sln("  She agrees to quit her job and take care of your house.")
	g.sln("")
	g.sln("  You look forward to inspiration from her every day.")
	g.sln("")
	g.lln("  `%YOU RECEIVE " + prettyInt(1000*p.Level) + " EXPERIENCE!")
	p.Exp += p.Level * 1000
	g.sln("")
	p.SeenViolet = false
	// DIFF (lord.js): The second line wasn't going into the log...
	g.logLine("`2  `#Violet`2 has `%MARRIED `0" + p.Name + "`2!!!!!\n`2  She `%QUITS`2 her job at the bar to the towns dismay!")
	g.tournamentCheck()
}

// singleVioletMarriage is single_violet()'s marriage().
func (g *Game) singleVioletMarriage() {
	p := g.player
	if p.Divorced {
		g.sln("  Since you've just gotten out of a serious relationship, you feel")
		g.sln("  it'd be moving to quickly to get married again.")
		return
	}
	g.sln("")
	g.lln("  `2You take Violet's hand and squeeze it gently.  `0\"My sweet Violet.")
	g.sln("  will you make me the happiest man alive?  Will you marry me?\"")
	g.sln("")
	g.lw("  `2With your heart on your sleeve, you await her answer.")
	g.mswait(1500)
	g.sw(".")
	g.mswait(1500)
	g.sw(".")
	g.mswait(1500)
	g.sln("")
	g.sln("")

	switch {
	case p.Cha < 3:
		g.lln("`#  \"Marry YOU?!\" `2 She laughs as if she just heard the funniest joke on")
		g.sln("   earth.")
	case p.Cha < 6:
		g.lln("`#  \"Is this a joke?\"`2 she asks seriously.")
	case p.Cha < 10:
		g.lln("`#  \"You're sweet, but I just don't like you in that way.\" `2she tells you")
		g.sln("   seriously.")
	case p.Cha < 50:
		g.lln("`#  \"I like you, I really do!  I just don't feel ready for that kind of")
		g.lln("   commitment right now!\" `2she explains sadly.")
	case p.Cha < 80:
		g.lln("`#  \"I think I'm in love with you!  But I can't marry you.  My last")
		g.lln("   marriage was a disaster.\" `2she explains sadly.")
	case p.Cha < 100:
		g.lln("`#  \"I...I can't say yes.  Ask me another time. Please.\" `2she begs.")
	default:
		g.lln("`#  \"Yes!  I WILL marry you!\" `2she laughs!")
		g.sln("")
	}

	if p.Cha < 100 {
		g.singleVioletNo()
	} else {
		g.singleVioletMarry()
	}
	g.moreNomail()
}

// singleVioletCarry is single_violet()'s carry().
func (g *Game) singleVioletCarry() {
	p := g.player
	g.sln("  You pick Violet up and roughly carry her upstairs and throw her on")
	g.sln("  the bed...")
	g.mswait(1500)
	g.sln("")
	p.SeenViolet = true
	if p.MarriedTo > -1 {
		g.sln("")
		g.lln("  `#Violet `0is appalled that you would even suggest such")
		g.sln("  a thing!")
		g.sln("")
		g.sln("  She calls you a dirty old man!")
		g.logLine("`#  Violet `2calls `0" + p.Name + " `2a `)" + pick(g.random(2) != 0, "dirty old man", "bastard") + "`2!")
	} else if p.Cha > 32 {
		if g.random(3) == 1 {
			g.lln("  `0She smiles invitingly...  ")
			g.mswait(1500)
			g.sln("")
			g.lln("  `2Half an hour later....You saunter downstairs. When a bearded drunk")
			g.sln("  asks you all the racket was, you smile knowingly and look away...")
			g.sln("  The drunks are mystified!")
			g.sln("")
			g.lln("  `0YOU GET `%" + prettyInt(p.Level*40) + "`0 EXPERIENCE!")
			p.Exp += p.Level * 40
			g.logLine("`5  " + p.Name + " `%Got laid by `#Violet`%!")
			p.Laid++
			g.tournamentCheck()
		} else {
			switch g.random(4) {
			case 0:
				g.sln("  She tells you she has a headache!")
				g.sln("  You saunter down to the bar disappointed.")
			case 1:
				g.sln("  She tells you she has to wash her hair!")
				g.sln("  You saunter down to the bar disappointed.")
			case 2:
				g.sln("  She tells you she is just too busy!")
				g.sln("  You saunter down to the bar disappointed.")
			case 3:
				g.sln("  She tells you she just doesn't have the urge today!")
				g.sln("  You saunter down to the bar - Wondering if she's been faithful!?")
			}
		}
	} else {
		g.sln("")
		g.foreground(4)
		g.sln("  She jumps off the bed screaming at you!  She savagely kicks you in the")
		g.sln("  groin!  You trudge down the stairs dejectedly..")
		g.sln("  The entire bar laughs at your misfortune!!")
		g.sln("")
		g.sln("  YOUR HITPOINTS GO DOWN TO 1!")
		p.HP = 1
		g.logLine("`0  " + p.Name + " `2got kicked in the groin by `#Violet`2!")
	}
	g.more()
}

// singleVioletWink is single_violet()'s wink().
func (g *Game) singleVioletWink() {
	p := g.player
	p.SeenViolet = true
	g.sln("")
	g.lln("  `%You wink at Violet seductively..")
	g.mswait(1500)
	if p.Cha >= 1 {
		p.Exp += p.Level * 5
		g.lln("  `2She blushes and smiles!!")
		g.sln("  You are making progress with her!")
		g.sln("")
		g.lln("  `%You receive " + prettyInt(p.Level*5) + " experience!")
		g.tournamentCheck()
	} else {
		g.lln("  `4She dumps a tankard of ale on your head!!")
		g.sln("  The entire bar laughs at your misfortune!!")
	}
	g.more()
}

// singleVioletKiss is single_violet()'s kiss().
func (g *Game) singleVioletKiss() {
	p := g.player
	p.SeenViolet = true
	if p.Cha >= 2 {
		g.lln("  `%You take Violet's hand and kiss it!")
		g.mswait(1500)
		g.lln("  `2She pulls her hand away and laughs merrily!!")
		g.sln("  Your relationship with her is taking off!")
		g.sln("")
		g.lln("  `%You receive " + prettyInt(p.Level*10) + " experience!")
		// BUG (lord.js): adds to player.ext, a property that doesn't
		// exist, so the experience announced is never awarded.  Kept.
	} else {
		g.lln("  `%You take her hand and kiss it!")
		g.mswait(1500)
		g.lln("  `2She pulls her hand away and slaps your face!")
		g.sln("  The entire bar laughs at your misfortune!!")
		g.sln("")
		g.lln("  `4YOU LOSE " + prettyInt(p.Level) + " HIT POINTS!")
		p.HP -= p.Level
		if p.HP < 1 {
			p.HP = 1
		}
	}
	g.more()
}

// singleVioletPeck is single_violet()'s peck().
func (g *Game) singleVioletPeck() {
	p := g.player
	g.lln("  `%You attempt to meekly kiss her on the lips, but...")
	g.mswait(1500)
	p.SeenViolet = true
	if p.Cha >= 4 {
		p.Exp += p.Level * 20
		g.lln("  `2She pulls you into a rough embrace, and kisses you")
		g.sln("  soundly!!   You are elated!")
		g.sln("")
		g.lln("  `%YOU RECEIVE " + prettyInt(p.Level*20) + " EXPERIENCE!")
		g.tournamentCheck()
	} else {
		g.lln("  `2She pushes you away and kicks you in the shin!")
		g.sln("  The entire bar laughs at your misfortune!!")
		g.sln("")
		g.lln("  `4YOU LOSE " + prettyInt(p.Level*3) + " HIT POINTS!")
		p.HP -= p.Level * 3
		if p.HP < 1 {
			p.HP = 1
		}
	}
	g.more()
}

// singleVioletSit is single_violet()'s sit().
func (g *Game) singleVioletSit() {
	p := g.player
	g.lln("  `%You grab her and throw her on your lap...")
	g.mswait(1500)
	p.SeenViolet = true
	if p.Cha >= 8 {
		g.lln("  `2She laughs and hugs you around the neck!!")
		g.sln("  You cherish her warmth!!")
		g.sln("")
		g.lln("  `%YOU RECEIVE " + prettyInt(p.Level*30) + " EXPERIENCE!")
		p.Exp += p.Level * 30
		g.tournamentCheck()
	} else {
		g.sln("  She jumps off your lap, and punches you hard!")
		g.sln("  Your chair tips over backwards!")
		g.sln("  The entire bar laughs at your misfortune!!")
		g.lln("")
		p.HP -= p.Level * 5
		if p.HP < 1 {
			p.HP = 1
		}
		g.lln("  `4YOU LOSE " + prettyInt(p.Level*5) + " HIT POINTS!")
	}
	g.more()
}

// singleVioletGrab is single_violet()'s grab().
func (g *Game) singleVioletGrab() {
	p := g.player
	g.lln("  `%You palm her shapely cheeks...`2")
	g.mswait(1500)
	p.SeenViolet = true
	if p.Cha >= 16 {
		p.Exp += p.Level * 40
		g.sln("  She doesn't object!  ")
		g.sln("  You become over excited and force your hand away!")
		g.sln("")
		g.lln("  `%YOU RECEIVE " + prettyInt(p.Level*40) + " EXPERIENCE!")
		g.tournamentCheck()
	} else {
		g.sln("  She breaks a tankard over your head!")
		g.sln("  You swoon and fall on your face!")
		g.sln("  The entire bar laughs at your misfortune!!")
		g.sln("")
		g.lln("  `4YOU LOSE " + prettyInt(p.Level*10) + " HITPOINTS!")
		p.HP -= p.Level * 10
		if p.HP < 1 {
			p.HP = 1
		}
	}
	g.more()
}

// flirtWithViolet is flirt_with_violet().  The result only mattered to
// the RIP menu redraw in red_dragon_inn().
func (g *Game) flirtWithViolet() bool {
	p := g.player

	if p.Sex == "F" {
		g.sln("")
		g.sln("")
		g.sln("  You would rather flirt with Seth Able.")
		return true
	}
	g.getState(false)
	if g.state.MarriedToViolet > -1 {
		if p.SeenViolet {
			g.sln("")
			g.sln("")
			g.lln("  `2You are still shaking from your last encounter with `4Grizelda`2!")
			g.sln("")
			g.more()
			return true
		}
		g.sln("")
		g.sln("")
		g.sln("  You whistle loudly for the barmaid, hardly containing your glee at the")
		g.sln("  thought of patting Violet's soft supple hips, but when you pat...")
		g.sln("")
		g.more()
		g.lln("  `4YOU FEEL HUGE LUMPS OF CELLULITE!")
		g.sln("")
		g.lln("`2  The obese barmaid intruduces her portly self as `4Grizelda`2!")
		g.sln("")
		if g.state.MarriedToViolet == int64(p.Record) {
			g.sln("  You remember now that Violet quit work when she married you!")
			g.sln("")
			g.sln("  You feel sorry for all the other warriors who must suffer at the hands of")
			g.lln("  `4Grizelda`2.")
		} else {
			op := g.playerGet(g.state.MarriedToViolet)
			g.lln("  You suddenly remember seeing something in the news about `0" + op.Name + "`2")
			g.sln("  marrying Violet!  As Grizelda grabs you for a kiss, her buckteeth jab")
			g.lln("  you painfully.  You curse `0" + op.Name + "`2 as you scream in horror.")
		}
		g.sln("")
		p.SeenViolet = true
		g.more()
		return true
	}
	if p.SeenViolet {
		g.sln("")
		g.sln("")
		g.sln("  You feel you had better not go too fast, maybe tomorrow.")
		g.sln("")
		return false
	}
	g.singleViolet()
	return true
}

// redDragonInn is red_dragon_inn().
//
// DIFF: lord.js returns -1 when the player got a room and 0 otherwise,
// but no caller uses the value, so it's dropped.
func (g *Game) redDragonInn() {
	p := g.player
	done := false

	menu := func() {
		if !p.Expert {
			// TODO (lord.js): RDI is ANSI, RDIF is lordtext!
			g.lrdfile(pick(p.Sex == "M", "RDI", "RDIF"))
		}
	}

	g.foreground(5)
	g.sln("")
	p.Inn = false
	g.sclrscr()
	menu()

	for !done {
		g.checkMail()
		g.foreground(5)
		g.sln("")
		g.lln("  The Red Dragon Inn   `8(? for menu)")
		g.sln("  (C,D,F,T,G,V,H,M,R)")
		g.commandPrompt()
		ch := strings.ToUpper(g.getkey())
		g.sln(ch)
		switch ch {
		case "R", "Q", "\r":
			done = true
		case "?":
			menu()
		case "V", "Y":
			g.showStats()
		case "W":
			g.composeMail()
		case "G":
			if g.getARoom() {
				done = true
			}
		case "D":
			g.showLog()
			menu()
		case "T":
			g.talkWithBartender()
			if p.Dead {
				done = true
			}
		case "C":
			g.converse(false)
		case "H":
			g.talkToBard()
		case "M":
			g.announce()
		case "F":
			g.flirtWithViolet()
		case "1", "2":
			// RIP-only pictures (LEGS, GUITAR).
		}
	}
}
