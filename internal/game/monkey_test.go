package game

import (
	"io"
	"math/rand/v2"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/pridkett/doorgames-common/door"
)

// keyFeeder is an endless stream of random keypresses, weighted towards
// the game's menu keys, Enter and space.
type keyFeeder struct {
	rng  *rand.Rand
	stop <-chan struct{}
}

const monkeyKeys = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890?" + "\r\r\r\r\r    YYYNNN" + "AAAALLLLFFFFRRRRQQ"

func (k *keyFeeder) Read(p []byte) (int, error) {
	select {
	case <-k.stop:
		return 0, io.EOF
	case <-time.After(time.Millisecond):
	}
	c := monkeyKeys[k.rng.IntN(len(monkeyKeys))]
	if k.rng.IntN(2) == 0 {
		c = byte(strings.ToLower(string(c))[0])
	}
	p[0] = c
	return 1, nil
}

// TestMonkey plays random keys against a fresh world for a while and fails
// on any panic other than the game's own exit.  LORD_MONKEY sets the
// duration (default 20s).
func TestMonkey(t *testing.T) {
	d := 20 * time.Second
	if v := os.Getenv("LORD_MONKEY"); v != "" {
		d, _ = time.ParseDuration(v)
	}
	if testing.Short() {
		d = 3 * time.Second
	}
	dir := t.TempDir()
	end := time.Now().Add(d)
	for seed := uint64(1); time.Now().Before(end); seed++ {
		promote(dir, rand.New(rand.NewPCG(seed, 99)))
		for _, alias := range []string{"Ann", "Ben", "Cat"} {
			ctx := door.Local(alias)
			ctx.PlayerID = int64(len(alias)) + int64(alias[0])
			ctx.SecondsLeft = 3600
			stop := make(chan struct{})
			sess, err := door.Open(ctx, &keyFeeder{rng: rand.New(rand.NewPCG(seed, uint64(alias[0]))), stop: stop}, io.Discard, door.Options{})
			if err != nil {
				t.Fatal(err)
			}
			g := New(dir, sess)
			done := make(chan error, 1)
			go func() {
				_, err := g.Run()
				done <- err
			}()
			select {
			case err := <-done:
				if err != nil {
					t.Fatalf("seed %d %s: %v", seed, alias, err)
				}
			case <-time.After(time.Duration(2+seed%10) * time.Second):
				// Still playing: hang up, as Renegade would.
				close(stop)
				if err := <-done; err != nil {
					t.Fatalf("seed %d %s after hangup: %v", seed, alias, err)
				}
				stop = nil
			}
			if stop != nil {
				close(stop)
			}
			sess.Close()
		}
	}
}

// promote gives the world's warriors random advanced stats so the monkey
// reaches the content that needs them (the Dark Horse tavern, skills, the
// dragon, gambling, romance).
func promote(dir string, r *rand.Rand) {
	g := New(dir, nil)
	for _, p := range g.allPlayers() {
		if r.IntN(3) == 0 {
			continue
		}
		p.Level = 1 + r.Int64N(12)
		p.Exp = r.Int64N(50_000_000)
		p.Gold = r.Int64N(100_000_000)
		p.Bank = r.Int64N(100_000_000)
		p.Gem = r.Int64N(50)
		p.Cha = r.Int64N(120)
		p.HPMax = 20 + r.Int64N(3000)
		p.HP = p.HPMax
		p.Str = 10 + r.Int64N(2000)
		p.Def = 1 + r.Int64N(1000)
		p.SkillW, p.SkillM, p.SkillT = r.Int64N(41), r.Int64N(41), r.Int64N(41)
		p.LevelW, p.LevelM, p.LevelT = r.Int64N(8), r.Int64N(8), r.Int64N(8)
		p.WeaponNum, p.ArmNum = 1+r.Int64N(15), 1+r.Int64N(15)
		p.Horse = r.IntN(2) == 0
		p.Amulet = r.IntN(4) == 0
		p.DragKills = r.Int64N(4)
		p.Laid = r.Int64N(30)
		p.ForestFights, p.PvPFights = 15, 3
		p.Dead, p.SeenMaster, p.SeenDragon, p.SeenViolet, p.SeenBard, p.Flirted = false, false, false, false, false, false
		p.Clss = 1 + r.Int64N(3)
		p.Put()
	}
	g.pfile.Close()
}
