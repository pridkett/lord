package game

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// turgons is turgons(): Turgon's Warrior Training.
func (g *Game) turgons() {
	var trainer *Player
	var ch string
	p := g.player

	menu := func() {
		g.sclrscr()
		g.lrdfile("TURGON")
	}

	prompt := func() bool {
		if p.Level < 12 {
			trainer = g.getTrainer(p.Level)
		} else {
			trainer = g.getTrainer(11)
		}
		g.lln("  `3`2Your master is `%" + trainer.Name + "`2.")
		g.sln("")
		g.lln("  `5Turgon's Warrior Training`2")

		if p.Level > 11 {
			g.sln("")
			g.sln("  You pay your respects to Turgon, and stroll around the grounds.  Lesser")
			g.sln("  warriors bow low as you pass.  Turgon's last words advise you to find and")
			g.lln("  kill the `4Red Dragon`2..")
			g.sln("")
			g.more()
			return true
		}
		g.lln("`2  (Q,A,V,R)  (`0? for menu`2)")
		return false
	}

	menu()
	for {
		if prompt() {
			return
		}

		g.commandPrompt()
		ch = strings.ToUpper(g.getkey())
		g.sln(ch)
		switch ch {
		case "Q":
			g.turgonsAsk(trainer)
		case "R", "\r":
		case "A":
			g.turgonsAttackMaster(trainer)
		case "?":
			menu()
		// DIFF (lord.js): Removed space which just skipped prompt.
		case "V":
			tf := g.gamedir("temp" + itoa(p.Record))
			g.turgonsRankKing(tf)
			g.showLordFile(tf, false, false)
			fileRemove(tf)
			g.sln("")
			g.moreNomail()
			g.foreground(2)
		}
		if strings.Contains("R\r", ch) {
			break
		}
	}
}

var pweRE = regexp.MustCompile(`(?i)&PWE`)

