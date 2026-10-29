package game

// Stubs for functions another group calls.  The owning agent replaces this
// file with the real implementations.

func (g *Game) obCleanup()                                                     { panic("TODO") }
func (g *Game) obSendResp(op *Player, str string)                              { panic("TODO") }
func (g *Game) onBattle(op *Player, first bool, action string) string          { panic("TODO") }
func (g *Game) battlePromptLines() int                                         { panic("TODO") }
func (g *Game) battle(op *Player, pfight, cantRun bool) string                 { panic("TODO") }
func (g *Game) say(fname string, enemy *Player, typ int, header, force string) { panic("TODO") }
func (g *Game) attackInInn()                                                   { panic("TODO") }
func (g *Game) fightDragon(cantRun bool)                                       { panic("TODO") }
func (g *Game) slaughterOthers()                                               { panic("TODO") }
