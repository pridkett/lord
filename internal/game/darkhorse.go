package game

import (
	"fmt"
	"slices"
	"strings"
)

// forestDarkhorseTavern is forest()'s darkhorse_tavern().  ch is forest()'s
// `ch`, which chance() shares (see forestDarkhorseChance).
func (g *Game) forestDarkhorseTavern(ch *string) {
	var dch string
	gamesPlayed := 0

	menu := func() {
		g.lrdfile("CLOAK")
	}

	menu()
	for {
		g.sln("")
		g.lln("  `0DarkCloak Tavern `2(C,E,T,V,D,G,W,R) (? for menu)")
		g.sln("")
		g.lw("  `2Your command? : ")
		dch = strings.ToUpper(g.getkey())
		if strings.Contains("\rQ", dch) {
			dch = "R"
		}
		g.sln(dch)
		switch dch {
		case "D":
			g.showLog()
		case "Y":
			g.showStats()
		case "G":
			g.forestDarkhorseGamble(&gamesPlayed)
		case "T":
			g.forestDarkhorseChance(ch)
		case "W":
			g.forestDarkhorseOldMan()
		case "C":
			g.converse(true)
		case "?":
			menu()
		case "E":
			g.forestDarkhorseRankLays(g.gamedir("temp" + itoa(g.player.Record)))
			g.displayFile(g.gamedir("temp"+itoa(g.player.Record)), true, true)
			g.sln("")
			fileRemove(g.gamedir("temp" + itoa(g.player.Record)))
			g.moreNomail()
		}
		if dch == "R" {
			break
		}
	}

	if g.random(25) == 17 {
		g.blackjack()

	}
}

