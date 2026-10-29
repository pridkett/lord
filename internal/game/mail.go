package game

import (
	"fmt"
	"strings"
)

// writeMail is write_mail(): the line editor for a message to player to,
// optionally quoting the message just read (the quote file).
func (g *Game) writeMail(to int64, quote bool) {
	p := g.player
	qf := g.quoteFile()
	mb := g.gamedir("mail" + itoa(to) + ".lrd")
	wrap := ""
	lines := 0
	derps := []string{"`.`%  Greetings.  How fare you, traveler?",
		"`.`%  Well met.  Any news you can share?",
		"`.`%  How goes it, fellow adventurer?",
		"`.`%  Didn't I see you on that table in the Dark Cloak tavern?",
		"`.`%  Sorry - I forgot what I was going to say, old bean."}
	bs := strings.Repeat("\x08", 76)
	ws := strings.Repeat(" ", 76)

	g.sln("")
	g.sln("  Enter message now..Blank line quits!")
	// DIFF: lord.js writes the message to msg<Record>.lrd and copies that
	// file into the mailbox; here it's collected in memory.
	var mf []string
	mf = append(mf, "`. ")
	mf = append(mf, "`.  `0"+p.Name+"`2 sent you this on "+g.formatDate())
	mf = append(mf, "`.`0-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-`2")
	if quote && fileExists(qf) {
		if ql, ok := readLines(qf); ok {
			for _, l := range ql {
				if strings.HasPrefix(l, "`% ") {
					l = "`0>" + l[3:]
					mf = append(mf, l)
					g.lln(l)
				}
			}
		}
	}
	fileRemove(qf)
	for {
		l := wrap
		wrap = ""
		g.sln("")
		g.foreground(2)
		g.sw(" >")
		g.foreground(15)
		g.sw(l)
		for {
			ch := g.getkey()
			g.sw(ch)
			if ch == "\x08" {
				if len(l) > 0 {
					l = l[:len(l)-1]
				}
				g.sw(" \x08")
			} else if ch != "\x1b" {
				if ch != "\r" {
					l += ch
				}
				if len(l) > 75 {
					if t := strings.LastIndex(l, " "); t != -1 {
						wrap = l[t+1:]
						l = l[:t]
						n := min(len(wrap), 76)
						g.sw(bs[:n])
						g.sw(ws[:n])
					}
					ch = "\r"
				}
			}
			if ch == "\r" {
				break
			}
		}
		if len(l) == 0 && lines == 0 {
			// DIFF (lord.js): Derps weren't passed through clean_str(),
			//       so kept the `. and therefore weren't quoted.
			mf = append(mf, cleanStr(derps[g.random(int64(len(derps)))]))
			break
		}
		lines++
		l = cleanStr("`.`%  " + l)
		mf = append(mf, l)
		if len(l) < 5 {
			break
		}
	}
	mf = append(mf, "`-"+fmt.Sprintf("%03d", p.Record))
	g.fmutex(mb)
	if err := appendLines(mb, mf...); err != nil {
		g.rfmutex(mb)
		g.fatalf("Failed to open %s", mb)
	}
	g.rfmutex(mb)
	if to == int64(p.Record) {
		g.sln("")
		g.sln("  You are a very stupid individual.")
		g.sln("")
		g.moreNomail()
	}
}

// mailTo is mail_to(): append mail (lines separated by "\n") to player
// to's mailbox.
func (g *Game) mailTo(to int64, mail string) {
	mb := g.gamedir("mail" + itoa(to) + ".lrd")
	g.fmutex(mb)
	if err := appendLines(mb, mail); err != nil {
		g.rfmutex(mb)
		g.fatalf("Failed to open %s", mb)
	}
	g.rfmutex(mb)
}

// showLooks is show_looks().
func (g *Game) showLooks(op *Player) {
	upronoun := pick(op.Sex == "M", "His", "Her")
	upronoun2 := pick(op.Sex == "M", "He", "She")

	if op.Cha < 30 {
		g.lln("  `4WARNING:  `2" + upronoun + " Charm Rating is only `%" + prettyInt(op.Cha) + "`2!")
	} else {
		g.lln("  `2" + upronoun + " Charm Rating is `%" + prettyInt(op.Cha) + "`2.")
	}
	g.sln("")
	g.lln("  `2" + upronoun2 + " " + g.looks(op.Sex, op.Cha))
	g.sln("")
}

