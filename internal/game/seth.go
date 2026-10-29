package game

import "strings"

// singleSeth is single_seth(): a woman flirting with an unmarried Seth
// Able.
func (g *Game) singleSeth() {
	p := g.player
	done := false

	if p.SeenViolet {
		g.sln("  The Bard seems occupied...Maybe tomorrow...")
		g.sln("")
		return
	}
	g.sln("")
	g.lrdfile("SETH")
	for {
		g.sln("")
		g.lw("  `0`2Your choice?  (`0? for menu`2) : ")
		ch := strings.ToUpper(g.getkey())
		if ch == "\r" {
			ch = "R"
		}
		g.sln(ch)
		g.sln("")
		switch ch {
		case "N", "R", "Q":
			done = true
		case "?":
			g.sln("")
			g.lrdfile("SETH")
		case "W":
			g.singleSethWink()
		case "F":
			g.singleSethFlutter()
		case "D":
			g.singleSethHanky()
		case "A":
			g.singleSethDrink()
		case "K":
			g.singleSethKiss()
		case "C":
			if !g.Settings.CleanMode {
				g.singleSethSeduce()
			} else {
				g.sln("")
				g.sln("")
				g.sln("  Sorry, your sysop has disabled that function...")
				g.sln("")
			}
		case "M":
			g.getState(false)
			if p.Divorced {
				g.sln("  Since you've just gotten out of a serious relationship, you feel")
				g.sln("  it'd be moving to quickly to get married again.")
			} else if p.MarriedTo > -1 {
				op := g.playerGet(p.MarriedTo)
				g.sln("")
				g.lln("  `0You kind of doubt you husband `)" + op.Name + " `0would")
				g.sln("  like that...")
				g.moreNomail()
			} else if g.state.MarriedToViolet == int64(p.Record) {
				g.sln("")
				g.lln("  `0You kind of doubt you wife `#Violet `0would")
				g.sln("  like that...")
				g.moreNomail()
			} else {
				g.singleSethMarriage()
			}
		}
		if done || p.SeenViolet {
			break
		}
	}
	g.sln("")
}

// singleSethSeduce is single_seth()'s seduce().
func (g *Game) singleSethSeduce() {
	p := g.player
	insults := []string{"filthy harlot", "slut"}

	g.sln("  You turn on the charm, moving your body seductively")
	g.sln("  against him, you ask him to come upstairs with you...")
	g.sln("")
	g.mswait(2000)
	p.SeenViolet = true
	if p.MarriedTo > -1 {
		g.sln("")
		g.foreground(10)
		g.sln("  Seth Able is appalled that you would even suggest such")
		g.sln("  a thing!")
		g.sln("")
		g.sln("  He calls you a filthy harlot!")
		g.logLine("`%  Seth Able `2calls `0" + p.Name + " `2a `)" + insults[g.random(int64(len(insults)))] + "`2!")
		return
	}
	if p.Cha >= 32 {
		switch g.random(4) {
		case 2:
			g.sln("  He smiles invitingly...  ")
			g.mswait(4000)
			g.sln("")
			g.foreground(2)
			g.sln("  Hours later.. You saunter downstairs. When a bearded drunk")
			g.sln("  asks you all the racket was, you smile knowingly and look away..")
			g.sln("  The drunks are mystified!")
			g.sln("")
			g.lln("  `0YOU GET `%" + itoa(p.Level*40) + " `0EXPERIENCE!")
			p.Exp += p.Level * 40
			g.logLine("`0  " + p.Name + " `2got laid by `%Seth Able`2!")
			p.Laid++
			g.tournamentCheck()
		case 0:
			g.sln("  He tells you he has a headache!")
			g.sln("  You are very disappointed.")
		case 1:
			g.sln("  He tells you he is not in the mood!")
			g.sln("  You saunter around the bar disappointed.")
		case 3:
			g.sln("  He tells you he is too exausted from last night!")
			g.sln("  You saunter around the bar disappointed.")
		}
	} else {
		g.sln("")
		g.foreground(2)
		g.sln("  He shoves you away harshly!")
		g.sln("  He calls you a filthy whore!")
		g.sln("  You trudge away from him dejectedly..")
		g.sln("  The entire bar laughs at your misfortune!!")
		g.sln("")
		g.foreground(4)
		g.sln("  YOUR HITPOINTS GO DOWN TO 1!")
		p.HP = 1
		g.logLine("`5  " + p.Name + " `2 was called a whore by `%Seth Able`2!")
	}
	g.more()
}

