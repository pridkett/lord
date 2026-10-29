package game

import (
	"strings"
)

// forestDeathKnightLevel is forest()'s death_knight_level().
func (g *Game) forestDeathKnightLevel() {
	var title string
	var ich string
	var guilty bool

	g.forestEventHeader()
	g.sln("  While trekking through the forest, you come to the hidden castle of")
	g.sln("  The Black Knights.  You are immediately greeted by a score of men in shiny")
	g.sln("  black armour.")
	g.sln("")
	title = pick(g.player.Sex == "F", "Lady", "Lord")

	if g.player.SkillW > 39 {
		g.lln("  `0\"Well met " + title + " `%" + g.player.Name + "`0!  A fellow Black Knight is")
		g.lln("  Always welcome here.\"`2  You walk the grounds, and eat with your")
		g.sln("  comrades.  You are fully refreshed.")
		g.lln("")
		if g.player.HP < g.player.HPMax {
			g.player.HP = g.player.HPMax
		}
		g.player.LevelW += 1
		g.lln("  `%HIT POINTS FILLED AND YOU RECEIVE THE ENERGY FOR 1 DEATH KNIGHT ATTACK!")
		g.sln("")
		g.more()
		return
	}
	if g.player.SkillW < 20 {
		g.lln("  `0\"Well met " + title + " `%" + g.player.Name + "`0!  We have been expecting you.  We")
		g.sln("   know you aspire to join us, and become a Death Knight.  We will")
		g.sln("   teach you a lesson today, but only if you pass a test of our device.\"")
	} else {
		g.lln("  `0\"Greetings " + title + " `%" + g.player.Name + "`0!  What are you doing around here?!")
		g.lln("   Even tho you are already a member, now is a great time to practice")
		g.lln("   your skills.  To become an even better warrior.  You know the routine...\"")
	}

	g.sln("")
	g.foreground(2)
	g.sln("  THEY LEAD YOU TO THE DEATH KNIGHT DUNGEON.")
	g.sln("")
	g.moreNomail()
	g.foreground(15)
	g.sclrscr()
	g.sln("")
	g.sln("")
	g.sln("                            ** THE TEST **")
	g.sln("")
	g.lln("  `2You ignore the shrieks of pain, the suffering of this dungeon's occupants.")
	g.sln("  You are shown a man kneeling over a stained chopping block.")
	g.sln("")
	g.lln("  `0\"This man is accused of a crime.  Is he innocent or guilty?\"")
	g.sln("")
	g.lln("  `2(`01`2) Decapitate Him")
	g.lln("  (`02`2) Release Him")
	g.sln("")
	g.commandPrompt()
	ich = strings.ToUpper(g.getkey())
	if ich != "2" {
		ich = "1"
	}
	g.sln(ich)
	g.sln("")
	guilty = g.random(2) == 1
	g.sln("")
	if ich == "1" {
		g.lln("  `2You take the axe and bring it down as hard as you can.  After")
		g.sln("  a sickening (but satisfying) crunch the deed is done.")
		g.sln("")
	} else {
		g.lln("  `%\"That man is innocent!  You shall not harm a hair on his")
		g.lln("  head, as long as I have a breath in me to fight!\" `2 you ")
		g.sln("  shout dramatically.")
	}
	g.lw("  `0\"You have chosen.")
	g.mswait(500)
	g.sw(".")
	g.mswait(500)
	g.sw(".")
	g.mswait(500)
	if guilty {
		if ich == "1" {
			g.lln(" `%WISELY!`0\"")
			g.sln("  \"You have done this country justice today.\"")
			g.foreground(2)
			g.sln("")
			g.more()
			g.lln("")
			g.raiseClass()
		} else {
			g.lln(" `4POORLY`0.\"")
			g.sln("  \"That man raped 6 women.  And you defend him?  Good God man!")
			g.sln("")
			g.more()
			g.sln("")
			return
		}
	} else {
		// BUG (lord.js): compares the key with the number 1, which never
		// matches, so beheading an innocent man is judged "WISELY" too.
		// Kept.
		if false {
			g.lln(" `4POORLY`0.\"")
			g.sln("  \"This man did no crime.  He was the father of 6 children.  \"Was\"")
			g.sln("   being the key word here.  Perhaps another time.\"")
			g.foreground(2)
			g.sln("")
			if g.player.OliviaCount > 7 {
				g.foreground(2)
				g.sln("  Thinking of Olivia, you pick up the severed head and check it for life.")
				g.sln("")
				g.foreground(15)
				g.sln("  IT'S DEAD, BUT YOUR KINDNESS MAKES YOU BEAUTIFUL. (1 CHARM ADDED)")
				g.sln("")
				g.player.Cha += 1
			}
			g.more()
			g.sln("")
			return
		}
		g.lln(" `%WISELY!`0\"")
		g.sln("  \"You speak eloquently.  Your words do not fall upon deaf ears.")
		g.sln("   We believe you.  You are wise today.\"")
		g.sln("")
		g.more()
		g.sln("")
		g.raiseClass()
	}
}

