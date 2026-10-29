package game

import (
	"strings"
)

// enemyAttack is enemy_attack(): the opponent strikes the player.
func (g *Game) enemyAttack(op *Player, pfight bool) {
	var son, his string
	player := g.player

	atk := g.randomf(float64(op.Str)/2) + op.Str/2
	if op.IsDragon {
		switch g.random(4) {
		case 0:
			op.Weapon = "Huge Claw"
		case 1:
			op.Weapon = "Swishing Tail"
		case 2:
			op.Weapon = "`4Flaming Breath"
			atk += atk
		case 3:
			op.Weapon = "Stomping The Ground"
		}
	}
	if g.random(30) == 1 {
		atk += atk / 2
		g.lln("  `4** `0" + op.Name + "`4 Executes A Power Move **")
		g.sln("")
	}
	atk = atk - player.Def
	if atk < 1 {
		g.lln("  `%** `0" + op.Name + "`2 misses you Completely! `%**")
		return
	}
	if player.Kids > 5 && ((g.random(100) + 1) < (player.Kids - 5)) {
		if g.random(2) == 1 {
			son = "son"
			his = "him"
		} else {
			son = "duaghter"
			his = "her"
		}
		player.Kids--
		g.lln("  `0Your enemy is startled - someone is attacking him from behind!")
		// DIFF: Originally, this had an off-by-one where hold and hit were set the same...
		tmp := op.HP / 2
		op.HP -= tmp
		g.sln("")
		g.moreNomail()
		g.lln("  `%** `2Your `0" + son + "`2 hits " + op.Name + " `2for `4" + prettyInt(tmp) + " `2points of damage! `%**")
		g.sln("")
		g.moreNomail()
		g.lln("  `0" + op.Name + " `2points its `4" + op.Weapon + " `2at " + his + "!")
		g.sln("")
		g.moreNomail()
		g.lln("  `2You jump to intercept the blow - but you are too late.  Your child")
		g.lln("  is struck down.  Your `0" + son + "`2 dies in your arms.")
		g.sln("")
		g.moreNomail()
	} else {
		if pfight {
			his = pick(op.Sex == "M", "his", "her")
			g.lw("  `4** `0" + op.Name + " `2hits with " + his + " `0" + op.Weapon + " `2for ")
		} else {
			g.lw("  `4** `0" + op.Name + " `2hits with its `0" + op.Weapon + " `2for ")
		}
		if player.LightShield {
			atk = atk / 2
		}
		player.HP -= atk
		if player.HP < 1 {
			player.Dead = true
		}
		g.foreground(4)
		g.lln(prettyInt(atk) + " `2damage! `4**")
		if player.HasFairy && player.Dead {
			g.sln("")
			g.lln("  `%YOU FEEL A BUZZING IN YOUR POUCH!")
			g.sln("")
			g.moreNomail()
			g.lln("  `2Knowing you are too weak to go on, you decide to release the")
			g.sln("  fairy in your pocket.")
			g.sln("")
			g.moreNomail()
			g.lln("  `0The tiny thing rises in the air - And blows you a kiss!")
			g.sln("")
			g.lln("  `%YOU FEEL MUCH BETTER!")
			g.sln("")
			g.moreNomail()
			player.HasFairy = false
			player.Dead = false
			player.HP = player.HPMax
		}
	}
	// DIFF: You could horse-kill someone after you die...
	if player.HP > 0 && player.Horse && !pfight && g.random(25) == 10 {
		if op.IsDragon {
			g.foreground(10)
			g.sln("")
			g.sln("In an unexpected move, the Dragon makes a quick swipe towards you!")
			g.sln("")
			g.moreNomail()
			g.sln("Your horse moves its huge body to intercept the attack")
			g.sln("")
			g.sln("You're horse is killed by the Dragon's attack!")
			// DIFF: Horse wasn't killed before...
			player.Horse = false
		} else {
			g.sln("")
			op.HP = 0
			player.Horse = false
			g.lln("  `0\"Prepare to die, fool!\" `2" + op.Name + " `2screams.")
			g.sln("")
			g.moreNomail()
			g.lln("  `2He takes a `)Death Crystal`2 from his cloak and throws it at")
			g.sln("  you!")
			g.sln("")
			g.moreNomail()
			g.lln("  `2Your horse moves its huge body to intercept the crystal.")
			g.sln("")
			g.lln("  `4YOUR HORSE IS VAPORIZED!")
			g.sln("")
			g.moreNomail()
			g.lln("  `2Tears of anger flow down your cheeks.  Your valiant steed must")
			g.lln("  be avenged.")
			g.sln("")
			g.moreNomail()
			g.lln("  `%YOU PUMMEL `0" + op.Name + " `%WITH BLOWS!")
			g.sln("")
			g.lln("  `2A few seconds later, your adversary is dead.")
			g.sln("")
			g.moreNomail()
			g.lln("  `0You bury your horse in a small clearing.  The best")
			g.sln("  friend you ever had.")
			g.sln("")
		}
	}
}