// forestDarkhorseGamble is darkhorse_tavern()'s gamble().
func (g *Game) forestDarkhorseGamble(gamesPlayed *int) {
	var bet int64
	var mynum int64
	var hisnum int64

	wager := func() int64 {
		var ret int64 = -1

		for {
			g.lln("  `2How much gold of your `%" + prettyInt(g.player.Gold) + "`2 will you hazard? (`00`2 to chicken out)")
			g.lw("  `2WAGER : `0")
			ret, _ = jsParseInt(g.consoleGetstr(GetStrOpts{Len: 11, Integer: true}))
			g.sln("")
			if ret > g.player.Gold {
				g.lln("  `2Betting what you don't have is `4NOT`2 a good idea.")
				g.sln("")
				ret = -1
			} else if ret < 0 {
				g.lln("  `2You don't think that will go over too big.")
				g.sln("")
			}
			if ret >= 0 {
				break
			}
		}
		g.sln("")
		return ret
	}

	g.lln("`c`%  ** GAMBLE TIME! **")
	g.sln("")
	g.foreground(2)
	g.sln("  You saunter over to the bar and demand that someone gamble with you.")
	g.sln("")
	if *gamesPlayed > 1 {
		g.sln("  No one seems too thrilled at the prospect.  Perhaps if you came back")
		g.sln("  another time.")
		g.sln("")
		if *gamesPlayed > 2 {
			g.sln("  (You wonder if it had anything to do with your making a joke out")
			g.sln("  of the word honor)")
			g.sln("")
		}
		return
	}

	*gamesPlayed += 1
	switch g.random(3) {
	case 0:
		g.sln("  The old man stops etching on the table he is at and walks over to you.")
		g.sln("")
		g.lln("  `0\"I'll play a game with ya, kid!  I'll bet I can guess what number")
		g.sln("  you are thinkin'!\"")
		g.sln("")
		g.more()
		bet = wager()
		if bet == 0 {
			g.lln("  `2The old man laughs in your face then continues his carving in the")
			g.sln("  table.")
			break
		}
		g.lln("  `0\"Fine!  `%" + prettyInt(bet) + "`0 it is!  Concentrate on a number.\"`2")
		g.sln("")
		mynum = g.random(100) + 1
		g.lln("  `2You concentrate on the number.")
		g.mswait(500)
		g.sw(".")
		g.mswait(500)
		g.sw(".")
		g.mswait(500)
		g.lln("`%" + itoa(mynum))
		g.sln("")
		g.lln("  `2The old man studies you quietly for a moment.  Then screams in delight.")
		g.sln("")
		g.moreNomail()
		// DIFF: This was clearly broken, and he always guess right.
		// it was hisnum !== mynum in the loop test.
		hisnum = g.random(100) + 1
		if hisnum > 55 {
			hisnum = g.random(100) + 1
			for hisnum == mynum {
				hisnum = g.random(100) + 1
			}
		} else {
			hisnum = mynum
		}
		g.lw("  `0\"The number is")
		g.mswait(500)
		g.sw(".")
		g.mswait(500)
		g.sw(".")
		g.mswait(500)
		g.lw(".`%")
		g.lln(itoa(hisnum) + "`0 isn't it?!!!!!!!!!!!")
		g.sln("")
		if hisnum != mynum {
			g.lln("  `2You sadly inform the Old Man of his mistake, and he grudgingly gives")
			g.lln("  you `%" + prettyInt(bet) + "`2 gold from his pouch.")
			g.player.Gold += bet
			if g.player.Gold > 2000000000 {
				g.player.Gold = 2000000000
			}
			break
		}
		g.lln("  `2You feel obliged to admit that he chose correctly.  You count out")
		g.lln("  `%" + prettyInt(bet) + "`2 gold and give it to him with a scowl.")
		g.player.Gold -= bet
		g.sln("")
		g.lln("  `2The old man dances a jig of joy!")
	case 1:
		g.sln("  The old man stops etching on the table he is at and walks over to you.")
		g.sln("")
		g.lln("  `0\"I'll play a game with ya, kid!\"")
		g.sln("")
		g.lln("  `2The old man grabs two wooden mugs from a table and slaps them down")
		g.lln("  in front of you upside down.  `0\"Guess which one I hid muh teeth in!\"")
		g.sln("")
		g.more()
		bet = wager()
		if bet == 0 {
			g.lln("  `2The old man laughs in your face then continues his carving in the")
			g.sln("  table.")
			break
		}
		g.lln("  `0\"Agreed!\" `2The old man waits for your response.")
		g.sln("")
		g.more()
		switch g.random(2) {
		case 0:
			g.lln("  `2You demand to see what's in the first mug!")
		case 1:
			g.lln("  `2You demand to see what's in the second mug!")
		}
		g.sln("")
		hisnum = g.random(100) + 1
		g.lw("  `2The old man slowly turns over the mug.")
		g.mswait(500)
		g.sw(".")
		g.mswait(500)
		g.sw(".")
		g.mswait(500)
		g.sw(".")
		if hisnum > 45 {
			g.lln("`%IT HAS HIS WOODEN TEETH IN IT!")
			g.player.Gold += bet
			if g.player.Gold > 2000000000 {
				g.player.Gold = 2000000000
			}
			g.sln("")
			g.lln("  `2The entire bar cheers at your success!")
			g.sln("")
			g.sln("  The old man groans and hands you the gold you've won.")
		} else {
			g.lln("`4IT IS EMPTY SAVE SOME STALE BEER!")
			g.player.Gold -= bet
			g.sln("")
			g.lln("  `2The old man howls in delight as you pay him.")
		}
	case 2:
		g.sln("  The old man stops etching on the table he is at and walks over to you.")
		g.sln("")
		g.lln("  `0\"I'll play a game with ya, kid!\"")
		g.sln("")
		g.lln("  `2The old man walks over to you and hands you a small dagger.  Then he")
		g.sln("  moves to the other side of the Tavern, and picks up a tankard of brew.")
		g.sln("")
		g.lln("  `0\"I'll bet you can't knock this off my head without getting me wet!\"")
		g.sln("")
		g.more()
		bet = wager()
		if bet == 0 {
			g.lln("  `2The old man laughs in your face then continues his carving in the")
			g.sln("  table.")
			break
		}
		g.lln("  `2The old man positions himself carefully, and places the Mug on his head.")
		g.sln("")
		g.lln("  `0\"You'll never hit it, sharpshooter!  Throw it already!\"")
		g.sln("")
		g.more()
		g.lln("  `2You give it your best shot.")
		g.sln("")
		g.lw("  `%** `0WHO")
		for mynum = 0; mynum < 15; mynum += 1 {
			g.mswait(100)
			g.sw("O")
		}
		g.lln("SH `%**")
		g.sln("")
		if (g.random(100) + 1) > 44 {
			g.lln("  `2YOU HAVE KNOCKED IT OFF LEAVING THE OLD MAN HIGH AND DRY!")
			g.sln("")
			g.lln("  The old man swears sourly, but pays you the `%" + prettyInt(bet) + "`2 gold.")
			g.sln("")
			g.player.Gold += bet
			if g.player.Gold > 2000000000 {
				g.player.Gold = 2000000000
			}
			break
		}
		g.foreground(4)
		g.sln("  YOU MISS YOUR TARGET!")
		g.foreground(2)
		g.sln("")
		g.mswait(500)
		switch g.random(6) {
		case 0:
			g.lln("  `2The dagger ricochetted off the wall and into `%Chance`2's drink!")
			g.sln("")
			g.sln("  He looks furious!")
		case 1:
			g.sln("  The dagger smoothly implants itself into the old mans forehead!")
			g.sln("")
			g.sln("  But he is ok!")
		case 2:
			g.sln("  The dagger stabs deeply into the oak behind the old mans head!")
		case 3:
			g.sln("  The dagger stabs near your foot!  The entire tavern laughs at you!")
		case 4:
			g.sln("  The dagger hits the bar, and slides to a screeching stop in front of a")
			g.sln("  young warrior.  His face is white as a sheet - The entire bar laughs!")
		case 5:
			g.sln("  The dagger flies through an open window!")
		}
		g.sln("")
		g.sln("  The old man laughs with glee.  You mumble curses as you pay him.")
		g.player.Gold -= bet
	}
	g.sln("")
}

