package game

import (
	"sort"
)

// tournamentOver is tournament_over(): a days-mode tournament has ended;
// crown the winner and leave.
func (g *Game) tournamentOver() {
	var a []*Player

	for _, p := range g.allPlayers() {
		if p.Name != "X" {
			a = append(a, p)
		}
	}

	// winstats = ['exp','drag_kills','pvp_kills','level','laid']
	// BUG (lord.js): 'pvp_kills' isn't a player property (it's `pvp`), so
	// the player-kills stat compares undefined to undefined and falls
	// through to the tie-breakers.  Kept.
	sk := func(p *Player) (int64, bool) {
		switch g.Tournament.WinStat {
		case 0:
			return p.Exp, true
		case 1:
			return p.DragKills, true
		case 3:
			return p.Level, true
		case 4:
			return p.Laid, true
		}
		return 0, false
	}

	sort.SliceStable(a, func(i, j int) bool {
		x, y := a[i], a[j]
		if vx, ok := sk(x); ok {
			vy, _ := sk(y)
			if vx != vy {
				return vy-vx < 0
			}
		}
		if y.Exp != x.Exp {
			return y.Exp-x.Exp < 0
		}
		if y.DragKills != x.DragKills {
			return y.DragKills-x.DragKills < 0
		}
		if y.PvP != x.PvP {
			return y.PvP-x.PvP < 0
		}
		if y.Level != x.Level {
			return y.Level-x.Level < 0
		}
		if y.Laid != x.Laid {
			return y.Laid-x.Laid < 0
		}
		// BUG (lord.js): compares b.gem with a.gen (undefined), which is
		// always unequal; the result is the same as comparing gems.
		return y.Gem-x.Gem < 0
	})

	g.sclrscr()
	g.sln("")
	g.sln("")
	g.foreground(10)
	// DIFF: lord.js throws on a[0] when there are no players; this can't
	// happen in practice (the caller is one), but guard anyway.
	name := ""
	rec := int64(-1)
	if len(a) > 0 {
		name = a[0].Name
		rec = int64(a[0].Record)
	}
	g.lln("  The tournament is now over! The winner is: " + name)
	g.sln("")
	g.sln("  Now taking you back to the BBS..")
	g.sln("")
	g.lln("  `2<`0MORE`2>")
	g.getkey()
	g.state.WonBy = rec
	g.state.LatestHero = name
	g.putState()
	if g.player != nil {
		g.player.OnNow = false
		g.player.Put()
	}
	g.exit(0)
}

// killmail is killmail(): throw away a player's waiting mail.
func (g *Game) killmail(to int) {
	fileRemove(g.gamedir("mail" + itoa(to) + ".lrd"))
}

// divorceSeth is divorce_seth().
func (g *Game) divorceSeth() {
	g.state.MarriedToSeth = -1
}

// divorceViolet is divorce_violet().
func (g *Game) divorceViolet() {
	g.state.MarriedToViolet = -1
}