// deadScreen is dead_screen().
func (g *Game) deadScreen(op *Player) {
	g.sln("")
	g.lln("  `4You have been killed by " + op.Name + "`2.")
	g.sln("")
	g.moreNomail()
	g.lln("  `2GOLD ON HAND WAS `4LOST`2.  ")
	g.sln("")
	g.lln("  `2TEN PERCENT OF EXPERIENCE `4LOST`2.")
	g.sln("  ")
	g.foreground(2)
	g.sln("  You have been defeated on your way to glory.  The road to success")
	g.lln("  is long and hard.  You have encountered a minor setback.  But do `0NOT`2")
	g.sln("  lose heart, you can continue your struggle tomorrow.")
	g.sln("")
	g.moreNomail()
}

// ---- mail actions (functions nested in lw) ---------------------------------

const romanticRule = "`2-=-=-=-=-=-=-=-=-=-=--=-=-=-=-=-=-=-=-=-=-=-=-=--=-=-=-=-=-=-\n"

// answerMail is lw's answer_mail().
func (g *Game) answerMail(to int64) {
	p := g.player
	op := g.playerGet(to)

	g.sln("")
	g.foreground(10)
	g.sw("  Answer " + pick(op.Sex == "M", "him?", "her?"))
	g.lw(" `2[`5Y`2] `0: ")
	ch := strings.ToUpper(g.getkey())
	if ch != "N" {
		ch = "Y"
	}
	g.sln(ch)
	if ch == "N" {
		fileRemove(g.quoteFile())
		return
	}
	if to == int64(p.Record) {
		g.sln("")
		g.sln("  ** THE MESSENGER REFUSES TO DELIVER THE MESSAGE! **")
		return
	}
	he := pick(op.Sex == "M", "he", "she")
	g.sln("")
	g.sln("")
	g.foreground(10)
	g.sw("  Quote what " + he + " said? : ")
	g.lw("`2[`5Y`2] `0: ")
	ch = strings.ToUpper(g.getkey())
	if ch != "N" {
		ch = "Y"
	}
	g.sln(ch)
	g.writeMail(to, ch == "Y")
	fileRemove(g.quoteFile())
}

// smileMail is lw's smile_mail().
func (g *Game) smileMail(to int64) {
	p := g.player
	op := g.playerGet(to)

	g.sln("")
	him := pick(op.Sex == "M", "him", "her")
	he := pick(p.Sex == "M", "He", "She")
	g.showLooks(op)
	g.lln("  `2(`0P`2)ointedly ignore " + him)
	g.lln("  `2(`0S`2)mile hugely")
	g.sln("")
	g.lw("  What do you do about it? `0[`2S`0]`2 : ")
	ch := strings.ToUpper(g.getkey())
	if ch != "P" {
		ch = "S"
	}
	g.sln(ch)
	g.sln("")
	g.lln("  `%** Please wait, writing " + him + " back **")
	// DIFF (lord.js): Returns didn't have `|, so the "Return from" bit was
	//       quoted in the NEXT message.  Starting the first line with `.
	//       would also fix the problem.
	if ch != "P" {
		g.mailTo(to, "`%  ** Romantic Mail Return From "+p.Name+"`% **\n"+
			romanticRule+
			"`0  "+he+" smiles back, encouragingly!\n"+
			"`0  \n"+
			"`2  You receive `%"+prettyInt(5*op.Level)+" `2experience points!\n"+
			"`E"+itoa(5*op.Level)+"\n"+
			"```|")
	} else {
		g.mailTo(to, "`%  ** Romantic Mail Return From "+p.Name+"`% **\n"+
			romanticRule+
			"`0  "+he+" ignores you!  You are snubbed good!\n"+
			"`0  \n"+
			"`2  You `4LOSE `%"+prettyInt(5*op.Level)+" `2experience points!\n"+
			"`V"+itoa(5*op.Level)+"\n"+
			"```|")
	}
	fileRemove(g.quoteFile())
}