// forestDarkhorseChance is darkhorse_tavern()'s chance().
//
// BUG (lord.js): the loop tests forest()'s `ch`, not `ich`, and 'R' sets
// that outer variable to 'Q'.  A second visit to Chance during the same
// tavern stay therefore only takes a single command.  Kept.
func (g *Game) forestDarkhorseChance(ch *string) {
	var ich string
	var op *Player
	var pl int64
	var him string
	var he string
	var prof string
	classStrs := []string{"Nobody", "`0Warrior", "`#Mystical Skills User`0", "`9Thief`0"}

	cmenu := func() {
		g.lrdfile("CHANCE")
	}

	cmenu()

	for {
		g.lw(" `2 Your command? (`0? `2for menu)  [`0R`2] : ")
		ich = strings.ToUpper(g.getkey())
		if ich == "\r" {
			ich = "R"
		}
		g.sln(ich)
		g.sln("")
		if ich == "L" {
			g.lln("  `0\"I know many things about many people.  Who is your enemy?\"`2")
			pl = g.findPlayer()
			if pl == -1 {
				g.lln("")
				g.lln("  `0\"I don't know anyone with a name even close to that.\"`2")
				g.sln("")
				cmenu()
			} else if pl == int64(g.player.Record) {
				if g.player.Sex == "M" {
					him = "him"
					he = "he"
				} else {
					him = "her"
					he = "she"
				}
				g.lln("  `0\"Yes..I know " + him + ".  " + he + " is a favorite customer of mine!\"")
				g.lln("  `2Chance laughs heartily.")
				g.lln("")
			} else {
				op = g.playerGet(pl)
				prof = classStrs[op.Clss]
				g.lln("  `2Chance's face turns somber.")
				g.lln("  `0\"" + op.Name + "`0 the " + prof + "?  I know who that is.\"")
				g.sln("")
				g.lln("  `0\"This information was not easily come by, and I am going to have to charge")
				g.lln("  two `%Gems`0 for it.  Now you know why this tavern is REALLY here.\"")
				g.lln("")
				if g.player.Gem < 2 {
					g.lln("`2  Not having two `%Gems`2, you decline.")
					g.sln("")
				} else {
					g.lw("  `2Pay Chance two `%Gems`2 for the info? [`0Pay 'Em`2] : ")
					ich = strings.ToUpper(g.getkey())
					// DIFF: Previously, 'N' was forced to 'Y', and anything except N or Y would be an effective "no"
					if ich != "N" {
						ich = "Y"
					}
					g.sln(ich)
					g.sln("")
					if ich != "Y" {
						g.lln("  `0\"No problem!  I know how it is these days.\"")
						g.sln("")
					} else {
						g.player.Gem -= 2
						g.lln("  `0\"Alright.  Come with me.\"")
						g.sln("")
						g.lln("  `2Chance leads you to a small room in back of the tavern, and you take")
						g.sln("  a comfortable seat.")
						g.sln("")
						g.lln("  `0\"Well...Here is everthing I know about " + op.Name + "`0.\"")
						g.sln("")
						g.more()
						if op.Sex == "M" {
							him = "him"
							he = "he"
						} else {
							him = "her"
							he = "she"
						}
						g.lln("  `0\"Fights with a " + op.Weapon + "`0 and has a total Strength of `%" + prettyInt(op.Str) + "`0.\"")
						g.lln("  `0\"Wears a " + op.Arm + "`0 and has a total Defense of `%" + prettyInt(op.Def) + "`0.\"")
						g.sln("")
						if op.Cha < 3 {
							g.lln("  `0\"" + op.Name + " is very ugly.\"")
						} else if op.Cha < 5 {
							g.lln("  `0\"" + op.Name + " is kind of blah looking.\"")
						} else if op.Cha < 10 {
							g.lln("  `0\"" + op.Name + " is fairly good looking.\"")
						} else if op.Cha < 50 {
							g.lln("  `0\"" + op.Name + " has a very fair countenance.\"")
						} else if op.Cha < 90 {
							// DIFF: This used the *players* sex!
							if op.Sex == "F" {
								g.lln("  `0\"" + op.Name + " is a very good looking woman.\"")
							} else {
								g.lln("  `0\"" + op.Name + " gets all the women...The lucky brute!\"")
							}
						} else {
							// DIFF: This used the *players* sex!
							if op.Sex == "F" {
								g.lln("  `0\"I have heard " + op.Name + " has the face and body of a Goddess.\"")
							} else {
								g.lln("  `0\"" + op.Name + " is a good looking bastard.\"")
							}
						}
						g.sln("")
						g.lln("  `0\"Total worth in gold is " + prettyInt(op.Gold+op.Bank) + ".\"")
						g.sln("")
						g.lln("  `0\"Last time we checked, " + he + " had " + prettyInt(op.Gem) + " `%Gems`0.\"")
						g.sln("")
						if op.Kids == 0 {
							g.lln("  `0\"" + op.Name + " has no offspring.\"")
						} else {
							g.lln("  `0\"" + op.Name + " has `%" + prettyInt(op.Kids) + " `2offspring.")
						}
						if op.Horse {
							g.lln("  `0That person owns a horse.")
							g.sln("")
						}
						// DIFF: Didn't check Seth and Violet marriages...
						g.getState(false)
						if op.MarriedTo > -1 || g.state.MarriedToSeth == int64(op.Record) || g.state.MarriedToViolet == int64(op.Record) {
							g.lln("  `0That person is married.")
							g.sln("")
						}
						g.sln("")
						g.moreNomail()
					}
				}
			}
		} else if ich == "?" {
			cmenu()
		} else if ich == "C" || ich == "S" {
			// DIFF: The 'S' key was not in here, but it was in the menu.
			g.lln("  `0\"Tired of what you do?  I know how it is.  To figure out what you")
			g.lln("  REALLY want to do in your life, think about your childhood.\"`2 ")
			g.sln("")
			g.lln("  `0\"Remember.  You NEVER forget what you learn in ANY profession.\"`2")
			g.player.Clss = g.chooseProfession(true)
			g.sln("")
			cmenu()
		} else if ich == "T" {
			g.sln("")
			g.lln("  `0\"C`3o`4l`5o`6r`7s`8?`0\"`2, Chance laughs, `0\"They are easy.\"")
			g.sln("")
			g.sln("  `1 `2 `3 `4 `5 `6 `7 `8 `9 `0 `! `@ `# `$ `%")
			g.lln("  `1^^ `2^^ `3^^ `4^^ `5^^ `6^^ `7^^ `8^^ `9^^ `0^^ `!^^ `@^^ `#^^ `$^^ `%^^")
			g.sln("")
			g.sln("  `5Colors are `%FUN`5! would look like...")
			g.sln("")
			g.lln("  `5Colors are `%FUN`5!")
			g.sln("")
			g.lln("  `0\"Using that symbol and a number, you can make any color.  They don't")
			g.sln("  work in mail, but try 'em talking to people or write it in the dirt.\"")
			g.lln("`2")
		} else if ich == "P" {
			g.lln("  `0\"Ok, enter your practice line.\"")
			g.lw("  `2Text `0: `%")
			line := cleanStr(g.consoleGetstr(GetStrOpts{Len: 69}))
			g.sln("")
			g.lln("  `0\"Here is what that would look like.\"")
			g.lln("`2  `%" + line)
			g.sln("")
		} else {
			if ich == "R" {
				*ch = "Q"
			}
		}
		if *ch == "Q" {
			break
		}
	}
}