// battlePrompt is battle_prompt().
func (g *Game) battlePrompt(op *Player) {
	player := g.player
	g.sln("")
	g.lln("  `2Your Hitpoints : `0" + prettyInt(player.HP))
	g.lln("  `2" + op.Name + "`2's Hitpoints : `0" + prettyInt(op.HP))
	g.sln("")
	g.lln("  `2(`5A`2)ttack")
	g.lln("  (`5S`2)tats")
	g.lln("  (`5R`2)un")
	if player.FairyLore {
		g.lln("  `2(`5H`2)eal")
	}
	if player.LevelW > 0 || player.LevelM > 0 || player.LevelT > 0 {
		g.sln("")
	}
	if player.LevelW > 0 {
		g.lln("  `2(`0D`2)`0eath Knight Attack (`%" + itoa(player.LevelW) + "`0)")
	}
	if player.LevelM > 0 {
		g.lln("  `2(`0M`2)`0ystical Skills     (`%" + itoa(player.LevelM) + "`0)")
	}
	if player.LevelT > 0 {
		g.lln("  `2(`0T`2)`0hieving Skills     (`%" + itoa(player.LevelT) + "`0)")
	}
	g.sln("")
	g.lw("  `2Your command, `0" + player.Name + "`2?  [`5A`2] : ")
}

// battlePromptLines is battle_prompt_lines(): how many lines
// battle_prompt() takes.
// DIFF: This is new for the training...
func (g *Game) battlePromptLines() int {
	player := g.player
	ret := 8
	if player.FairyLore {
		ret++
	}
	if player.LevelW > 0 || player.LevelM > 0 || player.LevelT > 0 {
		ret++
	}
	if player.LevelW > 0 {
		ret++
	}
	if player.LevelM > 0 {
		ret++
	}
	if player.LevelT > 0 {
		ret++
	}
	return ret
}

// battleKillBonus is the gem/gold roll shared by handle_hit and do_attack
// when a blow stronger than the player's strength kills.
func (g *Game) battleKillBonus(op *Player, pfight bool) {
	g.sln("")
	if !pfight {
		g.lln("  " + op.Death)
	}
	// DIFF: In 4.07, this was random(2) so more gold never happened.
	// rand(1,3) is what's used in the PHP version, and what we use here.
	switch g.random(3) {
	case 0:
		g.sln("")
		g.lln("  `2You find a `%Gem`2!")
		g.player.Gem++
	case 1:
		g.sln("")
		g.sln("  You find more gold than expected!")
		op.Gold *= 2
	}
}

// handleHit is handle_hit(): a special skill's blow lands.
func (g *Game) handleHit(atk int64, op *Player, pfight bool) {
	g.sln("")
	g.foreground(2)
	if pfight {
		atk -= op.Def
	}
	if atk < 1 {
		atk = 1
	}
	g.lln("  `2You hit `0" + op.Name + "`2 for `0" + prettyInt(atk) + " `2damage!")
	op.HP -= atk
	if atk > g.player.Str && op.HP < 1 {
		// DIFF: The death line was unconditional in original...
		g.battleKillBonus(op, pfight)
	}
	if op.HP > 0 {
		g.enemyAttack(op, pfight)
	}
}

