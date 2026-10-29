package game

import (
	"time"
)

// Player is a Player_Def record (player.bin) and, as in lord.js where every
// combatant is a duck-typed object, also any opponent: monsters, trainers
// and the dragon use the same struct with the fields they need.  Fields
// with a `rec` tag are stored; the rest live only for the session.
type Player struct {
	Name             string `rec:"String:20"`
	RealName         string `rec:"String:50"`
	HP               int64  `rec:"SignedInteger16,def=20"`
	Divorced         bool   `rec:"Boolean"`
	Amulet           bool   `rec:"Boolean"`
	Gone             int64  `rec:"SignedInteger16"`
	HPMax            int64  `rec:"SignedInteger16,def=20"`
	WeaponNum        int64  `rec:"SignedInteger8,def=1"`
	Weapon           string `rec:"String:20,def=Stick"`
	SeenMaster       bool   `rec:"Boolean"`
	ForestFights     int64  `rec:"SignedInteger16,def=15"`
	PvPFights        int64  `rec:"SignedInteger16"`
	Gold             int64  `rec:"SignedInteger,def=500"`
	Bank             int64  `rec:"SignedInteger"`
	Def              int64  `rec:"SignedInteger16,def=1"`
	Str              int64  `rec:"SignedInteger16,def=10"`
	Cha              int64  `rec:"SignedInteger16,def=1"`
	SeenDragon       bool   `rec:"Boolean"`
	SeenViolet       bool   `rec:"Boolean"`
	Level            int64  `rec:"SignedInteger8,def=1"`
	Time             int64  `rec:"Integer,def=2147483648"`
	Arm              string `rec:"String:20,def=Coat"`
	ArmNum           int64  `rec:"SignedInteger8,def=1"`
	Dead             bool   `rec:"Boolean"`
	Inn              bool   `rec:"Boolean"`
	Gem              int64  `rec:"SignedInteger16"`
	Exp              int64  `rec:"SignedInteger,def=1"`
	Sex              string `rec:"String:1,def=M"`
	SeenBard         bool   `rec:"Boolean"`
	LastReincarnated int64  `rec:"Integer"`
	Laid             int64  `rec:"SignedInteger16"`
	GotDelicious     bool   `rec:"Boolean"`
	OnNow            bool   `rec:"Boolean"`
	TimeOn           string `rec:"String:5,def=00:00"`
	Clss             int64  `rec:"SignedInteger8"`
	Horse            bool   `rec:"Boolean"`
	Kids             int64  `rec:"SignedInteger16"`
	DragKills        int64  `rec:"SignedInteger16"`
	SkillW           int64  `rec:"SignedInteger8"`
	SkillM           int64  `rec:"SignedInteger8"`
	SkillT           int64  `rec:"SignedInteger8"`
	LevelW           int64  `rec:"SignedInteger8"`
	LevelM           int64  `rec:"SignedInteger8"`
	LevelT           int64  `rec:"SignedInteger8"`
	MarriedTo        int64  `rec:"SignedInteger16,def=-1"`
	TransferredGold  int64  `rec:"SignedInteger"`
	PvP              int64  `rec:"SignedInteger16"`
	Weird            bool   `rec:"Boolean"`
	HighSpirits      bool   `rec:"Boolean,def=true"`
	Flirted          bool   `rec:"Boolean"`
	Olivia           bool   `rec:"Boolean"`
	Asshole          bool   `rec:"Boolean"`
	OliviaCount      int64  `rec:"SignedInteger8"`
	DoneTower        bool   `rec:"Boolean"`
	HasDes           bool   `rec:"Boolean"`
	Des1             string `rec:"String:255"`
	Des2             string `rec:"String:255"`
	LeftBank         bool   `rec:"Boolean"`
	Padding          string `rec:"String:99"`

	// Record is the record number in player.bin (-1 for a non-player).
	Record int

	// Session-only properties lord.js sets on the player object.
	Expert             bool
	HasFairy           bool
	FairyLore          bool
	LightShield        bool
	KilledAPlayer      bool
	MagicallyDelicious bool
	OliviaAsshole      bool
	RanAway            bool

	// Opponent properties (monsters, trainers, the dragon).
	Death     string
	IsArena   bool
	IsDragon  bool
	Swear     string
	NeedStr1  string
	NeedStr2  string
	HPGained  int64
	StrGained int64
	Need      int64

	pf *RecordFile[Player]
}

// State is State_Def (state.bin): the game-wide record.
type State struct {
	LatestHero      string `rec:"String:20,def=Master Turgon"`
	MarriedToSeth   int64  `rec:"SignedInteger16,def=-1"`
	MarriedToViolet int64  `rec:"SignedInteger16,def=-1"`
	Days            int64  `rec:"Integer"`
	WonBy           int64  `rec:"SignedInteger16,def=-1"`
	LogDate         int64  `rec:"Integer"`
	LastBar         int64  `rec:"SignedInteger16"`
	ForestGold      int64  `rec:"SignedInteger,def=100"`

	sf *RecordFile[State]
}

