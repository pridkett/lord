package game

import (
	"strconv"
	"strings"
)

// yeOldBank is ye_old_bank(): Ye Olde Bank.
func (g *Game) yeOldBank() {
	var ch string
	var amt int64
	var to int64
	p := g.player
	classList := []string{"stranger", "warrior", "magician", "thief"}

	menu := func() {
		g.sclrscr()
		g.lrdfile("BANK")
	}

	prompt := func() {
		g.sln("")
		g.lw("  `2Gold In Hand: `0" + prettyInt(p.Gold))
		g.lln("  `2Gold In Bank: `0" + prettyInt(p.Bank))
		g.lln("  `5The Bank `8(W,D,R,T,Q)  (? for menu)")
	}

	// getAmount reads parseInt(dk.console.getstr({len:11, integer:true})),
	// NaN being 0.
	getAmount := func() int64 {
		n, ok := jsParseInt(g.consoleGetstr(GetStrOpts{Len: 11, Integer: true}))
		if !ok {
			return 0
		}
		return n
	}

	menu()
	prompt()
	for {
		g.commandPrompt()
		ch = strings.ToUpper(g.getkey())
		g.sln(ch)
	sw:
		switch ch {
		case "W":
			g.lln("`c  `%Ye Olde Bank`#")
			g.sln("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
			g.lln("  `2Gold In Hand: `0" + prettyInt(p.Gold) + "  `2Gold In Bank: `0" + prettyInt(p.Bank))
			g.sln("")
			g.lln("  `2\"How much gold would you like to withdraw?\" `0(1 for ALL of it)")
			g.sln("")
			g.lw("  `0AMOUNT : ")
			amt = getAmount()
			if amt == 1 {
				amt = p.Bank
			}
			if p.Gold+amt > 2000000000 {
				amt = 2000000000 - p.Gold
			}
			switch {
			case amt == 0:
				g.sln("")
				g.sln("  \"Okay. Maybe another time.\"")
			case amt < 0:
				g.sln("")
				g.sln("  \"Uh...Wouldn't that be depositing?!\"")
			case amt > p.Bank:
				g.sln("")
				g.sln("  \"I'm afraid you don't have that much in your account, ")
				g.sln(pick(p.Sex == "M", "sir.\"", "ma'am.\""))
			default:
				g.sln("")
				g.sln("  Done! " + prettyInt(amt) + " withdrawn.")
				p.Bank -= amt
				p.Gold += amt
			}
		case "D":
			g.lln("`c  `%Ye Olde Bank`#")
			g.sln("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
			g.lln("  `2Gold In Hand: `0" + prettyInt(p.Gold) + "   `2Gold In Bank: `0" + prettyInt(p.Bank))
			g.sln("")
			g.lln("  `2\"How much gold would you like to deposit?\" `0(1 for ALL of it)")
			g.sln("")
			g.lw("  `0AMOUNT: ")
			amt = getAmount()
			if amt == 1 {
				amt = p.Gold
			}
			if p.Bank+amt > 2000000000 {
				amt = 2000000000 - p.Bank
			}
			switch {
			case amt == 0:
				g.sln("")
				g.sln("  \"Okay. Maybe another time.\"")
			case amt < 0:
				g.sln("")
				g.lln("  \"Uh...Wouldn't that be withdrawing?!\"")
			case amt > p.Gold:
				g.sln("")
				g.sln("  \"I'm afraid you don't have that much on you, ")
				g.sln(pick(p.Sex == "M", "sir.\"", "ma'am.\""))
			case p.Bank >= 2000000000:
				p.Bank = 2000000000
				g.sln("")
				g.sln("  \"I'm sorry, but we can only keep 2,000,000,000 gold at a time. ")
			default:
				g.sln("")
				g.sln("  Done! " + prettyInt(amt) + " deposited.")
				p.Gold -= amt
				p.Bank += amt
			}
		case "T":
			if !g.Settings.TransfersOn {
				g.sln("")
				g.sln("")
				g.lln("  `0I'm sorry, but all of our carriers are on permanent vacation.")
				g.sln("")
				break
			}
			if g.Settings.TransfersPerDay > 0 && p.TransferredGold >= g.Settings.TransfersPerDay {
				g.lln("  `0\"Very sorry, but all couriers are busy.")
				g.sln("")
				break
			}
			g.lln("`c  `%Ye Olde Bank`#")
			g.lln("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
			g.lw("  `2Gold In Hand: `0" + prettyInt(p.Gold))
			g.lln("  `2Gold In Bank: `0" + prettyInt(p.Bank))
			g.foreground(2)
			g.sln("  \"How much gold would you like to transfer?\"")
			g.lw("  `0AMOUNT : ")
			amt = getAmount()
			switch {
			case amt == 0:
				g.sln("")
				g.sln("  \"Okay. Maybe some other time.\"")
			case amt < 0:
				g.sln("")
				g.sln("  \"Uh...Wouldn't that be taking money from his account?!\"")
			case amt > p.Bank:
				g.sln("")
				g.sln("  \"I'm afraid you don't have that much in your account, ")
				g.sln(pick(p.Sex == "M", "sir.\"", "woman.\""))
			case amt > g.Settings.TransferAmount:
				g.sln("")
				g.sln("  \"Our messenger will not carry more than " + prettyInt(g.Settings.TransferAmount) + " gold pieces at a time!\"")
			default:
				g.sln("")
				g.sln("  \"And who would you like to send this gold to?\"")
				to = g.findPlayer()
				if to == -1 || to == int64(p.Record) {
					g.sln("")
					g.sln("  Gold Not Sent!")
					break sw
				}
				op := g.playerGet(to)
				p.Bank -= amt
				p.TransferredGold++
				p.Put()
				g.sln("")
				g.lln("  Done!  " + op.Name + " will be notified!")
				a := strconv.FormatInt(amt, 10)
				mline := " \n  `%BANK NOTICE:\n`0-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-`2\n  `0" + p.Name + " `2has transferred `%" + a + " `2gold\n  to your account.\n`b" + a
				g.mailTo(to, mline)
			}
		case "R", "Q", "\r":
		case "?":
			menu()
			prompt()
		case "1":
			g.sln("")
			g.sln("")
			g.sln("  Nice rock.")
			g.sln("")
		case "4":
			g.sln("")
			g.sln("")
			g.lln("  `2Look wise guy, these are `)NOT `2secret keys.  They say these things")
			g.sln("  when you click on pictures in RIP.  Nothing more. (Then again...)")
			g.sln("")
		case "2":
			g.sln("")
			g.sln("")
			if p.Clss == 3 {
				if p.HasFairy {
					g.lln("  `0** `%TRICKY EVENT! `0**`2")
					g.sln("")
					g.sln("  The buzzing in your pocket reminds you")
					g.sln("  that the little fairy can feel your emotions.")
					g.sln("")
					g.moreNomail()
					g.lln("  `%IT ESCAPES FROM YOUR POCKET AND PICKS THE LOCK!")
					g.sln("")
					g.moreNomail()
					if g.Settings.OldSteal {
						amt = ((g.random(500) + 500) * p.Level * p.Level * p.Level)
					} else {
						amt = ((g.random(500) + 500) * p.Level)
					}
					p.HasFairy = false
					g.lln("  `2You steal `0" + prettyInt(amt) + " `2while the banker isn't looking.")
					// BUG (lord.js): no 2,000,000,000 clamp here (the record
					// write clamps to the field's range).  Kept.
					p.Gold += amt
					g.sln("")
					g.moreNomail()
					prompt()
					break
				}
				g.sln("  You're a thief, find the key.")
			} else {
				g.sln("  Steal?  Never!")
			}
			g.sln("")
		case "3":
			g.sln("")
			g.sln("")
			g.sln("  Maybe there is a Genie in there.")
			g.sln("")
		}
		if strings.Contains("RQ\r", ch) {
			break
		}
	}
	if !p.LeftBank {
		p.LeftBank = true
		if g.random(30) == 25 {
			if !p.Amulet {
				// class_list[player.clss] is undefined for other classes.
				cls := "undefined"
				if p.Clss >= 0 && p.Clss < int64(len(classList)) {
					cls = classList[p.Clss]
				}
				amt = p.Level * 1000
				g.sclrscr()
				g.lln("  `%A random event!")
				g.lln("`#-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
				g.sln("")
				g.lln("  `2While leaving the bank, a stranger approaches you..")
				g.sln("")
				g.lln("  `2\"`0Hello there " + cls + ". My name is Aidan. I have something that might")
				g.sln("  interest you. It's just a little necklace I found in the forest, but")
				g.lln("  the old man in the Dark Cloak Tavern called it an '`%Amulet of Accuracy`2'")
				g.lln("  I have no use for it. I'll sell it to you - only " + strconv.FormatInt(amt, 10) + " gold`2\"")
				g.sln("")
				g.sln("")
				g.lw("  `2Buy the Amulet? `2[`0N`2] : ")
				ch = strings.ToUpper(g.getkey())
				if ch != "Y" {
					ch = "N"
				}
				g.sln(ch)
				if ch == "Y" {
					switch {
					case p.Gold < amt:
						g.lln("  `2\"`0I'm sorry, but it seems you dont have enough gold to make")
						g.lln("  the trade. Maybe some other time.`2\"")
					case p.Amulet:
						g.lln("  `2\"`0I'm sorry, but you already have an amulet. I cannot sell")
						g.lln("  you another.`2\"")
					default:
						g.lln("  `2\"`0Thank you. I believe you'll enjoy that little amulet.")
						g.lln("  It may come in handy in battle.`2\"")
						p.Gold -= amt
						p.Amulet = true
					}
				}
			}
		}
	}
}