// doAttack is do_attack(): a normal attack.  It returns text for the
// victim's mail in a player fight.
func (g *Game) doAttack(op *Player, pfight bool) string {
	player := g.player
	ret := ""

	g.sln("")
	atk := g.random(player.Str/2) + player.Str/2
	c := g.random(10) + 1
	if player.Amulet {
		c += 2
	}
	if pfight {
		atk -= op.Def
	}
	if player.Amulet && atk < 1 {
		atk = 1
	}
	if atk < 1 {
		g.sln("")
		g.lw("  `2You miss `0" + op.Name + "`2 completely!")
		g.sln("")
		g.enemyAttack(op, pfight)
		return ""
	}
	if c > 9 {
		atk *= 3
		g.sln("")
		g.foreground(15)
		g.sln("  **POWER MOVE**")
		if pfight {
			ret += "\n  `0" + player.Name + " `2did a power move for `4" + prettyInt(atk) + " `2damage!"
		}
	}
	g.sln("")
	g.lln("  `2You hit `0" + op.Name + "`2 for `0" + prettyInt(atk) + " `2damage!")
	op.HP -= atk
	if atk > player.Str && op.HP < 1 {
		g.battleKillBonus(op, pfight)
	}
	if op.HP > 0 {
		g.enemyAttack(op, pfight)
	}
	return ret
}

// tryRunning is try_running().
func (g *Game) tryRunning(op *Player, pfight, cantRun bool) string {
	if op.IsArena {
		g.sln("  You cannot run from an Arena!!!  You came here to prove")
		g.sln("  your worth, and that's what you are going to do.")
		g.battlePrompt(op)
		return ""
	}

	if cantRun || (g.random(9) == 1) {
		g.sln("")
		g.lln("  `0" + op.Name + " `2sees you!")
		g.enemyAttack(op, pfight)
		return ""
	}

	g.player.RanAway = true
	if pfight {
		g.sln("")
		g.lln("  `2You barely manage to escape!  `0" + op.Name + "`2 laughs as you")
		g.sln("  scurry away.")
		return "\n  `0" + g.player.Name + "`2 has run away like a scared rat!"
	}

	if op.IsDragon {
		g.sln("")
		g.lln("  `2You barely flip out of the way, as the `4Dragon`2 breathes huge")
		g.sln("  amounts of fire where you were a second ago!  You run towards")
		g.sln("  the forest, screaming all the way!")
	}
	return ""
}

// skillHonor is the honour check at the start of use_death_knight,
// use_thief_skill and use_mystical_skill.  It returns false when the skill
// may not be used.
func (g *Game) skillHonor(op *Player, pfight bool) bool {
	if pfight && g.player.Level >= op.Level {
		g.sln("  Your honor stops you from using the more unorthodox methods of ")
		g.sln("  battle.")
		g.sln("")
		g.battlePrompt(op)
		return false
	}
	if op.IsArena {
		g.sln("  Your honor stops you from using the more unorthodox methods of ")
		g.sln("  battle against your teacher.")
		g.sln("")
		g.battlePrompt(op)
		return false
	}
	if pfight && g.player.Level < op.Level {
		g.sln("  In a situation like this, you need every advantage you can get.")
		g.sln("")
	}
	return true
}

// useDeathKnight is use_death_knight().
func (g *Game) useDeathKnight(op *Player, pfight bool) {
	player := g.player

	g.sln("")
	g.sln("")
	if !g.skillHonor(op, pfight) {
		return
	}
	g.lln("  `%** ULTRA POWERFUL MOVE **")
	g.sln("")
	g.foreground(2)
	switch g.random(7) {
	case 0:
		g.lln("  In a swift move you have `0" + op.Name + "`2 by the neck and begin to")
		g.sln("  exert huge amounts of pressure.  You hear a sickening crack.")
	case 1:
		g.lln("  In a spectacular move, you drive your weapon up between `0" + op.Name + "`2's")
		g.lln("  legs.  You grin as your " + player.Weapon + " crunches noisily.")
	case 2:
		g.sln("  Seeing an opening, you feel your tendons strain as you swing your")
		g.lln("  " + player.Weapon + ". `0" + op.Name + "'s`2 left arm falls to the ground,")
		g.sln("  a gout of blood pulsating at the stump.")
	case 3:
		g.sln("  You duck an uncontrolled swing, and retaliate by slicing a wide gash")
		g.lln("  under `0" + op.Name + "'s`2 belly.  You feel sick as steamy intestines")
		g.sln("  slither out into a pile at your enemy's feet.")
	case 4:
		g.lln("  In a scream of rage, you brandish your " + player.Weapon + " wildly.  The")
		g.lln("  surprised `0" + op.Name + "`2 misses a parry, and you are able to send your")
		g.sln("  enemy's nose into the air with a smooth blow.")
	case 5:
		g.lln("  After exchanging blows & blocks for nearly a minute, `0" + op.Name)
		g.foreground(2)
		g.sln("  gets too fancy, and as a result an attempt to jump a low swing results")
		g.sln("  in the severing of your enemy's left leg.  His remaining leg slips in")
		g.sln("  blood.")
	case 6:
		g.lln("  You swing your " + player.Weapon + " as hard as you can, intending to")
		g.sln("  drive it into the fiend's side.  Instead he ducks, and your blade slides")
		g.sln("  right into the top part of his skull.  He stares at you in horror.  The")
		g.sln("  creature's warm brain slides smoothly from the cloven skull and lands")
		g.sln("  with a soft plop on your shoe.  You kick it away, disgusted.")
	}

	atk := g.random(player.Str/2) + player.Str/2
	if g.Settings.DKBoost {
		atk = int64(float64(atk) * 3.3)
	} else {
		atk *= 3
	}
	g.handleHit(atk, op, pfight)
	player.LevelW--
}