// turgonsAsk is turgons' ask(): question your master.
func (g *Game) turgonsAsk(trainer *Player) {
	p := g.player
	mc := g.morechk

	g.sclrscr()
	g.sln("")
	g.lln("  `%Questioning Your Master`2")
	g.sln("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
	g.foreground(10)
	g.morechk = false
	g.lrdfile("LEVEL"+itoa(p.Level)+pick(p.Sex == "M", "MALE", "FEMALE"), true)
	g.sln("")
	if p.Exp > trainer.Need {
		g.lln("  `0" + trainer.Name + " `2looks at your carefully and says: ")
		g.sln("")
		if trainer.NeedStr1 != "" {
			// String.replace with a regexp replaces the first match only.
			s := trainer.NeedStr1
			if loc := pweRE.FindStringIndex(s); loc != nil {
				s = s[:loc[0]] + p.Weapon + s[loc[1]:]
			}
			g.lln("  `2\"`0" + s + "`2\"")
		}
		if trainer.NeedStr2 != "" {
			g.lln("  `2\"`0" + trainer.NeedStr2 + "`2\"")
		}
	} else {
		g.sln("")
		g.lln("  `0" + trainer.Name + "`2 looks at you carefully.")
		g.sln("")
		g.lln("  `2\"`0You need about `%" + prettyInt(trainer.Need-p.Exp) + "`0 more experience before")
		g.lln("   you will be as good as I am.`2\"")
	}
	g.sln("")
	g.morechk = mc
	g.more()
}

// turgonsRankKing is turgons' rank_king(): write the heroes list to fname.
func (g *Game) turgonsRankKing(fname string) {
	var pl []*Player

	for _, op := range g.allPlayers() {
		if op.Name != "X" && op.DragKills > 0 {
			pl = append(pl, op)
		}
	}
	sort.SliceStable(pl, func(i, j int) bool {
		a, b := pl[i], pl[j]
		if a.DragKills != b.DragKills {
			return b.DragKills < a.DragKills
		}
		if a.Level != b.Level {
			return b.Level < a.Level
		}
		return b.Exp < a.Exp
	})
	lines := []string{
		"",
		"",
		"                           `%Heroes Of The Realm",
		"",
		"  `0Name                       Heroic Deeds Done         Current Level",
		"`2-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-",
	}
	for _, op := range pl {
		lines = append(lines, "`0  "+spacePad(op.Name, 22)+"           `%"+fmt.Sprintf("%3d", op.DragKills)+"                       `0"+fmt.Sprintf("%2d", op.Level))
	}
	if len(pl) == 0 {
		lines = append(lines, "  `0Sad times indeed, there are no heroes in this realm.")
	}
	if err := writeLines(fname, lines...); err != nil {
		g.fatalf("Unable to open %s", fname)
	}
}

// turgonsAttackMaster is turgons' attack_master(): fight your master.
func (g *Game) turgonsAttackMaster(trainer *Player) {
	p := g.player

	g.sclrscr()
	g.sln("")
	g.sln("")
	g.lln("  `%Fighting Your Master")
	g.foreground(2)
	g.sln("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
	g.foreground(10)
	if p.SeenMaster {
		son := pick(p.Sex == "M", "son", "daughter")
		g.sln("  \"I would like to battle again, but it is too late my " + son + ".\"")
		g.lln("  `2" + trainer.Name + " tells you.  You figure you will try again")
		g.sln("  tomorrow.")
		g.sln("")
		g.moreNomail()
		return
	}
	if p.Exp < trainer.Need {
		g.sln("  You are escorted down the hallway and into the battle arena.")
		g.sln("")
		g.foreground(10)
		g.sln("  THE BATTLE BEGINS!")
		g.sln("")
		g.moreNomail()
		g.lln("  `2You raise your `0" + p.Weapon + " `2to strike!  You wonder why everyone")
		g.sln("  is looking at you with grins on their faces...")
		g.sln("")
		g.moreNomail()
		g.sln("  Your weapon is gone!  You are holding air!")
		g.sln("")
		g.moreNomail()
		g.sln("  Your Master is holding it!  The entire crowd is laughing at you!")
		g.sln("")
		g.sln("  You meekly accept the fact that you are not ready for your testing.")
		g.sln("")
		g.moreNomail()
		p.SeenMaster = true
		return
	}
	g.lln("`2  You enter the fighting arena, ready with your `0" + p.Weapon + "`2.")
	g.sln("")
	g.sln("  When your name is called, you move to the proper position and take a")
	g.sln("  fighting stance against your master.")
	p.SeenMaster = true
	g.sln("")
	g.sln("")
	g.lln("  `2**`%MASTER FIGHT`2**")
	g.sln("")
	g.lln("  You have encountered " + trainer.Name + "`2!!")
	g.sln("")
	// DIFF (lord.js): Added this pause when the battle prompt would pause.
	if g.curlinenum+g.battlePromptLines() > (g.con.Rows - 1) {
		g.moreNomail()
	}
	g.battle(trainer, false, false)
	if p.Dead {
		p.Dead = false
		p.HP = p.HPMax
		g.sln("")
		g.lln("  `0" + trainer.Name + " `2raises his " + trainer.Weapon + "`2 to kill you!")
		g.sln("")
		g.moreNomail()
		g.sln("  At the last minute, he reaches down and helps you up.  He tells you not to")
		g.sln("  be discouraged, and for good gesture has you healed before you go.")
		g.sln("")
		g.more()
	} else if trainer.HP < 1 {
		p.Dead = false
		g.sln("")
		g.lln("  `%You have bested " + trainer.Name + "`%!")
		g.sln("")
		g.lln("`%  " + trainer.Swear)
		g.sln("")
		p.HPMax += trainer.HPGained
		p.HP = p.HPMax
		p.Str += trainer.StrGained
		p.Def += trainer.Def
		p.Level++
		g.lw("`2  You receive `0" + prettyInt(trainer.HPGained) + "`2 hitpoints,")
		g.lw(" `0" + prettyInt(trainer.StrGained) + "`2 strength")
		g.lln(" and `0" + prettyInt(trainer.Def) + "`2 defense points!")
		g.sln("")
		g.lln("  `%YOU ARE NOW LEVEL " + prettyInt(p.Level) + ".")
		g.sln("")
		mline := "  `0" + p.Name + " `2has beaten `%" + trainer.Name + "!"
		if p.Level == 12 {
			mline += "\n"
			if p.Sex == "M" {
				mline += "  He "
			} else {
				mline += "  She "
			}
			mline += "has become the Ultimate Warrior!"
		}
		g.logLine(mline)
		g.raiseClass()
		p.SeenMaster = false
		g.tournamentCheck()
	}
}
