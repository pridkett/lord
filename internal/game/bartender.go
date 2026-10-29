package game

import (
	"os"
	"strings"
)

// talkWithBartender is talk_with_bartender(): the bartender of the Red
// Dragon Inn (gems for elixirs, name changes, bribes).
func (g *Game) talkWithBartender() {
	p := g.player
	done := false
	classNames := []string{"lazy", "warrior", "mystical skills user", "thief"}

	uniqueName := func(name string) bool {
		plen := g.playerLength()
		for ui := 0; ui < plen; ui++ {
			if ui != p.Record {
				uop := g.playerGet(int64(ui))
				if name == strings.TrimSpace(dispStr(uop.Name)) {
					g.sln("")
					g.foreground(2)
					g.sln("  You recall hearing that another warrior currently goes by that name.  Not")
					g.lln("  wanting to dishonor " + uop.Name + "`2 you decide to pick another.")
					g.sln("")
					return false
				}
			}
		}
		return true
	}

	menu := func() {
		if p.Level > 11 {
			g.lrdfile(pick(p.Sex == "M", "BT1", "BT1F"))
		}
		if p.Level < 12 {
			g.lrdfile(pick(p.Sex == "M", "BT", "BTF"))
		}
	}

	g.sln("")
	g.sln("")
	g.foreground(2)
	g.sln("  You find the bartender and ask if he will talk")
	g.sln("  privately with you.")
	g.sln("")

	if p.Level == 1 {
		g.foreground(5)
		g.sln("  \"I don't recall ever hearing the name")
		g.lln("  `0" + p.Name + "`5 before!  Get outta my face!\"")
		g.foreground(2)
		g.sln("")
		return
	}

	menu()

	for !done {
		g.foreground(5)
		g.sln("")
		g.sw("  \"Well?\"")
		g.lln(" The bartender inquires.  (`0? for menu`2)")
		g.commandPrompt()
		ch := strings.ToUpper(g.getkey())
		if ch == "\r" {
			ch = "R"
		}
		g.sln(ch)
		switch ch {
		case "?":
			menu()
		case "R", "Q":
			done = true
		case "V":
			if p.Sex == "M" {
				g.sln("")
				g.foreground(5)
				g.lln("  \"Ya want to know about `#Violet`5 do ya?  She is every warrior's")
				g.sln("  wet dream...But forget it, Lad, she only goes for the type")
				g.sln("  of guy who would help old people...\"")
			}
		case "S":
			if p.Sex == "F" {
				g.sln("")
				g.foreground(5)
				g.lln("  \"Ya want to know about `%Seth Able`5 the Bard, eh?  Well... He")
				g.sln("  has a good voice, and can really play that mandolin.  I don't")
				g.sln("  think he would go for your type tho,  he likes he type of girl")
				g.sln("  that would help an old man or something... A lot of Charm is")
				g.sln("  is what you would need.  However, I could sure go for your type")
				g.sln("  Har har har!\"")
			}
		case "G":
			g.talkWithBartenderGems()
		case "C":
			g.sln("")
			g.foreground(5)
			g.sln("  \"Ya wanna change your name, eh?  Yeah..")
			cn := "undefined"
			if p.Clss >= 0 && p.Clss < int64(len(classNames)) {
				cn = classNames[p.Clss]
			}
			g.lln("  " + p.Name + " `5the " + cn + " does sound kinda funny..")
			g.sw("  it would cost ya " + prettyInt(p.Level*500) + " gold... Deal?\"")
			g.lw("  `2Change your name?  [`0N`2] : ")
			ch = strings.ToUpper(g.getkey())
			if ch != "Y" {
				ch = "N"
			}
			g.sln(ch)
			if ch != "Y" {
				g.foreground(5)
				g.sln("  \"Fine.. Keep your stupid name.. See if I care..\"")
				break
			}
			g.sln("")
			if p.Gold < 500*p.Level {
				g.foreground(4)
				g.sln("")
				g.sln("  \"Hey!  You stupid fool! You don't have that much gold!\"")
				g.sln("")
				break
			}
			p.Gold -= 500 * p.Level

			var nname string
			for {
				g.foreground(5)
				g.sln("  \"What would you like as an alias?\" ")
				g.foreground(2)
				g.sw("  NAME: ")
				nname = g.consoleGetstr(GetStrOpts{Len: 50})
				g.sln("")
				nname = strings.TrimSpace(nname)
				nname = cleanStr(nname)
				tmp := strings.TrimSpace(dispStr(nname))
				if len(nname) > 18 {
					g.foreground(4)
					g.sln("  \"Try a shorter name ya stupid, ugly troll!\"")
					g.sln("")
				} else if len(nname) < 3 {
					g.foreground(4)
					g.sln("  Try a longer name, bonehead!")
					g.sln("")
				} else if len(tmp) < 3 {
					g.foreground(4)
					g.sln("  You think yer pretty damn smart, don't ya.  Ha!")
					g.sln("")
				} else if uniqueName(tmp) && g.checkName(tmp) {
					g.lw("  `5\"" + nname + "`5?  That's kinda flaky, you sure?\"  `2[`0Y`2] : `%")
					ch = strings.ToUpper(g.getkey())
					if ch != "N" {
						ch = "Y"
					}
					g.sln(ch)
					if ch == "Y" {
						break
					}
					g.sln("")
				}
			}
			// DIFF (lord.js): This did not happen previously,
			// so people could steal the heros name
			// and write in the dirt.
			g.getState(true)
			if g.state.LatestHero == p.Name {
				g.state.LatestHero = nname
			}
			g.putState()
			p.Name = nname
			g.sln("")
			g.sln("")
			g.lln("  `%\"Done!  You are now " + p.Name + "`%!\"")
			g.sln("")
			g.more()
		case "B":
			g.foreground(5)
			g.sln("")
			g.sln("  \"Ahh.. Bribe.. Now you are speaking my language friend!")
			g.sln("  I will let you borrow my room keys.. on one condition..")
			g.sw("  That ya pay me " + prettyInt(p.Level*1600) + " gold!!\"")
			g.lw("   `2Deal?  [`0N`2] : `%")
			ch = strings.ToUpper(g.getkey())
			if ch != "Y" {
				ch = "N"
			}
			g.sln(ch)
			if ch == "N" {
				g.sln("")
				g.foreground(5)
				g.sln("  \"Fine.. Forget I offered that deal to you..\"")
				break
			}
			g.sln("")
			g.sln("")
			if p.Gold < (800 * (p.Level * 2)) {
				g.foreground(4)
				g.sln("  \"Hey! You slobbering idiot! You don't have that much gold!\"")
				break
			}
			p.Gold -= (800 * (p.Level * 2))
			g.attackInInn()
			if p.Dead {
				done = true
			}
		case "D":
			if p.Level == 12 {
				g.foreground(5)
				g.sln("")
				g.sln("  \"Ahhh...If anyone could kill that Dragon, it would be you..")
				g.sln("  I have heard rumors that strange noises have been heard in the")
				g.sw("  forest... You might try (")
				g.lw("`2S`5")
				g.sln(")earching the forest...My son did,\"")
				g.sln("  and never came back.")
				g.sln("")
				g.moreNomail()
			}
		}
	}

	g.sln("")
}

