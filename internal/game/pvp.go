package game

import (
	"strings"
)

// attackPlayer is attack_player(): fight another (offline) player, in the
// fields or (inn) in their room at the Inn.
func (g *Game) attackPlayer(op *Player, inn bool) {
	player := g.player
	outName := g.gamedir("out" + itoa(op.Record) + ".lrd")

	if player.PvPFights < 1 {
		g.sln("")
		g.sln("  You are too tired to look for that warrior.  Try again tomorrow.")
		g.sln("")
		g.more()
		return
	}
	pronoun := pick(op.Sex == "M", "his", "her")
	g.fmutex(outName)
	if fileExists(outName) {
		lines, ok := readLines(outName)
		if !ok {
			g.fatalf("Unable to open %s", outName)
		}
		where := "null"
		if len(lines) > 0 {
			where = lines[0]
		}
		g.rfmutex(outName)
		g.sln("")
		g.sln("  That warrior is currently visiting another section of the realm.")
		g.lln("  `2(" + op.Name + " " + where + "`2)")
		g.moreNomail()
		return
	}
	g.rfmutex(outName)
	g.fmutex(g.gamedir("battle"+itoa(op.Record)+".rec"), itoa(player.Record)+"\n")
	if op.OnNow {
		g.rfmutex(g.gamedir("battle" + itoa(op.Record) + ".rec"))
		g.sln("")
		g.sln("  That warrior is currently online!")
		g.sln("")
		g.onlineBattle(op)
		return
	}
	mail := " \n  `%YOU HAVE BEEN ATTACKED!\n`0-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-`2\n  `0" + player.Name + "`2 has attacked you!"
	player.PvPFights--
	g.sln("")
	g.sln("")
	g.lln("  `2** `%PLAYER FIGHT `2**")
	g.sln("")
	g.lln("  You have encountered " + op.Name + "`2!!")
	g.sln("")
	if op.HP < op.HPMax {
		op.HP = op.HPMax
	}
	mail += g.battle(op, true, false)
	if player.Dead {
		g.sln("")
		player.Gold = 0
		player.Exp -= player.Exp / 10
		mail += "\n`.  `2You have killed `0" + player.Name + " `2in self defense!\n`.  `2You receive `%" + prettyInt(player.Exp/2) + "`2 experience!"
		op.PvP++
		if op.PvP > 32000 {
			op.PvP = 32000
		}
		op.Exp += player.Exp / 2
		if op.Exp > 2000000000 {
			op.Exp = 2000000000
		}
		op.Put()
		g.mailTo(int64(op.Record), mail)
		g.rfmutex(g.gamedir("battle" + itoa(op.Record) + ".rec"))
		g.deadScreen(op)
		g.badSay(op, "`.  `0"+op.Name+" `2has killed `5"+player.Name+" `2in self defence!")
		g.sln("")
		return
	}

	if op.HP < 1 {
		g.foreground(15)
		g.sln("")
		g.lln("  You have killed " + op.Name + "`%!")
		g.sln("")
		g.foreground(2)
		tmp := player.Gold
		player.Gold += op.Gold
		if player.Gold >= 2000000000 {
			player.Gold = 2000000000
		}
		g.lw("  You receive `%" + prettyInt(player.Gold-tmp) + " `2gold, ")
		tmp = player.Exp
		player.Exp += op.Exp / 2
		if player.Exp >= 2000000000 {
			player.Exp = 2000000000
		}
		g.lln("and " + prettyInt(player.Exp-tmp) + " experience!")
		if op.Gem >= 2 {
			tmp = player.Gem
			player.Gem += op.Gem / 2
			g.sln("")
			tmp = player.Gem - tmp
			g.lln("  `2You also find `0" + prettyInt(tmp) + " `%" + pick(tmp == 1, "Gem", "Gems") + "`2!")
		}
		if player.Def < (32000 - 3) {
			if g.Settings.DefForPK {
				g.lw("  `2You also gain `%3 `2defense points")
				player.Def += 3
				if player.Str < (32000 - 2) {
					// In 4.07, you don't actually get the str points!
					g.lw(", and you gain `%2 `2strength points")
					if g.Settings.StrForPK {
						player.Str += 2
					}
				}
			}
			g.sln("")
		}
		if !g.Settings.DefForPK && g.Settings.StrForPK {
			if player.Str < (32000 - 2) {
				g.lw("  `2You also gain `%2 `2strength points")
				player.Str += 2
			}
		}
		op.Gem -= op.Gem / 2
		op.Exp = op.Exp - op.Exp/10
		op.Gold = 0
		op.Dead = true
		op.Inn = false
		mail += "\n`.  `0" + player.Name + " `2has killed you!"

		// TODO: This is a massive pain in the arse to get right in remote games.  Don't try.
		if inn {
			if g.random(10) == 1 {
				if player.Level < (op.Level - 1) {
					if op.WeaponNum > player.WeaponNum {
						g.sln("")
						g.lln("  `2Do you wish to trade your " + player.Weapon + " `2for " + pronoun)
						g.lln("  `%" + op.Weapon + "`2? `0[N] : `%")
						ch := strings.ToUpper(g.getkey())
						if ch != "Y" {
							ch = "N"
						}
						g.sln(ch)
						if ch == "Y" {
							// DIFF: lord.js calls get_weapon(op.weapon) with the
							// weapon's name, which throws; use its number.
							wep := g.getWeapon(op.WeaponNum)
							op.Str -= wep.Num
							player.Str += wep.Num
							wep = g.getWeapon(player.WeaponNum)
							op.Str += wep.Num
							player.Str -= wep.Num
							op.Weapon, player.Weapon = player.Weapon, op.Weapon
							op.WeaponNum, player.WeaponNum = player.WeaponNum, op.WeaponNum
							mail += "\n  `$" + player.Name + " took your weapon!"
							g.lln("  `2Done! You now have a " + player.Weapon + "`2!")
						}
					}
				}
			}
		}
		g.mailTo(int64(op.Record), mail)
		g.sln("")
		player.KilledAPlayer = true
		if player.PvP < 32000 {
			player.PvP++
		}
		op.Put()
		g.rfmutex(g.gamedir("battle" + itoa(op.Record) + ".rec"))
		g.goodSay(op, "`.  `0"+player.Name+" `2has killed `5"+op.Name+"`2!")
		g.tournamentCheck()
	}

	player.LightShield = false
}