// violetMarriage is violet_marriage(): the daily news of Violet's husband.
func (g *Game) violetMarriage() {
	mstr := "\n`%  Latest news about `#Violet`%, your wife.\n`0-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-\n"
	husband := g.playerGet(g.state.MarriedToViolet)

	if g.random(5) < 1 {
		g.divorceViolet()
		husband.Cha = husband.Cha / 2
		husband.Put()
		switch g.random(3) {
		case 0:
			g.logLine("`2  `#Violet`2 has `%DIVORCED `0" + husband.Name + "`2 for cheating on her with\n" +
				"`2  `4Grizelda!  `#Violet`2 got her barmaid job back!")
			mstr += "  `2Violet has divorced you because Grizelda said you kissed her!  You\n"
			mstr += "  curse Grizelda!  What will people think?\n"
		case 1:
			g.logLine("`2  `#Violet`2 has `%LEFT `0" + husband.Name + "`2 so she could get her old job\n" +
				"`2  back at the Inn!  `0" + husband.Name + " `2is heartbroken.")
			mstr += "  `2You hunt around the house for Violet, and all you find is a note!  She\n"
			mstr += "  left you!  She could not resist the temptation of working at the Inn\n"
			mstr += "  once again.  You try to control your convulsive sobs.\n"
		case 2:
			g.logLine("`2  `0" + husband.Name + " `2has `%DIVORCED `#Violet`2 because she refused to do\n" +
				"`2  any housework!  She got her old job back at the Inn!")
			mstr += "  `2You ask Violet to wash the dishes, and she refuses!  You have a big\n"
			mstr += "  fight!  You decide she isn't the women you married, and divorce her.\n"
			mstr += "  The whole experience has left you bitter.\n"
		}
		mstr += "\n"
		mstr += "  `4CHARM DROPS TO " + prettyInt(husband.Cha) + "\n"
		g.mailTo(int64(husband.Record), mstr)
	} else {
		switch g.random(4) {
		case 0:
			g.logLine("`2  `#Violet`2 has `%PMS!  `0" + husband.Name + "`2 is understanding, and peace\n" +
				"`2  is restored.  For now.")
			mstr += "  `2Violet gets angry over little things!  You realize she has PMS this\n"
			mstr += "  morning.  You treat her gently and calamity is avoided.\n"
			mstr += "\n"
			mstr += "  `%YOU RECEIVE " + prettyInt(100*husband.Level) + " EXPERIENCE.\n"
			husband.Exp += 100 * husband.Level
		case 1:
			g.logLine("`2  `#Violet`2 bears `0" + husband.Name + "`2 a male child.")
			mstr += "  `2Violet bears you a male child.  You have never loved her more.\n"
			mstr += "\n"
			mstr += "  `%YOU RECEIVE " + prettyInt(150*husband.Level) + " EXPERIENCE.\n"
			husband.Exp += 150 * husband.Level
			husband.Kids++
		case 2:
			g.logLine("`2  `#Violet`2 bears `0" + husband.Name + "`2 a female child.")
			mstr += "  `2Violet bears you a female child.  You are a little disappointed.\n"
			mstr += "  However, you are very pleased she retained her voluptuous figure.\n"
			mstr += "\n"
			mstr += "  `%YOU RECEIVE " + prettyInt(50*husband.Level) + " EXPERIENCE.\n"
			husband.Exp += 50 * husband.Level
			husband.Kids++
		case 3:
			g.logLine("`2  `#Violet`2 and `0" + husband.Name + "`2 didn't appear to get much sleep" +
				"`2  last night.  The town is mystified.")
			mstr += "  `2Violet pleases you in ways you had only dreamed about.  Nothing\n"
			mstr += "  has ever felt so good as your wife giving herself freely to you.\n"
			mstr += "\n"
			mstr += "  `%YOU RECEIVE " + prettyInt(100*husband.Level) + " EXPERIENCE.\n"
			husband.Exp += 100 * husband.Level
		}
		husband.Put()
		g.mailTo(int64(husband.Record), mstr)
	}
}

