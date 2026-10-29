package game

import (
	"fmt"
	"sort"
	"strings"
)

// beefUp is beef_up(): make a monster tougher for dragon slayers.
func (g *Game) beefUp(m *Player) {
	var drk int64

	if g.Settings.BeefUp {
		if g.player.DragKills > 0 && g.random(5) < 2 {
			if g.player.DragKills < 6 {
				drk = g.player.DragKills
			} else {
				drk = 5
			}
			max := 2000000000 - m.HP
			r := g.random(10) + g.random(10)
			adj := g.randomf(float64(drk*r*m.HP) / 100)
			if adj < max {
				m.HP += adj
			}
			max = 2000000000 - m.Str
			r = g.random(10) + g.random(10)
			adj = g.randomf(float64(drk*r*m.Str) / 100)
			if adj < max {
				m.Str += adj
			}
		}
	}
}

// getArmourweap is get_armourweap(): entry num of a weapon/armour table,
// from the sysop's fname if there is one.
func (g *Game) getArmourweap(num int64, stats []WeaponArmour, fname string) WeaponArmour {
	fn := g.gameOrExec(fname)

	if num < 1 {
		return stats[0]
	}
	if num > 15 {
		num = 15
	}
	if fn == "" {
		return stats[num]
	}
	ret, _, err := readRecordAt[WeaponArmour](fn, int(num-1))
	if err != nil {
		// f.length < num
		return stats[num]
	}
	ret.IsArena = true
	return *ret
}

// getTrainer is get_trainer(): a copy of the master for level num.
func (g *Game) getTrainer(num int64) *Player {
	fn := g.gameOrExec("trainers.bin")

	if num < 1 {
		t := trainerStats[1]
		return &t
	}
	if num >= int64(len(trainerStats)) {
		num = int64(len(trainerStats)) - 1
	}
	if fn == "" {
		t := trainerStats[num]
		return &t
	}
	rec, _, err := readRecordAt[trainerRec](fn, int(num-1))
	if err != nil {
		// f.length < num
		t := trainerStats[num]
		return &t
	}
	return &Player{
		Name:      rec.Name,
		Weapon:    rec.Weapon,
		Swear:     rec.Swear,
		NeedStr1:  rec.NeedStr1,
		NeedStr2:  rec.NeedStr2,
		Str:       rec.Str,
		HP:        rec.HP,
		HPGained:  rec.HPGained,
		StrGained: rec.StrGained,
		Def:       rec.Def,
		Need:      rec.Need,
		IsArena:   true,
		Record:    -1,
	}
}

// getArmour is get_armour().
func (g *Game) getArmour(num int64) WeaponArmour {
	return g.getArmourweap(num, armourStats, "armour.bin")
}

// getWeapon is get_weapon().
func (g *Game) getWeapon(num int64) WeaponArmour {
	return g.getArmourweap(num, weaponStats, "weapon.bin")
}

// formatDate is format_date(): today as MM/DD/YYYY.
func (g *Game) formatDate() string {
	now := g.now()
	return fmt.Sprintf("%02d/%02d/%04d", int(now.Month()), now.Day(), now.Year())
}

// tournamentCheck is tournament_check(): has the player won a stat-mode
// tournament?
func (g *Game) tournamentCheck() {
	t := &g.Tournament
	p := g.player

	if !t.Enabled {
		return
	}
	if t.Days > 0 {
		return
	}
	if p.Exp < t.XP {
		return
	}
	if p.DragKills < t.DKills {
		return
	}
	if p.PvP < t.PKills {
		return
	}
	if p.Level < t.Level {
		return
	}
	if p.Laid < t.Lays {
		return
	}

	g.getState(true)
	if g.state.WonBy < 0 {
		g.state.WonBy = int64(p.Record)
		g.state.LatestHero = p.Name
		g.putState()
		g.sclrscr()
		g.sln("")
		g.sln("")
		g.foreground(10)
		g.sln("  Congratulations! You have won the tournament!")
		g.sln("")
		g.lln("  `2<`0MORE`2>")
		g.getkey()
		p.OnNow = false
		p.Put()
	} else {
		g.putState()
		op := g.playerGet(g.state.WonBy)
		g.sclrscr()
		g.sln("")
		g.sln("")
		g.foreground(10)
		g.lln("  You have completed the tournament, but " + op.Name)
		g.lln("  got there first.")
		g.sln("")
		g.lln("  `2<`0MORE`2>")
		g.getkey()
		p.OnNow = false
		p.Put()
	}
	g.exit(0)
}