// singleSethWink is single_seth()'s wink().
func (g *Game) singleSethWink() {
	p := g.player
	p.SeenViolet = true
	if p.Cha >= 1 {
		g.sln("")
		g.foreground(15)
		g.sln("  You wink at Seth Able seductively..")
		g.mswait(5000)
		g.foreground(2)
		g.sln("  He blushes and smiles!!")
		g.sln("  You are making progress with him!")
		g.sln("")
		g.foreground(15)
		g.sln("  You receive " + prettyInt(p.Level*5) + " experience!")
		p.Exp += p.Level * 5
		g.tournamentCheck()
	} else {
		g.foreground(15)
		g.sln("  You wink at Seth Able seductively..")
		g.mswait(5000)
		g.foreground(4)
		g.sln("  He looks the other way!")
		g.sln("  You nearly die of embarrassment!")
	}
	g.more()
}

// singleSethFlutter is single_seth()'s flutter().
func (g *Game) singleSethFlutter() {
	p := g.player
	p.SeenViolet = true
	g.foreground(15)
	g.sln("  You turn to Seth Able and flutter your eyelashes violently..")
	g.mswait(5000)
	g.foreground(2)
	if p.Cha >= 2 {
		g.sln("  He smiles broadly and winks back!")
		g.sln("  Your relationship with him is taking off!")
		g.sln("")
		g.foreground(15)
		g.sln("  You receive " + prettyInt(p.Level*10) + " experience!")
		p.Exp += p.Level * 10
		g.tournamentCheck()
	} else {
		g.sln("  He doesn't seem interested and starts talking to Violet!")
		g.sln("  You nearly faint from embarresment!")
		g.sln("")
		g.foreground(4)
		g.sln("  YOU LOSE " + prettyInt(p.Level) + " HIT POINTS!")
		p.HP -= p.Level
		if p.HP < 1 {
			p.HP = 1
		}
	}
	g.more()
}

// singleSethHanky is single_seth()'s hanky().
func (g *Game) singleSethHanky() {
	p := g.player
	g.foreground(15)
	g.sln("  You innocently drop your lace hanky on the ground..")
	g.mswait(5000)
	p.SeenViolet = true
	g.foreground(2)
	if p.Cha >= 4 {
		g.sln("  Seth sees it, picks it up and offers it to you!")
		g.sln("  He smiles hugely!  You are elated!")
		g.sln("")
		g.foreground(15)
		g.sln("  YOU RECEIVE " + prettyInt(p.Level*20) + " EXPERIENCE!")
		p.Exp += p.Level * 20
		g.tournamentCheck()
	} else {
		g.sln("  He ignores your plea for attention, and looks away!")
		g.sln("  You nearly faint from embarresment!")
		g.sln("")
		g.foreground(4)
		g.sln("  YOU LOSE " + prettyInt(p.Level*3) + " HIT POINTS!")
		p.HP -= p.Level * 3
		if p.HP < 1 {
			p.HP = 1
		}
	}
	g.more()
}

// singleSethDrink is single_seth()'s drink().
func (g *Game) singleSethDrink() {
	p := g.player
	g.foreground(15)
	g.sln("  You wave Seth over and ask him to buy you a drink..")
	g.mswait(5000)
	p.SeenViolet = true
	g.foreground(2)
	if p.Cha >= 8 {
		g.sln("  He buys you the most expensive drink in the house, and gives")
		g.sln("  you a hug!  You cherish his warmth!")
		g.sln("")
		g.foreground(15)
		g.sln("  YOU RECEIVE " + prettyInt(p.Level*30) + " EXPERIENCE!")
		p.Exp += p.Level * 30
		g.tournamentCheck()
	} else {
		g.sln("  He looks startled at your advances and mumbles a lame excuse!")
		g.sln("  You nearly die of embarresment!")
		g.sln("")
		g.foreground(4)
		g.sln("  YOU LOSE " + prettyInt(p.Level*5) + " HIT POINTS!")
		p.HP -= p.Level * 5
		if p.HP < 1 {
			p.HP = 1
		}
	}
	g.more()
}