// WeaponArmour is a Weapon_Armor_Def record (weapon.bin, armour.bin) or an
// entry of the built-in tables.
type WeaponArmour struct {
	Name    string `rec:"String:20"`
	Price   int64  `rec:"Integer"`
	Num     int64  `rec:"Integer"`
	IsArena bool
}

// trainerRec is Trainer_Def (trainers.bin).
type trainerRec struct {
	Name      string `rec:"String:30"`
	Weapon    string `rec:"String:20,def=Fists"`
	Swear     string `rec:"String:50,def=Oh sugar!"`
	NeedStr1  string `rec:"String:50,def=You're really improving!"`
	NeedStr2  string `rec:"String:50,def=Come at me!"`
	Str       int64  `rec:"SignedInteger16,def=5"`
	HP        int64  `rec:"SignedInteger16,def=1"`
	HPGained  int64  `rec:"SignedInteger16,def=1"`
	StrGained int64  `rec:"SignedInteger16,def=1"`
	Def       int64  `rec:"SignedInteger16,def=1"`
	Need      int64  `rec:"SignedInteger,def=1"`
}

// monsterRec is Monster_Def (lenemy.bin, dragon.bin).
type monsterRec struct {
	Name   string `rec:"String:60"`
	Str    int64  `rec:"SignedInteger,def=5"`
	Gold   int64  `rec:"SignedInteger,def=500"`
	Weapon string `rec:"String:60,def=Fists"`
	Exp    int64  `rec:"SignedInteger,def=1"`
	HP     int64  `rec:"SignedInteger,def=1"`
	Death  string `rec:"String:100,def=Aw Nuts"`
}

// Put writes the player's record (player.put()).
func (p *Player) Put() {
	if p.pf == nil {
		return
	}
	if err := p.pf.Write(p.Record, p); err != nil {
		panic(err)
	}
}

// ReLoad re-reads the stored fields from player.bin (player.reLoad()),
// keeping the session-only ones.
func (p *Player) ReLoad() {
	if p.pf == nil {
		return
	}
	if err := p.pf.Read(p.Record, p); err != nil {
		panic(err)
	}
}

// Lock holds the player's record lock until UnLock (player.lock()).
func (p *Player) Lock() {
	if p.pf != nil {
		p.pf.Lock(p.Record, -1)
	}
}

// UnLock releases the record lock (player.unLock()).
func (p *Player) UnLock() {
	if p.pf != nil {
		p.pf.UnLock(p.Record)
	}
}

// ReInit resets every stored field to its default (player.reInit()).
func (p *Player) ReInit() {
	if p.pf != nil {
		p.pf.Defaults(p)
	}
}

// Put writes state.bin (state.put()).
func (s *State) Put() {
	if err := s.sf.Write(0, s); err != nil {
		panic(err)
	}
}

func (g *Game) pfileInit() {
	if g.pfile != nil {
		return
	}
	pf, err := OpenRecordFile[Player](g.gamedir("player.bin"))
	if err != nil {
		panic(err)
	}
	g.pfile = pf
}

// playerNew appends a fresh player record (player_new); with leaveLocked
// the record stays locked.
func (g *Game) playerNew(leaveLocked bool) *Player {
	g.pfileInit()
	rec := g.pfile.Length()
	if !g.pfile.Lock(rec, time.Second) {
		return nil
	}
	p := &Player{Record: rec, pf: g.pfile}
	g.pfile.Defaults(p)
	p.Put()
	if !leaveLocked {
		g.pfile.UnLock(rec)
	}
	return p
}

// playerGet reads player record rec (player_get), or nil when there is none.
func (g *Game) playerGet(rec int64) *Player {
	g.pfileInit()
	if rec < 0 || int(rec) >= g.pfile.Length() {
		return nil
	}
	p := &Player{Record: int(rec), pf: g.pfile}
	if err := g.pfile.Read(int(rec), p); err != nil {
		panic(err)
	}
	return p
}

// playerGetLocked reads rec and keeps its lock (player_get(rec, true)).
func (g *Game) playerGetLocked(rec int64) *Player {
	g.pfileInit()
	if rec < 0 || int(rec) >= g.pfile.Length() {
		return nil
	}
	g.pfile.Lock(int(rec), -1)
	return g.playerGet(rec)
}

// playerLength is the number of player records (player_length).
func (g *Game) playerLength() int {
	g.pfileInit()
	return g.pfile.Length()
}

// allPlayers reads every player record (all_players).
func (g *Game) allPlayers() []*Player {
	n := g.playerLength()
	out := make([]*Player, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, g.playerGet(int64(i)))
	}
	return out
}
