package game

import (
	"runtime"
	"strings"
)

// conjugalityList is conjugality_list(): who is married to whom.
func (g *Game) conjugalityList() {
	some := false
	l := 0
	phrase := []string{"hitched with", "attached to", "in love with", "a love slave to",
		"smitten with love for", "in matrimony with", "in marital bliss with", "wedded to"}

	g.lln("`c                     `%** CONJUGALITY LIST **")
	g.lln("                    `2-=-=-=-=-=-=-=-=-=-=-=-=-")
	g.player.Put()
	g.sln("")
	g.getState(false)
	for i, op := range g.allPlayers() {
		// DIFF: This check if name length was over 2 chars...
		if op.Name != "X" {
			if op.MarriedTo > -1 && int64(i) < op.MarriedTo {
				// DIFF: lord.js throws if the spouse's record is missing;
				// such a marriage is skipped here.
				if mt := g.playerGet(op.MarriedTo); mt != nil {
					some = true
					g.lln("  `0" + op.Name + " `2is " + phrase[l] + " `0" + mt.Name + "`2.")
					l++
					if l >= len(phrase) {
						l = 0
					}
				}
			}
			if int64(i) == g.state.MarriedToSeth {
				some = true
				g.lln("  `0" + op.Name + " `2is the property of `0Seth Able`2.")
			}
			if int64(i) == g.state.MarriedToViolet {
				some = true
				g.lln("  `0" + op.Name + " `2belongs to `0Violet`2.")
			}
		}
	}
	if !some {
		g.lln("  `%No one is married in this realm.")
	}
	g.sln("")
	g.moreNomail()
}

// showGameStats is show_game_stats().
func (g *Game) showGameStats() {
	s := &g.Settings

	g.getState(false)
	g.lln("`c                     `%** GAME STATISTICS " + Version + " **")
	g.sln("")
	g.lln(" `2 This game has been running for `%" + prettyInt(g.state.Days) + "`2 days.")
	g.lln("  `2You are playing a `0J`2A`0V`2A`0S`2C`0R`2I`0P`2T game.")
	if s.WinDeeds > 0 {
		g.lln("  `2Game can be `0FINISHED `2by getting " + prettyInt(s.WinDeeds) + " heroic deeds")
	} else {
		g.lln("  `2Game is set to run indefinitely.")
	}
	g.sln("")
	var tmp int64
	for _, op := range g.allPlayers() {
		if op.Name != "X" {
			tmp++
		}
	}
	g.lln("  `2There are currently `%" + prettyInt(tmp) + " `2people playing.")
	g.sln("")
	g.lln("  `2Players are deleted after `%" + prettyInt(s.DeleteDays) + " `2days of inactivity.")
	g.sln("")
	if s.CleanMode {
		g.lln("  `2Clean mode is `0ENABLED`2.")
	} else {
		g.lln("  `2Clean mode is OFF!  (Thank the stars!)")
	}
	g.sln("")
	if s.TransfersOn {
		g.lln("  `2Transferring Funds is enabled, max of `%" + prettyInt(s.TransferAmount) + "`2 per transfer.")
		if s.TransfersPerDay > 0 {
			g.lln("  `2You can do this up to `%" + prettyInt(s.TransfersPerDay) + "`2 times a day.")
		} else {
			g.lln("  `2You can do this `%unlimited`2 times per day.")
		}
	} else {
		g.lln("  `2Bank option Transferring Funds is disabled by Sysop")
	}
	g.sln("")
	g.lln("  `2Need an ego boost?  This is what you overhear being said about yourself.")
	if g.player.Sex == "M" {
		g.lln("  `0\"He " + g.looks(g.player.Sex, g.player.Cha) + "\"")
	} else {
		g.lln("  `0\"She " + g.looks(g.player.Sex, g.player.Cha) + "\"")
	}
	g.sln("")
	// DIFF: system.os_version is Synchronet's host OS string; the Go
	// runtime's OS and architecture stand in for it.
	g.lln("  `0Current user tasker: " + runtime.GOOS + "/" + runtime.GOARCH + " and you're on node " + prettyInt(int64(g.Ctx.Node)))
	g.lw("`2  ")
	if !s.Del1XP {
		g.lw("NK1 ") // No Kill 1XP
	}
	if s.ResDays != 3 {
		g.lw("R" + itoa(s.ResDays) + " ") // Resurection Days
	}
	if !s.Olivia {
		g.lw("NO ") // No Olivia
	}
	if !s.FunkyFlowers {
		g.lw("BF ") // Boring Flowers
	}
	if !s.ShopLimit {
		g.lw("NL ") // No Shop Limit
	}
	if s.OldSkillPoints {
		g.lw("5S ") // 5 skill points per use
	}
	if s.DefForPK {
		g.lw("DP ") // Def for pk
	}
	if s.StrForPK {
		g.lw("SP ") // Str for pk
	}
	if s.BeefUp {
		g.lw("BEEF ") // Beef up monsters
	}
	if !s.OldSteal {
		g.lw("CT ") // Crappy Thieves
	}
	if s.DKBoost {
		g.lw("DK ") // Great Death Knights
	}
	if s.SleepDragon {
		g.lw("SD ") // Anti-camping Dragon attacks
	}
	g.sln("")
	g.moreNomail()
}

// listPlayers is list_players(): the rankings with who is online.
func (g *Game) listPlayers() {
	f := g.gamedir("temp" + itoa(g.player.Record))

	g.sln("")
	g.generateRankings(f, true, true, false)
	g.displayFile(f, true, true)
	fileRemove(f)
}