// attackInInn is attack_in_inn(): the bartender's "kill someone in the
// Inn" service.
func (g *Game) attackInInn() {
	player := g.player
	done := false
	temp := g.gamedir("temp" + itoa(player.Record))

	g.generateRankings(temp, true, true, true)
	g.sclrscr()
	g.displayFile(temp, true, true)
	fileRemove(temp)
	g.sln("")
	for !done {
		g.sln("")
		g.foreground(5)
		g.sln("  \"What next, kid?\"")
		g.sln("")
		g.foreground(2)
		g.lln("  `2(`5L`2)ist Warriors in the Inn")
		g.lln("  `2(`5S`2)laughter Warriors in the Inn")
		g.lln("  `2(`5R`2)eturn to Bar")
		g.commandPrompt()
		ch := strings.ToUpper(g.getkey())
		g.sln(ch)
		switch ch {
		case "L":
			g.generateRankings(temp, true, true, true)
			g.displayFile(temp, true, true)
			fileRemove(temp)
		case "R":
			g.foreground(5)
			g.sln("")
			g.sln("  \"So ya aren't going to attack anyone?  Here is half your money back...")
			g.lln("'   The other half will cover my time!\" `2the bartender laughs harshly.")
			player.Gold += 400 * (player.Level * 2)
			g.sln("")
			g.moreNomail()
			done = true
		case "S":
			g.sln("  Who would you like to attack?")
			atk := g.findPlayer()
			if atk == -1 {
				g.sln("  No warriors found.")
				break
			}
			g.sln("")
			op := g.playerGet(atk)
			pronoun := pick(op.Sex == "M", "he", "she")
			if atk == int64(player.Record) {
				g.sln("  You wish to attack yourself?!!  You decide against it.")
				break
			}
			if op.Dead {
				g.sln("  That warrior isn't at the Inn at the moment.")
				g.sln("  You recall seeing in the news that " + pronoun + " was dead..")
				break
			}
			if !op.Inn {
				g.sln("  That warrior is not staying at the Inn today,")
				g.sln("  " + pronoun + " is probably in the fields.")
				break
			}
			if op.Level+1 < player.Level {
				g.sln("  A child could beat that wimp!  Attack someone else!")
				break
			}
			if op.Sex == "M" {
				g.lln("  You enter " + op.Name + "`2s room...He is sleeping.")
				g.lln("  You notice he has a dangerous looking " + op.Weapon + " `2by his bed..")
				g.lln("  Are you sure you want to attack him?")
			} else {
				g.lln("  You enter " + op.Name + "`2s room...She is sleeping.")
				g.lln("  You notice she has a dangerous looking " + op.Weapon + " `2by her bed..")
				g.lln("  Are you sure you want to attack her?")
			}
			g.sln("")
			g.lw("  `2Attack `5" + op.Name + " `2[`0Y`2] :`0")
			ch = strings.ToUpper(g.getkey())
			if ch != "N" {
				ch = "Y"
			}
			g.sln(ch)
			if ch == "N" {
				break
			}
			g.sln("")
			g.sln("  Being a trained warrior, " + pronoun + " jumps up suddenly, aware of your")
			g.sln("  presence!")
			g.mailTo(atk, "  `%YOU HAVE BEEN ATTACKED IN YOUR ROOM!\n"+
				"`0-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-`2\n"+
				"  `0"+player.Name+"`2 has broken into your room at the Inn!")
			g.more()
			g.attackPlayer(op, true)
			done = true
		}
	}
}