// singleSethKiss is single_seth()'s kiss().
func (g *Game) singleSethKiss() {
	p := g.player
	g.foreground(15)
	g.sln("  You reach over his mandolin, and kiss him soundly...")
	g.foreground(2)
	g.mswait(5000)
	p.SeenViolet = true
	if p.Cha >= 16 {
		g.sln("  He doesn't object!  He kisses you back! ")
		g.sln("  You become over excited and force yourself away!")
		g.sln("")
		g.foreground(15)
		g.sln("  YOU RECEIVE " + prettyInt(p.Level*40) + " EXPERIENCE!")
		p.Exp += p.Level * 40
		g.tournamentCheck()
	} else {
		g.sln("  He pushes your face away, and says he is just")
		g.sln("  not ready for that kind of relationship!")
		g.sln("  You nearly die of embarresment!")
		g.sln("")
		g.foreground(4)
		g.sln("  YOU LOSE " + prettyInt(p.Level*10) + " HITPOINTS!")
		p.HP -= p.Level * 10
		if p.HP < 1 {
			p.HP = 1
		}
	}
	g.more()
}

// singleSethMarriage is single_seth()'s marriage().
func (g *Game) singleSethMarriage() {
	p := g.player
	p.SeenViolet = true
	switch {
	case p.Cha < 10:
		g.sln("")
		g.sln("")
		g.sln("  Seth Able looks at you seriously.")
		g.sln("")
		g.lln("  `0\"" + p.Name + "`0!  I will not consider it!\"`2")
		g.sln("")
		g.sln("  His words carry through the air and everyone looks at you.")
		g.sln("")
		g.sln("  Your world crumbles to the ground!")
		g.sln("")
		g.moreNomail()
	case p.Cha < 20:
		g.sln("")
		g.sln("")
		g.sln("  Seth Able looks at you seriously.")
		g.sln("")
		g.lln("  `0\"" + p.Name + "`0..We are just not right for each other.. We are")
		g.sln("   so different - And in this case opposites DON'T attract.\"")
		g.sln("")
		g.foreground(2)
		g.sln("  His words carry through the air and everyone looks at you.")
		g.sln("")
		g.sln("  You feel so stupid.")
		g.sln("")
		g.moreNomail()
	case p.Cha < 50:
		g.sln("")
		g.sln("")
		g.sln("  Seth Able looks at you seriously.")
		g.sln("")
		g.lln("  `0\"" + p.Name + "`0.. I cannot say yes.  I do find you attractive,")
		g.sln("   but my life is here, serving the people.\"")
		g.sln("")
		g.lln("  `2His words carry through the air and everyone looks at you.")
		g.sln("")
		g.sln("  Your heart hurts.")
		g.sln("")
		g.moreNomail()
	case p.Cha < 80:
		g.sln("")
		g.sln("")
		g.sln("  Seth Able looks at you seriously.")
		g.sln("")
		g.lln("  `0\"" + p.Name + "`0.. I do love you - You are one of the most")
		g.sln("   wonderful women in the realm!  I still cannot say yes.\"")
		g.sln("")
		g.lln("  `2His words carry through the air and everyone looks at you.")
		g.sln("")
		g.sln("  Humiliation turns your face a deep scarlet.")
		g.sln("")
		g.moreNomail()
	case p.Cha < 100:
		g.sln("")
		g.sln("")
		g.sln("  Seth Able looks at you seriously.")
		g.sln("")
		g.lln("  `0\"" + p.Name + "`0.. You are truly a thing of beauty.  Please")
		g.sln("   have mercy on me and do not hate me when I tell I cannot.\"")
		g.sln("")
		g.lln("  `2His words carry through the air and everyone looks at you.")
		g.sln("")
		g.sln("  A fierce anger burns inside.")
		g.sln("")
		g.moreNomail()
	case p.Cha < 120:
		g.sln("")
		g.sln("")
		g.sln("  Seth Able looks at you very seriously.")
		g.sln("")
		g.lln("  `0\"" + p.Name + "`0.. You are so dear to me.  You are truly THE")
		g.sln("  most beautiful woman in the realm.  I pray to God you will")
		g.sln("  not hate me when I tell you I cannot mary you.  You see, I")
		g.sln("  was married once before.  My wife was killed many years ago.")
		g.sln("  I cannot go through that pain again.  I am so sorry.")
		g.sln("")
		g.lln("  `2His words carry through the air and everyone looks at you.")
		g.sln("")
		g.sln("  You see tears in the eyes of many onlookers.")
		g.sln("")
		g.moreNomail()
	default:
		g.sln("")
		g.sln("")
		g.sln("  Seth Able looks at you very seriously.")
		g.sln("")
		g.lln("  `0\"" + p.Name + "`0.. The many times you've visited me have")
		g.sln("   truly done my heart good.\"")
		g.sln("")
		g.moreNomail()
		g.sln("   \"I do love you.  And if you'll have me, I will be the best")
		g.sln("   husband I can be for you.\"")
		g.sln("")
		g.moreNomail()
		g.lln("  `2The bard then takes you in his arms and squeezes you tightly.")
		g.sln("")
		g.sln("  The news of the wedding travels like wildfire - shortly, a vast")
		g.sln("  assemblege of townspeople are sitting rather impatiently in the")
		g.sln("  the great church found near the outskirts of town.")
		g.sln("")
		g.moreNomail()
		g.lln("  `c`%                    ** THE WEDDING **")
		g.sln("")
		g.getState(true)
		if g.state.MarriedToSeth == -1 || g.state.MarriedToSeth == int64(p.Record) {
			g.state.MarriedToSeth = int64(p.Record)
			g.putState()
			if p.Laid > 0 {
				g.sln("  Wearing a slightly off-white dress you walk demurely down the aisle.")
				g.logLine("`)  ** `0" + p.Name + " `2has MARRIED `%Seth Able`2!`) **")
			} else {
				g.lln("  `0Wearing a PURE `%white `0dress you walk demurely down the asle.")
				g.sln("")
				g.lln("  `0Your heart is completely fulfilled.")
				g.logLine("`)  ** `0" + p.Name + " `2has MARRIED `%Seth Able`2 in `%WHITE`2!`) **")
				g.sln("")
				p.Exp += p.Level * 5000
				g.lln("  `%You feel heavenly!  You receive " + prettyInt(p.Level*5000) + " experience!")
			}
			g.sln("")
			g.lln("  `0You barely notice the ceremony taking place - your eyes and")
			g.sln("  your husbands never leave each other.  ")
			g.sln("")
			g.moreNomail()
			g.sln("  At the appropriate time, the bard kisses you deeply and long.")
			g.sln("")
			g.lln("  `%THE ENTIRE CHURCH ROARS ITS APPROVAL!")
			g.sln("")
			g.moreNomail()
			g.tournamentCheck()
		} else {
			g.putState()
			mop := g.playerGet(g.state.MarriedToSeth)
			if p.Laid > 0 {
				g.sln("  Wearing a slightly off-white dress you walk demurely to the aisle.")
			} else {
				g.lln("  `0Wearing a PURE `%white `0dress you walk demurely to the aisle.")
			}
			g.moreNomail()
			g.lln("  `0When you get there, you can't believe your ears...")
			g.lln("  `0\"I now pronounce you husband and wife... you may kiss the bride\"")
			g.sln("")
			g.lln("  `0Seth Able has just married " + mop.Name + "!")
			g.moreNomail()
			g.lln("  `0You run away sobbing.")
		}
	}
}