// useThiefSkill is use_thief_skill().
func (g *Game) useThiefSkill(op *Player, pfight bool) {
	player := g.player

	g.sln("")
	g.sln("")
	if !g.skillHonor(op, pfight) {
		return
	}
	g.lln("  `%** ULTRA SNEAKY MOVE **")
	g.sln("")
	g.foreground(2)
	switch g.random(7) {
	case 0:
		g.lln("  You point behind your enemy, and yell \"Whats that?!\"  `0" + op.Name)
		g.lln("  `2turns around stupidly, before the fool realizes its error, there are")
		g.sln("  two daggers planted in its back.")
	case 1:
		g.lln("  You suddenly find `0" + op.Name + "'s`2 right eye very unnatractive.")
		g.sln("  In a smooth motion you slide a dagger from your boot into the air. End")
		g.sln("  over end it flies, finding its mark.  Your enemy's right eye \"pops\".")
	case 2:
		g.sln("  A devious plan enters your mind.  You fall to the ground, clutching")
		g.lln("  at your chest.  The dumbfounded `0" + op.Name + "`2 drops his guard and")
		g.sln("  leans over you.  You scream \"Surprise you hoochie fiend!\" and drive your")
		g.lln("  " + player.Weapon + " into his neck.  It's nice to be covered with someone")
		g.sln("  else's blood for once.")
	case 3:
		g.sln("  You duck an uncontrolled swing, and retaliate by slicing a wide gash")
		g.lln("  under `0" + op.Name + "'s`2 belly.  You feel sick as steamy intestines")
		g.sln("  slither out into a pile at your enemy's feet.")
	case 4:
		g.lln("  In a scream of rage, you brandish your " + player.Weapon + " wildly.  The")
		g.lln("  surprised `0" + op.Name + "`2 misses a parry, and you are able to send your")
		g.sln("  enemy's nose into the air with a smooth blow.")
	case 5:
		g.lln("  After exchanging blows & blocks for nearly a minute, `0" + op.Name)
		g.lln("  `2gets too fancy, and as a result an attempt to jump a low swing results")
		g.sln("  in the severing of your enemy's left leg.  His remaining leg slips in")
		g.sln("  blood.")
	case 6:
		g.lln("  You swing your " + player.Weapon + " as hard as you can, intending to")
		g.sln("  drive it into the fiend's side.  Instead he ducks, and your blade slides")
		g.sln("  right into the top part of his skull.  He stares at you in horror.  The")
		g.sln("  creature's warm brain slides smoothly from the cloven skull and lands")
		g.sln("  with a soft plop on your shoe.  You kick it away, disgusted.")
	}

	atk := g.random(player.Str/2) + player.Str/2
	atk *= 3
	g.handleHit(atk, op, pfight)
	player.LevelT--
}