// sethMarriage is seth_marriage(): the daily news of Seth Able's wife.
func (g *Game) sethMarriage() {
	mstr := "\n`%  Latest news about `0Seth Able`%, your husband.\n`0-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-\n"
	wife := g.playerGet(g.state.MarriedToSeth)

	if g.random(5) < 1 {
		g.divorceSeth()
		wife.Cha = wife.Cha / 2
		wife.Put()
		switch g.random(3) {
		case 0:
			g.logLine("`2  `%Seth Able`2 has `%DIVORCED `0" + wife.Name + "`2 for cheating on him with\n" +
				"`2  the Bartender!  He now sings a song of woe!")
			mstr += "  `2Seth Able has divorced you because the Bartender said you kissed him!\n"
			mstr += "  You curse him!  What will people think?\n"
		case 1:
			g.logLine("`2  `0Seth Able `2has been `%BOOTED `2by `0" + wife.Name + "`2!  Artistic\n" +
				"`2  differences are to be blamed.")
			mstr += "  `2You are getting tired of seeing Seth hang around the house in his\n"
			mstr += "  underwear 'composing' music.  You tell him to get a real job.\n"
			mstr += "  There is a fight - And you end up booting this artist.\n"
		case 2:
			g.logLine("`2  `0" + wife.Name + " `2has `%DIVORCED `0Seth Able`2 because he refused to do\n" +
				"`2  any housework!  It is \"womans work\" was his reply!")
			mstr += "  `2You ask Seth to wash the dishes, and he refuses!  You have a big\n"
			mstr += "  fight!  You decide he isn't the man you married, and divorce him.\n"
			mstr += "  The whole experience has left you bitter.\n"
		}
		mstr += "\n"
		mstr += "  `4CHARM DROPS TO " + prettyInt(wife.Cha) + "\n"
		g.mailTo(int64(wife.Record), mstr)
	} else {
		switch g.random(4) {
		case 0:
			g.logLine("`2  `0" + wife.Name + "`2 screams at `%Seth Able `2for leaving\n" +
				"`2  the lid up!  Seth is able to calm her down - peace is restored.")
			mstr += "  `2You wake up to find the tiolet seat up - AGAIN!  You scream your\n"
			mstr += "  objections at your man - He promises to be more carefull in the future.\n"
			mstr += "\n"
			mstr += "  `%YOU RECEIVE " + prettyInt(100*wife.Level) + " EXPERIENCE.\n"
			wife.Exp += 100 * wife.Level
			wife.Put()
			g.mailTo(int64(wife.Record), mstr)
		case 1:
			g.logLine("`2  `0" + wife.Name + "`2 bears a male child - `%Seth Able`2 is proud.")
			mstr += "  You bear a male child!  Seth Able is extremely pleased with you.\n"
			mstr += "\n"
			mstr += "  `%YOU RECEIVE " + prettyInt(150*wife.Level) + " EXPERIENCE.\n"
			wife.Exp += 150 * wife.Level
			wife.Kids++
			wife.Put()
			g.mailTo(int64(wife.Record), mstr)
		case 2:
			g.logLine("`2  `0" + wife.Name + "`2 bears a female child - `%Seth Able `2approves.")
			mstr += "  `2You bear a female child!  You are puzzled why Seth is not as happy\n"
			mstr += "  as you.  He refuses to speak about it.\n"
			mstr += "\n"
			mstr += "  `%YOU RECEIVE " + prettyInt(50*wife.Level) + " EXPERIENCE.\n"
			wife.Exp += 50 * wife.Level
			wife.Kids++
			wife.Put()
			g.mailTo(int64(wife.Record), mstr)
		case 3:
			g.logLine("`2  `%Seth Able`2 and `0" + wife.Name + "`2 didn't appear to get much sleep\n" +
				"`2  last night.  The town is mystified.")
			mstr += "  `2Seth Able pleases you in ways you had only dreamed about.  Nothing\n"
			mstr += "  has ever felt so good as your husband doing the chores around the house.\n"
			mstr += "\n"
			mstr += "  `%YOU RECEIVE " + prettyInt(100*wife.Level) + " EXPERIENCE.\n"
			wife.Exp += 100 * wife.Level
			wife.Put()
			g.mailTo(int64(wife.Record), mstr)
		}
	}
}

// dailyMaintTalkNow is daily_maint()'s talk_now(): an NPC answers the last
// player who talked in the bar.
func (g *Game) dailyMaintTalkNow(file string, pl *Player) {
	if g.Settings.NoChat {
		return
	}
	n := g.random(24)
	if n > 11 {
		return
	}
	if pl.Name == "X" {
		return
	}
	all, ok := readLines(file)
	if !ok {
		return
	}
	switch n {
	case 0:
		all = append(all, pick(pl.Sex == "F", "  `%Seth Able:", "  `%Violet:"))
		if pl.Cha < 3 {
			all = append(all, "  `0"+pl.Name+"`0..  Have you showered recently?")
		} else if pl.Cha < 21 {
			all = append(all, "  `0Greetings, "+pl.Name+"`0.. Don't let them get you down.")
		} else {
			all = append(all, "  `0Hey "+pl.Name+"`0, come closer to me..")
		}
	case 1:
		all = append(all, "  `%Bartender:")
		all = append(all, "  `0You jabber too much, "+pl.Name+".")
	case 2:
		all = append(all, "  `%Barak:")
		if pl.Level > 2 {
			all = append(all, "  `0You think you're tough, "+pl.Name+" `0'cuz you got a "+pl.Weapon+"?")
		} else {
			all = append(all, "  `0You really got a big mouth, "+pl.Name+"`0.  Ugly one too.")
		}
	case 3:
		all = append(all, "  `%Turgon:")
		if pl.Level < 6 {
			all = append(all, "  `0Are you saying what I think you're saying "+pl.Name+"`0?")
		} else if pl.Sex == "M" {
			all = append(all, "  `0Man I need to get wasted...")
		} else {
			all = append(all, "  `0Heed the words of "+pl.Name+", `0For she is a great warrior.")
		}
	case 4:
		all = append(all, "  `%Aragorn:")
		all = append(all, "  `0Hey, "+pl.Name+"`0, whats the deal with Jennie Garth?")
	case 5:
		all = append(all, "  `%Seth Able:")
		all = append(all, "  `0Damn, I need a drink.")
	case 6:
		all = append(all, "  `%Grizelda:")
		if pl.Sex == "M" {
			all = append(all, "  `0"+pl.Name+"`0! You're a cute one!  Come to mama!")
		} else if g.Settings.CleanMode {
			all = append(all, "  `0"+pl.Name+"`0..Stay away from my man, wench!")
		} else {
			all = append(all, "  `0"+pl.Name+"`0..Stay away from my man, slut.")
		}
	case 7:
		all = append(all, "  `%Old Women From The Corner:")
		all = append(all, "  `0There is no `4Dragon`0...The children all just ran away.")
	case 8:
		all = append(all, "  `%Violet:")
		all = append(all, "  `0I'm so tired.  Will you escort me to my room,"+pl.Name+"`0?")
	case 9:
		all = append(all, "  `%Hooded Warrior:")
		all = append(all, "  `0Watch you're back, "+pl.Name+"`0.")
	case 10:
		all = append(all, "  `%Old Man:")
		if pl.Cha > 3 {
			all = append(all, "  `0Well met, "+pl.Name+"`0!")
		} else {
			all = append(all, "  `0"+pl.Name+"`0!  Why did you kick my ass the other day?!")
		}
	case 11:
		all = append(all, "  `%Barak:")
		if pl.Sex == "M" {
			all = append(all, "  `0"+pl.Name+"`0.  Interesting name...")
		} else if pl.Cha < 4 {
			all = append(all, "  `0"+pl.Name+"`0!  You have no breasts!  Are you a man!?  Har!")
		} else if pl.Cha < 8 {
			all = append(all, "  `0Finally!  A woman with some sense!")
		} else {
			all = append(all, "  `0"+pl.Name+"`0!  I think I'm in love!  Marry me honey!  Har!")
		}
	}
	for len(all) > 18 {
		all = all[1:]
	}
	_ = writeLines(file, all...)
}

