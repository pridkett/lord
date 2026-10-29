package game

import (
	"os"
	"strings"
)

// The online duel: two callers on the same world fight through small files
// in the game directory.  war<rec>.bin carries a battle response (damage,
// "A" or "R") to player <rec>, mess<rec>.bin a yelled message, and the
// fmutex fight<rec>.rec marks player <rec> as being in a duel.
//
// lord.js keeps `cantaunt` as a global; it only matters while on_battle
// runs (on_battle sets it before its loop), so here it is a flag owned by
// each on_battle call and passed down to ob_get_message.

// fileMutex is Synchronet's file_mutex(name, text): create name
// exclusively with text as its content.  It fails (false) when the file
// already exists.
func fileMutex(name, text string) bool {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return false
	}
	defer f.Close()
	_, _ = f.WriteString(text)
	return true
}

// obCleanup is ob_cleanup().
func (g *Game) obCleanup() {
	fileRemove(g.gamedir("war" + itoa(g.player.Record) + ".bin"))
	g.rfmutex(g.gamedir("fight" + itoa(g.player.Record) + ".rec"))
	fileRemove(g.gamedir("mess" + itoa(g.player.Record) + ".bin"))
}

// obReadMsg is ob_read_msg(): wait (up to 15 seconds) for a complete line
// in fname, then remove the file and return the line.
func (g *Game) obReadMsg(fname string) string {
	var msg string
	end := g.time() + 15

	for {
		var b []byte
		var err error
		for {
			b, err = os.ReadFile(fname)
			if err == nil || g.time() > end {
				break
			}
			g.mswait(100)
		}
		msg = string(b)
		if g.time() > end {
			break
		}
		if len(msg) > 0 && msg[len(msg)-1] == '\n' {
			break
		}
		// DIFF: lord.js re-reads a partly written file in a tight loop;
		// pause briefly between reads instead.
		g.mswait(10)
	}
	msg = strings.NewReplacer("\r", "", "\n", "").Replace(msg)

	fileRemove(fname)
	return msg
}

// obSendResp is ob_send_resp().
func (g *Game) obSendResp(op *Player, str string) {
	fileMutex(g.gamedir("war"+itoa(op.Record)+".bin"), str+"\n")
}

// obGetMessage is ob_get_message(): show a yelled message if there is one.
func (g *Game) obGetMessage(op *Player, cantaunt *bool) bool {
	if fileExists(g.gamedir("mess" + itoa(g.player.Record) + ".bin")) {
		g.sln("")
		g.sln("")
		g.lln("  `0" + op.Name + " `2screams something at you!`%")
		g.sln("")
		msg := g.obReadMsg(g.gamedir("mess" + itoa(g.player.Record) + ".bin"))
		g.lln(msg)
		*cantaunt = true
		g.sln("")
		g.moreNomail()
		return true
	}
	return false
}

// obGetResponse is ob_get_response(): wait for the opponent's battle
// response.  It returns "R" when the opponent left, "ABORT" when the
// player gave up waiting.
func (g *Game) obGetResponse(op *Player, cantaunt *bool) string {
	var dh string
	loops := 0
	done := false

	g.sln("")
	g.lw("  `2Waiting`0.")
	for {
		op.ReLoad()
		if !op.OnNow {
			return "R"
		}
		if g.con.WaitKey(0) {
			dh = strings.ToUpper(g.getkey())
		} else {
			dh = ""
		}
		if loops > 70 && dh == "R" {
			g.logLine("`.  `0" + g.player.Name + " `2has been eluded by `0" + op.Name + "`2.")
			g.obSendResp(op, "R")
			g.sln("")
			return "ABORT"
		}
		g.sw(".")
		if loops == 70 {
			g.sln("")
			g.lln("  `2(`0R`2)un Away")
			g.sln("")
			g.lw("  `2Continuing wait`0.")
		}
		if g.obGetMessage(op, cantaunt) {
			g.sln("")
			g.lw("  `2Continuing wait`0.")
		}
		loops++
		g.mswait(600)
		if fileExists(g.gamedir("war" + itoa(g.player.Record) + ".bin")) {
			done = true
		}
		if done {
			break
		}
	}
	msg := g.obReadMsg(g.gamedir("war" + itoa(g.player.Record) + ".bin"))
	g.sln("")
	return msg
}

