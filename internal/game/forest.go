package game

import (
	"strings"
)

// forest is forest(): the forest menu and its events.
func (g *Game) forest() {
	var ch string
	over := false

	g.sclrscr()
	g.forestMenu()
	for {
		g.forestPrompt()
		g.commandPrompt()
		ch = strings.ToUpper(g.getkey())
		if ch != "J" {
			g.sln(ch)
		}
		switch ch {
		case "B":
			if g.player.Gold > 0 {
				g.sln("")
				g.sln("")
				g.sln("  You throw your gold pouch up into the air gleefully.")
				g.sln("")
				g.lln("  `0AN UGLY VULTURE `)GRABS `0IT IN MID AIR!")
				g.player.Bank += g.player.Gold
				if g.player.Bank > 2000000000 {
					g.player.Bank = 2000000000
				}
				g.player.Gold = 0
			}
		case "R", "Q", "\r":
			over = true
		// DIFF: Removed explicit space handling...
		case "V":
			g.showStats()
		case "L":
			if g.player.ForestFights < 1 {
				g.sln("")
				g.sln("")
				g.sln("  You are too tired.")
				g.sln("")
				g.sln("  Try again tomorrow.")
				g.sln("")
				g.more()
			} else {
				g.forestLookToKill(&ch)
				if g.player.Dead {
					over = true
				}
			}
		case "T":
			if g.player.Horse {
				g.sln("")
				g.lln("  `2You nudge your horse deeper into the woods.")
				g.sln("")
				g.moreNomail()
				g.sln("")
				g.sln("")
				g.foreground(15)
				g.sln("  Event In The Forest")
				g.lln("`0-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
				g.foreground(2)
				g.sln("  In the gloom of the shady forest, you see smoke coming from a bright")
				g.sln("  chimney.")
				g.sln("")
				g.more()
				g.forestDarkhorseTavern(&ch)
				g.sclrscr()
				g.forestMenu()
				break
			}
			g.sln("")
			g.sln("")
			g.sln("  Your Thieving skills cannot help you here.")
			g.sln("")
		case "M":
			g.sln("")
			g.sln("")
			g.sln("  Your Mystical skills cannot help you here.")
			g.sln("")
		case "D":
			g.sln("")
			g.sln("")
			g.sln("  Your Death Knight skills cannot help you here.")
			g.sln("")
		case "A":
			g.sln("")
			g.sln("")
			g.sln("  You brandish your weapon dramatically.")
			g.sln("")
		case "H":
			g.healers()
		case "?":
			if g.player.Expert {
				g.lrdfile("FOREST")
			}
			g.forestMenu()
		case "S":
			if g.player.Level < 12 {
				g.showStats()
			} else {
				g.forestAttackDragon()
				if g.player.Dead {
					over = true
					break
				}
			}
			g.sclrscr()
			g.forestMenu()
		case "J":
			if g.player.HighSpirits {
				g.player.HighSpirits = false
				if strings.ToUpper(g.getkey()) == "E" &&
					strings.ToUpper(g.getkey()) == "N" &&
					strings.ToUpper(g.getkey()) == "N" &&
					strings.ToUpper(g.getkey()) == "I" &&
					strings.ToUpper(g.getkey()) == "E" {
					g.sln("")
					g.sln("")
					g.lln("  `0Jennie?  Jennie Garth?")
					g.lw("  `2Define her. ")
					ch = cleanStr(strings.ToUpper(g.getstr(0, 0, 4, 1, 15, "", GetStrOpts{})))
					g.curlinenum = 1
					g.sln("")
					g.sln("")
					switch ch {
					case "BABE":
						g.lln("  `0That is correct. `2(YOU RECIEVE AN EXTRA FOREST FIGHT!)")
						g.player.ForestFights += 1
						if g.player.ForestFights > 32000 {
							g.player.ForestFights = 32000
						}
					case "SEXY":
						g.lln("  `0Exellent. `2(YOU RECIEVE AN EXTRA USER FIGHT!)")
						g.player.PvPFights += 1
						if g.player.PvPFights > 32000 {
							g.player.PvPFights = 32000
						}
					case "LADY":
						g.lln("  `0Very true.  `2(YOU GET SOME GOLD!)")
						g.player.Gold += 1000 * g.player.Level
						if g.player.Gold > 2000000000 {
							g.player.Gold = 2000000000
						}
					case "DUMB":
						g.lln("  `0You idiot.  You will `)never`0 be a useful member of society.")
					case "STAR":
						g.lln("  `0A huge star, infant.")
						g.lln("  `4(YOU GET NOTHING, YOU STATED THE OBVIOUS)")
					case "DUNG":
						if g.player.Sex == "M" {
							g.lln("  `0You are a fool, sir.  `4(YOU ARE TURNED INTO A FROG)")
						} else {
							g.lln("  `0You are a fool, woman.  `4(YOU ARE TURNED INTO A FROG)")
						}
						g.lln("")
						g.moreNomail()
						for {
							g.lln("  `c`2The Forest Floor")
							g.sln("")
							g.lln("  `2(`0H`2)op Like Crazy")
							g.lln("  `2(`0A`2)pologize")
							g.sln("")
							g.lw("  `2Your command, greeny? : `%")
							ch = strings.ToUpper(g.getkey())
							if ch != "A" {
								ch = "H"
							}
							g.sln(ch)
							if ch != "A" {
								g.foreground(2)
								g.sln("  You hop around like a crazy frog.  What")
								g.sln("  is this accomplishing, pray tell?")
								g.sln("")
								g.moreNomail()
							}
							if ch == "A" {
								break
							}
						}
						g.sln("")
						g.lln("  `2You apologize humbly, knowing what you did was wrong.")
						g.lln("")
						g.lln("  `%(YOU ARE CHANGED BACK TO YOUR (MOSTLY) HUMAN FORM)")
					case "FOXY":
						g.lln("  `0Very wise. `%(YOU RECIEVE AN EXTRA GEM!)")
						g.player.Gem += 1
						if g.player.Gem > 32000 {
							g.player.Gem = 32000
						}
					// TODO: Toggle for old (opposite) behaviour
					case "FAIR":
						g.lln("  `0Very fair. `2(YOU FEEL EXCITED!)")
						g.player.Flirted = false
					case "UGLY":
						g.player.OnNow = false
						g.player.Put()
						g.lln("  `0You understand nothing.  `4(YOU ARE BITCH SLAPPED!)")
						g.player.HP = 1
						g.sln("")
						g.moreNomail()
						g.player.OnNow = false
						g.exit(0)
					case "HOTT":
						g.lln("  `0\"Hot\" is spelled with only one T.. But good job, nonetheless.")
						g.sln("")
						g.player.HP = int64(float64(g.player.HPMax) + float64(g.player.HPMax)/5)
						if g.player.HP > 32000 {
							g.player.HP = 32000
						}
						g.lln("  `%(YOU FEEL ENERGIZED!)")
					case "COOL":
						g.lln("  `0Why, you are cool to notice that.")
						g.sln("")
						if g.player.HP < g.player.HPMax {
							g.lln("  `%GOD NOTICES YOU ARE WOUNDED AND PITIES YOU.  YOU LOOK BETTER!")
							g.player.Cha += 1
							if g.player.Cha > 32000 {
								g.player.Cha = 32000
							}
						}
					// TODO: Option to turn off
					case "GIFT":
						g.lln("  `0Yes, she is this.  And now, a magical gift for you.")
						g.sln("")
						switch g.player.Clss {
						case 2:
							if g.player.SkillM < 1 || g.player.MagicallyDelicious {
								g.lln("  `%You are unable to accept the gift.")
							} else {
								g.lln("  `5YOU FEEL MAGICALLY DELICIOUS.")
								g.player.LevelM = g.player.SkillM
								g.player.MagicallyDelicious = true
							}
						case 1:
							if g.player.SkillW < 1 || g.player.MagicallyDelicious {
								g.lln("  `%You are unable to accept the gift.")
							} else {
								g.lln("  `5YOU FEEL MAGICALLY DELICIOUS.")
								g.player.LevelW = g.player.SkillW
								g.player.MagicallyDelicious = true
							}
						case 3:
							if g.player.SkillT < 1 || g.player.MagicallyDelicious {
								g.lln("  `%You are unable to accept the gift.")
							} else {
								g.lln("  `5YOU FEEL MAGICALLY DELICIOUS.")
								g.player.LevelT = g.player.SkillT
								g.player.MagicallyDelicious = true
							}
						}
					// TODO 'NICE' (Duh (nothing special) only in versions prior to 4)
					default:
						if g.player.Sex == "M" {
							g.sln("  You do not understand her, my son.")
						} else {
							g.sln("  Perhaps if you were male you might understand better.")
						}
					}
					g.sln("")
					g.moreNomail()
				}
			}
		}
		if over {
			break
		}
	}
}

// forestEventHeader is forest()'s event_header().
func (g *Game) forestEventHeader() {
	g.sln("")
	g.sln("")
	g.foreground(15)
	g.sln("  Event In The Forest")
	g.lln("`0-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
	g.foreground(2)
}

// forestLoadMonster is forest()'s load_monster(): lenemy.bin when the sysop
// supplies one, else a copy of the built-in table.
func (g *Game) forestLoadMonster(num int64) *Player {
	fname := g.gameOrExec("lenemy.bin")
	var ret *Player

	if fname != "" && fileExists(fname) && !g.Settings.NoLEnemyBin {
		m, _, err := readRecordAt[monsterRec](fname, int(num))
		if err != nil {
			g.fatalf("%v", err)
		}
		ret = &Player{
			Name:   m.Name,
			Str:    m.Str,
			Gold:   m.Gold,
			Weapon: m.Weapon,
			Exp:    m.Exp,
			HP:     m.HP,
			Death:  m.Death,
		}
	} else {
		m := monsterStats[num]
		ret = &m
	}
	ret.Record = -1
	g.beefUp(ret)
	return ret
}

// forestSaveForestGold is forest()'s save_forest_gold(): a fifteenth of a
// dead player's gold goes to the forest for find_lost_gold().
func (g *Game) forestSaveForestGold(amt int64) {
	amt = amt / 15
	g.getState(true)
	g.state.ForestGold += amt
	g.putState()
}

// forestLookToKill is forest()'s look_to_kill(): a forest event or a
// monster.  ch is forest()'s `ch` (see forestDarkhorseChance).
func (g *Game) forestLookToKill(ch *string) {
	var mnum int64
	var ich string
	var tmp int64

	g.sclrscr()
	if g.random(5) == 1 {
		// TODO: Allow JS events...
		switch g.random(15 + pick(g.player.Horse, int64(1), 0)) {
		case 0:
			g.lln("`c  `%Event In The Forest`0")
			g.lln("`l")
			g.sln("")
			g.lln("  `2You come across an old man.  He seems confused and asks if")
			g.sln("  you would direct him to the Inn.  You know that if you do,")
			g.sln("  you will lose time for one fight today.")
			g.sln("")
			g.lw("  Do you take the old man? [`0Y`2] : ")
			ich = strings.ToUpper(g.getkey())
			if ich != "N" {
				ich = "Y"
			}
			g.sln(ich)
			g.sln("")
			if ich == "Y" {
				g.sln("")
				tmp = g.player.Level * 500
				g.player.Gold += tmp
				g.lln("  `2You gladly take the old man to the Inn. He is pleased")
				g.lln("  with you, and gives you `%" + prettyInt(tmp) + " `2gold!")
				g.player.Cha += 1
				g.player.ForestFights -= 1
				g.sln("")
				g.lln("  `%**CHARM GOES UP BY 1**")
			} else {
				if g.player.Cha < 10 {
					g.lln("  `2\"`0I don't have time for you old man..Goodbye..`2\"")
					g.sln("  You tell him coldly.  The old man shakes his head very sadly.")
				} else {
					g.lln("  `2\"`0Old man, I think it's time you learned the way yourself.`2\"")
					g.lln("  You tell him coldly.  The old man shakes his head sadly.")
				}
			}
			g.sln("")
			g.moreNomail()
		case 1:
			g.lln("`c  `%Event In The Forest`0")
			g.lln("`l")
			g.sln("")
			if g.player.HP == g.player.HPMax {
				g.lln("  `2You come across an ugly old hag.  `5\"Give me a gem!  You won't be sorry")
				g.lln("  my pet!\"`2 she screeches.")
			} else {
				g.lln("  `2You come across an ugly old hag.  `5\"Give me a gem and I will completely")
				g.lln("  heal you warrior!\"`2 she screeches.")
			}
			g.sln("")
			g.lw("  Give her the gem? [`0N`2] : ")
			ich = strings.ToUpper(g.getkey())
			if ich != "Y" {
				ich = "N"
			}
			g.sln(ich)
			g.sln("")
			if ich == "Y" {
				if g.player.Gem < 1 {
					g.lln("  `5\"You have no gems, fool!\"`2  the old woman screams at you.  She then")
					g.lln("   hits you in a sensitive spot with her cane, and you feel VERY weak.")
					g.player.HP = 1
				} else {
					g.player.Gem -= 1
					g.sln("  You give her a gem.  She waves her wand strangely.")
					g.sln("")
					g.lln("  `%YOU FEEL BETTER")
					if g.player.HP < g.player.HPMax {
						g.player.HP = g.player.HPMax
					} else {
						g.player.HPMax += 1
					}
				}
			} else {
				g.lln("  `5\"Hurumph!\" `2 the old woman grunts sourly as you leave.")
			}
			g.sln("")
			g.moreNomail()
		case 2:
			g.lln("`c  `%Event In The Forest`0")
			g.lln("`l")
			g.sln("")
			tmp = (g.random(500) + 250) * g.player.Level * g.player.Level
			g.lln("  `2You find a sack with `%" + prettyInt(tmp) + " `2gold in it!")
			g.player.Gold += tmp
			g.sln("")
			g.moreNomail()
			g.sclrscr()
		case 3:
			g.player.HP = g.player.HPMax
			g.lln("`c  `%Event In The Forest`0")
			g.lln("`l")
			g.sln("")
			g.lln("  `2You stumble upon a group of Merry Men!  After partying with them")
			g.sln("  for 2 hours you feel totally refreshed.")
			g.sln("")
			g.moreNomail()
			g.sclrscr()
		case 4:
			g.lln("`c  `%Event In The Forest`0")
			g.lln("`l")
			g.sln("")
			g.lln("  `2Fortune smiles, and you find a gem!")
			g.sln("")
			g.player.Gem += 1
			g.moreNomail()
			g.sclrscr()
		case 5:
			if g.random(2) == 1 {
				g.forestFlowerGarden()
			} else {
				g.lln("`c  `%Event In The Forest`0")
				g.lln("`l")
				g.sln("")
				g.lln("  `2You find a `%Hammer Stone`2!")
				g.sln("")
				g.moreNomail()
				g.lln("  `2As is the tradition, you break it in two with your")
				g.lln("  `0" + g.player.Weapon + "`2. Your weapon tingles!")
				g.sln("")
				g.lln("  `%ATTACK STRENGTH RAISED.")
				g.sln("")
				g.player.Str += 1
				g.sln("")
				g.moreNomail()
				g.sclrscr()
			}
		case 6:
			switch g.player.Clss {
			case 1:
				g.forestDeathKnightLevel()
			case 2:
				g.forestMysticalLevel()
			case 3:
				g.forestThiefLevel()
			}
		case 7:
			/* for the longest time, this event has always given 5 charm and
			 * taken 1 charm. Doing that kinda throws off the balance of the game,
			 * since a player could get very very pretty and never have a chance
			 * of getting ugly. Randomizing it helps that a bit. */
			tmp = g.random(2) + 1
			if g.random(3) == 1 {
				if !g.Settings.NewUglyStick {
					tmp = 5
				}
				g.sln("")
				g.sln("")
				g.lln("  `%MEGA EVENT IN THE FOREST")
				g.lln("`0-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
				// DIFF: LADY indent was one char... changed to two 'cause that's what everything uses.
				g.lln("  `2You are whacked with a pretty stick by an old man!")
				g.sln("")
				g.sln("  He giggles and runs away!")
				g.sln("")
				g.lln("  `%YOU GET " + itoa(tmp) + " CHARM!")
				g.sln("")
				g.player.Cha += tmp
			} else {
				if !g.Settings.NewUglyStick {
					tmp = 1
				}
				g.sln("")
				g.sln("")
				g.lln("  `%MEGA EVENT IN THE FOREST")
				g.lln("`0-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
				g.lln("  `2You are wacked with an `4ugly`2 stick by an old man!")
				g.sln("")
				g.lln("  `2He giggles and runs away!")
				g.sln("")
				g.lln("  `4YOU LOSE " + itoa(tmp) + " CHARM!")
				g.sln("")
				g.player.Cha -= tmp
			}
			g.moreNomail()
			g.sclrscr()
		case 8:
		inner:
			for {
				tmp = g.player.Level * 10000
				g.lln("`c  `%Event In The Forest`0")
				g.lln("`l")
				g.sln("")
				g.lln("  `2Walking through the forest, you stumble across a clearing in the woods.")
				g.sln("  Looking around, you see a small man and a lot of horses. The man walks up")
				g.lln("  to you and asks you, \"`0What can I do for you, mmhmm?`2\"")
				g.sln("")
				g.lln("  `2(`0B`2)uy a horse")
				g.lln("  `2(`0S`2)ell your horse")
				g.lln("  `2(`0G`2)o back to the forest")
				g.sln("")
				g.lw("  `0Your command `2[`0GBS`2] : ")
				ich = strings.ToUpper(g.getkey())
				if !strings.Contains("GBS", ich) {
					ich = "G"
				}
				g.sln(ich)
				switch ich {
				case "B":
					if g.player.Horse {
						g.sln("")
						g.lln("  `2\"`0Now, now. You already have a horse and I cannot sell you another.`2\"")
						g.sln("")
						g.moreNomail()
						break
					}
					g.sln("")
					g.lw("  `2\"`0Buy one of my horses for`% " + prettyInt(tmp) + "`0?`2\" ")
					ich = strings.ToUpper(g.getkey())
					if ich != "N" {
						ich = "Y"
					}
					g.sln(ich)
					if ich == "N" {
						g.lln("  `2\"`0Okay, mmhmm, maybe some other time.`2\"")
						g.moreNomail()
						break
					}
					if g.player.Gold < tmp {
						g.lln("  `2Chuckling, the hairyfoot says to you `2\"`0I am sorry, mmhmm, but I cannot")
						g.lln("  make a deal with someone that does not have enough money.`2\"")
						g.sln("")
						g.moreNomail()
						break
					}
					g.player.Gold -= tmp
					g.player.Horse = true
					g.sln("")
					g.lln("  `2\"`0Good deal! Mmhmm, there you are.`2\"")
					g.sln("")
					g.moreNomail()
				case "S":
					if !g.player.Horse {
						g.sln("")
						g.lln("  `2The happy look on the gnome's face goes sour. \"`0I am sorry, mmhmm,")
						g.lln("  but I cannot make a deal with someone that does not have a horse.`2\"")
						g.sln("")
						g.moreNomail()
						break
					}
					tmp = g.player.Level * 5000
					g.sln("")
					g.lw("  `2\"`0Sell your horse for`% " + itoa(tmp) + "`0?`2\" ")
					ich = strings.ToUpper(g.getkey())
					if ich != "N" {
						ich = "Y"
					}
					g.sln(ich)
					if ich == "N" {
						g.sln("")
						g.lln("  `2\"`0Okay.. Maybe some other time..`2\"")
						g.sln("")
						g.moreNomail()
						break
					}
					g.player.Gold += tmp
					g.player.Horse = false
					g.sln("")
					g.lln("  `2\"`0Thank you, kind person. Here is the `%" + itoa(tmp) + "`0 gold I promised.`2\"")
					g.sln("")
				case "G":
					g.sln("")
					g.lln("  `0\"Good day, kind person. May the Beast not play a deadly part in")
					g.lln("  your future, mmhmm.\"")
					g.sln("")
					g.moreNomail()
					break inner
				}
			}
			g.sclrscr()
		case 9:
			g.lrdfile("FAIRY")
			g.more()
			g.sclrscr()
			g.sln("")
			g.sln("")
			g.foreground(15)
			g.sln("  YOU ARE NOTICED!")
			g.sln("")
			g.foreground(2)
			g.sln("  The small things encircle you.  A small wet female bangs your")
			g.lln("  shin.  `0\"How dare you spy on us, human!\" `2you can't help but")
			g.sln("  smile, the defiance in her silvery voice is truly a sight, you")
			g.sln("  think to yourself.  Further contemplation is interrupted by")
			g.sln("  another sharpfully painful prod.")
			g.sln("")
			g.lln("  `2(`0A`2)sk for a blessing")
			if !g.player.HasFairy {
				g.lln("  `2(`0T`2)ry to catch one to show your friends")
			}
			g.sln("")
			g.lw("  `0Your choice? [`2A`0] : `%")
			ich = strings.ToUpper(g.getkey())
			if ich != "T" {
				ich = "A"
			}
			if g.player.HasFairy && ich == "T" {
				ich = "A"
			}
			g.sln(ich)
			if ich == "A" {
				g.sln("")
				g.lln("  `%\"Bless me!\" `2you implore the small figure.")
				g.sln("")
				g.lln("  `0\"Very well.\" `2she agrees, `0\"But we're still angry at you!")
				g.sln("")
				g.lw("  `0Your blessing is")
				g.mswait(500)
				g.sw(".")
				g.mswait(500)
				g.sw(".")
				g.mswait(500)
				g.sw(".")
				g.mswait(500)
				g.sw(".")
				// BUG (lord.js): tests player.has_horse, which doesn't
				// exist, so the horse blessing is offered even to
				// players who already own one.  Kept.
				switch g.random(3 + 1) {
				case 0:
					g.sw("a kiss from ")
					if g.player.Sex == "M" {
						g.sln("Teesha!")
					} else {
						g.sln("Nolmar!")
					}
					g.sln("")
					g.moreNomail()
					g.foreground(2)
					g.sln("  A fairy near her wordlessly upstretches its arms to you.")
					g.sln("")
					g.lln("  `%THE KISS IS STRANGELY FULFILLING! `2(You're refreshed)")
					g.player.HP = g.player.HPMax
					g.sln("")
					g.moreNomail()
				case 1:
					g.sln("an incredibly sad story!\"")
					g.moreNomail()
					g.foreground(2)
					g.sln("  The small body beckons you to lift her.  As you bring her")
					g.sln("  closer, she begins to whisper.")
					g.sln("")
					g.moreNomail()
					g.sln("  You almost immediately begin to cry.")
					g.sln("")
					g.lln("  `%YOUR TEARS TURNS INTO GEMS AND FALL INTO YOUR HANDS!")
					g.sln("")
					g.player.Gem += 2
					if g.player.Gem > 32000 {
						g.player.Gem = 32000
					}
					g.moreNomail()
				case 2:
					g.sln("a forest melody!\"")
					g.moreNomail()
					g.foreground(2)
					g.sln("  Immediately a small figure in the back raises her tiny reed")
					g.sln("  flute to her lips and begins to play.")
					g.sln("")
					g.moreNomail()
					g.sln("  The strange sounds send thousands of images to your mind.")
					g.sln("")
					g.foreground(15)
					g.sln("  YOU LEARN FAIRY LORE")
					tmp = g.player.Exp
					g.player.Exp += (10 * (g.player.Level * g.player.Level))
					if g.player.Exp > 2000000000 {
						g.player.Exp = 2000000000
					}
					tmp = g.player.Exp - tmp
					if tmp > 0 {
						g.lln(" - AND GET " + prettyInt(tmp) + " EXPERIENCE!")
					}
					g.player.FairyLore = true
					g.sln("")
					g.moreNomail()
					g.tournamentCheck()
				case 3:
					g.sln("a companion for your travels!\"")
					g.sln("")
					g.moreNomail()
					g.foreground(2)
					if g.random(2) == 1 {
						g.sln("  A shiny black stallion surfaces its head in the lake!")
					} else {
						g.sln("  A pure white mare nudges your back!")
					}
					g.sln("")
					g.foreground(15)
					g.sln("  YOU FEEL THE DAY GROW LONGER.")
					g.player.ForestFights = int64(float64(g.player.ForestFights) * 1.25)
					g.player.Horse = true
					g.sln("")
					g.moreNomail()
				}
			} else if !g.player.HasFairy {
				g.sln("")
				g.foreground(15)
				g.sln("  YOU MAKE A WILD GRAB FOR THE SMALL FIGURES!")
				g.sln("")
				g.moreNomail()
				g.foreground(2)
				g.sw("  Your hand finally connects with")
				g.mswait(500)
				g.sw(".")
				g.mswait(500)
				g.sw(".")
				g.mswait(500)
				g.sw(".")
				g.mswait(500)
				g.sw(".")
				if g.random(2) == 1 {
					g.lln("`4a thornberry bush`2!")
					g.sln("")
					g.player.HP = 1
					g.moreNomail()
					g.lln(" `4 (HIT POINTS GO DOWN....WAAAAY DOWN)")
					g.sln("")
					g.lln("  `2The fairy leader laughs at you!  You slink away in defeat.")
					g.sln("")
					g.moreNomail()
				} else {
					g.foreground(15)
					g.sln("  A FAIRY!")
					g.sln("")
					g.player.HasFairy = true
					g.foreground(2)
					g.sln("  You throw the screaming creature into your pouch.  What hidden")
					g.sln("  powers could it have?")
					g.sln("")
					g.sln("  You think now would be splendid time to leave.")
					g.sln("")
					g.moreNomail()
				}
			}
		case 10:
			if g.Settings.Olivia {
				g.lrdfile("CREEPY")
				ich = g.getkey()
				g.forestOlivia()
			}
		case 11:
			g.forestRescueThePrincess()
		case 12:
			g.forestFindLostGold()
		case 13:
			// Ive always wondered why a random event couldnt kill you..
			// Ive also wondered why there arent as many bad guys in the random
			// events.. Another funny thought. This must be a huge forest.. How
			// else could there be soo many monsters? :)

			// The player will lose 1 to 10 hitpoints... except thieves...
			// DIFF: This used to override and always do damage...
			if g.player.Clss == 3 {
				tmp = 0
			} else {
				tmp = g.random(10) + 1
			}
			g.sln("")
			g.sln("")
			g.lln("  `%MEGA EVENT IN THE FOREST")
			g.lln("`0-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
			g.lln("  `2An ugly troll jumps out from behind and tries to steal")
			g.sln("  your coin purse!")
			g.sln("")
			g.player.HP -= tmp
			if g.player.HP > 0 {
				if g.player.Clss == 3 {
					g.sln("  But being the sly theif that you are, you")
					g.sln("  fend him off easily!")
				} else {
					g.lln("  You fend him off, but you lose `%" + itoa(tmp) + "`2 hitpoints!")
					g.sln("")
				}
			} else {
				// If the player is smart, he wont "look for something to kill" while
				// he's low on hitpoints.. If he isnt keen on the idea of protecting
				// himself, well, the troll will get him..
				g.player.Gold = 0
				g.player.Gem = 0
				g.player.Dead = true
				// Lord will save the player's data when he gets back from the event,
				// but some unscrupulous people will drop carrier when they see that theyve
				// been killed. So save the player data now so they cant cheat the game.
				g.player.Put()
				g.sln("  You valiantly try to fend off the little troll, but at the")
				g.sln("  last moment, he pulls a dagger out of his boot and gives you")
				g.sln("  final blow!")
				g.sln("")
				g.lln("  `2The troll made off with you coin purse and gem sack!")
				g.sln("")
				g.lln("  `4You have been killed by the troll!")
				g.sln("")
			}
			g.moreNomail()
			// psst: one day, this lil bugger of a troll is going to try and make off
			// with your weapon.. your armor? nah. its hard enough to take off as it is.
			// but not hard to grab your spear of gold and take off running...
		case 14:
			// TODO: Nothing.. annoy the player. :D
		case 15:
			g.sln("")
			g.sln("")
			g.foreground(15)
			g.sln("  Event In The Forest")
			g.lln("`0-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
			g.foreground(2)
			g.sln("  In the gloom of the shady forest, you see smoke coming from a bright")
			g.sln("  chimney.")
			g.sln("")
			g.more()
			g.forestDarkhorseTavern(ch)
		}
		return
	}

	if g.player.Level == 1 {
		mnum = g.random(10)
	} else {
		if g.random(6) != 2 {
			mnum = ((g.player.Level - 1) * 11) + g.random(10)
		} else {
			mnum = (g.random(g.player.Level) * 11) + g.random(10)
		}
	}
	enemy := g.forestLoadMonster(mnum)
	g.player.ForestFights -= 1
	g.sln("")
	g.sln("")
	g.lln("  `2**`%FIGHT`2**")
	g.sln("")
	g.lln("  You have encountered " + enemy.Name + "`2!!")
	g.sln("")
	g.battle(enemy, false, false)
	if g.player.Dead {
		g.player.Put()
		g.say("normsay.lrd", enemy, 1, "  `5"+g.player.Name+" `2has been killed by `0"+enemy.Name+"`2!", "")
		g.forestSaveForestGold(g.player.Gold)
		g.player.Gold = 0
		g.player.Exp = g.player.Exp - g.player.Exp/10
		g.player.OnNow = false
		g.player.Put()
		g.sln("")
		g.deadScreen(enemy)
		g.sln("")
		g.exit(0)
	}
	if enemy.HP < 1 {
		g.foreground(15)
		g.sln("")
		g.lln("  You have killed " + enemy.Name + "`%!")
		g.sln("")
		g.foreground(2)
		if enemy.Gold > 0 {
			tmp = g.player.Gold
			g.player.Gold += enemy.Gold
			if g.player.Gold > 2000000000 {
				g.player.Gold = 2000000000
			}
			if tmp == g.player.Gold {
				g.sln("  You don't find any gold, but you do get")
			} else {
				g.lw("  `2You receive `0" + prettyInt(g.player.Gold-tmp) + " `2gold, and")
			}
		} else {
			g.sln("  You don't find any gold, but you do get")
		}
		g.lln("`0 " + prettyInt(enemy.Exp) + "`2 experience!")
		g.player.Exp += enemy.Exp
		if g.player.Exp > 2000000000 {
			g.player.Exp = 2000000000
			g.foreground(10)
			g.sln("Wow, you have a LOT of experience!")
		}
		g.sln("")
		g.more()
		g.tournamentCheck()
	}
}

// forestSpecial is forest()'s forest_special(): the WIERD EVENT.
func (g *Game) forestSpecial() {
	var gems int64

	g.lln("`c`%                        ** WIERD EVENT **")
	g.lln("`2-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
	g.lln("  `0You are heading into the forest, when you hear the voice of")
	g.sln("  angels singing.")
	g.lln("")
	g.moreNomail()
	g.sln("  You follow the sound for some time - when you are about to give")
	g.sln("  up...")
	g.lln("")
	g.moreNomail()
	g.lw("  `2You find.")
	g.mswait(40)
	g.sw(".")
	g.mswait(40)
	g.sw(".")
	gems = g.random(g.player.Level) + 1
	g.sln(pick(gems > 1, itoa(gems)+" Gems!", "1 Gem!"))
	g.sln("")
	g.player.Gem += gems
	g.player.Weird = false
	g.moreNomail()
}

// forestAttackDragon is forest()'s attack_dragon().
func (g *Game) forestAttackDragon() {
	var ch string
	dragonScreen := func() {
		g.lrdfile("LAIRANS")
	}

	if g.player.SeenDragon {
		g.sln("")
		g.sln("")
		g.foreground(2)
		g.sln("  You are shaking so badly from your previous encounter,")
		g.sln("  you deem it wise to wait and gather your strength!")
		g.sln("")
		g.more()
		return
	}
	dragonScreen()
	for {
		for {
			g.commandPrompt()
			ch = strings.ToUpper(g.getkey())
			g.sw(ch)
			if strings.Contains("?AR", ch) {
				break
			}
		}
		g.sln(ch)
		switch ch {
		case "?":
			dragonScreen()
		case "A":
			g.fightDragon(false)
		case "R":
			g.sln("")
			g.foreground(2)
			g.sln("  You decide it would be wise to depart from this wicked place.")
			g.sln("")
			g.more()
			return
		}
		if g.player.SeenDragon {
			break
		}
	}
}

// forestMenu is forest()'s menu().  The special argument is only false in
// a RIP branch, so it is always true here.
func (g *Game) forestMenu() {
	if g.player.Weird {
		g.forestSpecial()
	}
	if !g.player.Expert {
		g.lrdfile("FOREST")
		g.foreground(2)
		g.sln("")
		g.lln("  `2(`0L`2)ook for something to kill       (`0H`2)ealer's Hut")
		g.lw("  `2(`0R`2)eturn to town                 ")
		if g.player.Horse {
			g.lln("  `2(`0T`2)ake Horse To DarkCloak Tavern")
		} else {
			g.sln("")
		}
	}
}

// forestPrompt is forest()'s prompt().
func (g *Game) forestPrompt() {
	g.sln("")
	g.lw("  `2HitPoints: (`0" + prettyInt(g.player.HP) + "`2 of `0" + prettyInt(g.player.HPMax) + "`2)")
	g.lw("  Fights: `0" + prettyInt(g.player.ForestFights) + "`2 Gold: `0" + prettyInt(g.player.Gold))
	g.lw("  `2Gems: `0" + prettyInt(g.player.Gem))
	g.sln("")
	g.lw("  `5The Forest")
	g.foreground(2)
	g.sw("   (L,H,R,")
	if g.player.Horse {
		g.sw("T,")
	}
	g.sln("Q)  (? for menu)")
}