// dailyMaint is daily_maint(): the first caller of the day runs the
// realm's maintenance.  The caller holds the state lock; state is written
// after it returns.
func (g *Game) dailyMaint() {
	st := g.state

	g.sln("")
	g.foreground(10)
	switch g.random(4) {
	case 0:
		g.sw("  You stop to notice a bird singing this bright morning.")
	case 1:
		g.sw("  You take the time to clean your weapon this morning.")
	case 2:
		g.sw("  Time seems to stand still as you contemplate life.")
	case 3:
		g.sw("  Waking up is slow this morning.  You blame lastnights ale.")
	}
	g.sw(".")
	st.Days++
	g.sw(".")
	if g.Tournament.Enabled && g.Tournament.Days > 0 && g.Tournament.Days <= st.Days {
		g.tournamentOver()
	}
	g.sw(".")
	// NPC talks to last person who chatted in the bar...
	if bar := g.gamedir("bar.lrd"); fileExists(bar) {
		g.fmutex(bar)
		if st.LastBar != -1 {
			if tp := g.playerGet(st.LastBar); tp != nil {
				g.dailyMaintTalkNow(bar, tp)
				st.LastBar = -1
			}
		}
		g.rfmutex(bar)
	}

	g.sln(".")
	// this is the start of the actual player maintenance
	plen := g.playerLength()
	for i := 0; i < plen; i++ {
		pl := g.playerGet(int64(i))
		if pl == nil {
			continue
		}
		if !pl.OnNow {
			fileRemove(g.gamedir("war" + itoa(pl.Record) + ".bin"))
			fileRemove(g.gamedir("mess" + itoa(pl.Record) + ".bin"))
		}
		if pl.OnNow {
			continue
		}
		switch {
		case pl.Name == "X":
			if st.MarriedToSeth == int64(i) {
				st.MarriedToSeth = -1
			}
			if st.MarriedToViolet == int64(i) {
				st.MarriedToViolet = -1
			}
		case g.Settings.Del1XP && pl.Exp == 1 && pl.DragKills == 0,
			g.Settings.DeleteDays > 0 && pl.Time < st.Days-g.Settings.DeleteDays:
			pl.Name = "X"
			pl.RealName = "X"
			pl.Put()
			g.killmail(pl.Record)
			if st.MarriedToSeth == int64(i) {
				st.MarriedToSeth = -1
			}
			if st.MarriedToViolet == int64(i) {
				st.MarriedToViolet = -1
			}
		case st.Days > 3:
			if pl.Time < st.Days-2 {
				if pl.LastReincarnated < st.Days-2 {
					if pl.Dead {
						if pl.Gone < g.Settings.ResDays {
							pl.Dead = false
							g.lln("  `2(`0" + pl.Name + " `2is raised from the dead)")
							pl.Inn = false
							pl.LastReincarnated = st.Days
							// BUG (lord.js): tests state.res_days, which
							// doesn't exist, so `gone` never increases.
							// Kept.
							pl.Divorced = false
							pl.Put()
						}
					}
				}
			}
		}
	}

	g.sln("")
	plen = g.playerLength()
	if st.MarriedToSeth > -1 {
		if st.MarriedToSeth > int64(plen) {
			g.divorceSeth()
		} else if tp := g.playerGet(st.MarriedToSeth); tp != nil {
			// DIFF: lord.js throws when married_to_seth == player_length()
			// (player_get returns null); nothing happens here.
			if tp.Name == "X" {
				g.divorceSeth()
			} else {
				g.sethMarriage()
			}
		}
	}
	if st.MarriedToViolet > -1 {
		if st.MarriedToViolet > int64(plen) {
			g.divorceViolet()
		} else if tp := g.playerGet(st.MarriedToViolet); tp != nil {
			// DIFF: as above for married_to_violet == player_length().
			if tp.Name == "X" {
				g.divorceViolet()
			} else {
				g.violetMarriage()
			}
		}
	}
	g.more()
}

