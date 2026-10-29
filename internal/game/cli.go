package game

import "os"

// LoadSettings is cmdline(): lord.ini, then the -P/-M overrides.
func (g *Game) LoadSettings(prefix, menu string) {
	if prefix != "" {
		g.Settings.GamePrefix = prefix
	}
	if menu != "" {
		g.Settings.MenuDir = menu
	}
	g.readINI()
	if prefix != "" {
		g.Settings.GamePrefix = prefix
	}
	if menu != "" {
		g.Settings.MenuDir = menu
	}
}

// Reset is cmdline()'s RESET: delete state.bin and player.bin.
func (g *Game) Reset() {
	_ = os.Remove(g.gamedir("state.bin"))
	_ = os.Remove(g.gamedir("player.bin"))
}
