package game

// mailCheck is mail_check(): is there mail waiting for the player?
func (g *Game) mailCheck() bool {
	return fileExists(g.gamedir("mail" + itoa(g.player.Record) + ".lrd"))
}

// checkMail is check_mail(): deliver waiting mail ("STOPPED BY A
// MESSENGER").
func (g *Game) checkMail() {
	mf := g.gamedir("mail" + itoa(g.player.Record) + ".lrd")
	g.fmutex(mf + ".tmp")
	fileRemove(mf + ".tmp")
	g.fmutex(mf)
	if fileExists(mf + ".tmp") {
		fileRemove(mf + ".tmp")
	}
	if fileExists(mf) {
		fileRename(mf, mf+".tmp")
		g.rfmutex(mf)
		g.foreground(15)
		g.sln("  ** YOU ARE STOPPED BY A MESSENGER WITH THE FOLLOWING NEWS: **")
		g.sln("")
		g.showLordFile(mf+".tmp", true, true)
		fileRemove(mf + ".tmp")
		g.sln("")
	} else {
		g.rfmutex(mf)
	}
	g.rfmutex(mf + ".tmp")
}