// forestMysticalLevel is forest()'s mystical_level().
func (g *Game) forestMysticalLevel() {
	var ich string
	var thenum int64
	var guesses int64
	var guess int64

	g.forestEventHeader()
	g.foreground(2)
	g.sln("  While trekking through the forest, you come upon a small hut.")
	g.sln("")
	g.lln("  `2(`0K`2)nock On The Door")
	g.lln("  (`0B`2)ang On The Door")
	g.lln("  (`0L`2)eave It Be")
	g.sln("")
	g.commandPrompt()
	ich = strings.ToUpper(g.getkey())
	if !strings.Contains("LB", ich) {
		ich = "K"
	}
	g.sln(ich)
	if ich == "L" {
		g.sln("  You walk away.  Who needs magical instruction anyway!")
		g.sln("")
		g.more()
		return
	}
	g.sln("")
	if ich == "K" {
		g.sln("  You politely knock on the knotted wooden door.")
	}
	if ich == "B" {
		g.sln("  You bang on the door as hard you can!")
	}
	g.sln("")
	if g.random(4) == 1 {
		g.sln("  You wait a while but no one is home.  Maybe next time.")
		g.sln("")
		g.more()
		return
	}
	g.lln("`2  You are about to leave, when you hear a voice from above:")
	if g.player.Sex == "M" {
		g.lln("  `0\"Watcha doin' down there Sonny?!\" `2 You look up and see a wizened")
	} else {
		g.lln("  `0\"Watcha doin' down there Miss?!\"  `2You look up and see a wizened")
	}
	g.sln("  old man.")
	g.sln("")
	g.lln("  `0\"Tell ya what!  I'll give ya a mystical lesson if you can pass my")
	g.lln("   test!\" `2the old man giggles.")
	g.sln("")
	g.more()
	g.lln("`c`%                            ** THE TEST **")
	g.sln("")
	g.lln("  `0\"All right now!  I'm thinking of a number between 1 and 100.  I'll")
	g.sln("   give ya six guesses.\"")
	g.sln("")
	g.lln("`2  (The old man leans even farther out the window in anticipation)")
	g.sln("")
	thenum = g.random(100) + 1
	for guesses < 6 {
		guesses += 1
		g.lw("  `2Guess `0" + itoa(guesses) + " `2: `%")
		guess, _ = jsParseInt(g.consoleGetstr(GetStrOpts{Len: 3, Integer: true}))
		if guess > thenum {
			g.lln("  `0\"The number is lower than that!\"")
		} else if guess < thenum {
			g.lln("  `0\"The number is higher than that!\"")
		} else if guess == thenum {
			g.sln("")
			g.lln("  `0\"That's right! That's the number I was thinking of!  You read my mind!\"")
			g.lln("  `2The old man nearly falls from his window in his excitement!")
			g.lln("")
			g.lln("`%                         ** YOU HAVE PASSED THE TEST **")
			g.lln("")
			g.moreNomail()
			if g.player.SkillM > 39 {
				g.sln("")
				g.sln("  The old man attempts to teach you something, but fails.  You know more")
				g.sln("  than him.  The only thing he is able to give you is hope.")
				g.sln("")
				g.foreground(15)
				g.sln("  YOU RECEIVE FOUR EXTRA MYSTICAL SKILLS USE POINTS!")
				g.sln("")
				g.player.LevelM += 4
				g.more()
				return
			}
			g.raiseClass()
			return
		}
	}

	if guess != thenum {
		g.sln("")
		g.lln("  `2The old man drops his head and shakes it sadly.  You notice small")
		g.lln("  dandruff flakes drifting down from the window. `0\"No, no, NO!  The number")
		g.sln("  was " + itoa(thenum) + "!  Geez!  I won't teach such an unpromising student!\"")
		g.sln("")
		if g.player.SkillM > 19 {
			g.lln("  `0\"I'm already a Master anyway, old man,\" `2you retort coldy.")
			g.sln("")
		}
		g.lln("  `2He slams his window shut, and leaves you no recourse but to leave.")
		g.sln("")
		g.more()
	}
}

