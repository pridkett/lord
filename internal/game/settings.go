package game

import (
	"os"
	"strconv"
	"strings"
)

// Settings are lord.js's `settings`, read from lord.ini's global section.
type Settings struct {
	DeleteDays      int64  `ini:"InactivityDays"`        // Inactive players are deleted after this many days
	Del1XP          bool   `ini:"Delete1XP"`             // Delete players with 1 xp daily
	CleanMode       bool   `ini:"CleanMode"`             // CLEAN MODE
	WinDeeds        int64  `ini:"DeedsToWin"`            // Heroic deeds to win
	TransfersOn     bool   `ini:"BankTransfersAllowed"`  //
	TransfersPerDay int64  `ini:"BankTransfersPerDay"`   //
	TransferAmount  int64  `ini:"MaximumTransferAmount"` // (lord.ini ships it as MaxTransferAmount; both are read)
	PvPFightsPerDay int64  `ini:"PlayerFightsPerDay"`    //
	ForestFights    int64  `ini:"ForestFightsPerDay"`    // Forest fights per day
	ResDays         int64  `ini:"ResurrectionDays"`      // Days to resurrect inactive players before leaving them dead
	NoTime          string `ini:"OutOfTimeFile"`         // Out of time filename
	SafeNode        bool   `ini:"SafeNode"`              // Creates node.xx files and if it exists, marks the player as off.
	UseFancyMore    bool   `ini:"UseFancyMore"`          // If false, does not erase the MORE prompt.
	NoChat          bool   `ini:"NoNPCChatter"`          // If true, NPCs don't chat in the inn.
	Olivia          bool   `ini:"Olivia"`                // If true, creepy event (olivia) occurs
	FunkyFlowers    bool   `ini:"FunkyFlowers"`          // Allow `rX and `Bx codes in flowers (only remove ``)
	ShopLimit       bool   `ini:"ShopLimit"`             // Don't restrict weapons and armour by base def/str
	OldSkillPoints  bool   `ini:"FiveSkillPoints"`       // Require 5 skill points for use points for Death Knight and Thief
	DefForPK        bool   `ini:"BonusDefForPlayerKill"` // Get bonus defence for killing a player.
	StrForPK        bool   `ini:"BonusStrForPlayerKill"` // Get bonus strength for killing player.
	GamePrefix      string `ini:"GamePrefix"`            // Directory/file prefix for game data files.
	MenuDir         string `ini:"MenuDir"`               //
	BeefUp          bool   `ini:"BeefUpForHeroes"`       // Beef up monsters after you kill the dragon
	OldSteal        bool   `ini:"OldStealAmounts"`       //
	SleepDragon     bool   `ini:"WakeUpDragon"`          // Anti-camping dragon attacks.
	DKBoost         bool   `ini:"DeathKnightBoost"`      // Make death knight attack multiple 3.3 instead of 3.
	NewUglyStick    bool   `ini:"NewUglyStick"`          // Fair ugly stick means the new one, old one was always 5 pretty/1 ugly
	NoLEnemyBin     bool   `ini:"NoLEnemy"`              // Never use the lenemy file
	NoIGMsAllowed   bool   `ini:"NoIGMs"`                // Never allow any IGMs (always true here: IGMs are not supported)
}

// DefaultSettings are lord.js's defaults.
func DefaultSettings() Settings {
	return Settings{
		DeleteDays: 15, Del1XP: true, WinDeeds: 3, TransfersOn: true,
		TransfersPerDay: 2, TransferAmount: 2000000, PvPFightsPerDay: 3,
		ForestFights: 15, ResDays: 3, SafeNode: true, UseFancyMore: true,
		Olivia: true, FunkyFlowers: true, ShopLimit: true, OldSteal: true,
		NoIGMsAllowed: true,
	}
}

// TournamentSettings are lord.ini's [Tournament] section.
type TournamentSettings struct {
	Enabled bool
	Days    int64 // zero for stat mode, non-zero for days mode
	WinStat int64 // tournamentWinStats index
	XP      int64
	DKills  int64
	PKills  int64
	Level   int64
	Lays    int64
}

var tournamentWinStats = []string{"xp", "dragonkills", "playerkills", "level", "lays"}

// ini is a parsed lord.ini: section -> lower-case key -> value.
type ini map[string]map[string]string