// kissMail is lw's kiss_mail().
func (g *Game) kissMail(to int64) {
	p := g.player
	op := g.playerGet(to)

	g.sln("")
	him := pick(op.Sex == "M", "him", "her")
	he := pick(p.Sex == "M", "He", "She")
	g.showLooks(op)
	g.lln("  `2(`0P`2)ointedly ignore " + him)
	g.lln("  `2(`0K`2)iss " + him + " soundly")
	g.sln("")
	g.lw("  What do you do about it? `0[`2K`0]`2 : ")
	ch := strings.ToUpper(g.getkey())
	if ch != "P" {
		ch = "K"
	}
	g.sln(ch)
	g.sln("")
	g.lln("  `%** Please wait, writing " + him + " back **")
	// DIFF (lord.js): Returns didn't have `|, so the "Return from" bit was
	//       quoted in the NEXT message.
	if ch != "P" {
		g.mailTo(to, "`%  ** Romantic Mail Return From "+p.Name+"`% **\n"+
			romanticRule+
			"`0  "+he+" gives you a big wet kiss!\n"+
			"`2  You receive `%"+prettyInt(10*op.Level)+" `2experience points!\n"+
			"`E"+itoa(10*op.Level)+"\n"+
			"```|")
	} else {
		g.mailTo(to, "`%  ** Romantic Mail Return From "+p.Name+"`% **\n"+
			romanticRule+
			"`0  "+he+" ignores you!  You are snubbed really good!\n"+
			"`0  \n"+
			"`2  You `4LOSE `%"+prettyInt(10*op.Level)+" `2experience points!\n"+
			"`V"+itoa(10*op.Level)+"\n"+
			"```|")
	}
	fileRemove(g.quoteFile())
}

// dinnerMail is lw's dinner_mail().
func (g *Game) dinnerMail(to int64) {
	p := g.player
	op := g.playerGet(to)

	g.sln("")
	him := pick(op.Sex == "M", "him", "her")
	he := pick(p.Sex == "M", "He", "She")
	g.showLooks(op)
	g.lln("  `2(`0P`2)ointedly ignore " + him)
	g.lln("  `2(`0G`2)o to dinner with " + him)
	g.sln("")
	g.lw("  What do you do about it? `0[`2G`0]`2 : ")
	ch := strings.ToUpper(g.getkey())
	if ch != "P" {
		ch = "G"
	}
	g.sln(ch)
	g.sln("")
	g.lln("  `%** Please wait, writing " + him + " back **")
	// DIFF (lord.js): Returns didn't have `|, so the "Return from" bit was
	//       quoted in the NEXT message.
	if ch != "P" {
		// lord.js leaves out the "\n" after the fourth `0 line here.
		g.mailTo(to, "`%  ** Romantic Mail Return From "+p.Name+"`% **\n"+
			romanticRule+
			"`0  "+he+" agrees to go to dinner with you!\n"+
			"`0  \n"+
			"`0  You both have a wonderful time, and learn a lot from each other.\n"+
			"`0  "+
			"`2  You receive `%"+prettyInt(20*op.Level)+" `2experience points!\n"+
			"`E"+itoa(20*op.Level)+"\n"+
			"```|")
		g.logLine("`%  " + p.Name + " `2agrees to have dinner with `0" + op.Name + "`2!")
		g.foreground(2)
		g.sln("")
		g.sln("  You have a wonderful time at dinner, and both learn a lot about")
		g.sln("  each other.  You both vow to do it again sometime.")
		g.sln("")
		g.moreNomail()
	} else {
		g.mailTo(to, "`%  ** Romantic Mail Return From "+p.Name+"`% **\n"+
			romanticRule+
			"`0  "+he+" laughs in your face!  You are snubbed so bad.\n"+
			"`0  \n"+
			"`2  You `4LOSE `%"+prettyInt(20*op.Level)+" `2experience points!\n"+
			"`V"+itoa(20*op.Level)+"\n"+
			"```|")
		g.logLine("`%  " + p.Name + " `2refuses to dine with `0" + op.Name + "`2!")
	}
	fileRemove(g.quoteFile())
}