// forestDarkhorseOldMan is darkhorse_tavern()'s old_man().
func (g *Game) forestDarkhorseOldMan() {
	var ich string
	var pl int64
	var him string
	var op *Player
	var l string

	g.lrdfile("OLDMAN")
	for {
		g.sln("")
		g.lw(" `2 Your command? (`0? `2for menu)  [`0R`2] : ")
		ich = strings.ToUpper(g.getkey())
		if !strings.Contains("V?E", ich) {
			ich = "R"
		}
		g.sln(ich)
		g.sln("")
		switch ich {
		case "?":
			g.lrdfile("OLDMAN")
		case "V":
			g.lln("  `0\"Who would you like to know more about?\"`2")
			pl = g.findPlayer()
			if pl == -1 {
				g.sln("")
				g.lln("  `0\"I don't know anyone with a name even close to that.\"`2")
				g.sln("")
			} else if pl == int64(g.player.Record) {
				g.lln("  `0\"Why, Id hope you know what you've said about yourself..\"`2,")
				g.sln("  the old man cackles.")
				g.lln("")
			} else {
				op = g.playerGet(pl)
				him = pick(op.Sex == "M", "him", "her")
				g.lln("  `2The Old Man thinks for a minute..")
				if !op.HasDes {
					g.lln("  `0\"" + op.Name + "?  I haven't heard anything about " + him + ".\"")
				} else {
					g.lln("  `0\"" + op.Name + "?  I know who that is. Last I heard..\"")
					g.sln("")
					g.lln("    `2" + op.Des1)
					g.lln("    `2" + op.Des2)
				}
			}
		case "E":
			g.lln("  `0\"What do you want me to remember about you?\"")
			g.lw(" `2-> `%")
			// TODO: Is a different string input than everything else...
			l = cleanStr(g.consoleGetstr(GetStrOpts{Len: 70}))
			if l != "" {
				g.player.Des1 = l
				// BUG (lord.js): sets player.des, which isn't a
				// record field, instead of has_des, so the Old
				// Man never repeats what you told him.  Kept.
				g.lw(" `2-> `%")
				l = cleanStr(g.consoleGetstr(GetStrOpts{Len: 70}))
				g.player.Des2 = l
				g.sln("")
				g.lln("  `2\"`0It has been noted..`2\"")
				g.player.Put()
			}
		}
		if ich == "R" {
			break
		}
	}
}