func parseINI(b []byte) ini {
	out := ini{"": {}}
	sec := ""
	for _, ln := range splitLines(b) {
		ln = strings.TrimSpace(ln)
		if ln == "" || ln[0] == ';' {
			continue
		}
		if ln[0] == '[' && strings.HasSuffix(ln, "]") {
			sec = strings.ToLower(strings.TrimSpace(ln[1 : len(ln)-1]))
			if out[sec] == nil {
				out[sec] = map[string]string{}
			}
			continue
		}
		k, v, ok := strings.Cut(ln, "=")
		if !ok {
			continue
		}
		out[sec][strings.ToLower(strings.TrimSpace(k))] = strings.TrimSpace(v)
	}
	return out
}

func (in ini) get(sec, key string) (string, bool) {
	v, ok := in[strings.ToLower(sec)][strings.ToLower(key)]
	return v, ok
}

func iniBool(v string, def bool) bool {
	switch strings.ToLower(v) {
	case "true", "yes", "on":
		return true
	case "false", "no", "off":
		return false
	}
	if n, err := strconv.ParseInt(v, 0, 64); err == nil {
		return n != 0
	}
	return def
}

func iniInt(v string, def int64) int64 {
	if n, err := strconv.ParseFloat(v, 64); err == nil {
		return int64(n)
	}
	return def
}

// readINI is cmdline()'s readini(): load lord.ini from the game directory.
func (g *Game) readINI() {
	b, err := os.ReadFile(g.gamedir("lord.ini"))
	if err != nil {
		return
	}
	in := parseINI(b)
	s := &g.Settings
	str := func(key string, p *string) {
		if v, ok := in.get("", key); ok {
			*p = v
		}
	}
	bl := func(key string, p *bool) {
		if v, ok := in.get("", key); ok {
			*p = iniBool(v, *p)
		}
	}
	num := func(sec, key string, p *int64) {
		if v, ok := in.get(sec, key); ok {
			*p = iniInt(v, *p)
		}
	}
	num("", "InactivityDays", &s.DeleteDays)
	bl("Delete1XP", &s.Del1XP)
	bl("CleanMode", &s.CleanMode)
	num("", "DeedsToWin", &s.WinDeeds)
	bl("BankTransfersAllowed", &s.TransfersOn)
	num("", "BankTransfersPerDay", &s.TransfersPerDay)
	// DIFF: lord.js maps MaximumTransferAmount, but its lord.ini says
	// MaxTransferAmount; read both.
	num("", "MaxTransferAmount", &s.TransferAmount)
	num("", "MaximumTransferAmount", &s.TransferAmount)
	num("", "PlayerFightsPerDay", &s.PvPFightsPerDay)
	num("", "ForestFightsPerDay", &s.ForestFights)
	num("", "ResurrectionDays", &s.ResDays)
	str("OutOfTimeFile", &s.NoTime)
	bl("SafeNode", &s.SafeNode)
	bl("UseFancyMore", &s.UseFancyMore)
	bl("NoNPCChatter", &s.NoChat)
	bl("Olivia", &s.Olivia)
	bl("FunkyFlowers", &s.FunkyFlowers)
	bl("ShopLimit", &s.ShopLimit)
	bl("FiveSkillPoints", &s.OldSkillPoints)
	bl("BonusDefForPlayerKill", &s.DefForPK)
	bl("BonusStrForPlayerKill", &s.StrForPK)
	str("GamePrefix", &s.GamePrefix)
	str("MenuDir", &s.MenuDir)
	bl("BeefUpForHeroes", &s.BeefUp)
	bl("OldStealAmounts", &s.OldSteal)
	bl("DeathKnightBoost", &s.DKBoost)
	bl("WakeUpDragon", &s.SleepDragon)
	bl("NewUglyStick", &s.NewUglyStick)
	bl("NoLEnemy", &s.NoLEnemyBin)
	// GameServer, GameUser, GamePassword (the remote server) and NoIGMs
	// are not used: neither feature is supported.

	t := &g.Tournament
	// DIFF: lord.js stores the [Tournament] values into `settings` instead
	// of `tournament_settings`, so they never take effect; they do here.
	if v, ok := in.get("Tournament", "Enabled"); ok {
		t.Enabled = iniBool(v, t.Enabled)
	}
	num("Tournament", "Days", &t.Days)
	num("Tournament", "Experience", &t.XP)
	num("Tournament", "DragonKills", &t.DKills)
	num("Tournament", "PlayerKills", &t.PKills)
	num("Tournament", "Level", &t.Level)
	num("Tournament", "Lays", &t.Lays)
	if v, ok := in.get("Tournament", "WinStat"); ok {
		found := false
		for i, n := range tournamentWinStats {
			if strings.EqualFold(n, v) {
				t.WinStat, found = int64(i), true
			}
		}
		if !found {
			if n, err := strconv.Atoi(v); err == nil && n >= 0 && n < len(tournamentWinStats) {
				t.WinStat = int64(n)
			}
		}
	}
}
