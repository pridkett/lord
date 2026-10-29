package game

import "strings"

// healers is healers(): the Healers' hut.
func (g *Game) healers() {
	var ch string
	p := g.player

	if p.HP < 0 {
		p.HP = 1
	}

	menu := func() bool {
		g.sclrscr()
		g.lrdfile("HEAL")
		// DIFF (lord.js): Was just ===
		if p.HP >= p.HPMax {
			g.sln("")
			g.sln("")
			g.foreground(15)
			g.lln("  `%Healers`#")
			g.sln("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
			g.lln("  `4`0\"You look fine to us!\"`2 the healers tell you.")
			g.sln("")
			g.more()
			return true
		}
		return false
	}

	if menu() {
		return
	}

	for {
		g.sln("")
		g.lln("  `3`2HitPoints: (`0" + prettyInt(p.HP) + " `2of`0 " + prettyInt(p.HPMax) + "`2)")
		g.lln("  `2Gold: `0" + prettyInt(p.Gold))
		g.lln("  `2(it costs `%" + prettyInt(5*p.Level) + "`2 to heal 1 hitpoint)")
		g.sln("")
		g.lln("  `5The Healers`2   (H,C,R)  (`0? for menu`2)")
		g.commandPrompt()
		ch = strings.ToUpper(g.getkey())
		g.sln(ch)
		switch ch {
		case "1":
			g.sln("")
			g.sln("  I wonder what's cooking..")
			g.sln("")
		case "2":
			g.sln("")
			g.sln("  Eye am not sure what you're")
			g.sln("  looking at.")
			g.sln("")
		case "3":
			g.sln("")
			g.sln("  Used to be a patient..  Hmmm.")
			g.sln("")
		case "C":
			g.healersHealSome()
		case "?":
			if menu() {
				return
			}
		case "H":
			if g.healersHealAll() {
				return
			}
			// DIFF (lord.js): space used to continue without show HitPoints etc...
		}
		if strings.Contains("RQ\r", ch) {
			break
		}
	}
}

// healersHealAll is healers' heal_all().
func (g *Game) healersHealAll() bool {
	p := g.player
	ret := false

	g.sclrscr()
	g.sln("")
	g.sln("")
	g.lln("  `%Healers`#")
	g.sln("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
	g.foreground(2)
	g.sln("")
	// DIFF (lord.js): was just =
	if p.HP >= p.HPMax {
		g.lln("  `0\"You look fine to us!\"")
	} else {
		afford := p.Gold / 5 / p.Level
		need := p.HPMax - p.HP
		// DIFF (lord.js): Was >, not >=
		if p.Gold >= (need * 5 * p.Level) {
			p.HP += need
			g.lln(" `0 " + prettyInt(need) + "`2 hit points are healed and you feel much better.")
			p.Gold -= (need * 5 * p.Level)
			ret = true
		} else if afford < need {
			p.HP += afford
			g.sln("  " + prettyInt(afford) + " hit points are healed and you feel much better.")
			p.Gold -= afford * 5 * p.Level
		}
	}
	g.sln("")
	g.more()
	return ret
}

// healersHealSome is healers' heal_some().
func (g *Game) healersHealSome() {
	p := g.player

	g.sclrscr()
	g.sln("")
	g.sln("")
	g.sln("")
	g.lw("  `2HitPoints: (`0" + prettyInt(p.HP) + " `2of `0" + prettyInt(p.HPMax) + "`2)")
	g.lln("  `2Gold: `0" + prettyInt(p.Gold))
	g.lln("  `2(it costs `%" + prettyInt(5*p.Level) + " `2to heal 1 hitpoint)")
	g.sln("")
	g.sln("  \"How many hit points would you like healed?\"")
	g.lw("  `0AMOUNT : `%")
	amt, ok := jsParseInt(g.consoleGetstr(GetStrOpts{Len: 5, Integer: true}))
	if !ok {
		amt = 0
	}
	g.sln("")
	switch {
	case amt == 0:
		g.sln("  \"Maybe some other time..\"")
	case amt < 0:
		g.sln("  \"Uh...Wouldn't that be hurting yourself?!\"")
	case amt*(5*p.Level) > p.Gold:
		g.sln("  \"I'm afraid you don't have enough gold to cover that.\"")
	case amt > (p.HPMax - p.HP):
		g.sln("  \"It would be deadly to over heal yourself!!\"")
	default:
		g.sln("  Done!")
		p.Gold = (p.Gold - (amt * 5 * p.Level))
		p.HP += amt
	}
}
