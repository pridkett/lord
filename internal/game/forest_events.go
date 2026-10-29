package game

import (
	"os"
	"strings"
)

// forestFindLostGold is forest()'s find_lost_gold().
func (g *Game) forestFindLostGold() {
	var found int64
	left := g.player.Gold / 15

	if left < 100 {
		left = 100
	}
	g.getState(true)
	found = g.state.ForestGold
	g.state.ForestGold = left
	g.putState()
	if found < 100 {
		found = 100
	}
	if g.player.Gold+found > 2000000000 {
		found = 2000000000 - g.player.Gold
	}
	g.lln("`c  `%Event In The Forest`0")
	g.lln("`l")
	g.sln("")
	g.lln("  `2Fortune smiles, and you find `%" + prettyInt(found) + " `2gold!")
	g.player.Gold += found
	g.sln("")
	g.more()
	g.sclrscr()
}

// forestRescueThePrincess is forest()'s rescue_the_princess().
func (g *Game) forestRescueThePrincess() {
	var him string
	var lad string
	var ich string
	var castle int64

	g.sclrscr()
	g.sln("")
	g.sln("")
	g.foreground(15)
	g.sln("                         FOREST EVENT")
	g.foreground(2)
	g.sln("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
	g.foreground(10)
	if g.player.Horse {
		g.sln("  Your horse stumbles.")
	} else {
		g.sln("  You stumble.")
	}
	g.sln("")
	g.foreground(2)
	if g.player.DoneTower {
		g.sln("  Sure enough.  Another dead bird with a scroll.")
		g.sln("")
	} else {
		g.sln("  This would normally not be noted, since this happens all the")
		g.sln("  time - but what is interesting here is WHAT you stumbed over.")
		g.sln("")
		g.more()
		g.foreground(10)
		g.sln("  The object in question is a large dead bird.  Its stink is great.")
		g.sln("")
		g.foreground(2)
		g.sln("  This would normally not be noted either, since the forest is full of")
		g.sln("  dead things, and sometimes they rot before being eaten.")
		g.sln("")
		g.foreground(10)
		g.sln("  But the bird has a small scroll carefully tied to one leg.")
		g.sln("")
	}
	g.more()
	g.lln("  `r1  `%THE SCROLL READS:  `r0")
	g.sln("")
	g.foreground(7)
	switch g.random(3) {
	case 0:
		g.sln("  Dear whomever:")
		g.sln("")
		g.sln("  My ruthless uncle has trapped me in his castle.  He is angered")
		g.sln("  because I will not submit to him - in the way that he wants me")
		g.sln("  to.  ")
		g.sln("")
		g.sln("  Please come save me,")
		g.sln("        -Sleepless in a tower")
	case 1:
		g.sln("  Dear Brave Heart:")
		g.sln("")
		g.sln("  I am to wed one against my will.  My father tells me I am")
		g.sln("  selfish, because this political marriage will bring peace.")
		g.sln("")
		g.sln("  Get me out of here,")
		g.sln("        -a prisoner of war")
	case 2:
		g.sln("  To whom it may concern:")
		g.sln("")
		g.sln("  Geez am I bored.  I've been locked in the highest peak")
		g.sln("  of this castle for an ever so long time.  Crikey,  I would")
		g.sln("  explain why, but I'm running out of ink.  You see, I've been")
		g.sln("  up here for such a long time, it's drying up.  Well what do you")
		g.sln("  know!  Guess I have enough to explain after all!  Isn't that a")
		g.sln("  lucky break?  Ok, I'm up here because I wa..")
		g.sln("")
		g.sln("  The note isn't signed.")
	}
	g.sln("")
	g.foreground(15)
	g.sln("  You quickly swipe a tear from your eye as you put down the note.")
	g.sln("")
	g.more()
	if g.player.Sex == "M" {
		him = "her"
		lad = "girl"
	} else {
		him = "him"
		lad = "lad"
	}
	g.lln("  `2(`0S`2)ave " + him)
	g.lln("  `2(`0I`2)gnore the " + lad)
	g.sln("")
	g.lw("  `2Well? [`%S`2] `8: `%")
	ich = strings.ToUpper(g.getkey())
	if ich != "I" {
		ich = "S"
	}
	g.sln(ich)
	g.sln("")
	if ich == "I" {
		g.sln("  You look at the note a moment - then blithely toss it in a")
		g.sln("  nearby creek.  You skip away merrily! (evil can be fun!)")
		g.sln("")
		g.more()
		return
	}
	g.lln("  `0You're quite a hero.  Unfortunately, the " + lad + " seems to have")
	g.sln("  forgotten the return address.  You'll have to guess.")
	g.sln("")
	g.lln("  `2(`0C`2)astle Coldrake")
	g.lln("  `2(`0F`2)ortress Liddux")
	g.lln("  `2(`0G`2)annon Keep")
	g.lln("  `2(`0P`2)enyon Manor")
	g.lln("  `2(`0D`2)ema's Lair")
	g.sln("")
	g.lw("  Where do we go now? `8: `%")
	for {
		ich = strings.ToUpper(g.getkey())
		if strings.Contains("CFGPD", ich) {
			break
		}
	}
	castle = int64(strings.Index("CFGPD", ich) + 1)
	g.sln(ich)
	g.sln("")
	g.lln("  `2You set your jaw resolutely and journey to `%" + castles[castle] + "`2.")
	g.sln("")
	g.more()
	g.sclrscr()
	g.lrdfile("TOWER")
	_ = g.getkey()
	g.sclrscr()
	g.sln("")
	g.sln("")
	g.lln("                         `%THE RESCUE")
	g.lln("`2-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
	g.player.DoneTower = true
	g.lln("  `%THE DOOR SWINGS WIDE OPEN!")
	g.sln("")
	g.more()
	if castle == g.whichCastle {
		g.whichCastle = g.random(5) + 1
		g.lln("  `0\"You've come for me!\" `2shouts an overjoyed (and darn good looking)")
		g.lln("  " + lad + ".")
		g.sln("")
		g.lln("  `2You breathe a sigh of relief.  This was the right place.")
		g.sln("")
		if g.player.Sex == "M" {
			if !g.Settings.CleanMode {
				g.lln("  `2The girl eyes you dreamily.  `0\"I can never repay you, and I..\"")
				g.sln("")
				g.lln("  `%\"Oh but you can.  Is that your bed?\" `2you interrupt.")
			}
			g.lln("")
			g.more()
			g.lln("  `2" + prettyInt(g.player.Cha) + " minute" + pick(g.player.Cha > 1, "s", "") + " later, you feel quite repaid.")
		}
		g.sln("")
		g.lln("  `0YOU GET `%" + prettyInt((g.player.Level*g.player.Level)*20) + " `0EXPERIENCE.")
		g.player.Exp += g.player.Level * g.player.Level * 20
		if g.player.Exp > 2000000000 {
			g.player.Exp = 2000000000
		}
		g.sln("")
		g.more()
		if g.player.Sex == "M" {
			if !g.Settings.CleanMode {
				g.lln("  `%\"Well, that was fun, gotta go,\" `2you mutter as you throw your tunic")
				g.sln("  back on.")
				g.sln("")
				g.player.Laid += 1
				g.lln("  A sweet voice from the bed stops you as you hit the stairs. `0\"Wait!")
				g.sln("  I have something else for you too.\"")
				g.lln("")
				g.more()
			}
			g.lln("  `2Your blank face turns to joy as she hands you a pouch.")
		} else {
			g.lln("`2  Your blank face turns to joy as he hands you a pouch.")
		}

		g.sln("")
		g.lln("  `0YOU GET `%" + prettyInt(3*g.player.Level) + " `0GEMS FOR YOUR TROUBLE.")
		g.player.Gem += 3 * g.player.Level
		if g.player.Gem > 2000000000 {
			g.player.Gem = 2000000000
		}
		g.sln("")
		if g.player.Sex == "M" {
			g.logLine("  `0" + g.player.Name + " `2saved a princess today!")
		} else {
			g.logLine("  `0" + g.player.Name + " `2saved a prince today!")
		}
		g.tournamentCheck()
	} else {
		switch g.random(2) {
		case 0:
			g.lln("  `2The room is empty, save a giant chest in the middle.")
			// DIFF: This didn't decrease your HP if it was less than 4...
			g.player.HP = g.player.HP / 4
			if g.player.HP < 1 {
				g.player.HP = 1
			}
			g.sln("")
			g.lln("  `%\"Hello?  Anybody home?\" `2you call softly.")
			g.sln("")
			g.sln("  You jump in surprise when you hear a voice answer - from the chest.")
			g.sln("")
			g.lln("  `0\"Help me!  I'm in here!  I wrote the note!\"`2 the voice pleads.")
			g.sln("")
			g.more()
			g.lln("  `%\"Um, how did you write the note from in there?\" `2you wonder out loud.")
			g.sln("")
			g.lln("  `0\"Nevermind that!  Just help me!\"")
			g.sln("")
			g.more()
			g.lln("  `%\"Fine.\" `2you open the chest.  Surprisingly, there is not a fair")
			if g.player.Sex == "M" {
				g.sln("  maiden inside.  Instead, you find a..")
			} else {
				g.sln("  prince inside.  Instead, you find a..")
			}
			g.sln("")
			g.more()
			if g.player.Sex == "M" {
				g.lln("  `%\"`bMY GOD A HOLL!  `%A HALF HUMAN HALF TROLL WOMAN!  NOOO!\" `2you ")
			} else {
				g.lln("  `%\"`bMY GOD A HOLL!  `%A HALF HUMAN HALF TROLL MAN!  NOOO!\" `2you ")
			}
			g.sln("  scream. ")
			g.sln("")
			g.more()
			g.lln("`c                         `%THE UGLY PART")
			g.foreground(2)
			g.sln("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
			if !g.Settings.CleanMode {
				if g.player.Sex == "M" {
					g.sln("  You are violated by this hairy beast.  Apparently she needs love")
				} else {
					g.sln("  You are violated by this hairy beast.  Apparently he needs love")
				}
				g.sln("  like anyone else.")
			} else {
				g.sln("  The thing beats you up.  Yeah.  Real bad.")
			}
			g.sln("")
			g.foreground(4)
			if g.player.Sex == "M" {
				g.sln("  YOU WEAKLY CRAWL AWAY FROM HER SOMETIME LATER.")
			} else {
				g.sln("  YOU WEAKLY CRAWL AWAY FROM HIM SOMETIME LATER.")
			}
			g.sln("")
			g.logLine("  `0" + g.player.Name + " `2had a nasty run in with a Holl!")
			g.more()
			return
		case 1:
			g.foreground(10)
			if g.player.Sex == "M" {
				g.sln("  You see two beautiful women playing chess.")
				g.sln("")
				g.lln("  `%\"Hello, ladys.  Which one of you needs rescuing?\" `2you ask politely.")
			} else {
				g.sln("  You see two handsome men playing chess.")
				g.sln("")
				g.lln("  `%\"Hello, boys.  Which one of you needs rescuing?\" `2you ask politely.")
			}
			g.sln("")
			g.more()
			g.lln("  `0\"Neither!\" `2they chime.")
			g.sln("")
			if g.player.Sex == "M" {
				g.lln("  `2You scratch your chin in confusion.  `%\"So where is the damn damsel")
			} else {
				g.lln("  `2You scratch your chin in confusion.  `%\"So where is the stinkin' prince")
			}
			g.sln("  in distress?!\"")
			g.sln("")
			g.more()
			g.lln("  `0At this moment, a messager burst through the broken doorway.")
			g.sln("")
			g.lln("  `5\"My friends, I bear terrible news - `%" + castles[castle] + "`5 has been")
			g.lln("  attacked.  Your father the King is dead.\" `2the messenger is now audibly")
			g.sln("  sobbing.")
			g.sln("")
			g.more()
			g.lln("  `%\"Ah.  Yes.  Well, this is all very tragic, but I uh, need to be going.\" `2you")
			g.sln("  studder uncomfortably.")
			g.sln("")
			g.foreground(10)
			if g.player.Sex == "M" {
				g.sln("  The now ashen white faced women look at you dumbfounded as you make")
			} else {
				g.sln("  The now ashen white faced men look at you dumbfounded as you make")
			}
			g.sln("  your exit.")
			g.sln("")
			if g.player.Sex == "M" {
				g.logLine("  `0" + g.player.Name + " `2showed courage today by trying to save a princess.")
			} else {
				g.logLine("  `0" + g.player.Name + " `2showed courage today by trying to save a prince.")
			}
		}
	}
	g.more()
}

// forestFlowerGarden is forest()'s flower_garden().
func (g *Game) forestFlowerGarden() {
	var ich string
	var message string
	var pname string
	fname := g.gamedir("garden.lrd")

	g.sln("")
	g.sln("")
	g.foreground(15)
	g.sln("  Event In The Forest")
	g.lln("`0-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
	g.foreground(2)
	g.sln("  You find a beautiful Garden, and decide to take a little rest...")
	g.sln("  You drink from the brook, and smell the flowers.")
	g.sln("")
	g.foreground(15)
	g.sln("  YOU ARE REFRESHED, AND GET ONE MORE FOREST FIGHT FOR TODAY!")
	g.sln("")
	g.player.ForestFights += 1
	if g.player.ForestFights > 32000 {
		g.player.ForestFights = 32000
	}
	if g.player.HP < g.player.HPMax {
		g.player.HP = g.player.HPMax
	}
	g.more()
	g.lln("  `0You notice the flowers seem to be arranged.")
	g.lw("  Study them? [`%Y`2] : `%")
	ich = strings.ToUpper(g.getkey())
	// DIFF: This wasn't done originally...
	if ich != "N" {
		ich = "Y"
	}
	g.sln(ich)
	if ich == "Y" {
		g.sclrscr()
		g.sln("")
		g.sln("")
		g.lln("  `%Understanding Pedals & Dew`2")
		g.sln("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
		g.fmutex(fname)
		if !fileExists(fname) {
			// game_or_exec('gstart.lrd'): the sysop's copy, else
			// the one built into the binary.
			if b, err := g.assetBytes("gstart.lrd"); err == nil {
				_ = os.WriteFile(fname, b, 0o644)
			}
		}
		if !fileExists(fname) {
			_ = os.WriteFile(fname, []byte("  `0Fairy Tisha         `%-\"`0Oooh!  I love flowers!  And I love to kiss.`%\"\n"), 0o644)
		}
		// DIFF: No more prompt since we're holding a lock.
		g.displayFile(fname, true, false)
		g.rfmutex(fname)
		g.sln("")
		g.lw("  `0Arrange them yourself? [`%N`2] : `%")
		ich = strings.ToUpper(g.getkey())
		if ich != "Y" {
			ich = "N"
		}
		g.sln(ich)
		if ich == "Y" {
			g.lln("  `2What would you like them to symbolize?")
			g.sw("  ")
			message = g.getstr(0, 0, 53, 1, 15, "", GetStrOpts{})
			if g.Settings.FunkyFlowers {
				message = strings.ReplaceAll(message, "``", "")
			} else {
				message = cleanStr(message)
			}
			g.sln("")
			g.sln("")
			// DIFF: This was after the length check...
			message = strings.TrimRight(message, " ")
			if dispLen(message) < 3 {
				g.sln("  You decide to arrange later.")
				g.sln("")
				g.more()
			} else {
				pname = spacePad(g.player.Name, 20)
				g.fmutex(fname)
				lines, ok := readLines(fname)
				if !ok {
					g.fatalf("Unable to open %s", fname)
				}
				lines = append(lines, "  `0"+pname+"`%- \"`0"+message+"`r0`%\"")
				for len(lines) > 20 {
					lines = lines[1:]
				}
				if err := writeLines(fname, lines...); err != nil {
					g.fatalf("Unable to open %s", fname)
				}
				g.rfmutex(fname)
				g.sln("")
				g.sln("")
			}
		}
	}
}