// quitMsg is quit_msg().
func (g *Game) quitMsg() {
	g.sclrscr()
	g.sln("")
	switch g.random(5) {
	case 0:
		g.lln("`2  The very earth groans at your depature.")
		return
	case 1:
		g.lln("`2  The very trees seem to moan as you leave.")
		return
	case 2:
		g.lln("`2  Echoing screams fill the wastelands as you close your eyes.")
		return
	case 3:
		g.lln("`2  The black thing inside rejoices at your departure.")
		return
	case 4:
		g.lln("`2  Your very soul aches as you wake up from your favorite dream.")
		return
	}
	g.sln("")
}

// instructions is instructions().
func (g *Game) instructions() {
	g.lrdfile("HINTS", true)
	g.moreNomail()
	g.sln("")
}

// checkName is check_name(): refuse reserved names.
func (g *Game) checkName(str string) bool {
	resp := ""

	str = strings.ToUpper(dispStr(str))

	switch str {
	case "BARAK":
		resp = "Naw, the real Barak would decapitate you if he found out."
	case "SETH":
		resp = "You are not Seth Able!  Don't take his name in vain!"
	case "SETH ABLE":
		resp = "You are not God!"
	case "TURGON":
		resp = "Haw.  Hardly - Turgon has muscles."
	case "VIOLET":
		resp = "Haw.  Hardly - Violet has breasts."
	case "RED DRAGON":
		resp = "Oh go plague some other land!"
	case "DRAGON":
		resp = "You ain't Bruce Lee, so get out!"
	case "JENNIE GARTH":
		resp = "You are not a goddess, don't use her name!"
	case "KIRSTEN DUNST":
		resp = "Hardly! You only wish you were in a movie with Wynona!"
	case "BUSH":
		resp = "Lower my taxes!"
	case "BAGGIO":
		resp = "Darius sucks!"
	case "DAVID FOLLEY":
		resp = "You rule, dude - but use a handle or you will mobbed!"
	case "ARNOLD PALMER":
		resp = "Ha!  You're too old to be playing games."
	case "BARTENDER":
		resp = "Nah, the bartender is smarter than you!"
	case "CHANCE":
		resp = "Why not go take a chance with a rattlesnake?"
	case "MICHAEL PRESLAR":
		resp = "You want to be a small town kid?"
	case "GOD", "JESUS":
		resp = "Why arent you in church?"
	}

	if resp != "" {
		g.sln("")
		g.lln("  `)** `%" + resp + " `)**`2")
		g.sln("")
		return false
	}
	return true
}

// chooseProfession is choose_profession(): 1 Death Knight, 2 Mystical,
// 3 Thieving.
func (g *Game) chooseProfession(dhTavern bool) int64 {
	var ch string

	if dhTavern {
		g.lln("  `2You concentrate deeply.")
	} else {
		g.lln("  `%As you remember your childhood, you remember...")
		g.sln("")
		g.lln("  `0(`5K`0)illing A Lot Of Woodland Creatures")
		g.lln("  `0(`5D`0)abbling In The Mystical Forces")
		g.lln("  `0(`5L`0)ying, Cheating, And Stealing From The Blind")
	}

	for {
		g.sln("")
		g.lw("  `2Pick one.  (`0K`2,`0D`2,`0L`2) : `%")
		ch = strings.ToUpper(g.getkey())
		g.sln(ch)
		g.foreground(10)
		g.sln("")
		g.sln("")
		switch ch {
		case "K":
			g.sln("  Now that you've grown up, you have decided to study the ways of the")
			g.sln("  the Death Knights.  All beginners want the power to use their body")
			g.sln("  and weapon as one.  To inflict twice the damage with the finesse only")
			g.sln("  a warrior of perfect mind can do.")
		case "D":
			g.sln("  You have always wanted to explain the unexplainable.  To understand the")
			g.sln("  powerful forces that rule the earth.  To tame the beast that oversees")
			g.sln("  all things.  Of course, having the power to burn someone by making a")
			g.sln("  gesture wouldn't hurt.")
		case "L":
			g.sln("  You decide to follow your instincts.  To get better at what you've")
			g.sln("  always done best.  So you decide to lead a dishonest lifestyle.")
			g.sln("  Of course, your ultimate goal will always be to join the Master Thieves")
			g.sln("  Guild.  To arrive, remember one thing, \"Even thieves have honor.\"")
		default:
			g.sln("")
			g.sln("  This is a very important choice.  Pay attention idiot!")
		}
		if strings.Contains("KDL", ch) {
			break
		}
	}
	g.moreNomail()
	return int64(strings.Index(" KDL", ch))
}