// bardSong is bard_song(): Seth Able sings (once a day).
func (g *Game) bardSong() {
	p := g.player
	g.sln("")
	g.sln("")
	if p.SeenBard {
		g.sln("  \"I'm sorry, but my throat is too dry..  Perhaps tomorrow.\"")
		g.sln("")
		return
	}
	p.SeenBard = true
	g.sln("")
	g.sln("")
	g.lln("  `%(Seth Able clears his throat)`2")
	g.sln("")
	if g.random(2) == 1 {
		switch g.random(5) {
		case 0:
			g.lln("  `2...\"The children are missing, the children are gone\"...")
			g.sln("  ...\"They have no pillow to lay upon,\"...")
			g.sln("  ...\"The hopes of the people are starting to dim\"...")
			g.lln("  ...\"They are gone, because the `4Dragon `2has eaten them\"...")
			g.sln("")
			g.sln("  Tears run down your face.  You swear to avenge the children.")
			g.sln("")
			g.lln("  `#YOU RECEIVE AN EXTRA USER BATTLE FOR TODAY!")
			p.PvPFights++
		case 1:
			g.lln("  `2...\"The Gods have powers, the Gods are just\"...")
			g.sln("  ...\"The Gods help us people, when they must\"...")
			g.sln("  ...\"Gods can heal the sick, even the cancered\"...")
			g.sln("  ...\"Pray to the Gods, and you will be answered\"...")
			g.sln("")
			g.sln("  You find yourself wishing for more money.")
			g.sln("")
			g.lln("  `#SOMEWHERE, MAGIC HAS HAPPENED!")
			p.Bank *= 2
			if p.Bank > 2000000000 {
				p.Bank = 2000000000
			}
		case 2:
			g.lln("  `2...\"This is the story of a bard\"...")
			g.sln("  ...\"Stabbed in the hand with a shard\"...")
			g.sln("  ...\"They said he'd never play, but they would rue the day\"...")
			g.sln("  ...\"He practiced all the time, becoming good before long\"...")
			g.sln("  ...\"Able got his revenge, by proving them wrong\"...")
			g.sln("")
			g.sln("  The song makes you feel sad, yet happy.")
			g.bardSongForestFights(2, "  `#YOU RECEIVE TWO MORE FOREST FIGHTS FOR TODAY!")
		case 3:
			g.lln("  `2...\"Let me tell you the Legend, of the Red Dragon\"...")
			g.sln("  ...\"They said he was old, that his claws were saggin\"...")
			g.sln("  ...\"It's not so, don't be fooled\"...")
			g.sln("  ...\"He's alive and well, still killing where he ruled\"...")
			g.sln("")
			g.sln("  The song makes you feel a strange wonder, an awakening.")
			if p.ForestFights < 32000 {
				p.ForestFights++
				g.sln("")
				g.lln("  `#YOU RECEIVE ONE MORE FOREST FIGHT FOR TODAY!")
			}
		case 4:
			g.lln("  `2...\"&hero `2has a story, that must be told\"...")
			g.sln("  ...\"He is already a legend, and he ain't even old\"...")
			g.sln("  ...\"He can drink a river of blood and not burst\"...")
			g.sln("  ...\"He can swallow a desert and never thirst\"...")
			g.sln("")
			g.sln("  The hero has inspired you, and in doing so, made you a better warrior.")
			g.sln("")
			g.lln("  `#YOUR HITPOINTS INCREASE BY ONE!")
			p.HPMax++
			if p.HPMax > 32000 {
				p.HPMax = 32000
			}
		}
		return
	}
	if p.Sex == "M" {
		switch g.random(4) {
		case 0:
			g.lln("  `2...\"`0" + p.Name + " `2was a warrior, a man\"...")
			g.sln("  ...\"When he wants to do a thing, he can\"...")
			g.sln("  ...\"To live, to die, it's all the same\"...")
			g.sln("  ...\"But it is better to die, than live in shame\"...")
			g.sln("")
			g.sln("  The song makes you feel your heritage.")
			g.sln("")
			g.lln("  `#YOUR HIT POINTS ARE MAXED OUT!")
			p.HP = p.HPMax
		case 1:
			g.lln("  `2...\"There once was a warrior, with a beard\"...")
			g.lln("  ...\"He had a power, yes `0" + p.Name + " `2was feared\"...")
			g.sln("  ...\"Nothing he did, could ever be wrong\"...")
			g.sln("  ...\"He was quick, and he was strong\"...")
			g.sln("")
			g.sln("  The song makes you feel powerful!")
			g.bardSongForestFights(3, "  `#YOU RECEIVE THREE MORE FOREST FIGHTS FOR TODAY!")
		case 2:
			g.lln("  `2...\"Waiting in the forest waiting for his prey\"...")
			g.lln("  ...\"`0" + p.Name + " `2didn't care what they would say\"...")
			g.sln("  ...\"He killed in the town, the lands\"...")
			g.sln("  ...\"He wanted evil's blood on his hands\"...")
			g.sln("")
			g.sln("  The song makes you feel powerful!")
			g.bardSongForestFights(3, "  `#YOU RECEIVE THREE MORE FOREST FIGHTS FOR TODAY!")
		case 3:
			g.lln("  `2...\"A true man was `0" + p.Name + "`2, a warrior proud\"...")
			g.sln("  ...\"He voiced his opinions meekly, never very loud\"...")
			g.lln("  ...\"But he ain't no wimp, he took `#Violet`2 to bed\"...")
			g.sln("  ...\"He's definately a man, at least that's what she said!\"...")
			g.sln("")
			g.sln("  The song makes you glad you are male!")
			g.bardSongForestFights(2, "  `#YOU RECEIVE TWO MORE FOREST FIGHTS FOR TODAY!")
		}
	} else {
		switch g.random(2) {
		case 0:
			g.lln("  `2...\"`0" + p.Name + " `2was a warrior, a queen\"...")
			g.sln("  ...\"She was a beauty, and she was mean\"...")
			g.sln("  ...\"She could melt a heart, at a glance\"...")
			g.sln("  ...\"And men would pay, to see her dance!\"...")
			g.sln("")
			g.sln("  The song makes you feel pretty!")
			g.sln("")
			g.lln("  `#YOU RECEIVE A CHARM POINT!")
			p.Cha++
			if p.Cha > 32000 {
				p.Cha = 32000
			}
		case 1:
			g.lln("  `2...\"There once was a woman, of exceeding fame\"...")
			g.lln("  ...\"She had a power, `0" + p.Name + " `2was her name\"...")
			g.sln("  ...\"Nothing she did, could ever be wrong\"...")
			g.sln("  ...\"She was quick, and she was strong\"...")
			g.sln("")
			g.sln("  The song makes you feel powerful!")
			g.bardSongForestFights(3, "  `#YOU RECEIVE THREE MORE FOREST FIGHTS FOR TODAY!")
		}
	}
}