// sleepMail is lw's sleep_mail().
func (g *Game) sleepMail(to int64) {
	p := g.player
	op := g.playerGet(to)

	g.sln("")
	him := pick(op.Sex == "M", "him", "her")
	he := pick(p.Sex == "M", "He", "She")
	ohe := pick(op.Sex == "M", "he", "she")
	g.showLooks(op)
	g.lln("  `2(`0P`2)ointedly tell " + him + " absolutely not!")
	g.lln("  `2(`0S`2)leep with " + him)
	g.sln("")
	g.lw("  What do you do about it? `0[`2P`0]`2 : ")
	ch := strings.ToUpper(g.getkey())
	if ch != "S" {
		ch = "P"
	}
	g.sln(ch)
	g.sln("")
	g.lln("  `%** Please wait, writing " + him + " back **")
	// DIFF (lord.js): Returns didn't have `|, so the "Return from" bit was
	//       quoted in the NEXT message.
	if ch == "S" {
		p.Laid++
		// BUG (lord.js): the `E amount goes through pretty_int(), so the
		// receiver's parseInt() stops at the first comma.  Kept.
		g.mailTo(to, "`%  ** Romantic Mail Return From "+p.Name+"`% **\n"+
			romanticRule+
			"`0  "+he+" agrees to sleep with you!\n"+
			"\n"+
			"  You both have a wonderfully sweaty time, and learn very little\n"+
			"  from each other.  (but you had fun!)\n"+
			"\n"+
			"`2  You receive `%"+prettyInt(50*op.Level)+" `2experience points!\n"+
			"`E"+prettyInt(50*op.Level)+"\n"+
			"`{\n"+
			"```|")
		g.logLine("`%  " + p.Name + " `2got laid by `0" + op.Name + "`2!")
		g.lln("`c`%  ** UPSTAIRS IN THE INN **")
		g.sln("")
		g.lln("  `2When you finally arrive at " + op.Name + "`2's bedroom at the appointed")
		g.sln("  time, you are nervous as ever!")
		g.sln("")
		g.moreNomail()
		g.sln("  Your knocks are unanswered.  Since the door is slightly ajar, you ")
		g.sln("  see no harm in letting yourself in.  You push through the beads into a")
		g.sln("  smokey room.  The smell of herbs fills your nostrils, and the sound of")
		g.sln("  soft music fills your ears.")
		g.sln("")
		g.moreNomail()
		g.lln("  And then a skimpily clad " + op.Name + "`2 walks into the room!")
		g.sln("")
		g.moreNomail()
		g.lln("  Without a word, " + ohe + " points at the heart-shaped bed.")
		g.sln("")
		g.moreNomail()
		g.sln("  A good while later you are drenched in sweat, tired, but feeling")
		g.lln("  good.  " + op.Name + "`2 holds you close, as if " + ohe + " never wanted")
		g.sln("  this moment to end.")
		g.sln("")
		if g.random(2) == 0 {
			g.moreNomail()
			if p.Sex == "M" {
				g.lln("  `%Seth Able hi-fives you on the way out!")
			} else {
				g.lln("  `%Violet gives you a dirty look on the way out!")
			}
		}
		g.sln("")
		g.moreNomail()
		g.tournamentCheck()
	} else {
		g.mailTo(to, "`%  ** Romantic Mail Return From "+p.Name+"`% **\n"+
			romanticRule+
			"`0  "+he+" laughs in your face!  You are snubbed the worst you\n"+
			"  have ever been snubbed before in your life.\n"+
			"  \n"+
			"`2  You `4LOSE `%"+prettyInt(50*op.Level)+" `2experience points!\n"+
			"`V"+itoa(50*op.Level)+"\n"+
			"```|")
		g.logLine("`%  " + op.Name + " `2got seriously snubbed by `0" + p.Name + "`2!")
		g.foreground(15)
		g.sln("")
		g.sln("  Ouch!")
		g.sln("")
	}
	fileRemove(g.quoteFile())
}

