package game

import (
	"strings"
)

// forestOliviaGetch is olivia()'s olivia_getch(): wait for one of chars
// (Enter picks the first).
func (g *Game) forestOliviaGetch(chars string) string {
	var ich string

	for {
		ich = strings.ToUpper(g.getkey())
		if ich == "\r" {
			ich = chars[:1]
		}
		if strings.Contains(chars, ich) {
			break
		}
	}
	g.sln(ich)
	return ich
}

// forestOliviaAsshole is olivia()'s asshole(): Olivia after you booted her.
func (g *Game) forestOliviaAsshole() {
	var ich string
	var cannot bool

	g.lln("`c  `%THE WAILING GROWS LOUDER.")
	g.sln("")
	g.lln("  You investigate - only to find a womans head on the ground.")
	g.lln("")
	if g.player.Sex == "M" {
		g.lln("  `0\"I see you, foolish boy.  Leave me alone!\" `2the head screams savagely.")
	} else {
		g.lln("  `0\"I see you, foolish girl.  Leave me alone!\" `2the head shouts.")
	}
	g.sln("")
	g.lln("  `2(`0A`2)pologize for what you did last time")
	g.lln("  `2(`0P`2)lay some \"head ball\"")
	g.sln("")
	g.lw("  `2Well? [`0A`2] `2: `0")
	ich = g.forestOliviaGetch("AP")
	g.sln("")
	if ich == "P" {
		// BUG (lord.js): there are only two "head ball" outcomes for
		// random(3); a 2 falls through to the apology below.  Kept.
		switch g.random(3) {
		case 0:
			g.lln("  `2You are distracted by a huge spider climbing up the cave wall.")
			g.sln("")
			g.lln("  `0Thinking quickly, you pick up the head and throw it at it!")
			g.sln("")
			g.more()
			g.lln("  `2You hear a satisfying squishing sound.  As you leave you hear the")
			g.lln("  head screaming after you.  `0\"How dare you!!  You're going to pay!\"")
			g.sln("")
			g.lln("  `%THE EXCERCISE GIVES YOU STRENGTH FOR ANOTHER FOREST FIGHT!")
			g.player.ForestFights += 1
			if g.player.ForestFights > 32000 {
				g.player.ForestFights = 32000
			}
			g.sln("")
			g.more()
			return
		case 1:
			g.lln("  `2Seeing the fine texture of the heads hair and the rip in your")
			g.sln("  garment, you decide to do something about it.")
			g.sln("")
			g.more()
			g.lln("  `0\"PUT ME DOWN!\" `2the head screams as you cut off a few strands of")
			g.sln("  its long hair and mend the tear.")
			g.sln("")
			cannot = false
			if g.player.Clss == 3 {
				if g.player.Level < 2 {
					cannot = true
				}
			} else {
				if g.player.Level < 6 {
					cannot = true
				}
			}
			if cannot {
				g.lln("  `4You lack the skill needed to sew your " + g.player.Arm + ".")
			} else {
				g.lln("  `2There!  Your " + g.player.Arm + " is better than ever!")
				g.sln("")
				g.lln("  `%YOU FEEL SO CLEVER YOU GAIN THE STRENGTH FOR ANOTHER FOREST FIGHT!")
				g.player.ForestFights += 1
				if g.player.ForestFights > 32000 {
					g.player.ForestFights = 32000
				}
			}
			g.sln("")
			g.more()
			return
		}
	}
	switch g.random(3) {
	case 0:
		g.lln("  `0\"I'm very sorry, ma'am.  I was a jerk.\"")
		g.sln("")
		g.lln("  `%She narrows her eyes at you.")
		g.sln("")
		g.more()
		g.lln("  `0\"And I uh, think you are perfectly beheading.  I mean, becoming!\"")
		g.sln("")
		g.lln("  `#\"AWK!! GO DIE YOU PIECE OF <choking spasm>\" `2she screams.")
		g.sln("")
		g.lln("  `4YOU FEEL SO WOEBEGONE YOU LOSE A FOREST FIGHT FOR TODAY.")
		g.sln("")
	case 1:
		g.lln("  `0\"Hey?  My old pal!  What are you doing in this 'neck' of the")
		g.lln("  woods?\" `2you inquire sincerely.")
		g.sln("")
		g.lln("  `%She narrows her eyes at you.")
		g.sln("")
		g.more()
		g.lln("  `0\"Don't be mad!  Please!  Geez, time to try decapit..I mean,")
		g.sln("  decaffinated!\"")
		g.sln("")
		g.lln("  `#\"GO EAT BUGS AND DIE, YOU HEARTLESS TROLL!\" `2she screams.")
		g.sln("")
		g.lln("  `4YOU FEEL SO CRAPPY YOU LOSE A FOREST FIGHT FOR TODAY.")
		g.sln("")
	case 2:
		g.lln("  `0\"You know before?  I was just kidding 'round!  Use your head")
		g.lln("  woman! `2you say earnestly.")
		g.sln("")
		g.lln("  `%She narrows her eyes at you.")
		g.sln("")
		g.more()
		g.lln("  `0Yeah, just joking!  Man you are so guillotine..I mean, er")
		g.lln("  gullible!\"")
		g.sln("")
		g.lln("  `#\"I HATE YOU!\" `2she screams.")
		g.sln("")
		g.lln("  `4YOU FEEL SO WOEBEGONE YOU LOSE A FOREST FIGHT FOR TODAY.")
		g.sln("")
	}
	g.player.ForestFights -= 1
	g.more()
}