// bardSongForestFights is the forest fight reward repeated in
// bard_song(): add n fights (capped at 32000) and announce it.
func (g *Game) bardSongForestFights(n int64, msg string) {
	p := g.player
	if p.ForestFights < 32000 {
		p.ForestFights += n
		if p.ForestFights > 32000 {
			p.ForestFights = 32000
		}
		g.sln("")
		g.lln(msg)
	}
}

// talkToBard is talk_to_bard(): Seth Able's corner of the inn.
func (g *Game) talkToBard() {
	p := g.player
	done := false

	menu := func() {
		g.lrdfile(pick(p.Sex == "M", "BARD", "BARDF"))
	}

	g.getState(false)
	menu()
	for !done {
		g.sln("")
		if g.state.MarriedToSeth == int64(p.Record) {
			g.lln("  `2Seth Able looks at you lovingly. (`0? for menu`2)")
		} else {
			g.lln("  `2Seth Able looks at you expectantly.  (`0? for menu`2)")
		}
		g.commandPrompt()
		ch := strings.ToUpper(g.getkey())
		g.sln(ch)
		switch ch {
		case "?":
			menu()
		case "R", "Q", "\r":
			done = true
		case "A":
			g.bardSong()
		case "F":
			if p.Sex == "M" {
				g.sln("")
				g.sln("  You would rather flirt with Violet.")
				break
			}
			g.getState(false)
			if g.state.MarriedToSeth > -1 {
				if g.state.MarriedToSeth == int64(p.Record) {
					g.sln("")
					if p.SeenViolet {
						g.lln("  `2You don't wish to interrupt his song.")
					} else {
						p.SeenViolet = true
						switch g.random(8) {
						case 0:
							g.lln("  `0\"I love you too, honey!\"")
						case 1:
							g.lln("  `2Seth gives you a quick peck on the cheek.")
						case 2:
							g.lln("  `0\"I'd love to sweetie, but duty calls.\"")
						case 3:
							g.lln("  `0\"We'll have plenty of time for that tonight!\"")
						case 4:
							g.lln("  `2In the roar of the crowd, he doesn't hear.")
						case 5:
							g.lln("  `2You don't wish to bother your busy husband.")
						case 6:
							g.lln("  `0\"After work, you're mine!\" `2he laughs.")
						case 7:
							g.lln("  `0\"Later, honey!\" `2Seth winks hugely at you.")
						}
					}
				} else {
					g.talkToBardMarried()
				}
			} else {
				g.singleSeth()
			}
		}
	}
	g.sln("")
}

