package game

// fightDragonStory is fight_dragon's nested story(): the epilogue for the
// player's class.
func (g *Game) fightDragonStory() {
	player := g.player
	if player.Clss == 1 {
		g.lln("`c                   `%EPILOGUE `2- `0The Warrior's Ending")
		g.lln("`0-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-==-=-=-=-=-=-=-=-=-=-=-")
		g.lln("  `2After your bloody duel with the huge Dragon, your first inpulse is to rip")
		g.lln("  its head off and bring it town.  Careful thought reveals it is much to")
		g.lln("  big for your horse, so that plan is moot.  Your second notion is bring back")
		g.lln("  the childrens bones.  Bags and bags of them for proper burial, but you")
		g.lln("  realize this would only cause the town's inhabitants `0MORE`2 pain.  You")
		g.lln("  finally decide on the Dragon's heart.  After adding ten years to your")
		g.lln("  sword's life, you finally chip off enough scales to wallow in the huge")
		g.lln("  beast's insides.")
		g.sln("")
		g.lln("  When you are finished, and fit the still heart in a gunny sack you brought,")
		g.lln("  (who would have thought this would be its use?) you make your way back to")
		g.lln("  town.  As you share your story to a crowd of excited onlookers, this crowd")
		g.lln("  becomes a gathering, and this gathering becomes an assemblage, and this")
		g.lln("  assemblage becomes a multitude!")
		g.sln("")
		g.lln("  This multitude nearly becomes a mob, but thinking quick, you make a")
		g.lln("  speech.")
		g.sln("")
		g.more()
		g.lln("  `0\"PEOPLE!\" `2your voice booms.  `0\"It is true I have ridden this town of")
		g.lln("  its curse, the `4Red Dragon`0.  And this is his heart.\"")
		g.sln("")
		g.lln("  `2You dump the bloody object onto the ground.  From the back, Barak's")
		g.lln("  voice is heard.  `5\"How do we know where you got that thing?  It looks")
		g.lln("  like you skinned a sheep!\"  `2A flicker of annoyance crosses your face,")
		g.lln("  but you force a smile.  `0\"Why Barak, would you doubt me?  A LEVEL 12")
		g.lln("  warrior?  If I am not mistaken, you are quite a bit lower, still at level")
		g.lln("  two, eh?\"")
		g.sln("")
		g.lln("  `2Barak gives you no more trouble, and you are declared a hero by all.")
		if player.Sex == "M" {
			g.lln("  `#Violet `2tops off the evening by giving you a kiss on the cheek, and a ")
		} else {
			g.lln("  `%Seth `2tops off the evening by giving you a kiss on the cheek, and a ")
		}
		g.lln("  whisper of things to come later that night makes even you almost")
		g.lln("  blush.  Almost.")
		g.sln("")
		g.moreNomail()
	}
	if player.Clss == 2 {
		g.lln("`c            `%EPILOGUE `2- `#The Mystical Skills User Ending")
		g.lln("`0-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-==-=-=-=-=-=-=-=-=-=-=-")
		g.lln("`2  Still shaking from the battle, you decide it is time to return to town")
		g.lln("  and share the good tidings.  You close your eyes and concentrate.")
		g.lln("  `0\"There is no place like Town\"`2...A when you open your eyes, you are")
		g.lln("  under a cow in a farm outside - Not far from Abduls Armour.  You trek the")
		g.lln("  distance to the Armoury, cursing as you go.  Wizards were not made for this")
		g.lln("  kind of hardship you tell yourself.")
		g.sln("")
		g.lln("  Miss Abdul is estatic when you tell her about your escapades.  She agrees")
		g.lln("  to accompany you to the Inn.  She lends you a horse when you tell her of how")
		g.lln("  your feet ache.  ")
		g.sln("")
		g.lln("  `2When you enter the smokey bar with a loud clatter, merry makers stop their")
		g.lln("  carousing, people hunched over their meals stop chewing and Seth Able stops")
		g.lln("  playing in mid-strum.")
		g.sln("")
		g.more()
		g.lln("  `0\"People!  I have slain the beast!\"`2 you shout triumphantly. A voice is")
		g.lln("  heard from the back. `0\"What beast?  A large rat or something?\"`2  You")
		g.lln("  scowl.  It is Barak.  You nonchalantly make a gesture with one hand. A")
		g.lln("  few moments later, Barak stands up suprised.  He looks wildly around, then")
		g.lln("  makes a bolt for the door.")
		g.sln("")
		g.lln("  `0\"What happened to him?\" `2Miss Adbul asks you in puzzlement.")
		g.sln("")
		g.lln("  `0\"Nature called.\" `2you smile.  Your face becomes serious as you address")
		g.lln("  the bar.  `0\"The `4Red Dragon`0 is no more.\"`2 `#Violet `2sets her serving")
		g.lln("  tray down to hear better. `0\"We cannot bring back those dead, but I have")
		g.lln("  stopped this from happening again.\" ")
		g.sln("")
		g.lln("  `0\"And he did it with Armour from Abduls Armour!\"`2 Miss Abdul adds.")
		g.sln("")
		g.lln("  `0\"Er, thank you.  Anyway, make sure you put something in the Daily")
		g.lln("  happenings about this!\"`2  The crowd gives you a standing ovation.")
		g.sln("")
		g.more()
	}
	if player.Clss == 3 {
		g.lln("`c                      `%EPILOGUE `2- `9The Thief's Ending")
		g.lln("`0-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-==-=-=-=-=-=-=-=-=-=-=-")
		g.lln("`2  You breathe a sigh of relief, retrieve your daggers and carefully clean")
		g.lln("  them.  Although you realized some may think it cold hearted if they ever")
		g.lln("  found out, you pick through the childrens bones, picking up a gold piece")
		g.lln("  here, a silver there.  Afterall, these children didn't need it to buy a")
		g.lln("  meal anymore... They WERE a meal!  You smile at your own dry wit and realize")
		g.lln("  you had better get to town and share the news.")
		g.sln("")
		g.lln("  The hike to town is long, but you are used to it, and rather enjoy the")
		g.lln("  peace it brings.  You find the town deserted.  You enter the Inn, hoping")
		if player.Sex == "M" {
			g.lln("  to find some clue, but all you find is `#Violet`2.")
			g.sln("")
			g.lln("  `0\"Via, Where has everyone gone?  Have they finally given up hope of ever")
			g.lln("  stopping the `4Red Dragon`0, and have gone to seek a new lifestyle?")
			g.lln("  ")
			g.lln("`2  There is a pause, then she responds.")
			g.sln("")
			g.more()
			g.lln("  `0\"No, they are having a feast at Turgons Place.  Naturally SOMEONE had")
			g.lln("  to stay here, and OF COURSE it would be me...I never get to have any fun!\"")
			g.sln("")
			g.lln("  `2You give her a wink.  `0\"That just isn't so...Remember last night?\"")
			g.lln("  `2She giggles, and after a few more naughty sayings you have her cheeks as")
			g.lln("  deep a scarlet as a rose.")
		} else {
			g.lln("  to find some clue, but all you find is `%The Bard`2.")
			g.sln("")
			g.lln("  `0\"Seth, Where has everyone gone?  Have they finally given up hope of ever")
			g.lln("  stopping the `4Red Dragon`0, and have gone to seek a new lifestyle?")
			g.lln("  ")
			g.lln("`2  There is a pause, then he responds.")
			g.sln("")
			g.more()
			g.lln("  `0\"No, they are having a feast at Turgons Place.  I didn't feel like ")
			g.lln("  going to a party.  At least not by myself.  I'm just not the party type.\"")
			g.sln("")
			g.lln("  `2You give him a wink.  `0\"That just isn't so...Remember last night?\"")
			g.lln("  `2He laughs, and after a few more naughty sayings you have him ready to")
			g.lln("  take a cold shower.  (But you do something else instead) ;> ")
		}
		g.sln("")
		g.lln("`2  When you finally meet up at Turgon's, you share your story.  You can't help")
		g.lln("  but be pleased at seeing so many faces in awe over your doings, so you")
		g.lln("  'spice' up the story in a few places... As you are finishing, a woman in the")
		g.lln("  back cries out.")
		g.sln("")
		g.lln("  `0\"That Thief is wearing the ring I gave my Ellie the last birthday before")
		g.lln("  she disappeared!\"`2")
		g.sln("")
		g.lln("  You decide now would be the perfect time to make your departure.")
		g.sln("")
		g.more()
		g.sln("")
		g.more()
	}
}

