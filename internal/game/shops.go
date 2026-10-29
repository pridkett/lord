package game

import (
	"fmt"
	"strings"
)

// The town's weapon and armour shops and raise_class() (lord.js lines
// 10016-10791).

// kingArthurs is king_arthurs(): the weapons shop.
func (g *Game) kingArthurs() {
	over := false
	var ch string
	p := g.player

	menu := func() {
		if !p.Expert {
			g.lrdfile("ARTHUR")
		}
	}

	prompt := func() {
		g.sln("")
		g.lln("  `2Current weapon: `0" + p.Weapon)
		g.lln("  `5King Arthur's Weapons `8(B,S,Y,R)  (? for menu)")
	}

	g.sclrscr()
	menu()

	for {
		prompt()
		g.commandPrompt()
		ch = strings.ToUpper(g.getkey())
		g.sln(ch)
		switch ch {
		case "L":
			g.sclrscr()
			g.lrdfile("BUYWEP")
			g.kingArthursDisplayWeapons()
		case "Y":
			g.showStats()
			menu()
		case "S":
			g.kingArthursSellWeapon()
		case "B":
			g.kingArthursBuyWeapon()
		case "R", "Q", "\r":
			over = true
		// DIFF (lord.js): ' ' would cycle, but not call prompt()...
		case "?":
			if p.Expert {
				g.lrdfile("ARTHUR")
			}
			menu()
		}
		if over {
			break
		}
	}
}

// kingArthursDisplayWeapons is king_arthurs' display_weapons().
func (g *Game) kingArthursDisplayWeapons() {
	mc := g.morechk
	g.morechk = false
	for i := int64(1); i < 16; i++ {
		w := g.getWeapon(i)
		l := w.Name
		pr := prettyInt(w.Price)
		n := fmt.Sprintf("%2d", i)
		for dispLen(l)+len(pr) < 42 {
			l += "."
		}
		g.con.GotoXY(17, int(i)+3)
		g.lln("`2" + n + ". " + l + "`0" + pr)
	}
	g.morechk = mc
}