// talkToBardMarried is talk_to_bard()'s flirting with Seth Able when he
// is married to someone else.
func (g *Game) talkToBardMarried() {
	p := g.player
	g.sln("")
	g.foreground(2)
	op := g.playerGet(g.state.MarriedToSeth)
	if p.SeenViolet {
		if p.Cha > op.Cha {
			g.sln("  You decide you've done enough ")
			g.sln("  homewrecking for one day.")
		} else {
			g.sln("  You are too shamed to even think")
			g.sln("  such a thought.")
		}
		return
	}
	g.sln("  You are about to give Seth your best, when")
	g.lln("  he turns away.  You see a shiny new `0ring `2on")
	g.sln("  his hand. ")
	g.sln("")
	g.lw("  Attempt to seduce him anyway? `![`0N`!] : ")
	ch := strings.ToUpper(g.getkey())
	if ch != "Y" {
		ch = "N"
	}
	g.sln(ch)
	g.sln("")
	if ch != "Y" {
		return
	}
	g.sln("  You rub your body against the man in")
	g.sln("  the most disarming way you can.")
	g.sln("")
	g.moreNomail()
	p.SeenViolet = true
	if p.Cha > op.Cha {
		g.lln("`2  Seth looks at you hungrily.")
		g.sln("")
		g.lln("  `0\"Dear " + p.Name + "`0..  You make me")
		g.lln("  think unworthy thoughts!\" `2the bard")
		g.sln("  cries in agony.")
		g.sln("")
		g.logLine("  `0" + p.Name + " `2has made `%Seth Able`2 falter.")
		g.mailTo(g.state.MarriedToSeth, "\n"+
			"  `0**** `)UGLY RUMOR `0****\n"+
			"`0-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-\n"+
			"  `2Anonymous sources tell you "+p.Name+"`2 was seen"+
			"  `2attempting to seduce your husband!")
		g.moreNomail()
		g.lln("  `2You have caused the bard to falter!")
		g.sln("")
		g.lln("  `%YOU RECEIVE " + prettyInt(p.Level*100) + " EXPERIENCE!")
		p.Exp += p.Level * 100
		g.tournamentCheck()
	} else {
		g.lln("`2  Seth looks at you angrily.")
		g.sln("")
		g.lln("  `0\"" + p.Name + "`0..  You are an extremely")
		g.sln("  silly wench.  My God woman, you")
		g.sln("  would need at least " + prettyInt(op.Cha-p.Cha) + " more")
		g.sln("  charm to be considered female!\"")
		g.sln("")
		g.logLine("  `%Seth Able`2 called `0" + p.Name + " `2a silly wench!")
	}
}