// useMysticalSkill is use_mystical_skill().
func (g *Game) useMysticalSkill(op *Player, pfight bool) {
	player := g.player

	tooTired := func() {
		g.foreground(10)
		g.sln("  Your face contorts in concentration.  You reach deep within yourself")
		g.sln("  for the power, but gasp as you realize is just isn't there.  You need")
		g.sln("  rest before your mind will be ready for this.")
		g.sln("")
		g.more()
	}

	g.sln("")
	g.sln("")
	g.lln("  `%** MYSTICAL SKILLS **")
	g.sln("")
	if !g.skillHonor(op, pfight) {
		return
	}
	g.sln("")
	g.foreground(2)
	switch g.random(7) {
	case 0:
		g.lln("  `2You quickly decide which mystical skill to use.")
	case 1:
		g.lln("  `2Your enemy gives you a chance to contemplate your next action.")
	case 2:
		g.lln("  `2You decide a little magic might change the outcome of this battle.")
	case 3:
		g.lln("  `2You struggle to keep your anger under control.")
	case 4:
		g.lln("  `2Your mind carefully goes over what you have learned.")
	case 5:
		g.lln("  `2You feel power dancing in your mind, maybe it's time to use it.")
	case 6:
		g.lln("  `0" + op.Name + " `2looks dumbfounded as you sheathe your " + player.Weapon + ".")
	}
	g.sln("")
	g.lln("  `5(`#P`5)inch Real Hard                 `5(`%1`5)")
	valid := "P"
	if player.LevelM > 3 && player.SkillM > 3 {
		g.lln("  `5(`#D`5)isappear                       `5(`%4`5)")
		valid += "D"
	}
	if player.LevelM > 7 && player.SkillM > 7 {
		g.lln("  `5(`#H`5)eat Wave                       `5(`%8`5)")
		valid += "H"
	}
	if player.LevelM > 11 && player.SkillM > 11 {
		g.lln("  `5(`#L`5)ight Shield                    `5(`%12`5)")
		valid += "L"
	}
	if player.LevelM > 15 && player.SkillM > 15 {
		g.lln("  `5(`#S`5)hatter                         `5(`%16`5)")
		valid += "S"
	}
	if player.LevelM > 19 && player.SkillM > 19 {
		g.lln("  `5(`#M`5)ind Heal                       `5(`%20`5)")
		valid += "M"
	}
	g.sln("")

	g.lw("  `5You Have `%" + itoa(player.LevelM) + "`5 Use Points.  Choose.  [`#Nothing`5] : ")
	ch := strings.ToUpper(g.getkey())
	if !strings.Contains(valid, ch) {
		ch = "Nothing"
	}
	g.sln(ch)
	if ch == "P" {
		if player.LevelM < 1 {
			tooTired()
			return
		}
		g.lln("  `0You whisper the word.  You smile as `2" + op.Name + "`0 screams out in pain.")
		atk := g.random(player.Str/2) + player.Str/2
		atk = (atk * 2) - (atk*2)/4
		g.handleHit(atk, op, pfight)
		player.LevelM--
		// TODO: This seems to be the only one that redraw the menu?
		if player.HP > 0 && op.HP > 0 {
			g.battlePrompt(op)
		}
	}
	if ch == "D" {
		if player.LevelM < 4 {
			tooTired()
			return
		}
		g.foreground(10)
		g.sln("  You imagine yourself being in a different part of the forest, and begin")
		g.sln("  to concentrate.  The next instant you are standing in a cool glade, nowhere")
		g.lln("  near your enemy.  You almost laugh remembering `2" + op.Name + "'s`0 befuddled")
		player.LevelM -= 4
		player.RanAway = true
	}
	if ch == "H" {
		if player.LevelM < 8 {
			tooTired()
			return
		}
		g.lln("  `0Your face contorts in anger.  How DARE `2" + op.Name + "`0 presume that IT")
		g.sln("  can beat YOU?  The air in front of you begins to shimmer.  A few seconds")
		g.sln("  later, the fiend is engulfed by flames.")
		player.LevelM -= 8
		atk := g.random(player.Str/2) + player.Str/2
		atk = (atk * 2) + (atk*2)/4
		g.handleHit(atk, op, pfight)
	}
	if ch == "L" {
		if player.LevelM < 12 {
			tooTired()
			return
		}
		g.lln("  `0You look into the clouds.  A beam of light covers your face.  You hear")
		g.lln("  `2" + op.Name + "`0 scream as the light expands to cover your whole body.")
		g.sln("  The beam disappears, but you are still glowing.")
		player.LevelM -= 12
		player.LightShield = true
	}
	if ch == "S" {
		if player.LevelM < 16 {
			tooTired()
			return
		}
		g.lln("  `0You visualize the bones in your enemy, then imagine them imploding.")
		g.lln("  `2" + op.Name + "`0 pierces the air with a scream of horror that makes your")
		g.lln("  blood run cold.  Even so, you feel no remorse.")
		player.LevelM -= 16
		atk := g.random(player.Str/2) + player.Str/2
		atk = (atk * 4) + (atk*4)/4
		g.handleHit(atk, op, pfight)
	}
	if ch == "M" {
		if player.LevelM < 20 {
			tooTired()
			return
		}
		g.lln("  `0You look down at your bloody body.  This will never do.  With an")
		g.lln("  inhuman scream, you will yourself to be healed.  `2" + op.Name)
		g.lln("  `0gasps as missing pieces of your flesh crawl back into place.")
		g.sln("")
		g.lln("  `%YOU ARE HEALED.`2")
		player.LevelM -= 20
		player.HP = player.HPMax
	}
}