// onBattle is on_battle(): the online duel itself.  action is the
// opponent's first strike when the opponent attacks first ("" is
// undefined: wait for it).  It returns "" (lord.js returns nothing).
func (g *Game) onBattle(op *Player, first bool, action string) string {
	player := g.player
	cantaunt := true

	onStruck := func(op *Player, amountStr string) {
		his := pick(op.Sex == "M", "his", "her")

		cantaunt = true
		amount, ok := jsParseInt(amountStr)
		// DIFF: lord.js turns a garbled response into NaN damage (and NaN
		// hitpoints); treat it as a miss.
		if !ok {
			amount = 0
		}

		if amount < 1 {
			g.lln("  `0" + op.Name + " `2misses you!")
			g.sln("")
			return
		}
		g.lln("  `0" + op.Name + " `2hits you for `4" + prettyInt(amount) + "`2 points of damage!")
		player.HP -= amount
		g.sln("")
		if player.HP < 1 {
			g.lln("  You fall to the ground.  You think " + op.Name + " `2is going to let you")
			g.lln("  go... Just before you feel " + his + " " + op.Weapon + " sliding")
			g.sln("  through your chest. ")
			g.sln("")
			g.sln("  You distantly feel your lungs filling with blood, when everything becomes")
			g.sln("  black.")
			player.Dead = true
			g.logLine("`.  `0" + op.Name + " `2has killed `5" + player.Name + " `2in an Online Duel!")
			player.Gold = 0
			player.Exp -= player.Exp / 10
			player.HP = 0
			player.OnNow = false
			player.Put()
			g.moreNomail()
			g.deadScreen(op)
			g.sln("")
		}
	}

	waitForHit := func(op *Player) string {
		g.lln("  You wait as " + op.Name + " `2strikes.")
		waction := g.obGetResponse(op, &cantaunt)
		if waction == "ABORT" {
			g.sln("")
			g.sln("  You run as fast as your legs will carry you!")
			g.sln("")
			g.moreNomail()
			return "ABORT"
		}
		if waction == "R" {
			g.lln("  `0The sniveling " + op.Name + " `0has run away.  Unfortunately, this battle")
			g.sln("  is OVER.")
			g.logLine("  `0" + player.Name + " `2& `0" + op.Name + " `2had an uneventful online duel.")
			g.sln("")
			g.moreNomail()
			return "ABORT"
		}
		onStruck(op, waction)
		return ""
	}

	if first {
		g.lln("  `%** YOUR SKILL ALLOWS YOU TO ATTACK FIRST **")
		g.sln("")
		g.moreNomail()
	} else {
		g.lln("  `%** YOUR ENEMY STRIKES YOU FIRST **")
		if action == "" {
			// DIFF: lord.js has `if (wait_for_hit(op) === 'ABORT');` with a
			// stray semicolon, so this side always left the duel after the
			// first strike; it now leaves only on ABORT.
			if waitForHit(op) == "ABORT" {
				return ""
			}
		} else {
			onStruck(op, action)
		}
	}

	cantaunt = true
outer:
	for {
		op.ReLoad()
		if player.HP < 1 || op.HP < 1 {
			break
		}
		g.lln("  `2Your Hitpoints:`0 " + prettyInt(player.HP))
		g.lln("  `0" + op.Name + "`2's Hitpoints:`0 " + prettyInt(op.HP))
		g.sln("")
		g.lln("  `2(`0A`2)ttack Your Enemy")
		if cantaunt {
			g.lln("  `2(`0Y`2)ell Something To Your Enemy")
		}
		g.lln("  `2(`0R`2)un For Your Life")
		g.sln("")
		g.lw("  `2Your Command ?  :`0 ")
		ch := strings.ToUpper(g.getkey())
		g.sln(ch)
		if g.obGetMessage(op, &cantaunt) {
			continue
		}
		switch ch {
		case "A":
			g.sln("")
			tmp := g.random(player.Str/2) + player.Str/2
			tmp -= op.Def
			if tmp < 1 {
				g.lln("  `2You miss `0" + op.Name + "`2 Completely!")
				g.sln("")
				g.obSendResp(op, "0")
				if waitForHit(op) == "ABORT" {
					break outer
				}
				break
			}
			g.lln("  `2You HIT `0" + op.Name + "`2 for `4" + prettyInt(tmp) + "`2 points of damage!")
			g.obSendResp(op, itoa(tmp))
			op.HP -= tmp
			if op.HP < 1 {
				g.sln("")
				g.lln("  You have killed " + op.Name + "`2.")
				player.KilledAPlayer = true
				g.sln("")
				tmp = player.Gold
				player.Gold += op.Gold
				if player.Gold > 2000000000 {
					player.Gold = 2000000000
				}
				g.lw("  `2You receive `%" + prettyInt(player.Gold-tmp) + "`2 gold, and " + prettyInt(op.Exp/2) + " experience!")
				if op.Gem > 2 {
					tmp = player.Gem
					player.Gem += op.Gem / 2
					if player.Gem > 32000 {
						player.Gem = 32000
					}
					// DIFF: You got an extra pvp fight here?!?!
					g.lw("  You also find " + prettyInt(player.Gem-tmp) + " Gems!")
				}
				g.sln("")
				player.Exp += op.Exp / 2
				player.Put()
				g.moreNomail()
				g.tournamentCheck()
				break outer
			}
			if waitForHit(op) == "ABORT" {
				break outer
			}
		case "Y":
			if cantaunt {
				g.sln("")
				g.sln("  Yell what?")
				g.sln("")
				g.lw(" `0>`2")
				tmp := "  " + g.consoleGetstr(GetStrOpts{})
				fileMutex(g.gamedir("mess"+itoa(op.Record)+".bin"), tmp+"\n")
				g.sln("")
				g.sln("")
				g.lln("  `%** MESSAGE YELLED **")
				g.sln("")
				cantaunt = false
			} else {
				g.sln("")
				g.lln("  `2You feel you had better stop taunting " + op.Name + "`2, or you will")
				g.sln("  lose your attack!")
				g.sln("")
			}
		case "R":
			g.obSendResp(op, "R")
			g.sln("")
			g.lln("  `2You run like the dickens!  You manage to elude `0" + op.Name + "`2!")
			break outer
		}
	}
	return ""
}