// forestOlivia is forest()'s olivia(): the severed head in the cave.
func (g *Game) forestOlivia() {
	var och string
	var tmp int64

	if g.player.OliviaAsshole {
		g.forestOliviaAsshole()
		return
	}
	if !g.player.Olivia {
		g.lln("`c  `%THE WAILING GROWS LOUDER")
		g.sln("")
		g.lln("  `2You catch a glimpse of something moving in the mouth of the")
		g.sln("  cave.")
		g.sln("")
		g.lln("  `2(`0I`2)nvestigate further")
		g.lln("  `2(`0G`2)et smart and leave it alone")
		g.sln("")
		g.lw("  `2Well? [`0I`2] : `0")
		och = g.forestOliviaGetch("IG")
		if och == "G" {
			g.sln("")
			g.lln("  `2You hurry away from this evil place.")
			g.sln("")
			g.more()
			return
		}
		g.sln("")
		g.sln("  You slowly make your way inside the the damp cave.  You feel a cold")
		g.sln("  breeze - or perhaps you just shivered for another reason.")
		g.sln("")
		g.lln("  `4YOU TRIP OVER SOMETHING!")
		g.sln("")
		g.moreNomail()
		g.lln("  `2You blindly reach down `8- `2your fingers find something wet `8- `2You have")
		g.lln("  put your hand inside the mouth of a severed head.  `0You scream like a child!")
		g.sln("")
		g.moreNomail()
		g.lln("  `#\"Oh, do shut up!\" `2the head implores you, scowling.")
		g.sln("")
		if g.player.Sex == "M" {
			g.lln("  `2You stare at the head in shock.  (which really isn't bad looking)")
		} else {
			g.lln("  `2You stare at the head in shock.")
		}
		g.sln("")
		g.lln("  `2(`0A`2)sk the head who she is")
		g.lln("  `2(`0B`2)oot her a distance")
		g.sln("")
		g.lln("  `2After careful consideration you.. [`0A`2] `8: `%")
		och = g.forestOliviaGetch("AB")
		g.sln("")
		if och == "B" {
			g.lln("  `0\"C'mere you freak of nature!\" `2you taunt, grabbing her by the hair.")
			g.sln("")
			g.lln("  `#\"NooOoooOooo!\" `2she screams as you toss her deep into the underbrush.")
			g.sln("")
			if g.player.Sex == "M" {
				g.lln("  `2You can't stop thinking about what she would be like in bed.")
			} else {
				g.lln("  `2You can't help but wonder if you made the right choice.")
			}
			g.sln("")
			g.player.OliviaAsshole = true
			g.more()
			return
		}
		g.player.Olivia = true
		g.lrdfile("OLIVIA")
		return
	}
	g.lln("`c  `%WAIT A SEC!")
	g.sln("")
	g.lln("  `2It's just your old pal `%Olivia `2the bodyless woman.")
	g.sln("")
	if g.player.OliviaCount == 0 {
		g.lln("  `0\"I was wondering if you were coming back,\" `2she informs you.")
		g.sln("")
	}
	if g.player.OliviaCount == 1 {
		g.lln("  `0\"You want more information I suppose.\" `2she says rather sadly.")
		g.sln("")
	}
	if g.player.OliviaCount == 2 {
		g.lln("  `0\"Life is worthless, I am worthless.\" `2she drones.")
		g.sln("")
		g.sln("  You realize she is dangerously depressed.")
		g.sln("")
	}
	if g.player.OliviaCount == 3 {
		g.lln("  `0\"I wish I had a friend.  You just use me.\" `2she rebukes.")
		g.sln("")
		g.lln("  She's still depressed apparently.")
		g.sln("")
	}
	if g.player.OliviaCount == 4 {
		g.lln("  `0\"It's about time.  I do get bored you know,\" `2she scolds.")
		g.lln("")
	}
	if g.player.OliviaCount == 5 {
		g.lln("  `0\"You're looking sharp today, " + g.player.Name + "`0,\" `2she complements.")
		g.lln("")
	}
	if g.player.OliviaCount == 6 {
		g.lln("  `0\"" + g.player.Name + "`0! I've missed you!\" `2she exclaims.")
		g.sln("")
		g.sln("  Odd, considering you were just here not too long ago.")
		g.sln("")
	}
	if g.player.OliviaCount > 7 && g.player.Sex == "F" {
		g.player.OliviaCount = 7
	}
	if g.player.OliviaCount == 7 {
		g.lln("  `0\"Hello, " + g.player.Name + "`0.  I'm glad you could visit.\"`2 she smiles.")
		g.sln("")
		g.sln("  She seems happy.")
		g.sln("")
	}
	if g.player.OliviaCount == 8 {
		g.lln("  `0\"Do you find me attractive?\" `2she asks pensively.")
		g.sln("")
		g.lln("  `%\"Why?  Are you unattached?\" `2you smirk.")
		g.sln("")
		g.lln("  `0\"Do shutup!\" `2she screeches.")
		g.sln("")
	}
	if g.player.OliviaCount == 9 {
		g.lln("  `0\"Am I pretty?\" `2she asks hopefully.")
		g.sln("")
		g.lln("  `%\"Well, uh, you don't weigh too much.\" `2you smirk.")
		g.sln("")
		g.lln("  `0\"I'm being serious!\" `2she screeches.")
		g.lln("")
	}
	// DIFF: The test for !== 25 wasn't here... always show the history now.
	if g.player.OliviaCount == 10 || (g.player.OliviaCount != 25 && g.player.OliviaCount > 10 && g.random(7) == 1) {
		g.player.OliviaCount += 1
		if !g.Settings.CleanMode {
			// BUG (lord.js): olivia_count was already incremented, so it is
			// never 10 here and the first proposal is never shown.  Kept.
			if g.player.OliviaCount == 10 {
				g.lln("  `0\"I'm through with the hints!  Do you want to mess around?\" `2she")
				g.sln("  states bluntly.")
				g.sln("")
				g.lln("  `%\"Olivia, I'm flattered - But .. you don't have a body.\"")
				g.sln("")
				g.more()
				g.lln("  `0\"If we put our minds to it, we can figure out a way.\" `2she grins")
				g.lln("  wickedly.")
				g.sln("")
				g.more()
			} else {
				g.lln("  `0\"I'm ready for love again.    Do you want to mess around?\" `2she")
				g.sln("  states bluntly.")
				g.lln("")
			}
			g.lln("  `2(`0B`2)e pleasured by Olivia")
			g.lln("  `2(`0T`2)urn her down.")
			g.sln("")
			g.lw("  What will it be? [`0B`2] `8:`% ")
			och = g.forestOliviaGetch("BT")
			if och == "B" {
				g.sln("")
				g.lln("  `0You stroke her stump lovingly.")
				g.sln("")
				g.more()
				g.lln("  `2Olivia gives all she has.  And in fact, what she is.")
				g.lln("")
				g.more()
				if g.player.HasFairy {
					g.lln("  `b** `%UH OH! `b**")
					g.lln("")
					g.lln("`2  You feel a buzzing in your pocket.  You smile, your old twig &")
					g.lln("  berries still works.  `0The buzzing grows uncomfortably loud.")
					g.lln("")
					g.logLine("  `0" + g.player.Name + "'s `2fairy escapes.")
					g.player.HasFairy = false
					g.player.OliviaCount = 0
					g.more()
					g.sln("")
					g.lln("`%  A tiny light shoots out of your trousers!")
					g.sln("")
					g.lln("  `2Olivia blinks in surprise.")
					g.sln("")
					g.more()
					g.sln("  You wince - Your fairy is loose!  You attempt to suavely cover up")
					g.sln("  your 'accident'.")
					g.sln("")
					g.lln("  `%\"Wow baby!  Uh, that was great, gotta go!\"")
					g.sln("")
					g.more()
					g.lln("  `2You leave an angry Olivia while you search around her cave for the")
					g.sln("  little creature.  Looks like it's gone for good.")
					g.sln("")
					g.lln("  `4FOR BLOWING IT AND MAKING OLIVIA MAD, YOU GET NO EXPERIENCE.")
					g.sln("")
					g.more()
					return
				}
				g.logLine("  `0" + g.player.Name + " `2walks out of the forest acting chipper!")
				g.lln("  `2A very short time later...")
				g.sln("")
				g.lln("  `%\"That was wonderfull!  Can I put you in my napsack and take you")
				g.lln("  with me?\" `2you ask hopefully.")
				g.sln("")
				g.lln("  `0\"No, I don't think so.  My place is here, I would just be a freak")
				g.lln("  and an oddity anywhere else.\" `2she answers sadly.")
				g.sln("")
				g.more()
				g.lln("  `%\"But.. that's what you are here too!  Keep your chin up, honey.\" `2you")
				g.sln("  clumsly encourage.")
				g.sln("")
				g.lln("  `2She smiles wanly.")
				g.sln("")
				tmp = g.player.Level * g.player.Level * 25
				g.lln("  `%YOU RECIEVE " + prettyInt(tmp) + " EXPERIENCE.")
				g.player.Exp += tmp
				if g.player.Exp > 2000000000 {
					g.player.Exp = 2000000000
				}
				g.player.OliviaCount += 1
				g.sln("")
				g.more()
				g.tournamentCheck()
				return
			}
		}
		if g.Settings.CleanMode {
			g.lln("  `0\"I'll never find a husband!\" `2Olivia cries for no apparent reason.")
			g.lln("")
		}
		g.sln("")
		g.lln("  `%\"Lets not get a ..head of ourselves.  You need someone who you")
		g.lln("  can relate too.  Someone who... well also is a severed head?\" `2you")
		g.sln("  let her down as easy as you can.")
		g.sln("")
		g.lln("  `0\"But I'll NEVER find him!!\" `2Olivia wails.")
		g.sln("")
		g.more()
		g.lln("  `0\"Oh yes you will!  Someday.\" `2you lie glibly.")
		g.sln("")
		g.lln("  `2Feeling uncofortably, you decide to take off.")
		g.sln("")
		g.more()
		return
	}
	if g.player.OliviaCount == 25 {
		g.lln("  `2Olivia is perched on a rock waiting for you.")
		g.sln("")
		g.lln("  `0\"Hello.  I need to tell you something.  I need to tell you")
		g.sln("  why they beheaded me.\"")
		g.sln("")
		g.more()
		g.lln("  `%\"Fine.  If you need to get this off your che, <cough>, uh, go ahead.\"")
		g.lln("  `2you encourage.")
		g.sln("")
		g.more()
		g.lln("`c                         `%OLIVIAS HISTORY")
		g.lln("`2-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
		g.lln("  `2Olivia was not beheaded publicly.  She was not even beheaded")
		g.sln("  by a decree of the king.")
		g.sln("")
		g.lln("  She was decapitated by a man named `0Earnest Drinklewip`2. ")
		g.sln("")
		g.sln("  Why did he do it?  Olivia tells you because she would not find")
		g.sln("  favor with him.  His reason was she was a witch.  She assures you")
		g.sln("  she is not.")
		g.sln("")
		g.sln("  You don't know what to believe - all you know for sure is this man")
		g.sln("  slew her, and put her head in this cave to rot.  The problem is it")
		g.sln("  didn't. ")
		g.sln("")
		g.more()
		g.lln("  `%\"So why aren't you dead and rotting, Olivia?\" `2you ask, confused.")
		g.sln("")
		g.lln("  `0\"I don't know!  Maybe I'm charmed, maybe I'm cursed.  I don't know.\"")
		g.sln("")
		g.more()
		g.lln("  `%\"I'm going to find 'em!  And find out what happened to the rest of")
		g.sln("  you!  Goodbye for now!\"")
		g.sln("")
		g.player.OliviaCount += 1
		g.more()
		return
	}
	if g.player.OliviaCount > 10 && g.player.OliviaCount < 18 {
		g.lln("  `0Olivia warmly greets you.")
		g.sln("")
	}
	if g.player.OliviaCount > 17 && g.player.OliviaCount < 21 {
		g.lln("  `0Olivia greets you with a butterfly kiss.")
		g.sln("")
	}
	if g.player.OliviaCount > 20 {
		g.lln("  `0Olivia greets you with a head hug.")
		g.lln("")
	}
	if g.player.OliviaCount > 100 {
		g.player.OliviaCount = 100
	}
	choices := "G"
	g.lln("  `2(`0G`2)et inside her head")
	if g.player.Sex == "M" {
		choices += "A"
		g.lln("  `2(`0A`2)sk for a kiss.")
	} else {
		g.lln("  `2(`0C`2)onsole her.")
		g.lln("  `2(`0D`2)o her hair.")
		choices += "CD"
	}
	g.sln("")
	g.lw("  What will it be? [`0G`2] `8:`% ")
	och = g.forestOliviaGetch(choices)
	g.player.OliviaCount += 1
	g.player.Put()
	switch och {
	case "C":
		if g.player.Sex != "F" {
			break
		}
		g.sln("")
		g.lln("`2  You feel sorry for the head, and attempt to solace her.")
		g.lw("  `%\"")
		switch g.random(5) {
		case 0:
			g.lln("Don't feel bad Olivia, I'll help you cope.\" `2you offer.")
		case 1:
			g.lln("We girls will always stick together!\" `2you swear.")
		case 2:
			g.lln("You're lucky, you don't have to watch your weight.\"`2you tell her.")
		case 3:
			g.lln("Why Olivia, I love that uh, hair style you have today!\" `2you compliment.")
		case 4:
			g.lln("I'm your friend always.  Lets do each others hair.\" `2you offer.")
		}
		g.sln("")
		if g.player.OliviaCount < 6 {
			g.lln("  `2She thanks you, and feels better!")
			g.sln("")
			g.lln("  `%HELPING ANOTHER LIFTS YOUR SPIRITS.")
			g.player.HighSpirits = true
		} else {
			g.lln("  `0\"But " + g.player.Name + "`0, I'm not depressed anymore!\" `2she laughs.")
			g.sln("")
		}
		g.sln("")
		g.more()
		return
	case "D":
		if g.player.Sex != "F" {
			break
		}
		g.sln("")
		g.lln("`2  You take your brush to Olivia's tangled locks.")
		g.sln("")
		if g.player.OliviaCount < 6 {
			g.lln("  `0\"Don't touch me!\" `2Olivia screeches.")
			g.sln("")
		} else {
			g.lln("  `2She thanks you, and feels better!")
			g.sln("")
			g.lln("  `%HELPING ANOTHER LIFTS YOUR SPIRITS.")
			g.player.HighSpirits = true
			g.sln("")
		}
		g.more()
		return
	case "A":
		if g.player.Sex != "M" {
			break
		}
		g.lln("`c                        `%LOVE IN YOUR HAND")
		g.lln("`2-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
		if g.random(2) == 0 {
			g.lln("`2  You pick Olivia up by the ears and kiss her soundly.")
		} else {
			g.lln("`2  You carefully wipe the dirt of Olivia's lips and kiss her.")
		}
		// BUG (lord.js): tests player.olivia_acount, which doesn't
		// exist (undefined < 6 is false), so the kiss always works.
		// Kept.
		if false {
			g.lln("  `0\"Don't touch me!\" `2Olivia screeches.")
			g.sln("")
		} else {
			g.sln("")
			g.lln("  `%YOU LEAVE IN HIGH SPIRITS!")
			g.sln("")
			g.player.HighSpirits = true
		}
		g.more()
		return
	case "G":
		g.lln("`c                        `%A FREUDIAN SLIP")
		g.lln("`2-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
		g.sw("  \"")
		g.foreground(15)
		switch g.random(5) {
		case 0:
			g.lln("Penny for your thoughts?\" `2you inquire.")
		case 1:
			g.lln("Watcha thinkin'?\" `2you ask politely.")
		case 2:
			g.lln("Will you tell me about yourself?\" `2you ask Olivia.")
		case 3:
			g.lln("What's on your mind?\" `2you inquire.")
		case 4:
			g.lln("You must be thinking hard, since thats all you can do.\" `2you laugh.")
		}
		g.sln("")
		g.lln("  `2Olivia nudges herself to a better speaking position.")
		g.sln("")
		if g.random(2) == 0 {
			g.lln("  `0\"I was just thinking about my past.  I grew up in a great palace.")
			g.sln("  I was pampered and treated like a queen.  In fact, I would have")
			g.lln("  been queen too - if not for the evil Duke of..`%" + castles[g.whichCastle] + "`0 was")
			g.sln("  it?  Something like that.  And then..\"")
			g.sln("")
			g.more()
		} else {
			g.lln("  `0\"Well..  I was thinking about when I had a body.  A very beautiful")
			g.sln("  one if I do say myself.  Men followed me like flies on honey.  Especially")
			g.lln("  one man.  At the time he was a town crier at `%" + castles[g.whichCastle] + "`0. ")
			g.sln("")
			g.sln("  It is because of him I am like this..You see, he..\"")
			g.sln("")
			g.more()
		}
	}
	g.lln("  `%\"Getting bored here.  Gotta go, see ya,\" `2you break in rudely.")
	g.sln("")
	g.sln("  You travel back to the forest entrance.")
	g.lln("")
	g.more()
}