// mainGoodbye is main()'s goodbye(): sleep in the fields and leave.
func (g *Game) mainGoodbye() {
	g.sclrscr()
	g.sln("")
	g.sln("")
	g.lln("  `%Quitting To The Fields...")
	g.lln("`0-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
	g.lln("`2  You find a comfortable place to sleep under a small tree...")
	g.sln("")
	g.player.OnNow = false
	g.player.Put()
	g.generateRankings(g.gamedir("scores.lrd"), true, false, false)
	g.exit(0)
}

// mainDoQuit is main()'s do_quit().
func (g *Game) mainDoQuit() bool {
	g.sln("")
	g.sln("")
	g.lw("`2  Quit game?  [`0Y`2] : ")
	ch := strings.ToUpper(g.getkey())
	if ch != "N" {
		ch = "Y"
	}
	g.sln(ch)
	return ch == "Y"
}

// mainPrompt is main()'s prompt().
func (g *Game) mainPrompt() {
	g.sln("")
	g.foreground(5)
	g.sln("  The Town Square")
	g.foreground(8)
	g.sln("  (? for menu)")
	g.sln("  (F,S,K,A,H,V,I,T,Y,L,W,D,C,O,X,M,P,Q)")
	g.commandPrompt()
}

// mainMenu is main()'s menu().
func (g *Game) mainMenu() {
	p := g.player
	if p.Bank < 0 {
		p.Bank = 0
	}
	if p.Gold < 0 {
		p.Gold = 0
	}
	if p.Exp < 0 {
		p.Exp = 0
	}
	if !p.Expert {
		g.lrdfile("MAIN")
	}
}

// mainOtherPlacesMenu is main()'s other_places_menu(), as shown with no
// IGMs installed (3rdalt.lrd/3rdparty.lrd are IGM lists).
func (g *Game) mainOtherPlacesMenu() {
	g.lln(" `%`r0`c  Other Options")
	g.lln("`2-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
	g.lln("  `0(`2Currently Undeveloped Land`0)")
}

// main is main(): the Town Square.
func (g *Game) main() {
	quit := false
	var ch string

	g.foreground(5)
	g.sln("")

	if g.player.Inn {
		g.redDragonInn()
		if g.player.Dead {
			return
		}
	}

	g.mainMenu()
	for {
		g.foreground(5)
		g.checkMail()
		g.mainPrompt()
		ch = strings.ToUpper(g.getkey())
		g.sln(ch)
		switch ch {
		case "V":
			g.showStats()
			g.mainMenu()
		case "D":
			g.showLog()
			g.mainMenu()
		case "R":
			if g.mailCheck() {
				g.checkMail()
			} else {
				g.sln("")
				g.sln("  You have no mail.")
			}
			g.sln("")
		case "?":
			// DIFF: Couldn't see menu in expert mode before...
			if g.player.Expert {
				g.lrdfile("MAIN")
			}
			g.mainMenu()
		case "K":
			g.kingArthurs()
			g.mainMenu()
		case "S":
			g.slaughterOthers()
			if g.player.Dead {
				quit = true
			}
			g.sln("")
		case "A":
			g.abdulsArmour()
			g.mainMenu()
			g.sln("")
		case "T":
			g.turgons()
			g.sln("")
		case "W":
			g.composeMail()
		case "L":
			g.listPlayers()
		case "O":
			g.foreground(10)
			g.background(0)
			// IGMs are not supported: there are no other places.
			if g.player.Expert {
				g.sclrscr()
				g.sln("")
				g.sln("")
			} else {
				g.mainOtherPlacesMenu()
			}
			g.checkMail()
			for {
				g.sln("")
				g.lw("  `2What's your pleasure? [`0Q`2] (`0? for menu`2) : `%")
				ch = strings.ToUpper(g.getkey())
				if ch == "\r" {
					ch = "Q"
				}
				if digitsRE.MatchString(ch) {
					// No place number is valid without IGMs.
					g.mainOtherPlacesMenu()
					g.checkMail()
				} else if ch == "?" {
					g.mainOtherPlacesMenu()
					g.checkMail()
				}
				if ch == "Q" {
					break
				}
			}
			g.sln("")
			g.sln("")
		case "F":
			g.forest()
			if g.player.Dead {
				quit = true
			}
			g.sln("")
		case "H":
			g.healers()
			g.sln("")
		case "I":
			g.redDragonInn()
			if g.player.Dead || g.player.Inn {
				quit = true
				g.sln("")
				break
			}
			g.mainMenu()
			g.sln("")
		case "2":
			// DIFF: toggles RIP mode in lord.js; RIP is not supported.
			g.mainMenu()
		case "M":
			g.announce()
		case "P":
			g.warriorsOnNow(false)
			g.sln("")
		case "1":
			g.showGameStats()
			g.sln("")
		case "C":
			g.conjugalityList()
			g.sln("")
		case "3": // TODO (ansi)
			g.sln("")
		case "4": // TODO (time left)
			g.sln("")
		case "X":
			g.player.Expert = !g.player.Expert
			g.sln("")
			g.sln("")
			if g.player.Expert {
				g.sln("  EXPERT MODE ON")
			} else {
				g.sln("  EXPERT MODE OFF")
			}
		case "Q":
			if !g.mainDoQuit() {
				break
			}
			g.mainGoodbye() // goodbye() exits...
			g.sln("")
		case " ":
			g.sln("")
		case "Y":
			g.yeOldBank()
			g.sln("")
		default:
			g.sln("")
		}
		if quit {
			break
		}
	}
}
