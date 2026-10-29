package game

import (
	"strings"
)

// romance is romance(): ask for a romantic line (suggesting one at random)
// and mail it to player to with the mail action code kind.
func (g *Game) romance(prompt string, mopts, fopts []string, to int64, shorttxt, mailtxt, kind string) bool {
	p := g.player
	var line string

	g.sw(prompt)
	if p.Sex == "M" {
		line = mopts[g.random(int64(len(mopts)))]
	} else {
		line = fopts[g.random(int64(len(fopts)))]
	}

	line = cleanStr(g.getstr(0, 0, 53, 1, 15, line, GetStrOpts{}))
	g.sln("")
	g.sln("")

	if len(line) < 5 {
		g.sln(shorttxt)
		return false
	}

	g.foreground(15)
	g.sln("  ** WRITING ROMANTIC MAIL, PLEASE WAIT **")
	line = "`0  \"" + line + "`0\""
	g.mailTo(to, " \n"+
		"  `2Romantic Message From `0"+p.Name+"`2!\n"+
		"`0-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-`2\n"+
		mailtxt+"\n"+
		"  `2\n"+
		line+"\n"+
		"`"+kind+itoa(p.Record))
	return true
}

// composeMail is compose_mail(): write mail (or romantic mail) to another
// player.
func (g *Game) composeMail() {
	p := g.player
	sent := false
	proposeLog := []string{"`0  " + p.Name + " `2has been smitten - with love.", "`0  " + p.Name + " `2is sick - lovesick, that is."}

	g.sln("")
	g.sln("  Who would you like to send mail to?")
	to := g.findPlayer()
	if to == -1 {
		g.foreground(15)
		g.sln("  No matching names found.")
		return
	}
	op := g.playerGet(to)
	if p.Sex != op.Sex && !p.Flirted {
		pronoun := pick(op.Sex == "M", "him", "her")
		ppronoun := pick(p.Sex == "M", "him", "her")
		ppronoun2 := pick(p.Sex == "M", "his", "her")
		uppronoun := pick(p.Sex == "M", "He", "She")
		upronoun := pick(op.Sex == "M", "Him", "Her")
		g.lw("  `2Say something `0Romantic`2 to " + pronoun + "`2? : `2[`5N`2] `0: ")
		ch := strings.ToUpper(g.getkey())
		if ch != "Y" {
			ch = "N"
		}
		g.sln(ch)
		g.sln("")
		if ch == "Y" {
			g.lln("`c                     `%** ROMANTIC MESSAGE **")
			g.lln("`2-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
			g.lln("  `0Your hand trembles as you try to conjure up something to stir")
			g.lln("  your beloved " + op.Name + ".")
			g.sln("")
			g.showLooks(op)
			g.moreNomail()
			g.lln("  `2(`0N`2)ever mind.")
			g.lln("  `2(`0F`2)latter " + upronoun)
			g.lln("  `2(`0A`2)sk For A Kiss")
			g.lln("  `2(`0B`2)uy " + upronoun + " Dinner")
			if !g.Settings.CleanMode {
				g.lln("  `2(`0I`2)nvite " + upronoun + " To Your Room At The Inn")
			}
			if p.Cha > 99 {
				g.lln("  `2(`0P`2)ropose To " + upronoun)
			}
			g.sln("")
			g.lw("  `2Which way would you like to show your love? : ")
			ch = strings.ToUpper(g.getkey())
			if ch == "" || !strings.Contains("FABIP", ch) {
				ch = "N"
			}
			if g.Settings.CleanMode && ch == "I" {
				ch = "N"
			}
			if p.Cha < 100 && ch == "P" {
				ch = "N"
			}
			g.sln(ch)
			g.sln("")
			switch ch {
			case "F":
				sent = g.romance("  Flattering Text? :", []string{
					"I like, um, think you are, uh, nice?",
					"You are just as cute as my dog!",
					"You have like, big uh, thingies.",
					"Your lips are like roses.",
					"May I drink in your beauty?  You're a tankard of joy!",
					"I'd pick up yer hankie anywheres! Hyuck!"},
					[]string{"I like, um, think you are, uh, nice?",
						"You are just as cute as my dog!",
						"I like the way you kill people.",
						"I've been watching you for a while...",
						"Might I let you know that you are handsome?",
						"Ya wanna do something sometime, somewhere?"},
					to, "  Somehow you don't think that will turn "+pronoun+" on.",
					"  `2"+p.Name+" is flirting with you!", "K")
				p.Flirted = sent
			case "A":
				sent = g.romance("  By saying what? :", []string{
					"Gimmie a kiss, woman!",
					"Git over here, you!",
					"Please give me the honor of touching your lips.",
					"I'm an explorer.  Can I explore your mouth?",
					"Kiss me, I use a mouth wash!",
					"Uh, if you kiss me, I'll turn back into a prince."},
					[]string{"Kiss me you big hunk oh man!",
						"Pleeeeeeeeease kiss me!",
						"I'm an adventurer.  Can I explore your mouth?",
						"If you kiss me, our relationship might grow.",
						"Kiss me, lover!  I brush regularly.",
						"Can I teach you a french custom I know?"},
					to, "  Somehow you don't think that will turn "+pronoun+" on.",
					"  `2"+p.Name+" wants you to kiss "+ppronoun+"!", "Y")
				p.Flirted = sent
			case "B":
				sent = g.romance("  And just how? :", []string{
					"Please let me buy you dinner.  It will be good.",
					"Let me take you out - I won't expect you to put out!",
					"No obligations - Just dinner, baby.",
					"You ever eat at the Red Dragon Inn?  It's on me.",
					"Come on sweet thing, I know you wanna eat!",
					"Please?  I'll cook Dragon for ya!"},
					[]string{"I'm not a feminist, but I'll pay for it!",
						"Pleeeeeeeeease let me buy you a warm meal.",
						"I'll cook it myself!  A woman's place is the kitchen!",
						"If you're not busy, I'd REALLY apreciate the company.",
						"I'll give you a real meal - Not crap.",
						"If you clean your plate, I'll even dance for you!"},
					to, "  Somehow you don't think that will turn "+pronoun+" on.",
					"  `2"+p.Name+" wants to treat you to dinner!", "U")
				p.Flirted = sent
			case "I":
				sent = g.romance("  Worded how? :", []string{
					"Um, do you uh, um, think you could, uh.. wanna do it?",
					"Glorious girl.  Please be my valentine.  Tonight.",
					"Come take a ride on the wild stallion, baby!",
					"Make me a man tonight, honey!  Pleeeeeease?!",
					"I know you want my body, I know you think I'm sexy...",
					"I've been waiting to ask you all my life for this."},
					[]string{"I'll pleasure you like no one has ever pleasured a man.",
						"Say yes, honey.  I swear I'll be here in the morning.",
						"Say yes.  Your body says yes, listen to it.",
						"Do it!  I swear I've been tested recently..I think?",
						"You'll never regret a night spent with me, studmuffin.",
						"If you say no, I'll never ask again!"},
					to, "  Somehow you don't think that will turn "+pronoun+" on.",
					"  `2"+p.Name+"`2 wants you to join "+ppronoun+" in a night of\n  `2unbridled passion, in "+ppronoun2+" room at the Inn.", "I")
				p.Flirted = sent
			case "P":
				if p.Cha > 99 {
					g.getState(false)
					st := g.state
					if p.MarriedTo > -1 || st.MarriedToSeth == int64(p.Record) || st.MarriedToViolet == int64(p.Record) {
						g.sln("")
						g.sln("  Er - you sort of ARE married.  That kind of puts a crink in your")
						g.sln("  romantic desires, now doesn't it?")
						g.sln("")
						g.moreNomail()
					} else if op.MarriedTo > -1 || st.MarriedToSeth == int64(op.Record) || st.MarriedToViolet == int64(op.Record) {
						g.sln("")
						g.lln("  `)PROBLEM! `2 You cannot marry a married person!  No way!")
						g.sln("")
						g.moreNomail()
					} else {
						sent = g.romance("  Worded how? :",
							[]string{"Would you do me the honor of being my wife?",
								"I love you, " + op.Name + ".  Marry me.",
								"Be my wife.  Make me the happiest man in the world.",
								"Please - I've watched you for the longest time..",
								"Say yes!  It's such a wonderful easy word!",
								"I cannot live another day without you."},
							[]string{"Would you do me the honor of being my husband?",
								"I love you, " + op.Name + ".  Marry me.",
								"Be mine.  Make me the happiest woman in the world.",
								"Please - I've watched you for the longest time..",
								"Marry me...So I can say \"I gotta man!\"...",
								"I've put my heart on my sleeve.  Don't break it."},
							to, "  Somehow you don't think that will convince "+pronoun+"...",
							"  `2"+p.Name+" has publicly declared "+ppronoun2+" love for\n  `2you.  "+uppronoun+" is asking for your hand in marriage.", "P")
						p.Flirted = sent
						if sent {
							g.logLine(proposeLog[g.random(int64(len(proposeLog)))])
						}
					}
				}
			}
		}
		p.Put()
	}
	g.sln("")
	if !sent {
		g.writeMail(to, false)
	}
}

// haveBaby is have_baby().
func (g *Game) haveBaby() {
	p := g.player
	var op *Player

	g.sln("")
	g.lln("  `0** `%YOU FEEL SICK THIS MORNING! `0**")
	g.sln("")
	g.moreNomail()
	g.lln("  `2You head to the healers in dismay, what can it be?")
	g.sln("")
	g.moreNomail()
	g.sclrscr()
	g.sln("")
	g.sln("")
	g.foreground(15)
	g.sln("  Healers")
	g.foreground(1)
	g.sln("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
	g.foreground(2)
	g.lln("  `0\"Madam!  You are with child, and it's due right now!\"")
	g.foreground(2)
	g.sln("  Nathan, an elderly healer yells.")
	g.sln("")
	g.moreNomail()
	g.sln("  At this moment, you feel something ripping your insides out.  You")
	g.sln("  realize that it's childbirth pain.")
	g.sln("")
	g.moreNomail()
	g.sln("  Then the fun begins.  Three hours of sweat and screaming later, it")
	g.sln("  is over.  Nathan hands you a small bundle with tears in his eyes.")
	g.sln("")
	g.moreNomail()
	g.lw("  `0The baby is `%")
	n := g.random(20) + 1
	if n == 20 {
		g.lln("`4Not breathing`%.")
		g.sln("")
		g.moreNomail()
		g.foreground(2)
		g.sln("  Your pain cannot be expressed in words.  You feel as though you")
		g.sln("  have lost your own heart.  From this day on, everytime you think of")
		g.sln("  this small helpless child who never got a chance to live, you cry.")
		g.sln("")
		g.sln("  Live is precious.  Period.")
		g.sln("")
		if p.MarriedTo == -1 {
			g.logLine("  `0" + p.Name + " `2is grief stricken over a great loss.")
		} else {
			op = g.playerGet(p.MarriedTo)
			g.logLine("  `0" + p.Name + " `2and `0" + op.Name + " `2are grief stricken.")
			g.mailTo(p.MarriedTo, "  `0TRAGEDY STRIKE.\n"+
				"  `2-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-\n"+
				"  `0Your wife had a miscarriage today.  You feel you\n"+
				"  `0should console her.\n"+
				"`|")
		}
	} else if n < 10 {
		g.lln("`%A beautiful baby boy!")
		g.sln("")
		g.moreNomail()
		p.Kids++
		g.foreground(2)
		g.sln("  What would you like to name this child?")
		g.sw("  Name :")
		nm := strings.TrimSpace(cleanStr(g.getstr(0, 0, 40, 1, 15, "", GetStrOpts{})))
		mstr := "  `0" + p.Name + " `2gives birth to a boy!  His name is `%" + nm + "`2."
		if p.MarriedTo > -1 {
			op = g.playerGet(p.MarriedTo)
			mstr += "\n  `0" + op.Name + " `2is extremely proud."
		}
		g.logLine(mstr)
		g.sln("")
		g.sln("")
		if p.MarriedTo == -1 {
			g.lln("  `2Oddly, he seems to  look a little like `%The Bard`2...")
		} else {
			g.lln("  `2He looks just like `0" + op.Name + "`2.")
			mstr = "  `0JOYOUS OCCASION!\n"
			mstr += "  `%-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-\n"
			mstr += "  `2Your wife had a baby `%boy `2today!  Her face glows\n"
			mstr += "  `2like only a mothers can.  This baby is as much yours,\n"
			mstr += "  `2so you aspire to teach little `%" + nm + "`2 much.\n"
			mstr += "`K\n"
			mstr += "`|"
			g.mailTo(p.MarriedTo, mstr)
		}
		g.sln("")
	} else {
		g.lln("`%A beautiful baby girl!")
		g.sln("")
		g.moreNomail()
		p.Kids++
		g.foreground(2)
		g.sln("  What would you like to name this child?")
		g.lw("  `2Name :")
		nm := strings.TrimSpace(cleanStr(g.getstr(0, 0, 40, 1, 15, "", GetStrOpts{})))
		mstr := "  `0" + p.Name + " `2gives birth to a girl!  Her name is `%" + nm + "`2."
		if p.MarriedTo > -1 {
			op = g.playerGet(p.MarriedTo)
			mstr += "\n  `0" + op.Name + " `2is pleased."
		}
		g.logLine(mstr)
		if p.MarriedTo > -1 {
			mstr = "  `0JOYOUS OCCASION!\n"
			mstr += "  `%-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-\n"
			mstr += "  `2Your wife had a baby `)girl `2today!  Her face glows\n"
			mstr += "  `2like only a mothers can.  This baby is as female, but\n"
			mstr += "  `2you hope to make a warrior out of `%" + nm + "`2 yet.\n"
			mstr += "`K\n"
			mstr += "`|"
			g.mailTo(p.MarriedTo, mstr)
		}
		g.sln("")
		g.sln("")
		g.foreground(2)
		g.sln("  You can tell she will be a heartbreaker.")
		g.sln("")
	}
}

// checkMarriage is check_marriage(): the morning with your spouse.
func (g *Game) checkMarriage() {
	p := g.player
	st := g.state
	goodMood := []string{"Happy", "Wonderful", "Joyous", "Excited", "Loving",
		"Lucky", "Providential", "Felicitous", "Glad", "Light-hearted"}
	badMood := []string{"Angry", "Resentful", "Disgusted", "Revolted", "Repulsed",
		"Sickened", "Nauseated", "Betrayed", "Let Down", "Cheated On"}

	if st.MarriedToSeth == int64(p.Record) {
		return
	}
	if st.MarriedToViolet == int64(p.Record) {
		return
	}
	if p.MarriedTo < 0 {
		return
	}
	op := g.playerGet(p.MarriedTo)
	if op.Name == "X" {
		g.sln("")
		g.lln("  `2Your lover has been declared MISSING!")
		g.sln("")
		g.lln("  `%YOU FEEL HEARTBROKEN!  (CHARM DROPS 50%)")
		g.sln("")
		p.MarriedTo = -1
		p.Cha = p.Cha / 2
		g.moreNomail()
		return
	}
	g.lln("  `2You notice `0" + op.Name + "`2 is asleep in your bed.")
	g.sln("")
	g.lln("  `%You feel...")
	g.sln("")
	good := goodMood[g.random(int64(len(goodMood)))]
	bad := badMood[g.random(int64(len(badMood)))]
	g.lln("  `2(`01`2) `%" + good)
	g.lln("  `2(`02`2) `%" + bad)
	g.sln("")
	g.lw("  `0Choose.  [`21`0] :`%")
	ch := g.getkey()
	if ch != "2" {
		ch = "1"
	}
	g.sln(ch)
	g.sln("")
	if ch == "1" {
		g.lln("  `2You smile at `%" + op.Name + "`2.  Life is good.")
		g.sln("")
	} else {
		g.lln("  `2You scowl at the snoring figure in your bed.  Suddenly")
		g.sln("  you feel your rage boiling over!")
		g.sln("")
		g.moreNomail()
		g.lln("  You kick " + op.Name + " `2into consciousness.")
		g.sln("")
		g.moreNomail()
		if p.Sex == "M" {
			g.lln("  The woman awakes in confusion.  `0\"Get out, woman.\" `2you tell her")
			g.sln("  harshly.")
			g.sln("")
			g.sln("  She slaps you in the face and walks out.")
			g.sln("")
			g.lln("  `%YOU FEEL FED UP WITH WOMEN!  (CHARM `)DROPS`% TO 50)")
		} else {
			g.lln("  The pig awakes in confusion.  `0\"Get out of my house.\" `2you tell him")
			g.sln("  coldly.")
			g.sln("")
			g.sln("  A short while later, he stands outside your door, naked, gawking in")
			g.sln("  disbelief. ")
			g.sln("")
			g.lln("  `%YOU FEEL FED UP WITH MEN!  (CHARM `)DROPS`% 50%)")
		}
		g.sln("")
		g.logLine("`0  " + p.Name + " `2has `)DIVORCED `0" + op.Name + "`2!")
		p.Divorced = true
		p.Cha = p.Cha / 2
		op.Lock()
		op.ReLoad()
		op.MarriedTo = -1
		op.Put()
		op.UnLock()
		p.MarriedTo = -1
		mstr := "`*\n"
		mstr += "  `)BAD NEWS!\n"
		mstr += "`2-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-\n"
		mstr += "  `0" + p.Name + " `2has `4DIVORCED `2you!\n"
		mstr += "\n"
		mstr += "  `0(CHARM DROPS 50%)"
		g.mailTo(int64(op.Record), mstr)
	}
}

// wakeUp is wake_up(): the start of a new day for the player.
func (g *Game) wakeUp() {
	p := g.player
	vds := []string{"AIDS!", "HERPES!", "SYPHILIS!", "GONORRHEA!", "the HIV virus!"}
	sethSupport := []string{"  `0I love you, and still do not regret what we did. Nor will I ever.\n",
		"  `0You must be strong - Do not let the gossip hurt you.\n",
		"  `0Teach our young to be strong - To be like her mother.\n",
		"  `0I am sending what I can - Not because of the law, I want to do it.\n"}
	salutations := []string{"  `0My dearest ", "  `0My lovely ", "  `0My goddess ", "  `0Dear "}
	origDead := p.Dead || (p.HP < 1)

	g.getState(false)
	st := g.state
	p.HP = p.HPMax
	p.SeenMaster = false
	p.TransferredGold = 0
	p.Flirted = false
	p.ForestFights = 15
	p.SeenBard = false
	p.LastReincarnated = st.Days
	p.Time = st.Days
	p.PvPFights = g.Settings.PvPFightsPerDay
	if g.random(5)+1 == 3 {
		p.Weird = true
	}
	p.SeenDragon = false
	p.SeenViolet = false
	p.Dead = false
	p.LeftBank = false
	if g.Settings.ForestFights > 32000 {
		g.Settings.ForestFights = 32000
	}
	if g.Settings.ForestFights > 0 && g.Settings.ForestFights <= 32000 {
		p.ForestFights = g.Settings.ForestFights
	}
	// do some other player maint ..  but not if theyve been gone for too long
	if p.Gone < g.Settings.ResDays {
		if g.Settings.OldSkillPoints {
			p.LevelW = p.SkillW / 5
		} else {
			p.LevelW = p.SkillW / 4
		}
		p.LevelM = p.SkillM
		if g.Settings.OldSkillPoints {
			p.LevelT = p.SkillT / 5
		} else {
			p.LevelT = p.SkillT / 4
		}
		if p.Clss == 1 {
			p.LevelW++
		}
		if p.Clss == 2 {
			p.LevelM++
		}
		if p.Clss == 3 {
			p.LevelT++
		}
		if p.Clss == 1 || p.Clss == 3 {
			g.foreground(15)
			if g.Settings.OldSkillPoints {
				g.sln("                 (Remember, 5 Skill Points = 1 Use Point)")
			} else {
				g.sln("                 (Remember, 4 Skill Points = 1 Use Point)")
			}
			g.sln("")
		}
		g.foreground(10)
		switch p.Clss {
		case 1:
			g.lln("  For being a `%Death Knight`0, you get an extra `%Death Knight`0 use point!")
		case 2:
			g.lln("  For being a `5Mystical Skills`0 student, you get an extra use point for today!")
		case 3:
			g.lln("  For being a Pilferer, you get an extra `%Thieving Skills`0 use point!")
		}

		g.sln("")
		if p.Kids > 0 {
			g.foreground(2)
			g.sw("  You are inspired by your child")
			if p.Kids == 1 {
				g.sln(".")
			} else {
				g.sln("ren.")
			}
			g.sln("")
			g.lln("  `%YOU RECEIVE `0" + prettyInt(p.Kids) + " `%EXTRA FOREST FIGHTS FOR TODAY.")
			g.sln("")
			p.ForestFights += p.Kids
			if p.ForestFights > 32000 {
				p.ForestFights = 32000
			}
		}
		if p.Bank >= 2000000000 {
			g.lln("  `$No interest earned today.  The coins won't fit in the bank!")
			p.Bank = 2000000000
		} else {
			p.Bank += p.Bank / 10
			if p.Bank > 2000000000 {
				p.Bank = 2000000000
				g.lln("  `0Wow! You have a lot of gold!")
			}
		}
	} else {
		p.Gone = 0
		g.lln("  `0The townspeople thought you had abandoned them!")
		g.sln("")
	}
	if g.Settings.SleepDragon && !origDead && (p.Level > 12 || p.Exp > 15000000) {
		g.lln("`c  `2You awake to odd movement in the middle of the night...")
		g.sln("")
		g.lln("  It feels like you're on a boat, a slow rolling side to side lulling")
		g.lln("  you back to sleep...")
		g.sln("")
		g.moreNomail()
		g.lln("  Suddenly, you bolt awake!  You went to sleep " + pick(p.Inn, "at the inn", "in the fields") + "!")
		g.sln("")
		g.moreNomail()
		g.lln("  The `4Red Dragon `2spits you out onto the ground in front of it.")
		g.lln("  Unprepared, you fumble for your " + p.Weapon + "`2.")
		g.sln("")
		g.fightDragon(true)
		g.moreNomail()
	}
	if p.Dead || p.HP < 1 {
		// BUG (lord.js): sets player.online, which isn't a Player_Def
		// field (on_now is), so only the put() has an effect.  Kept.
		p.Put()
		return
	}
	if p.Sex == "F" && p.Kids > 3 && p.Laid > 13 {
		if p.MarriedTo < 0 && st.MarriedToSeth != int64(p.Record) {
			if g.random(5) == 3 {
				mstr := "  `0Seth Able `2sent you this...\n"
				mstr += "`0-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-`2\n"
				mstr += salutations[g.random(int64(len(salutations)))]
				mstr += p.Name + ",\n"
				mstr += sethSupport[g.random(int64(len(sethSupport)))]
				mstr += "  `0\n"
				mstr += "  `0Visit me today, please.\n"
				mstr += " \n"
				mstr += "  `%BANK NOTICE:\n"
				mstr += "`0-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-`2\n"
				mstr += "  `0Seth Able `2has transfered `%" + prettyInt(p.Kids*1000*p.Level) + " `2gold\n"
				mstr += "  to your account.\n"
				mstr += "`b" + itoa(p.Kids*1000*p.Level) + "\n"
				g.mailTo(int64(p.Record), mstr)
			}
		}
	}
	if g.random(3)+1 > 1 {
		g.lln("  `2You are in `0high `2spirits today.")
		p.HighSpirits = true
		g.sln("")
	} else {
		g.lln("  `2You are in `0low `2spirits today.")
		g.sln("")
		p.HighSpirits = false
	}
	if p.Horse {
		g.lln("  `)** `%You ready your horse `)**`2")
		g.sln("")
		p.ForestFights += p.ForestFights / 4
		if p.ForestFights > 32000 {
			p.ForestFights = 32000
		}
	}
	if p.Inn {
		g.sln("  You wake up early, sheathe your weapon, and head down to the bar.")
	} else {
		g.sw("  You wake up early, strap your ")
		g.lw(p.Weapon)
		g.foreground(2)
		g.sln(" to your back, and")
		g.sln("  head out to the Town Square, seeking adventure, fame, and honor.")
		g.sln("")
	}
	if p.Sex == "F" && g.random(34)+1 == 11 && p.Laid*5 > p.Kids {
		g.haveBaby()
	}
	g.checkMarriage()
	if p.Laid > 50 {
		if g.random(23) == 13 {
			g.sln("")
			g.lln("  `0** `%YOU FEEL SICK THIS MORNING! `0**")
			g.sln("")
			p.HP = 1
			g.logLine("`0  " + p.Name + " `2went to the healers to get cured of `4VD`2!")
			g.moreNomail()
			g.lln("  `2You head to the healers in dismay, what can it be?")
			g.sln("")
			g.moreNomail()
			g.sclrscr()
			g.sln("")
			g.sln("")
			g.foreground(15)
			g.sln("  Healers")
			g.foreground(1)
			g.sln("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
			g.foreground(2)
			g.lln("  `0\"My God!\" `2Nathan, the lead healer screams.")
			g.sln("")
			g.moreNomail()
			g.lw("`%  You are diagnosed with ")
			g.sln(vds[g.random(int64(len(vds)))])
			g.sln("")
			g.lln("  `2You are fairly sure you didn't get it off a toilet seat.")
			g.sln("")
			g.moreNomail()
			g.lln("  `2Luckily, Nathan has a potion that can cure it.  This time.")
			g.sln("  However, you still feel quite weak from the experience.")
			g.sln("")
		}
	}
}

// getARoom is get_a_room(): rent a room at the inn (and leave the game).
func (g *Game) getARoom() bool {
	p := g.player

	g.sln("")
	g.sln("")
	g.foreground(2)
	g.sln("  The bartender approaches you at the mention of a room.")
	g.foreground(5)
	g.sln("  \"You want a room, eh?  That'll be " + prettyInt(400*p.Level) + " gold!\"")
	g.lw("  `2Do you agree? [`0Y`2] : `%")
	ch := strings.ToUpper(g.getkey())
	if ch != "N" {
		ch = "Y"
	}
	g.sln(ch)
	if ch != "Y" {
		g.sln("")
		g.sln("  The bartender grunts, then walks away uninterested.")
		return false
	}

	g.sln("")
	g.sln("")
	if p.Cha > 99 {
		if p.Sex == "M" {
			g.sln("  \"You seem like a nice guy, and I hear you're tough, so tell you what,")
			g.sln("  I'll just give you the room for free....\"")
		} else {
			g.sln("  \"Hey good lookin'!  I'll be glad to give you a freebie...")
			g.sln("  A free room I mean, of course.  Har!\"")
		}
		g.sln("")
		g.moreNomail()
	} else {
		if p.Gold < 400*p.Level {
			g.foreground(4)
			g.sln("  \"Hey!  You stupid fool! You don't have that much gold!\"")
			return false
		}
		if p.Cha < 100 {
			p.Gold -= 400 * p.Level
		}
	}

	p.Inn = true
	g.sclrscr()
	g.sln("")
	g.sln("")
	g.foreground(15)
	g.sln("  The Room At The Inn")
	g.foreground(10)
	g.sln("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
	g.foreground(2)
	g.sln("  You are escorted to a small but cozy room in the Inn.")
	g.sln("  You relax on the soft bed, and soon fall asleep, still")
	g.sln("  wearing your armour.")
	g.sln("")
	g.sln("")
	p.OnNow = false
	p.Put()
	g.exit(0)

	return true
}