// onlineBattle is online_battle(): challenge an online player to a duel.
func (g *Game) onlineBattle(op *Player) {
	player := g.player
	cantaunt := true

	g.lln("`c  `%                           ** ONLINE BATTLE **")
	g.lln("`0-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-==-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
	// DIFF: lord.js checks file_exists() on the fight file, which only
	// ever exists as fmutex's ".lock"; check the lock.
	if fileMutexExists(g.gamedir("fight" + itoa(op.Record) + ".rec")) {
		g.sln("")
		g.sln("  Sorry!  That user is currently already in battle.  Try again later.")
		g.sln("")
		g.moreNomail()
		g.obCleanup()
		return
	}
	g.fmutex(g.gamedir("fight" + itoa(player.Record) + ".rec"))
	g.mailTo(int64(op.Record), "`o"+itoa(player.Record))
	// DIFF: You got free full healing here!
	// player.hp = player.hp_max
	player.Put()
	g.lln("  `2Issuing an Online Duel challenge to " + op.Name + "`2.")
	action := g.obGetResponse(op, &cantaunt)
	if action == "ABORT" {
		g.sln("")
		g.sln("  You decide you are really not ready for battle anyway.")
		g.sln("")
		g.moreNomail()
		g.obCleanup()
		return
	}
	g.sln("")
	if strings.HasPrefix(action, "R") {
		g.lln("  `2You call out, but `0" + op.Name + "`2 hides from you. ")
		g.sln("")
		g.moreNomail()
		g.obCleanup()
		return
	}
	player.PvPFights--
	g.lln("  `0" + op.Name + " `2draws his " + op.Weapon + " for battle!")
	g.sln("")
	if strings.HasPrefix(action, "A") {
		g.onBattle(op, true, "")
	} else {
		if action == "" {
			// DIFF: an empty response (the opponent's reply timed out) is
			// NaN damage in lord.js; treat it as a miss.
			action = "0"
		}
		g.onBattle(op, false, action)
	}
	g.obCleanup()
}