// battle is battle(): fight op until someone dies or the player runs.
// It returns text for the victim's mail in a player fight (pfight).
func (g *Game) battle(op *Player, pfight, cantRun bool) string {
	player := g.player
	mail := ""

	tmp := g.random(99) + 1
	if pfight {
		if op.Level > player.Level {
			if tmp > 60 {
				tmp = 95
			}
		}
	}

	if tmp > 90 {
		g.lln("  `0" + op.Name + " `2surprises you.")
		g.enemyAttack(op, pfight)
	} else {
		g.sln("  Your skill allows you to get the first strike.")
	}

	player.RanAway = false
	for {
		if !player.Dead && player.HP > 0 {
			g.battlePrompt(op)
		}
		if player.Dead {
			break
		}
		if player.HP < 1 {
			break
		}
		// NOTE: This prompt is being left non-normalized...
		ch := strings.ToUpper(g.getkey())
		if ch == "\r" {
			ch = "A"
		}
		g.sw(ch)
		switch ch {
		case "Q":
			g.sln("")
			g.sln("")
			g.sln("  You are in combat!  Try running.")
			g.sln("")
		case "A":
			mail += g.doAttack(op, pfight)
		case "R":
			mail += g.tryRunning(op, pfight, cantRun)
		case "S":
			g.showStats()
		case "L":
			g.sln("")
			g.sln("")
			g.sln("  What?!  You want to fight two at once?")
			g.sln("")
		case "D":
			if player.LevelW < 1 {
				g.battleNoSkill(player.SkillW)
			} else {
				g.useDeathKnight(op, pfight)
			}
		case "T":
			if player.LevelT < 1 {
				g.battleNoSkill(player.SkillT)
			} else {
				g.useThiefSkill(op, pfight)
			}
		case "M":
			if player.LevelM < 1 {
				g.battleNoSkill(player.SkillM)
			} else {
				g.useMysticalSkill(op, pfight)
			}
		case "H":
			g.sln("")
			g.sln("")
			if !player.FairyLore {
				g.sln("  You are in combat, and they don't make house calls!")
			} else {
				g.foreground(10)
				g.sln("  Thanks to your fairy lore training today, you are able to heal yourself.")
				g.lw("  `2Do it? [`%Y`2] : `%")
				ch = strings.ToUpper(g.getkey())
				if ch != "N" {
					ch = "Y"
				}
				g.sln(ch)
				g.foreground(10)
				if ch == "Y" {
					player.HP = player.HPMax
					player.FairyLore = false
					g.sln("  Concentrating only on your wounds, they heal themselves!")
				}
			}
			g.sln("")
		}
		if player.RanAway {
			break
		}
		if player.Dead {
			break
		}
		if player.HP < 1 {
			break
		}
		if op.HP < 1 {
			break
		}
	}

	player.RanAway = false
	return mail
}

// battleNoSkill is battle()'s answer when a skill has no use points left.
func (g *Game) battleNoSkill(skill int64) {
	g.sln("")
	g.sln("")
	if skill < 1 {
		g.sln("  You don't know any at the moment!")
	} else {
		g.sln("  You will need rest before you can use those again.")
	}
	g.sln("")
}