// slaughterOthers is slaughter_others(): the player-fight menu.
func (g *Game) slaughterOthers() {
	player := g.player
	dirt := g.gamedir("dirt.lrd")
	temp := g.gamedir("temp" + itoa(player.Record))

	topKiller := func() *Player {
		// An object that only has pvp: its name is undefined.
		ret := &Player{Name: "undefined", PvP: -1}
		for _, pl := range g.allPlayers() {
			if pl.Name != "X" && pl.PvP > ret.PvP {
				ret = pl
			}
		}
		return ret
	}

	g.sclrscr()
	if !player.Expert {
		g.lrdfile("WAR")
	}
	for {
		g.sln("")
		g.foreground(5)
		g.sln("  Slaughter Other Players")
		g.foreground(2)
		g.sln("")
		g.sln("  (S,L,E,W,R)  (? for menu)")
		g.checkMail()
		g.commandPrompt()
		ch := strings.ToUpper(g.getkey())
		g.sln(ch)
		switch ch {
		case "L":
			g.generateRankings(temp, false, true, false)
			g.displayFile(temp, true, true)
			fileRemove(temp)
		case "R", "Q", "\r":
			return
		case "?":
			g.lrdfile("WAR")
		case "V":
			g.showStats()
		case "S":
			g.sln("  Who would you like to attack?")
			atk := g.findPlayer()
			if atk == -1 {
				g.sln("  No warriors found.")
			} else if atk == int64(player.Record) {
				g.sln("  You wish to attack yourself?!!  You decide against it.")
			} else {
				enemy := g.playerGet(atk)
				outName := g.gamedir("out" + itoa(enemy.Record) + ".lrd")
				g.fmutex(outName)
				if fileExists(outName) {
					lines, ok := readLines(outName)
					if !ok {
						g.fatalf("Unable to open %s", outName)
					}
					where := "null"
					if len(lines) > 0 {
						where = lines[0]
					}
					g.rfmutex(outName)
					g.sln("")
					g.sln("  That warrior is currently visiting another section of the realm.")
					g.lln("  `2(" + enemy.Name + " " + where + "`2)")
					g.sln("")
					g.moreNomail()
				} else if enemy.Dead {
					g.rfmutex(outName)
					g.sw("  You look for that warrior...And you find ")
					g.sln(pick(enemy.Sex == "M", "him...", "her..."))
					g.sln("  A rotting corpse...Looks like you were a little late..")
				} else if enemy.Inn {
					g.rfmutex(outName)
					g.sln("  You search the fields but do not find that warrior.")
					g.sw("  You conclude ")
					if enemy.Sex == "F" {
						g.sw("s")
					}
					g.sln("he is staying at the Inn.")
				} else {
					g.rfmutex(outName)
					player.Put()
					g.lw("  `2You hunt around for `0" + enemy.Name + "`2...")
					if enemy.Sex == "M" {
						g.sln("YOU FIND HIM!")
						g.sw("  He")
					} else {
						g.sln("YOU FIND HER!")
						g.sw("  She")
					}
					g.lln(" is brandishing a dangerous looking " + enemy.Weapon + ".")
					g.sln("")
					g.lw("  `2Attack `5" + enemy.Name + " `2[`0Y`2] : `%")
					ch = strings.ToUpper(g.getkey())
					if ch != "N" {
						ch = "Y"
					}
					g.sln(ch)
					g.foreground(2)
					if ch == "Y" {
						g.attackPlayer(enemy, false)
					}
					if player.Dead {
						return
					}
				}
			}
		case "E":
			g.lln("`c`%                         ** EXAMINING THE DIRT **")
			g.sln("")
			tk := topKiller()
			g.lln("  `2The town elders have named `0" + tk.Name + " `2the most dangerous")
			g.lln("  `2warrior in the realm with a whopping `%" + prettyInt(tk.PvP) + "`2 kills.")
			g.sln("")
			if fileExists(dirt) {
				g.fmutex(dirt)
				g.showLordFile(dirt, false, false)
				g.rfmutex(dirt)
			} else {
				g.lln("  `2The dirt appears soft and malleable.")
			}
		case "W":
			g.lln("`c `%                        ** WRITING IN THE DIRT **")
			g.sln("")
			if (!player.KilledAPlayer) && player.Name != g.state.LatestHero {
				g.lln("  `2You are about to inscribe something, when you hear `0" + g.state.LatestHero + "'`2s")
				g.sln("  voice in your head!")
				g.sln("")
				g.lln("  `0\"My son, you have not challenged and killed another warrior this day.\"")
				g.sln("")
				g.lln("  `2You lower your eyes in homage to your better.")
				g.sln("")
				g.lln("  `0\"You may share wisdom when you have killed, and have it.\"")
				g.lln("")
				g.more()
			} else {
				// BUG (lord.js): compares with state.lastesthero (a typo, always
				// undefined), so the dragon slayer's greeting never shows.
				if false {
					g.lln("  `2The crowd parts in awe, as you approach the dirt patch.  You smile, maybe")
					g.sln("  slaying that Dragon did more for your reputation than you know.")
					g.sln("")
				}
				g.lln("  `2You pick up a nearby stick and inscribe some of your wisdom.")
				g.sln("")
				var lines []string
				// DIFF: Previously, no max length on dirt...
				for {
					g.lw("`. `0>`%")
					line := "  " + cleanStr(g.consoleGetstr(GetStrOpts{}))
					lines = append(lines, line)
					if line == "  " || len(lines) >= 20 {
						break
					}
				}
				g.fmutex(dirt)
				_ = writeLines(dirt, lines...)
				g.rfmutex(dirt)
			}
		}
	}
}
