package game

import (
	"strconv"
	"strings"
)

// bjCard is a card of blackjack()'s deck ({suit, val} in lord.js).
type bjCard struct{ suit, val int }

var (
	bjSuits  = []string{"\x06", "\x03", "\x05", "\x04"}
	bjScol   = []string{"`r2", "`r2`4", "`r2", "`r2`4"}
	bjFaces  = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	bjValues = []int{11, 2, 3, 4, 5, 6, 7, 8, 9, 10, 10, 10, 10}
)

// TODO (lord.js): Option to disable?

// blackjack is blackjack(): the dwarf's game of chance.
func (g *Game) blackjack() {
	startgold := g.player.Gold
	var ch string

	for {
		g.lln("`c  `%A meeting of chance.`2")
		g.sln("")
		g.sln("  A sly looking dwarf hops out of the brush.")
		g.sln("")
		g.lln("  `2\"`0How about a game, friend?`2\"")
		g.sln("")
		g.lln("  `2(`0G`2)ive the game a chance")
		g.lln("  (`0T`2)ell the Dwarf to screw off")
		for {
			g.commandPrompt()
			ch = strings.ToUpper(g.getkey())
			if strings.Contains("GT?", ch) {
				break
			}
		}
		g.sln("")
		g.foreground(0)
		switch ch {
		case "?":
		case "T":
			g.sln("")
			g.lln("  `2\"`0Your loss, kid.  Forget you ever saw me!`2\"")
			g.sln("")
			g.more()
			return
		case "G":
			g.sln("")
			g.lln("  `0\"Excellent!\" `2the dwarf exclaims as he pulls out a makeshift black jack")
			g.sln("  table!")
			g.sln("")
			g.lln("  `2\"`0As for the rules... I have to stay on a 17.  You can double down on a")
			g.lln("  9, 10 or 11, and I do allow splitting pairs multiple times.`2\"")
			g.sln("")
			g.sln("  You rub your chin - with any luck you'll double your money...")
			g.sln("")
			g.more()
			g.sclrscr()
			g.lrdfile("BJTABLE")
			g.blackjackPlayHand(startgold)
			return
		}
	}
}