// customSaying is custom_saying(): the winner's or loser's words for the
// daily news.
func (g *Game) customSaying() string {
	ret := ""

	g.sln("")
	g.lw("  `2Say something to the press? [`0N`2] : `%")
	ch := strings.ToUpper(g.getkey())
	if ch != "Y" {
		ch = "N"
	}
	g.sw(ch)
	if ch == "N" {
		return ""
	}
	for {
		g.sln("")
		g.sln("")
		g.lln(" `2Share your feelings now.. (Max 50 char!)")
		g.sln("                                             Max")
		g.sln("                                              ^")
		g.lw("  `0\"`2")
		ret = g.consoleGetstr(GetStrOpts{})
		ret = cleanStr(ret)
		if len(ret) > 50 {
			g.sln("")
			g.lln("  `2The press thinks you are boring and leaves!")
			g.sln("")
			g.moreNomail()
			return ""
		}
		if len(ret) < 5 {
			g.sln("")
			g.lln("  `2That just isn't interesting enough.  Sorry man.")
			g.sln("")
			g.moreNomail()
			return ""
		}
		ret = "`0  \"" + ret + "\"`2 "
		g.sln("")
		g.lln("  `2(`0H`2)appy")
		g.lln("  `2(`0S`2)ad")
		g.lln("  `2(`0M`2)oaning")
		g.lln("  `2(`0C`2)ursing")
		g.lln("  `2(`0E`2)ffing Mad")
		g.sln("")
		g.lw("  Your emotional state? `2[`0H`2] : `%")
		ch = strings.ToUpper(g.getkey())
		if !strings.Contains("HSMCE", ch) {
			ch = "H"
		}
		g.sw(ch)
		switch ch {
		case "H":
			ret += "laughs `%" + g.player.Name + "`2."
		case "S":
			ret += "`%" + g.player.Name + "`2 declares sadly."
		case "M":
			ret += "moans `%" + g.player.Name + "`2."
		case "C":
			ret += "curses `%" + g.player.Name + "`2."
		case "E":
			ret += "screams `%" + g.player.Name + "`2 in anger."
		}
		g.sln("")
		g.sln("")
		g.lln(ret)
		g.sln("")
		g.lw("  `2Does this look ok? [`0Y`2] :`% ")
		ch = strings.ToUpper(g.getkey())
		if ch != "N" {
			ch = "Y"
		}
		g.sln(ch)
		if ch != "N" {
			g.sln("")
			g.lln("  Your comment has been noted.`2")
			return ret
		}
	}
}

// say is say(): log header and a saying to the daily news.  fname is the
// name of a saying file (badsay.lrd, goodsay.lrd, normsay.lrd); when force
// is empty a random saying is taken from it, with `e and `g replaced by
// the enemy's and the player's names.
func (g *Game) say(fname string, enemy *Player, typ int, header, force string) {
	if force == "" {
		lines, err := g.assetLines(fname)
		if err != nil {
			g.fatalf("Unable to open %s", fname)
		}
		readln := func(i int) string {
			if i < len(lines) {
				return lines[i]
			}
			return "null" // File.readln() past the end is null.
		}
		c, _ := jsParseInt(readln(1))
		c = g.random(c)
		for i := int64(0); i < c; i++ {
			force = "`2  " + readln(int(2+i))
		}
		var b strings.Builder
		for i := 0; i < len(force); i++ {
			if force[i] == '`' && i+1 < len(force) && strings.IndexByte("neg", force[i+1]) >= 0 {
				b.WriteString(g.sayCode(force[i+1], enemy, typ))
				i++
				continue
			}
			b.WriteByte(force[i])
		}
		force = b.String()
	}
	g.logLine(header + "\n" + force)
}

// sayCode is say()'s replacement for `n, `e and `g.
func (g *Game) sayCode(c byte, enemy *Player, typ int) string {
	if c == 'n' {
		return "\n"
	}
	switch typ {
	case 1, 3:
		if c == 'e' {
			return "`0" + enemy.Name + "`2"
		}
		if c == 'g' {
			return "`5" + g.player.Name + "`2"
		}
	case 2:
		if c == 'e' {
			return "`5" + enemy.Name + "`2"
		}
		if c == 'g' {
			return "`0" + g.player.Name + "`2"
		}
	}
	return "undefined"
}

// badSay is bad_say().
func (g *Game) badSay(enemy *Player, header string) {
	g.say("badsay.lrd", enemy, 3, header, g.customSaying())
}

// goodSay is good_say().
func (g *Game) goodSay(enemy *Player, header string) {
	g.say("goodsay.lrd", enemy, 2, header, g.customSaying())
}