// talkWithBartenderGems is the 'G' (gems for elixirs) case of
// talk_with_bartender().
func (g *Game) talkWithBartenderGems() {
	p := g.player
	g.sln("")
	g.foreground(5)
	g.lln("  \"You have `%Gems`5, eh?  I'll give ya a pint of magic elixir for two.\"")
	g.foreground(2)
	g.lw("  Buy how many elixirs?  [`0" + prettyInt(p.Gem/2) + "`2] : ")
	g.foreground(15)
	ch := g.consoleGetstr(GetStrOpts{Edit: itoa(p.Gem / 2), Len: 10, Integer: true})
	if ch == "" {
		ch = "0"
	}
	i, ok := jsParseInt(ch)
	if !ok {
		i = 0
	}
	if i == 0 {
		g.foreground(5)
		g.sln("")
		g.sln("  \"They are vary valuable, please reconsider it..\"")
		return
	}
	if i > 2000000 {
		g.foreground(5)
		g.sln("")
		g.sln("  \"Ha!  There isn't enough brew in all the world for that much!\"")
		return
	}
	if float64(p.Gem)/2 < float64(i) {
		g.sln("")
		g.foreground(5)
		g.lln("  \"You don't have `%" + prettyInt(i*2) + " `5Gems.  Find some...\"")
		return
	}
	if i < 0 {
		g.sln("")
		g.foreground(5)
		g.sln("  \"Ha.. Like I'm really gonna give you a negetive amount, idiot.\"")
		return
	}
	p.Gem -= i * 2
	g.foreground(2)
	g.sln("")
	switch i {
	case 1:
		ch = "small cup"
	case 2, 3, 4, 5:
		ch = "steaming mug"
	case 6, 7, 8, 9:
		ch = "huge tankard"
	default:
		ch = "keg 'o brew"
	}
	g.foreground(2)
	g.sln("  The bartender retrieves a " + ch + " from the back room.")
	g.sln("  Before you drink it, what do you wish for?")
	g.sln("")
	g.lln("  `2(`0H`2)it Points")
	g.lln("  `2(`%S`2)trength")
	g.lln("  `2(`#V`2)itality")
	g.commandPrompt()
	ch = strings.ToUpper(g.getkey())
	if !strings.Contains("HSV", ch) {
		ch = "V"
	}
	g.sln(ch)
	g.foreground(15)
	g.sln("")
	switch ch {
	case "H":
		g.sln("  YOU DRINK THE BREW AND FEEL REFRESHED!")
		g.sln("")
		g.addHP(i)
	case "V":
		g.sln("  YOU DRINK THE BREW AND YOUR SOUL REJOICES!")
		g.sln("")
		g.addDef(i)
	case "S":
		g.sln("  YOU DRINK THE BREW AND FEEL STRONGER!")
		g.sln("")
		g.addStr(i)
	}
}