// tday is create_log()'s day number: unique for each calendar day
// (Date.getYear() is years since 1900, getMonth() is 0-based).
func (g *Game) tday() int64 {
	now := g.now()
	return int64(now.Year()-1900)*366 + int64(now.Month()-1)*31 + int64(now.Day())
}

var happenings = []string{
	"`4  A Child was found today!  But scared deaf and dumb.",
	"`4  More children are missing today.",
	"`4  A small girl was missing today.",
	"`4  The town is in grief.  Several children didn't come home today.",
	"`4  Dragon sighting reported today by a drunken old man.",
	"`4  Despair covers the land - more bloody remains have been found today.",
	"`4  A group of children did not return from a nature walk today.",
	"`4  The land is in chaos today.  Will the abductions ever stop?",
	"`4  Dragon scales have been found in the forest today..Old or new?",
	"`4  Several farmers report missing cattle today.",
}

// createLog is create_log(): start today's lognow.lrd (moving yesterday's
// to logold.lrd) and run the daily maintenance once per day.
func (g *Game) createLog(havelock bool) {
	f := g.gamedir("lognow.lrd")
	tday := g.tday()

	if !havelock {
		g.fmutex(f)
	}

	g.getState(true)
	if fileExists(f) {
		if g.state.LogDate != tday {
			fileRemove(g.gamedir("logold.lrd"))
			fileRename(g.gamedir("lognow.lrd"), g.gamedir("logold.lrd"))
		} else {
			g.putState()
			if !havelock {
				g.rfmutex(f)
			}
			return
		}
	} else {
		// BUG (lord.js): this drops the state lock, so the maintenance
		// below runs without it.  Kept.
		g.putState()
	}
	err := appendLines(f,
		"`2  The Daily Happenings....",
		"`0-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-",
		happenings[g.random(int64(len(happenings)))],
		"`.                                `2-`0=`2-`0=`2-`0=`2-")
	if err != nil {
		g.rfmutex(f)
		g.fatalf("Unable to open %s", f)
	}
	g.state.LogDate = tday
	// TODO: IGM Maint and out check things.
	if !havelock {
		g.rfmutex(f)
	}
	g.dailyMaint()
	g.putState()
}

// logLine is log_line(): add an entry to today's happenings.
//
// BUG (lord.js): when lognow.lrd doesn't exist, create_log(true) runs the
// daily maintenance while this function holds the lognow.lrd mutex, and a
// marriage event there calls log_line again, which waits for that mutex
// and fails after 15 seconds.  Kept.
func (g *Game) logLine(text string) {
	f := g.gamedir("lognow.lrd")

	g.fmutex(f)
	if !fileExists(f) {
		g.createLog(true)
	}
	if err := appendLines(f, text, "`.                                `2-`0=`2-`0=`2-`0=`2-"); err != nil {
		g.fatalf("Unable to open %s", f)
	}
	g.rfmutex(f)
}