// blackjackPlayHand is blackjack's play_hand() ("sam2").
func (g *Game) blackjackPlayHand(startgold int64) {
	p := g.player
	var cards []bjCard
	var dealerHands, playerHands [][]bjCard
	var bet int64
	var res string
	var tmp int64
	var ch string
	mc := g.morechk
	shuffleStrs := []string{"/\\", "==", "\xcd\xcd"}
	gotoxy := g.con.GotoXY

	handTotal := func(hand []bjCard, mn bool) int {
		aces11 := 0
		ret := 0

		for _, c := range hand {
			if c.val == 0 {
				aces11++
			}
			ret += bjValues[c.val]
			if ret > 21 {
				if aces11 > 0 {
					ret -= 10
					aces11--
				}
			}
		}
		if mn {
			for aces11 > 0 {
				ret -= 10
				aces11--
			}
		}
		return ret
	}

	shuffle := func() {
		cards = nil
		for j := 0; j < 4; j++ {
			for si := 0; si < 13; si++ {
				cards = append(cards, bjCard{suit: j, val: si})
			}
		}
		for si := len(cards); si > 1; si-- {
			j := g.random(int64(si))
			cards[si-1], cards[j] = cards[j], cards[si-1]
		}
	}

	showTotal := func(play bool, total string, split int) {
		if split == 0 && len(playerHands) > 1 {
			split = 1
		}
		if play {
			switch split {
			case 1:
				gotoxy(0, 22)
				g.lw("`r0`2Hand1: `0" + total)
			case 2:
				gotoxy(13, 22)
				g.lw("`r0`2Hand2: `0" + total)
			default:
				gotoxy(0, 22)
				g.lw("`r0`2You: `0" + total)
			}
		} else {
			gotoxy(68, 22)
			g.lw("`r0`2Dealer: `0" + total)
		}
	}

	wizclrScr := func(num int) {
		switch num {
		case 1:
			g.background(0)
			for y := 11; y <= 19; y++ {
				gotoxy(1, y)
				g.sw("            ")
			}
		case 2:
			g.background(0)
			for y := 11; y <= 19; y++ {
				gotoxy(67, y)
				g.sw("            ")
			}
		case 3:
			g.background(2)
			for y := 1; y <= 7; y++ {
				gotoxy(24, y)
				g.sw("                              ")
			}
		case 4:
			g.background(2)
			for y := 10; y <= 16; y++ {
				gotoxy(16, y)
				g.sw("                       ")
			}
		case 5:
			g.background(2)
			for y := 10; y <= 16; y++ {
				gotoxy(40, y)
				g.sw("                       ")
			}
		case 6:
			g.background(1)
			gotoxy(37, 18)
			g.sw("        ")
		}
	}

	wizmore := func() {
		gotoxy(4, 19)
		g.lw("`2<`0MORE`2>")
		g.getkey()
		for wi := 0; wi < 6; wi++ {
			g.sw("\b")
		}
		for wi := 0; wi < 6; wi++ {
			g.sw(" ")
		}
	}

	getCpic := func(card bjCard) string {
		return bjScol[card.suit] + bjFaces[card.val] + bjSuits[card.suit]
	}

	// hands returns the dealer's or the player's hands.
	//
	// DIFF: lord.js picks the hand with `who === '1'`, comparing the number
	// it is always passed with a string, so the dealer's cards went into
	// the player's hand and the dealer's hidden card (hands.dealer[0][1])
	// was undefined when it was turned over, which throws.  The dealer's
	// hand is used for who == 1 here.
	hands := func(who int) *[][]bjCard {
		if who == 1 {
			return &dealerHands
		}
		return &playerHands
	}

	deal := func(who, count int, show, split bool) {
		var cpic string
		var card bjCard
		var x, y int
		hs := hands(who)
		hand := (*hs)[pick(split, 1, 0)]
		var full bool
		first := false
		var dtmp int

		for di := 0; di < count; di++ {
			card = cards[0]
			cards = cards[1:]
			cpic = getCpic(card)
			g.mswait(500)
			hand = append(hand, card)
			if who == 1 && len(hand) == 2 {
				first = true
			}
			full = (len(hand) % 2) == 1
			dtmp = len(hand)
			if split {
				dtmp = 14 - len(hand)
				if len(hand)%2 == 1 {
					dtmp -= 2
				}
			}
			x = 17 + ((dtmp-1)/2)*8
			if who == 1 {
				x += 8
			}
			y = pick(who == 1, 1, 10) + ((dtmp % 2) * 2)
			g.background(2)
			g.foreground(0)
			gotoxy(x-1, y)
			g.sw("\xd5\xcd\xcd\xcd\xcd\xcd\xb8")
			gotoxy(x-1, y+1)
			if first {
				g.sw("\xb3\xb1\xb1\xb1\xb1\xb1")
				g.foreground(0)
				gotoxy(x+5, y+1)
				g.sw("\xb3")
			} else {
				g.lw("\xb3" + cpic)
				g.foreground(0)
				gotoxy(x+5, y+1)
				g.sw("\xb3")
			}
			if full {
				gotoxy(x-1, y+2)
				g.sw("\xb3     \xb3")
				gotoxy(x-1, y+3)
				g.sw("\xb3     \xb3")
				gotoxy(x-1, y+4)
				g.sw("\xc0\xc4\xc4\xc4\xc4\xc4\xd9")
			}
		}
		if show {
			showTotal(who == 2, strconv.Itoa(handTotal(hand, false)), pick(split, 2, 0))
		} else {
			showTotal(who == 2, "?", 0)
		}
		(*hs)[pick(split, 1, 0)] = hand
	}

	cardCheck := func() string {
		var cch string
		candouble := true
		var cctmp int
		var ptotal, dtotal int
		dealcards := true
		hand := 0

		if playerHands[0][0].val == playerHands[0][1].val {
			wizclrScr(1)
			g.foreground(10)
			gotoxy(1, 11)
			g.sw("Would you")
			gotoxy(1, 12)
			g.sw("like to")
			gotoxy(1, 13)
			g.sw("split?")
			gotoxy(67, 11)
			g.lw("`2(`0Y`2)es (`0N`2)o. ")
			cch = strings.ToUpper(g.getkey())
			if cch != "Y" {
				cch = "N"
			}
			if cch == "Y" {
				if p.Gold < bet {
					wizclrScr(1)
					g.foreground(10)
					gotoxy(1, 11)
					g.sw("Hey! you")
					gotoxy(1, 12)
					g.sw("don't have")
					gotoxy(1, 13)
					g.sw("enough gold!")
					wizmore()
					candouble = false
				} else {
					last := playerHands[0][len(playerHands[0])-1]
					playerHands[0] = playerHands[0][:len(playerHands[0])-1]
					playerHands = append(playerHands[:1], []bjCard{last})
					candouble = false
					p.Gold -= bet
					gotoxy(0, 21)
					g.lw("`r0`2Gold on hand: `$" + prettyInt(p.Gold))
					gotoxy(37, 18)
					g.lw("`r2`4" + strconv.FormatInt(bet*2, 10))
					gotoxy(16, 10)
					g.lw("`r2`2       ")
					gotoxy(16, 11)
					g.sw("       ")
					g.foreground(0)
					gotoxy(56, 12)
					g.sw("\xd5\xcd\xcd\xcd\xcd\xcd\xb8")
					gotoxy(56, 13)
					g.lw("\xb3" + getCpic(playerHands[1][0]))
					g.foreground(0)
					gotoxy(62, 14)
					g.sw("\xb3")
					gotoxy(56, 14)
					g.sw("\xb3     \xb3")
					gotoxy(56, 15)
					g.sw("\xb3     \xb3")
					gotoxy(56, 16)
					g.sw("\xc0\xc4\xc4\xc4\xc4\xc4\xd9")
					gotoxy(0, 22)
					g.lw("`r0`2Hand1:       Hand2:  ")
					deal(2, 1, true, false)
					deal(2, 1, true, true)
				}
			}
		}
		if handTotal(playerHands[0], false) == 21 && len(playerHands[0]) == 2 {
			candouble = false
		}
		if candouble {
			cctmp = handTotal(playerHands[0], true)
			if cctmp > 8 && cctmp < 12 {
				for {
					wizclrScr(1)
					g.foreground(10)
					gotoxy(1, 11)
					g.sw("Would you")
					gotoxy(1, 12)
					g.sw("like to")
					gotoxy(1, 13)
					g.sw("double down?")
					gotoxy(67, 11)
					g.lw("`2(`0Y`2)es (`0N`2)o. ")
					cch = strings.ToUpper(g.getkey())
					if cch == "N" {
						break
					}
					if cch == "Y" {
						if p.Gold < bet {
							wizclrScr(1)
							g.foreground(10)
							gotoxy(1, 11)
							g.sw("Hey! you")
							gotoxy(1, 12)
							g.sw("don't have")
							gotoxy(1, 13)
							g.sw("enough gold!")
							wizmore()
							break
						}
						p.Gold -= bet
						bet *= 2
						gotoxy(0, 21)
						g.lw("`r0`2Gold on hand:                  ")
						gotoxy(0, 21)
						g.lw("`r0`2Gold on hand: `$" + prettyInt(p.Gold))
						gotoxy(37, 19)
						g.lw("`r2`4" + strconv.FormatInt(bet, 10))
						deal(2, 1, true, false)
						ptotal = handTotal(playerHands[0], false)
						if ptotal > 21 {
							return "BUST"
						}
						dealcards = false
						break
					}
				}
			}
		}
		for dealcards {
			ptotal = handTotal(playerHands[hand], false)
			if len(playerHands) > 1 {
				if ptotal > 21 {
					wizclrScr(1)
					gotoxy(1, 11)
					g.lw("`0Sorry!")
					gotoxy(1, 12)
					g.sw("That hand")
					gotoxy(1, 13)
					g.sw("busted.")
					wizmore()
					if hand == 1 {
						break
					}
					hand++
					ptotal = handTotal(playerHands[hand], false)
				}
			} else {
				if ptotal > 21 {
					return "BUST"
				}
				if len(playerHands[0]) == 2 && ptotal == 21 {
					return "BLACKJACK"
				}
			}
			// start3
			wizclrScr(1)
			wizclrScr(2)
			g.foreground(10)
			gotoxy(1, 11)
			g.sw("Would you")
			gotoxy(1, 12)
			g.sw("like to hit")
			gotoxy(1, 13)
			g.sw("or stay?")
			if len(playerHands) > 1 {
				gotoxy(1, 15)
				g.sw("Hand" + strconv.Itoa(hand+1))
			}
			gotoxy(67, 11)
			g.lw("`2(`0H`2)it")
			gotoxy(67, 12)
			g.lw("`2(`0S`2)tay")
			cch = strings.ToUpper(g.getkey())
			if cch == "S" {
				if len(playerHands) == 1 {
					break
				}
				if hand == 1 {
					break
				}
				hand++
			} else if cch == "H" {
				deal(2, 1, true, !(hand == 0))
			}
		}
		g.lw("`r2")
		g.foreground(0)
		gotoxy(24, 1)
		g.sw("\xd5\xcd\xcd\xcd\xcd\xcd\xb8")
		gotoxy(24, 2)
		g.sw("\xb3\xb1\xb1\xb1\xb1 \xb3")
		g.mswait(500)
		gotoxy(24, 2)
		g.sw("\xb3\xb1\xb1\xb1  \xb3")
		g.mswait(500)
		gotoxy(24, 2)
		g.sw("\xb3\xb1\xb1   \xb3")
		g.mswait(500)
		gotoxy(24, 2)
		g.sw("\xb3     ")
		gotoxy(24, 2)
		g.lw("\xb3" + getCpic(dealerHands[0][1]))
		g.foreground(0)
		gotoxy(30, 2)
		g.sw("\xb3")
		showTotal(false, strconv.Itoa(handTotal(dealerHands[0], false)), 0)
		for {
			dtotal = handTotal(dealerHands[0], false)
			if dtotal > 21 {
				if len(playerHands) > 1 {
					// BUG (lord.js): a dealer bust wins both split hands,
					// even one that busted.  Kept.
					return "WW"
				}
				return "DBUST"
			}
			if dtotal == 21 && len(dealerHands[0]) == 2 {
				return "DBLACKJACK"
			}
			if dtotal < 17 {
				deal(1, 1, true, false)
			} else {
				break
			}
		}
		if len(playerHands) == 1 {
			if dtotal == ptotal {
				return "PUSH"
			}
			if dtotal > ptotal {
				return "DWIN"
			}
			return "WIN"
		}
		r := ""
		ptotal = handTotal(playerHands[0], false)
		if dtotal == ptotal {
			r += "P"
		} else if dtotal > ptotal || ptotal > 21 {
			r += "D"
		} else {
			r += "W"
		}
		ptotal = handTotal(playerHands[1], false)
		// BUG (lord.js): unlike the first hand, a busted second hand
		// isn't checked here, so it wins against a lower dealer total.
		// Kept.
		if dtotal == ptotal {
			r += "P"
		} else if dtotal > ptotal {
			r += "D"
		} else {
			r += "W"
		}
		return r
	}

	g.morechk = false
hand:
	for {
		shuffle()
		gotoxy(1, 11)
		g.lw("`r0`%Shuffling.")
		for i := 0; i < 9; i++ {
			g.mswait(100)
			gotoxy(1, 12)
			g.sw(shuffleStrs[i%len(shuffleStrs)])
		}
		dealerHands = [][]bjCard{{}}
		playerHands = [][]bjCard{{}}
		// DIFF (lord.js): startgold was reset to player.gold here.
		// startgold = player.gold;
		gotoxy(0, 21)
		g.lw("`r0`2Gold on hand: `$" + prettyInt(p.Gold))
		gotoxy(0, 22)
		g.lw("`r0                                                                              ")
		showTotal(true, "0", 0)
		showTotal(false, "0", 0)
		for {
			wizclrScr(1)
			g.foreground(10)
			gotoxy(1, 11)
			g.sw("How much")
			gotoxy(1, 12)
			g.sw("ya gonna")
			gotoxy(1, 13)
			g.sw("wager?")
			g.foreground(4)
			g.background(2)
			gotoxy(37, 19)
			b, ok := jsParseInt(strings.TrimSpace(g.getstr(38, 19, 8, 1, 15, "200", GetStrOpts{Integer: true, InputBox: true, Select: true})))
			bet = b
			if !ok || bet == 0 {
				wizclrScr(1)
				g.foreground(10)
				gotoxy(1, 11)
				g.sw("Fine, maybe")
				gotoxy(1, 12)
				g.sw("later.")
				wizmore()
				break hand
			}
			if bet < 200 {
				gotoxy(1, 14)
				g.lw("`0`r0Too little!")
				wizmore()
				wizclrScr(1)
			} else if bet > p.Level*1000 {
				gotoxy(1, 14)
				g.lw("`0`r0Too much!")
				wizmore()
				wizclrScr(1)
			} else if p.Gold < bet {
				wizclrScr(1)
				g.foreground(10)
				gotoxy(1, 11)
				g.sw("Hey! you")
				gotoxy(1, 12)
				g.sw("don't have")
				gotoxy(1, 13)
				g.sw("enough gold!")
				wizmore()
				break hand
			} else {
				break
			}
		}
		p.Gold -= bet
		gotoxy(0, 21)
		g.lw("`r0`2Gold on hand: `$" + prettyInt(p.Gold) + "            ")
		deal(1, 1, false, false)
		deal(2, 1, true, false)
		deal(1, 1, false, false)
		deal(2, 1, true, false)
		res = cardCheck()
		switch res {
		case "BUST":
			g.foreground(10)
			wizclrScr(1)
			gotoxy(1, 11)
			switch g.random(3) {
			case 0:
				g.sw("You busted!")
				gotoxy(1, 12)
				g.sw("Better luck")
				gotoxy(1, 13)
				g.sw("next time.")
			case 1:
				g.sw("Oh too bad!")
				gotoxy(1, 12)
				g.sw("Maybe next")
				gotoxy(1, 13)
				g.sw("hand? (haw!)")
			case 2:
				g.sw("You busted!")
				gotoxy(1, 12)
				g.sw("Have you")
				gotoxy(1, 13)
				g.sw("played this")
				gotoxy(1, 14)
				g.sw("game before?")
			}
			wizmore()
		case "BLACKJACK":
			g.foreground(10)
			bet *= 3
			p.Gold += bet
			if p.Gold > 2000000000 {
				p.Gold = 2000000000
			}
			wizclrScr(1)
			gotoxy(1, 11)
			g.sw("You got a")
			gotoxy(1, 12)
			g.sw("Blackjack!?")
			gotoxy(1, 13)
			g.sw("Are you")
			gotoxy(1, 14)
			g.sw("cheating?!")
			gotoxy(1, 16)
			g.sw("You win")
			gotoxy(1, 17)
			g.lw("`$" + prettyInt(bet) + "`0")
			gotoxy(1, 18)
			g.sw("Gold.")
			wizmore()
		case "DBUST":
			g.foreground(10)
			bet *= 2
			p.Gold += bet
			if p.Gold > 2000000000 {
				p.Gold = 2000000000
			}
			wizclrScr(1)
			gotoxy(1, 11)
			switch g.random(3) {
			case 0:
				g.sw("I busted!")
				gotoxy(1, 12)
				g.sw("See? You")
				gotoxy(1, 13)
				g.sw("win ")
				gotoxy(1, 15)
				g.lw("`$" + prettyInt(bet) + "`0")
				gotoxy(1, 16)
				g.sw("Gold.")
			case 1:
				g.sw("Looks like")
				gotoxy(1, 12)
				g.sw("I busted.")
				gotoxy(1, 14)
				g.sw("You win")
				gotoxy(1, 15)
				g.lw("`$" + prettyInt(bet) + "`0")
				gotoxy(1, 16)
				g.sw("Gold... Arg.")
			case 2:
				g.sw("I busted!")
				gotoxy(1, 12)
				g.sw("Damnit!")
				gotoxy(1, 14)
				g.sw("You win")
				gotoxy(1, 15)
				g.lw("`$" + prettyInt(bet) + "`0")
				gotoxy(1, 16)
				g.sw("Gold.")
			}
			wizmore()
		case "DBLACKJACK":
			g.foreground(10)
			wizclrScr(1)
			gotoxy(1, 11)
			switch g.random(3) {
			case 0:
				g.sw("Dealer gets")
				gotoxy(1, 12)
				g.sw("Blackjack!")
				gotoxy(1, 13)
				g.sw("Aw, too bad")
				gotoxy(1, 14)
				g.sw("for the big")
				gotoxy(1, 15)
				g.sw("human.")
			case 1:
				g.sw("Blackjack!")
				gotoxy(1, 12)
				g.sw("I'm hot")
				gotoxy(1, 13)
				g.sw("today!")
			case 2:
				g.sw("Dealer gets")
				gotoxy(1, 12)
				g.sw("Blackjack!")
				gotoxy(1, 13)
				g.sw("Not your")
				gotoxy(1, 14)
				g.sw("day is it?")
			}
			wizmore()
		case "PUSH":
			g.foreground(10)
			wizclrScr(1)
			gotoxy(1, 11)
			g.sw("It's a push!")
			gotoxy(1, 12)
			g.sw("I guess it's")
			gotoxy(1, 13)
			g.sw("better than")
			gotoxy(1, 14)
			g.sw("losing.")
			p.Gold += bet
			if p.Gold > 2000000000 {
				p.Gold = 2000000000
			}
			wizmore()
		case "DWIN":
			g.foreground(10)
			wizclrScr(1)
			gotoxy(1, 11)
			g.sw("Looks like")
			gotoxy(1, 12)
			g.sw("you lose.")
			gotoxy(1, 13)
			g.sw("Oh well!")
			gotoxy(1, 14)
			g.sw("better luck")
			gotoxy(1, 15)
			g.sw("next time.")
			wizmore()
		case "WIN":
			g.foreground(10)
			wizclrScr(1)
			gotoxy(1, 11)
			bet *= 2
			p.Gold += bet
			if p.Gold > 2000000000 {
				p.Gold = 2000000000
			}
			switch g.random(3) {
			case 0:
				g.sw("That beats")
				gotoxy(1, 12)
				g.sw("my hand!")
				gotoxy(1, 14)
				g.sw("You win")
				gotoxy(1, 15)
				g.lw("`$" + prettyInt(bet) + "`0")
				gotoxy(1, 16)
				g.sw("Gold, Loser.")
			case 1:
				g.sw("You win!")
				gotoxy(1, 12)
				g.sw("Are you")
				gotoxy(1, 13)
				g.sw("counting?")
				gotoxy(1, 15)
				g.sw("You win")
				gotoxy(1, 16)
				g.lw("`$" + prettyInt(bet) + "`0")
				gotoxy(1, 17)
				g.sw("Gold.")
			case 2:
				g.sw("You win!")
				gotoxy(1, 12)
				g.sw("Not too bad")
				gotoxy(1, 14)
				g.sw("for a kid...")
				gotoxy(1, 15)
				g.sw("You win")
				gotoxy(1, 16)
				g.lw("`$" + prettyInt(bet) + "`0")
				gotoxy(1, 17)
				g.sw("Gold.")
			}
			wizmore()
		case "PP", "PD", "PW", "DP", "DD", "DW", "WP", "WD", "WW":
			tmp = 0
			wizclrScr(1)
			gotoxy(1, 11)
			g.lln("`0I have " + strconv.Itoa(handTotal(dealerHands[0], false)) + ".")
			g.mswait(1500)
			gotoxy(1, 12)
			switch res[0] {
			case 'W':
				g.sw("Hand1 Wins!")
				tmp += bet * 2
			case 'P':
				g.sw("Hand1 Push!")
				tmp += bet
			case 'D':
				g.sw("Hand1 Loses!")
			}
			gotoxy(1, 14)
			switch res[1] {
			case 'W':
				g.sw("Hand2 Wins!")
				tmp += bet * 2
			case 'P':
				g.sw("Hand2 Push!")
				tmp += bet
			case 'D':
				g.sw("Hand2 Loses!")
			}
			gotoxy(1, 15)
			g.sw("You win")
			gotoxy(1, 16)
			g.lw("`$" + strconv.FormatInt(tmp, 10) + "`0")
			gotoxy(1, 17)
			g.sw("Gold.")
			p.Gold += tmp
			if p.Gold > 2000000000 {
				p.Gold = 2000000000
			}
			wizmore()
		}
		wizclrScr(1)
		wizclrScr(2)
		gotoxy(1, 11)
		g.foreground(10)
		g.sw("Play again?")
		gotoxy(67, 11)
		g.lw("`2(`0Y`2)es (`0N`2)o. ")
		ch = strings.ToUpper(g.getkey())
		if ch != "Y" {
			ch = "N"
		}
		if ch == "N" {
			break
		}
		wizclrScr(1)
		wizclrScr(2)
		wizclrScr(3)
		wizclrScr(4)
		wizclrScr(5)
		wizclrScr(6)
	}
	g.foreground(10)
	if p.Gold > startgold {
		wizclrScr(1)
		gotoxy(1, 11)
		g.sw("Quitting")
		gotoxy(1, 12)
		g.sw("while your")
		gotoxy(1, 13)
		g.sw("ahead?")
		gotoxy(1, 14)
		g.sw("Smart move.")
	} else if p.Gold < startgold {
		wizclrScr(1)
		gotoxy(1, 11)
		g.sw("Come back")
		gotoxy(1, 12)
		g.sw("soon. We")
		gotoxy(1, 13)
		g.sw("enjoyed")
		gotoxy(1, 14)
		g.sw("your money!")
		gotoxy(1, 15)
		g.sw("::snicker::")
	} else {
		wizclrScr(1)
		gotoxy(1, 11)
		g.sw("It's been")
		gotoxy(1, 12)
		g.sw("nice doing")
		gotoxy(1, 13)
		g.sw("business")
		gotoxy(1, 14)
		g.sw("with you...")
	}
	wizmore()
	gotoxy(0, 23)
	g.morechk = mc
	g.curlinenum = 1
}