// initBar is init_bar(): create the bar conversation file from one of the
// starting files if it doesn't exist yet.
func (g *Game) initBar(darkhorse bool) {
	bfnames := []string{"start1.lrd", "start2.lrd", "start3.lrd", "start4.lrd", "start5.lrd"}
	fname := g.gamedir(pick(darkhorse, "darkbar.lrd", "bar.lrd"))

	g.fmutex(fname)
	if !fileExists(fname) {
		if darkhorse {
			g.innFileCopy("dstart.lrd", fname)
		} else {
			g.innFileCopy(bfnames[g.random(int64(len(bfnames)))], fname)
		}
	}
	g.rfmutex(fname)
}

// innFileCopy is file_copy(game_or_exec(src), dst), reading the
// distribution file through the built-in assets when needed.
func (g *Game) innFileCopy(src, dst string) bool {
	b, err := g.assetBytes(src)
	if err != nil {
		return false
	}
	return os.WriteFile(dst, b, 0o644) == nil
}

// converse is converse(): the conversation at the bar (darkhorse: the
// Dark Cloak Tavern's).
func (g *Game) converse(darkhorse bool) {
	p := g.player
	fname := g.gamedir(pick(darkhorse, "darkbar.lrd", "bar.lrd"))

	g.initBar(darkhorse)
	g.sclrscr()
	g.lln("  `%Conversation at the Bar`#")
	g.lln("`l")
	g.fmutex(fname)
	lines, ok := readLines(fname)
	if !ok {
		g.rfmutex(fname)
		g.fatalf("Unable to open %s", fname)
	}
	for _, l := range lines {
		g.lln(l)
	}
	g.rfmutex(fname)
	g.sln("")
	g.lw("  `2(`5C`2)ontinue  (`5A`2)dd to Conversation `0[`5C`0] : ")
	ch := strings.ToUpper(g.getkey())
	if ch != "A" {
		ch = "C"
	}
	g.sln(ch)
	if ch == "A" {
		g.sln("")
		g.foreground(2)
		g.sln(" Share your feelings now.. (Max 75 char!)")
		g.sln("                                                                       Max")
		g.sln("                                                                         ^")
		g.lw(" `0>`2")
		l := cleanStr(g.consoleGetstr(GetStrOpts{}))
		if len(l) < 2 {
			g.sln("  You decide not to speak..You really don't have anything to say.")
			g.sln("  (ENTRY NOT ENTERED)")
			return
		}
		if len(l) > 75 {
			g.sln("  You talk so long, that it bores the other warriors!")
			g.sln("  (ENTRY NOT ENTERED)")
			return
		}
		g.fmutex(fname)
		all, ok := readLines(fname)
		if !ok {
			g.rfmutex(fname)
			g.fatalf("Unable to open %s", fname)
		}
		all = append(all, "  `%"+p.Name+":")
		all = append(all, "  `2"+l)
		for len(all) > 18 {
			all = all[1:]
		}
		if err := writeLines(fname, all...); err != nil {
			g.rfmutex(fname)
			g.fatalf("Unable to open %s", fname)
		}

		if !darkhorse {
			g.getState(true)
			g.state.LastBar = int64(p.Record)
			g.putState()
		}

		g.rfmutex(fname)
	}
}