// forestThiefLevel is forest()'s thief_level().
func (g *Game) forestThiefLevel() {
	var ich string

	g.forestEventHeader()

	if g.player.SkillT > 39 {
		g.lln("`1`2  You are carefully moving through the forest, making absolutely no noise")
		g.sln("  when your sensitive ears pick up a twig breaking.  You circle around")
		g.lln("  towards the noise, and find it's not an animal, but The Master Thieves!")
		g.sln("")
		g.lln("  As they pass under a tree you are in, you call out. `0\"Ahh... Master")
		g.lln("  Thieves!  Do you think it would be possible to make even MORE noise?!\"`2")
		g.sln("")
		g.sln("  The group is very embarrassed, but they overcome it to chew the fat with")
		g.sln("  you.")
		g.sln("")
		g.lln("  `%EXTRA THIEVING USE FOR TODAY!")
		g.lln("")
		g.player.LevelT += 1
		g.more()
		return
	}

	g.lln("`2  You are innocently skipping through the forest, when you suddenly")
	g.sln("  notice you are surrounded by a group of rogues!")
	g.sln("")

	if g.player.SkillT < 20 {
		g.lln("  `0\"Greetings, `%" + g.player.Name + "`0.  We are members of the Master Thieves")
		g.sln("   Guild.  We know you are struggling to learn our ways.  We will")
		g.sln("   give you a lesson, for the price of one Gem.\"")
	} else {
		g.lln("  `0\"Greetings, `%" + g.player.Name + "`0.  We are proud that you have mastered")
		g.lln("   the skills.  But, we will continue to teach you anyway..For  ")
		g.lln("   the usual price!\" `2The scarred man laughs gleefully.")
	}

	for {
		g.sln("")
		g.lln("  `2(`0G`2)ive Them A Hard Earned Gem")
		g.lln("  `2(`0S`2)pit In Their Faces")
		g.lln("  `2(`0M`2)umble Apologies And Run")
		g.sln("")
		g.commandPrompt()
		ich = strings.ToUpper(g.getkey())
		g.sln(ich)
		if ich == "M" {
			g.sln("")
			g.lln("  `0\"You are nothing but a coward!  We will never let you join us!\"")
			g.sln("")
			g.more()
			return
		}
		if ich == "S" {
			g.sln("")
			g.lln("  `2You hawk a good sized piece of phlegm into the leader's face.")
			g.sln("")
			g.moreNomail()
			if g.player.Sex == "F" {
				g.lln("  `0\"You have spirit girl!  Maybe next time!\"`2")
			} else {
				g.lln("  `0\"You have spirit boy!  Maybe next time!\"`2")
			}
			g.sln("  The man laughs.")
			g.sln("")
			g.moreNomail()
			return
		}
		if ich == "G" {
			if g.player.Gem < 1 {
				g.sln("")
				g.sln("")
				g.lln("  `2You fumble through your pockets and find you don't posses a Gem.  You")
				g.sln("  don't think it would be wise to try to pull one over on the Master")
				g.sln("  Thieves' Guild.  You have a reputation to worry about!")
				g.sln("")
				g.moreNomail()
			} else {
				g.player.Gem -= 1
				g.sln("")
				g.sln("  You nonchalantly flip them a sparkling Gem.  The Thieves look impressed.")
				g.sln("")
				g.lln("  `0\"Nice rock.  Alright...True to our word, we will instruct you.\"`2")
				g.sln("")
				g.more()
				g.raiseClass()
				return
			}
		}
	}
}