// kingArthursSellWeapon is king_arthurs' sell_weapon().
func (g *Game) kingArthursSellWeapon() {
	p := g.player
	var price int64

	g.lln("`c  `%King Arthurs Weapons")
	g.foreground(2)
	g.sln("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
	g.sln("")
	if p.WeaponNum == 0 {
		g.lln("  `2\"`0What the...?!!`2\" the stout man shouts. `2\"`0You don't have")
		g.lln("  a weapon to sell!`2\"")
		g.sln("")
		g.sln("")
		g.more()
		return
	}
	oldw := g.getWeapon(p.WeaponNum)
	mult := p.Level * p.Cha * p.Level
	if mult > 0 && mult < 65530 {
		price = int64(float64(oldw.Price)/2 + float64(g.random(mult)))
	} else {
		price = int64(float64(oldw.Price)/2 + float64(g.random(65535)))
	}
	if float64(price) > float64(oldw.Price)-(float64(price)/3) {
		price = oldw.Price - oldw.Price/3
	}
	g.lln("  `2\"`0Hmmm I will buy your `%" + p.Weapon + "`0 for `%" + prettyInt(price) + "`0, Agreed?`2\"")
	g.sln("")
	g.lw("   Sell it?  [`0N`2] : `%")
	ich := strings.ToUpper(g.getkey())
	if ich != "Y" {
		ich = "N"
	}
	g.sln(ich)
	if ich != "Y" {
		g.sln("")
		g.lln("  `2\"`0You don't want to sell?!  Fine!  I don't want your stinken' weapon!`2\"")
		g.sln("")
		g.moreNomail()
		return
	}
	g.sln("")
	g.lln("  `2\"`0Great!`2\" The fat man takes your weapon, and gives you the")
	g.sln("  money.")
	p.WeaponNum = 0
	neww := g.getWeapon(p.WeaponNum)
	p.Weapon = neww.Name
	p.Gold += price
	if p.Gold > 2000000000 {
		p.Gold = 2000000000
		g.lln("Wow, you have a lot of money!")
	}
	p.Str -= oldw.Num
	if p.Str < 5 {
		p.Str = 5
	}
	g.moreNomail()
}

// kingArthursStrNeeded is king_arthurs' str_needed().
func (g *Game) kingArthursStrNeeded(weap int64) int64 {
	ret := int64(10)

	if !g.Settings.ShopLimit {
		return 0
	}
	if weap < 3 {
		return 0
	}
	for i := int64(1); i < weap-1 && i < int64(len(trainerStats)); i++ {
		t := g.getTrainer(i)
		ret += t.StrGained
	}
	return ret
}

// kingArthursBuyWeapon is king_arthurs' buy_weapon().
func (g *Game) kingArthursBuyWeapon() {
	p := g.player

	g.sclrscr()
	if !p.Expert {
		g.sclrscr()
		g.lrdfile("BUYWEP")
		g.kingArthursDisplayWeapons()
	}
	g.sln("")
	g.sln("")
	g.lln("  `2(`0Gold: `%" + prettyInt(p.Gold) + "`2)  (`00 to exit`2)")
	g.lw("  `0Number Of Weapon `2: `%")
	n, ok := jsParseInt(g.consoleGetstr(GetStrOpts{Len: 2, Integer: true}))
	if !ok {
		n = 0
	}
	if n > 0 && n < 16 {
		oldw := g.getWeapon(p.WeaponNum)
		neww := g.getWeapon(n)
		need := g.kingArthursStrNeeded(n)
		g.sclrscr()
		g.sln("")
		g.sln("")
		g.foreground(15)
		g.sln("  King Arthurs Weapons")
		g.foreground(2)
		g.sln("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
		g.sln("")
		g.lln("  `2\"`0Hmmm I will sell you my FAVORITE `%" + neww.Name + "`0 for `%" + prettyInt(neww.Price) + " `0gold!`2\"")
		g.sln("")
		g.sln("")
		g.lln("  `2Note: It takes `%" + prettyInt(need) + " `2strength points to weild this weapon.")
		g.lln("  `2You currently have `%" + prettyInt(p.Str-oldw.Num) + " `2strength points.")
		g.sln("")
		g.lw("  `2Buy it?  [`0N`2] : `%")
		ich := strings.ToUpper(g.getkey())
		if ich != "Y" {
			ich = "N"
		}
		g.sln(ich)
		if ich == "N" {
			g.sln("")
			g.lln("  `2\"`0Fine..You will come back...`2\" the man grunts.")
			g.sln("")
			g.moreNomail()
			return
		}
		if p.Str < need {
			g.sln("")
			g.lln("  `2\"`0You silly fool! You aren't strong enough to carry")
			g.lln("  that weapon!`2\"")
			g.sln("")
			g.moreNomail()
			return
		}
		if p.WeaponNum > 0 {
			g.sln("")
			g.lln("  `2\"`0You fool!  You already have a weapon, and you can't carry")
			g.lln("  two!`2\"  You realize he is right.")
			g.sln("")
			g.moreNomail()
			return
		}
		if p.Gold < neww.Price {
			g.sln("")
			g.lln("  `2\"`0You stupid fool!  You don't have that much gold!")
			g.lln("  I knew you were up to no good the moment I saw you!`2\"")
			g.sln("")
			g.sln("")
			// DIFF (lord.js): This more wasn't here...
			g.moreNomail()
			return
		}
		g.sln("")
		g.lln("  `2\"`0Great!`2\" The fat man takes your money, and gives you the")
		g.sln("  weapon.")
		p.WeaponNum = n
		p.Weapon = neww.Name
		p.Gold -= neww.Price
		p.Str += neww.Num
		if p.Str > 32000 {
			p.Str = 32000
		}
		g.sln("")
		g.moreNomail()
	}
}

// abdulsArmour is abduls_armour(): the armour shop.
func (g *Game) abdulsArmour() {
	over := false
	var ch string
	p := g.player

	menu := func() {
		g.sclrscr()
		if !p.Expert {
			g.lrdfile("ABDUL")
		}
	}

	menu()

	for {
		g.sln("")
		g.lln("  `2Current armour: `0" + p.Arm)
		g.lln("  `5Abduls Armour `8(B,S,Y,R)  (? for menu)")
		g.commandPrompt()
		ch = strings.ToUpper(g.getkey())
		g.sln(ch)
		switch ch {
		case "L":
			g.sclrscr()
			g.lrdfile("BUYARM")
			g.abdulsArmourDisplayArmour()
		case "V", "Y":
			g.showStats()
		case "S":
			g.abdulsArmourSellArmour()
		case "B":
			g.abdulsArmourBuyArmour()
		case "R", "Q", "\r":
			over = true
		case "?":
			if p.Expert {
				g.lrdfile("ABDUL")
			}
			menu()
		}
		if over {
			break
		}
	}
}

// abdulsArmourDisplayArmour is abduls_armour's display_armour().
func (g *Game) abdulsArmourDisplayArmour() {
	mc := g.morechk
	g.morechk = false
	for i := int64(1); i < 16; i++ {
		a := g.getArmour(i)
		l := a.Name
		pr := prettyInt(a.Price)
		n := fmt.Sprintf("%2d", i)
		for dispLen(l)+len(pr) < 42 {
			l += "."
		}
		g.con.GotoXY(17, int(i)+3)
		g.lln("`2" + n + ". " + l + "`0" + pr)
	}
	g.morechk = mc
}

// abdulsArmourDefNeeded is abduls_armour's def_needed().
func (g *Game) abdulsArmourDefNeeded(arm int64) int64 {
	ret := int64(0)

	if !g.Settings.ShopLimit {
		return 0
	}
	if arm < 3 {
		return 0
	}
	for i := int64(1); i < arm-1 && i < int64(len(trainerStats)); i++ {
		t := g.getTrainer(i)
		ret += t.Def
	}
	return ret
}

// abdulsArmourBuyArmour is abduls_armour's buy_armour().
func (g *Game) abdulsArmourBuyArmour() {
	p := g.player

	if !p.Expert {
		g.sclrscr()
		g.lrdfile("BUYARM")
		g.background(0)
		g.abdulsArmourDisplayArmour()
	}
	g.sln("")
	g.sln("")
	g.lln("`2  (`0Gold: `%" + prettyInt(p.Gold) + "`2)  (`00 to Exit`2) ")
	g.lw("  `0Number Of Armour `2: `%")
	n, ok := jsParseInt(g.consoleGetstr(GetStrOpts{Len: 2, Integer: true}))
	if !ok {
		n = 0
	}
	if n > 0 && n < 16 {
		newa := g.getArmour(n)
		olda := g.getArmour(p.ArmNum)
		need := g.abdulsArmourDefNeeded(n)
		g.sclrscr()
		g.sln("")
		g.sln("")
		g.foreground(15)
		g.sln("  Abduls Armour Shop")
		g.foreground(2)
		g.sln("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
		g.sln("")
		g.lln("  \"`0Hmmm I will sell you a nice `%" + newa.Name + "`0 for `%" + prettyInt(newa.Price) + "`0.")
		g.lln("   Agreed, friend?`2\"")
		g.sln("")
		g.sln("")
		g.foreground(2)
		g.lln("  Note: It takes `%" + prettyInt(need) + "`2 defense points to wear this armor.")
		g.lln("  `2You currently have `%" + prettyInt(p.Def-olda.Num) + " `2defense points.")
		g.sln("")
		g.lw("   `2Buy it?  [`0N`2] : `%")
		ich := strings.ToUpper(g.getkey())
		if ich != "Y" {
			ich = "N"
		}
		g.sln(ich)
		g.sln("")
		switch {
		case ich == "N":
			g.lln("`2  \"`0Ok!  No rush!`2\" the girl smiles.")
		case p.Def < need:
			g.lln("  `2\"`0I'm sorry, but you are not strong enough to wear")
			g.lln("  that armor.`2\"")
		case p.ArmNum > 0:
			g.lln("  `2\"`0You already have armour, and you can't wear")
			g.lln("  two!`2\" You realize she is right.")
		case p.Gold < newa.Price:
			g.lln("  `2\"`0I'm sorry, but you seem to be lacking funds at the")
			g.lln("  moment.`2\"  the girl tells you.")
		default:
			g.lln("  `2\"`0Wonderful!`2\" The girl takes your money, and helps you")
			g.sln("  into your new armour.")
			p.Arm = newa.Name
			p.ArmNum = n
			p.Gold -= newa.Price
			p.Def += newa.Num
			// DIFF (lord.js): Used to compare to 3200 and set to 32000!
			if p.Def > 32000 {
				p.Def = 32000
			}
		}
		g.sln("")
		g.moreNomail()
	}
}

// abdulsArmourSellArmour is abduls_armour's sell_armour().
func (g *Game) abdulsArmourSellArmour() {
	p := g.player
	var price int64

	g.sclrscr()
	g.foreground(15)
	g.sln("  Abduls Armour Shop")
	g.foreground(2)
	g.sln("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
	g.sln("")
	if p.ArmNum == 0 {
		g.lln("  `2\"`0You silly kidder!!`2\" Paula laughs, \"`0You don't have")
		g.lln("  any armour to sell!`2\"  ")
		g.sln("")
		g.sln("")
		g.more()
		return
	}
	olda := g.getArmour(p.ArmNum)
	mult := p.Level * p.Cha * p.Level
	if mult >= 0 && mult <= 65530 {
		price = int64(float64(olda.Price)/2 + float64(g.random(mult)))
	} else {
		price = int64(float64(olda.Price)/2 + float64(g.random(65530)))
	}
	if float64(price) > float64(olda.Price)-(float64(price)/3) {
		price = olda.Price - olda.Price/3
	}
	g.lln("  `2\"`0Hmmm I will buy your `%" + p.Arm + " `0for `%" + prettyInt(price) + "`0 gold.")
	g.lln("  Agreed, friend?`2\"")
	g.sln("")
	g.lw("  `2Sell it?  [`0N`2] : `%")
	ich := strings.ToUpper(g.getkey())
	// DIFF (lord.js): This was after the print...
	if ich != "Y" {
		ich = "N"
	}
	g.sln(ich)
	if ich == "N" {
		g.sln("")
		g.lln("  `2\"`0Thats ok.  Your armour probably has sentimental value to you.`2\"")
		return
	}
	g.sln("")
	g.lln("  `2\"`0Good doing business with you!`2\" The girl takes your armour")
	g.sln("  and gives you the money.")
	p.Arm = "Nothing!"
	p.ArmNum = 0
	p.Gold += price
	if p.Gold > 2000000000 {
		p.Gold = 2000000000
		g.sln("Wow, you have a lot of money!")
	}
	p.Def -= olda.Num
	if p.Def < 1 {
		p.Def = 0
	}
	g.moreNomail()
}

// raiseClass is raise_class(): a lesson in the player's class skill.
func (g *Game) raiseClass() {
	p := g.player

	g.sln("")

	if p.SkillW > 39 && p.SkillM > 39 && p.SkillT > 39 {
		g.lln("  `%** `0YOU HAVE ALREADY MASTERED ALL SKILLS `%**")
		return
	}
	switch p.Clss {
	case 1:
		if p.SkillW > 39 {
			g.lln("  `%** `0YOU HAVE ALREADY MASTERED THIS CLASS `%**")
			g.sln("")
			return
		}
	case 2:
		if p.SkillM > 39 {
			g.lln("  `%** `0YOU HAVE ALREADY MASTERED THIS CLASS `%**")
			g.sln("")
			return
		}
	case 3:
		if p.SkillT > 39 {
			g.lln("  `%** `0YOU HAVE ALREADY MASTERED THIS CLASS `%**")
			g.sln("")
			return
		}
	}

	g.lln("  `%** `0YOUR CLASS SKILL IS RAISED BY ONE! `%**")
	g.sln("")
	switch p.Clss {
	case 1:
		p.SkillW++
		if g.Settings.OldSkillPoints {
			switch p.SkillW % 5 {
			case 0:
				g.lln("  `2You now have `0" + prettyInt(p.SkillW/4) + "`2 uses of Death Knight Skills a day.")
				p.LevelW++
				g.sln("")
				if p.SkillW == 40 {
					g.sln("  You have mastered The Death Knight Skills Completely.  You may choose to")
					g.sln("  learn a NEW skill now.")
					g.sln("")
					p.Clss = g.chooseProfession(false)
				} else {
					g.sln("  (five more lessons needed for next raise in uses per day)")
				}
			case 4:
				g.sln("  You need one more lesson to also raise your Death Knight Uses Per Day.")
			case 3:
				g.sln("  You need two more lesson to also raise your Death Knight Uses Per Day.")
			case 2:
				g.sln("  You need three more lessons to also raise your Death Knight Uses Per Day.")
			case 1:
				g.sln("  You need four more lessons to also raise your Death Knight Uses Per Day.")
			}
		} else {
			switch p.SkillW % 4 {
			case 0:
				g.lln("  `2You now have `0" + prettyInt(p.SkillW/4) + "`2 uses of Death Knight Skills a day.")
				p.LevelW++
				g.sln("")
				if p.SkillW == 40 {
					g.sln("  You have mastered The Death Knight Skills Completely.  You may choose to")
					g.sln("  learn a NEW skill now.")
					g.sln("")
					p.Clss = g.chooseProfession(false)
				} else {
					g.sln("  (four more lessons needed for next raise in uses per day)")
				}
			case 3:
				g.sln("  You need one more lesson to also raise your Death Knight Uses Per Day.")
			case 2:
				g.sln("  You need two more lessons to also raise your Death Knight Uses Per Day.")
			case 1:
				g.sln("  You need three more lessons to also raise your Death Knight Uses Per Day.")
			}
		}
	case 2:
		p.SkillM++
		g.sln("")
		g.lln("  `2You now have `0" + prettyInt(p.SkillM) + "`2 Mystical Skill points a day.")
		p.LevelM++
		switch p.SkillM % 4 {
		case 0:
			g.sln("")
			g.more()
			g.lln("`c  `%                    **  MYSTICAL INSTRUCTION **")
			g.sln("")
			g.foreground(2)
			g.sln("  An old man with a long white beard suddenly appears next to you.")
			g.sln("")
			g.lln("  `0\"You're ready for your next lesson " + pick(p.Sex == "F", "my girl!", "boy!"))
			g.sln("")
			g.more()
			if p.SkillM == 4 {
				g.lln("  `0\"You ever been chased by a monster that just wouldn't quit?  I have,")
				g.sln("   Muh wife!  Heehee! In a case like that there is only one thing to do!")
				g.sln("   Disappear!\"")
			}
			if p.SkillM == 8 {
				g.lln("  `0\"Ok, you're still new at this, but I think you're ready for a little")
				g.sln("   trick I call The Heat Wave.  This little 'beaut will blow a wind that")
				g.sln("   not only warms your enemy, it cooks him.\"")
			}
			if p.SkillM == 12 {
				g.lln("  `0\"Ok, pardner, lemmie give it to ya straight.  Sometimes yer enemy is")
				g.sln("   stronger then ya.  Ya need an advantage.  For instance, having a")
				g.sln("   protective Light Shield wrapped around ya.  Half damage!\"")
			}
			if p.SkillM == 16 {
				g.lln("  `0\"Ok, sometimes ya lose your temper, and you want to vent your anger")
				g.lln("   in a contructive way?  Am I right?  Causing someone's bones to shatter")
				g.lln("   is a great way to relieve stress.\"")
			}
			if p.SkillM == 20 {
				g.lln("  `0\"This is also your LAST lesson!\" `2the old man's face turns sober.  `0\"You")
				g.sln("   have finally shown enough control to master the most complicated thing.")
				g.sln("   Your body.  Healing yourself with your mind is an awesome power.\"")
			}
			g.sln("")
			g.lln("  `2The old man vanishes as quickly as he appeared.")
			if p.SkillM == 40 {
				g.sln("")
				g.sln("  You have mastered the Mystical skills Completely.  You may choose to")
				g.sln("  learn a NEW skill now.")
				g.sln("")
				p.Clss = g.chooseProfession(false)
			}
			g.sln("  (four more lessons needed to learn a new Mystical Skill)")
		case 3:
			g.sln("  You need only one more lesson to learn a new Mystical Skill.")
		case 2:
			g.sln("  You need two more lessons to learn a new Mystical Skill.")
		case 1:
			g.sln("  You need three more lessons to learn a new Mystical Skill.")
		}
	case 3:
		p.SkillT++
		if g.Settings.OldSkillPoints {
			switch p.SkillT % 5 {
			case 0:
				g.lln("  `2You now have `0" + prettyInt(p.SkillT/4) + "`2 uses of The Thieving Skills a day.")
				p.LevelT++
				g.sln("")
				if p.SkillT == 40 {
					g.sln("")
					g.sln("  You have mastered The Thieving Knight Skills Completely.  You may choose to")
					g.sln("  learn a NEW skill now.")
					g.sln("")
					p.Clss = g.chooseProfession(false)
				}
				g.sln("  (five more lessons needed for next raise in uses per day)")
			case 4:
				g.sln("  You need only one more lesson to also raise your Thieving Uses Per Day.")
			case 3:
				g.sln("  You need two more lesson to also raise your Thieving Uses Per Day.")
			case 2:
				g.sln("  You need three more lessons to also raise your Thieving Uses Per Day.")
			case 1:
				g.sln("  You need four more lessons to also raise your Thieving Uses Per Day.")
			}
		} else {
			switch p.SkillT % 4 {
			case 0:
				g.lln("  `2You now have `0" + prettyInt(p.SkillT/4) + "`2 uses of The Thieving Skills a day.")
				p.LevelT++
				g.sln("")
				if p.SkillT == 40 {
					g.sln("")
					g.sln("  You have mastered The Thieving Knight Skills Completely.  You may choose to")
					g.sln("  learn a NEW skill now.")
					g.sln("")
					p.Clss = g.chooseProfession(false)
				}
				g.sln("  (four more lessons needed for next raise in uses per day)")
			case 3:
				g.sln("  You need only one more lesson to also raise your Thieving Uses Per Day.")
			case 2:
				g.sln("  You need two more lessons to also raise your Thieving Uses Per Day.")
			case 1:
				g.sln("  You need three more lessons to also raise your Thieving Uses Per Day.")
			}
		}
	}
	g.sln("")
}