// fightDragon is fight_dragon(): the Red Dragon encounter.
func (g *Game) fightDragon(cantRun bool) {
	player := g.player
	var dragon *Player

	// TODO: This could be rolled into lenemy.bin easily enough...
	var rec *monsterRec
	if fn := g.gameOrExec("dragon.bin"); fn != "" && !g.Settings.NoLEnemyBin {
		r, _, err := readRecordAt[monsterRec](fn, 0)
		if err != nil {
			panic(err)
		}
		rec = r
	}
	if rec != nil {
		dragon = &Player{
			Name:     rec.Name,
			Str:      rec.Str,
			Gold:     rec.Gold,
			Weapon:   rec.Weapon,
			Exp:      rec.Exp,
			HP:       rec.HP,
			Death:    rec.Death,
			IsDragon: true,
			Record:   -1,
		}
	} else {
		dragon = &Player{
			Name:     "`4The Red Dragon`2",
			Str:      2000,
			Gold:     1000,
			Weapon:   "`)Scorching Flame`2",
			Exp:      1000,
			HP:       15000,
			Death:    "  The earth shakes as the mighty beast falls.",
			IsDragon: true,
			Record:   -1,
		}
	}
	player.SeenDragon = true
	g.beefUp(dragon)
	g.sln("")
	g.sln("")
	g.lln("`%  **`4DRAGON ENCOUNTER`%**")
	g.sln("")
	g.lln("  `2The Red Dragon approaches.")
	g.sln("")
	g.battle(dragon, false, cantRun)
	g.foreground(2)
	if player.Dead {
		player.HP = 0
		player.Dead = true
		player.Gold = 0
		player.Put()
		g.sln("")
		g.sclrscr()
		g.sln("")
		g.sln("")
		g.sln("  The Dragon pauses to look at you, then snorts in a Dragon laugh, and")
		g.sln("  delicately rips your head off, with the finesse only a Dragon well practiced")
		g.sln("  in the art could do.  ")
		g.sln("")
		g.logLine("  `2The `4Red Dragon `2has killed `5" + player.Name + "`2!")
		g.more()
	}
	if dragon.HP < 1 {
		g.sln("")
		g.sln("  You have defeated The Red Dragon!")
		g.sln("")
		g.more()
		g.sclrscr()
		g.sln("")
		g.sln("")
		g.sln("  You have defeated the Dragon, and saved the town.  Your stomach churns")
		g.sln("  at the site of stacks of clean white bones - Bones of small children.")
		g.sln("")
		g.sln("  THANKS TO YOU, THE HORROR HAS ENDED!")
		g.sln("")
		g.logLine("`.  `%" + player.Name + " `2has slain the `4Red Dragon`2 and become a hero.")
		g.more()
		g.fightDragonStory()
		g.foreground(15)
		g.sln("             Thanks for being tough enough to win the game, ")
		//sln('                 and thank your sysop for registering!');
		g.sln("")
		g.foreground(10)
		g.sln("                              -Seth Able")
		g.sln("")
		g.more()
		player.Level = 1
		player.HPMax = 20
		player.HP = player.HPMax
		player.WeaponNum = 1
		player.Weapon = "Stick"
		player.Gold = 500
		player.Bank = 0
		player.Def = 1
		player.Str = 10
		player.Gem = 10
		player.Arm = "Coat"
		player.ArmNum = 1
		player.Dead = false
		player.Inn = false
		player.Exp = 10
		if g.Settings.ForestFights > 32000 {
			g.Settings.ForestFights = 32000
		}
		if g.Settings.ForestFights+player.Kids > 32000 {
			player.ForestFights = 32000
		} else {
			player.ForestFights = g.Settings.ForestFights + player.Kids
		}
		player.PvPFights = g.Settings.PvPFightsPerDay
		player.Flirted = false
		player.HighSpirits = true
		player.DragKills++
		player.Put()
		g.getState(true)
		g.state.LatestHero = player.Name
		g.tournamentCheck()
		if g.Settings.WinDeeds > 0 && player.DragKills >= g.Settings.WinDeeds {
			g.state.WonBy = int64(player.Record)
			g.putState()
			g.lln("`c  `%** YOUR QUEST IS OVER **")
			g.sln("")
			g.lln("  `2You must indeed be the chosen one.  The ancient magic that")
			g.lln("  kept the `4dragon `2alive is now truly no more.")
			g.sln("")
			g.moreNomail()
			g.lln("  `0Now begone, blessed among warriors - Your fight is over.")
			player.OnNow = false
			player.Put()
			g.exit(0)
		}
		g.putState()
		g.lln("`c  `%                ** YOUR QUEST IS NOT OVER **")
		g.sln("")
		g.lln("  `2You are a hero.  Bards will sing of your deeds, but that doesn't")
		g.lln("  mean your life doesn't go on.  ")
		g.sln("")
		g.lln("  `%YOUR CHARACTER WILL NOW BE RESET.  `2But you will keep a few things")
		g.lln("  you have earned.  Like the following.")
		g.sln("")
		g.lln("  `%ALL SPECIAL SKILLS.")
		g.lln("  CHARM.")
		g.lln("  A FEW OTHER THINGS.")
		g.sln("")
		g.more()
		g.lln("`c  `%YOU FEEL STRANGE.")
		g.sln("")
		g.lln("  `2Apparently, you have been sleeping.  You dust yourself off, and")
		g.lln("  regain your bearings.  You feel like a new person!")
		g.sln("")
		g.more()
	}
}