// forestDarkhorseRankLays is darkhorse_tavern()'s rank_lays().
func (g *Game) forestDarkhorseRankLays(fname string) {
	var pl []*Player

	for _, op := range g.allPlayers() {
		if op.Name != "X" && op.Laid > 0 {
			pl = append(pl, op)
		}
	}
	slices.SortStableFunc(pl, func(a, b *Player) int {
		if a.Laid != b.Laid {
			return int(b.Laid - a.Laid)
		}
		return int(b.Cha - a.Cha)
	})
	lines := []string{
		"",
		"",
		"                          `%The Old Man's Ranking",
		"",
	}
	if g.Settings.CleanMode {
		lines = append(lines, "  `0Name                             Evil Deeds             Player Kills")
	} else {
		lines = append(lines, "  `0Name                                Lays                Player Kills")
	}
	lines = append(lines, "`#-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
	for i := 0; i < len(pl); i++ {
		lines = append(lines, "  `0"+spacePad(pl[i].Name, 22)+"            `%"+fmt.Sprintf("%5d", pl[i].Laid)+"                      `4 "+fmt.Sprintf("%6.6s", prettyInt(pl[i].PvP)))
	}
	if len(pl) == 0 {
		lines = append(lines, "  `0Sad times indeed, no one has managed to make this list.")
	}
	if err := writeLines(fname, lines...); err != nil {
		g.fatalf("Unable to open %s", fname)
	}
}