// findPlayer is find_player(): pick a player by full or partial name;
// -1 if none was chosen.
func (g *Game) findPlayer() int64 {
	g.sln("")
	g.lln("  `2(full or `0PARTIAL`2 name)")
	g.lw("  NAME: `%")
	name := cleanStr(g.consoleGetstr(GetStrOpts{}))
	ucname := strings.ToUpper(dispStr(name))
	plen := g.playerLength()
	for i := 0; i < plen; i++ {
		op := g.playerGet(int64(i))
		if op.Name == "X" {
			continue
		}
		tname := strings.ToUpper(dispStr(op.Name))
		if strings.Contains(tname, ucname) {
			g.sln("")
			g.lw("  `2You mean \"`0" + op.Name + "\"? `2[`%Y`2] `: ")
			ch := strings.ToUpper(g.getkey())
			if ch != "N" {
				ch = "Y"
			}
			g.sln(ch)
			if ch == "Y" {
				g.sln("")
				return int64(op.Record)
			}
		}
	}
	return -1
}

// commandPrompt is command_prompt().
func (g *Game) commandPrompt() {
	// Subtract 30 seconds from time remaining.
	tl := g.secondsRemaining() - 30
	remain := itoa(tl/60) + ":" + itoa(tl%60)

	g.sln("")
	g.lw("`2  Your command,`0 " + g.player.Name + "`2? `2[`%" + remain + "`2] : `%")
}

// addHP is add_hp().
func (g *Game) addHP(num int64) {
	g.player.HPMax += num
	if g.player.HPMax > 32000 {
		g.player.HPMax = 32000
	}
}

// addStr is add_str().
func (g *Game) addStr(num int64) {
	g.player.Str += num
	if g.player.Str > 32000 {
		g.player.Str = 32000
	}
}

// addDef is add_def().
func (g *Game) addDef(num int64) {
	g.player.Def += num
	if g.player.Def > 32000 {
		g.player.Def = 32000
	}
}

// generateRankings is generate_rankings(): write the player rankings to
// fname.
func (g *Game) generateRankings(fname string, all, on, inn bool) {
	var a []*Player

	for _, p := range g.allPlayers() {
		if p.Name != "X" {
			if all || !p.Dead {
				if inn {
					if p.Inn {
						a = append(a, p)
					}
				} else if all || !p.Inn {
					a = append(a, p)
				}
			}
		}
	}
	sort.SliceStable(a, func(i, j int) bool { return a[i].Exp > a[j].Exp })
	g.fmutex(fname)
	lines := []string{
		"",
		"",
		"`%" + center("Legend Of The Red Dragon "+Version+" - Player Rankings"),
		"",
		"`2    Name                    Experience    Level    Mastered    Status",
		"`0-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-",
	}
	for _, p := range a {
		l := ""
		if p.Sex == "F" {
			l += "`#F "
		} else {
			l += "  "
		}
		switch p.Clss {
		case 1:
			l += "`0D "
		case 2:
			l += "`#M "
		case 3:
			l += "`9T "
		default:
			l += "  "
		}
		l += "`2" + spacePad(p.Name, 22) + "`2"
		l += spacePad(prettyInt(p.Exp), 13, true)
		l += "    "
		l += "`%"
		l += spacePad(prettyInt(p.Level), 2)
		l += "        "
		if p.SkillW > 39 {
			l += "`%D"
		} else if p.SkillW > 19 {
			l += "`0D"
		} else {
			l += " "
		}
		if p.SkillM > 39 {
			l += "`%M"
		} else if p.SkillM > 19 {
			l += "`#M"
		} else {
			l += " "
		}
		if p.SkillT > 39 {
			l += "`%T"
		} else if p.SkillT > 19 {
			l += "`9T"
		} else {
			l += " "
		}
		l += "     "
		if p.Dead {
			l += "`4Dead "
		} else {
			l += "`%Alive"
		}
		if on && p.OnNow {
			l += "  `%On"
		}
		lines = append(lines, l)
	}
	if err := writeLines(fname, lines...); err != nil {
		g.rfmutex(fname)
		g.fatalf("Unable to open %s", fname)
	}
	g.rfmutex(fname)
}

// checkFields is check_fields(): clamp the player's numbers.
func (g *Game) checkFields() {
	p := g.player
	clampField := func(v *int64, max int64) {
		if *v < 0 {
			*v = 0
		}
		if *v > max {
			*v = max
		}
	}
	clampField(&p.HP, 32000)
	clampField(&p.HPMax, 32000)
	clampField(&p.ForestFights, 32000)
	clampField(&p.PvPFights, 32000)
	clampField(&p.Gold, 2000000000)
	clampField(&p.Bank, 2000000000)
	clampField(&p.Def, 32000)
	clampField(&p.Str, 32000)
	clampField(&p.Cha, 32000)
	clampField(&p.Exp, 2000000000)
	clampField(&p.Laid, 32000)
	clampField(&p.Kids, 32000)
}