// marryMail is lw's marry_mail().
func (g *Game) marryMail(to int64) {
	p := g.player
	mop := g.playerGet(to)

	him := pick(mop.Sex == "M", "him", "her")
	she := pick(p.Sex == "M", "He", "She")
	her := pick(p.Sex == "M", "Her", "His")

	g.lln("  `2(`0S`2)ay YES!")
	g.lln("  `2(`0T`2)urn " + him + " down!")
	g.sln("")
	g.lw("  What do you do about it? `0[`2T`0]`2 : ")
	ch := strings.ToUpper(g.getkey())
	if ch != "S" {
		ch = "T"
	}
	g.sln(ch)
	g.sln("")
	g.lln("  `%** Please wait, writing " + him + " back **")
	if ch == "S" {
		p.MarriedTo = to
		mop.Lock()
		mop.ReLoad()
		// DIFF (lord.js): See what happens here... this check is new.
		if mop.MarriedTo == -1 || mop.MarriedTo == int64(p.Record) {
			mop.MarriedTo = int64(p.Record)
			mop.Put()
			mop.UnLock()
			p.Put()
			g.mailTo(to, "`%  ** Romantic Mail Return From "+p.Name+"`% **\n"+
				romanticRule+
				"`0  "+she+" agrees to marry you!\n"+
				"\n"+
				"  The wedding is simple, but beautiful.   Your significant other\n"+
				"  stands elequently beside you, as you echange vows.)\n"+
				"\n"+
				"  `2Knowing the elders don't take marriage lightly, you hope it\n"+
				"  will last...\n"+
				"`?"+itoa(p.Record))
			g.logLine("`%** `0" + p.Name + " `2said YES to `0" + mop.Name + "`2's marriage proposal! `%**")
			g.sclrscr()
			g.foreground(15)
			g.sln("")
			g.sln("")
			g.sln("  ** THE WEDDING **")
			g.sln("")
			g.foreground(2)
			g.sln("  Your spouse is overjoyed at your answer.  " + her + " face beams so")
			g.sln("  brightly,  you know you've made the right choice.")
			g.sln("")
			g.moreNomail()
			if p.Sex == "M" {
				g.sln("  The vows are short and to the point - but heartfelt.  You are especially ")
				g.sln("  proud of the one \"To have; To hold..  And to satisfy anytime she needs it\"")
				g.lln("  (you made that one up yourself)  `%Seth Able `2laughs at the back of the")
				g.sln("  church at this.  Many faces turn in his direction - making him turn")
				g.sln("  a bright purple.")
				g.sln("")
				g.moreNomail()
				if g.random(2) == 0 {
					g.lln("  `0A figure dressed in black attracts your attention in the back...")
					g.sln("")
					g.moreNomail()
					g.lln("  `2You recognize the voloptous curves to be `#Violet`2's.")
					g.sln("")
					g.sln("  A single tear glistens on her cheek.")
				}
				g.sln("")
				g.moreNomail()
				g.lln("  `2Turgon solomnly gives you permission to kiss the bride.")
				g.sln("")
				g.moreNomail()
			} else {
				g.sln("  The vows are short and to the point - but heartfelt.  Violet giggles as she")
				g.sln("  hands you the ring.  (You wonder why)  Even through the veil covering")
				g.sln("  your face - your man looks so handsome - so strong.  ")
				g.sln("")
				g.moreNomail()
				if g.random(2) == 0 {
					g.lln("  `0A figure dressed in black attracts your attention in the back...")
					g.sln("")
					g.sln("")
					g.moreNomail()
					g.lln("  `2You recognize the individual by the tightly clutched mandolin.")
					g.sln("")
					g.sln("  A single tear glistens on his cheek.")
					g.sln("")
					g.moreNomail()
				}
				g.sln("  Turgon solomnly gives your husband his permission to kiss you.")
				g.sln("")
				g.moreNomail()
			}
			g.lln("  `%THE ENTIRE CROWD ROARS ITS APPROVAL!")
			g.sln("")
		} else {
			// BUG (lord.js): player.married_to stays set to the proposer,
			// who is married to someone else.  Kept.
			mop.UnLock()
			bae := g.playerGet(mop.MarriedTo)
			g.lln("`0  " + mop.Name + " laughs in your face!")
			g.lln("  You never answered, so I married " + bae.Name + " instead!")
			g.sln("  " + she + " can at least answer a question.")
			g.logLine("`%** `0" + p.Name + " `2said YES to `0" + mop.Name + "`2's marriage proposal! `%**\n`2  (But not until after `0" + mop.Name + "`2 already married `0" + bae.Name + "`2!")
			g.lln("`4  YOU ARE `)ULTRA `4SNUBBED!")
			g.lln("`0  ")
			g.lln("`2  You `4LOSE `%" + prettyInt(100*mop.Level) + " `2experience points!")
			p.Exp -= 100 * mop.Level
		}
	} else {
		// DIFF (lord.js): Returns didn't have `|, so the "Return from" bit
		//       was quoted in the NEXT message.
		g.mailTo(to, "`%  ** Romantic Mail Return From "+p.Name+"`% **\n"+
			romanticRule+
			"`0  "+she+" laughs in your face!\n"+
			"`4  YOU ARE `)ULTRA `4SNUBBED!\n"+
			"`0  \n"+
			"`2  You `4LOSE `%"+prettyInt(100*mop.Level)+" `2experience points!\n"+
			"`V"+itoa(100*mop.Level)+"\n"+
			"```|")
		g.logLine("`%  " + mop.Name + " `2got `)TURNED DOWN `2by `0" + p.Name + "`2!")
		g.sln("")
		g.lln("  `%Yowzers you're cold!")
		g.sln("")
	}
	fileRemove(g.quoteFile())
}
