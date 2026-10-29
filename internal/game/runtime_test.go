package game

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestPlayerRecordLength(t *testing.T) {
	_, n := parseRecFields(reflect.TypeFor[Player]())
	// Player_Def in recorddefs.js adds up to 807 bytes.
	if n != 807 {
		t.Fatalf("Player record is %d bytes, want 807", n)
	}
	if _, n := parseRecFields(reflect.TypeFor[State]()); n != 40 {
		t.Fatalf("State record is %d bytes, want 40", n)
	}
}

func TestRecordRoundTrip(t *testing.T) {
	g := New(t.TempDir(), nil)
	p := g.playerNew(false)
	if p.HP != 20 || p.Weapon != "Stick" || p.MarriedTo != -1 || !p.HighSpirits || p.Time != 2147483648 {
		t.Fatalf("defaults not applied: %+v", p)
	}
	p.Name, p.Gold, p.Exp, p.HP = "Tester", 3_000_000_000, -5, 40000
	p.Put()
	q := g.playerGet(0)
	if q.Name != "Tester" || q.Gold != 2147483647 || q.Exp != -5 || q.HP != 32767 {
		t.Fatalf("round trip: %+v", q)
	}
	if g.playerLength() != 1 {
		t.Fatalf("length %d", g.playerLength())
	}
	if filepath.Base(g.pfile.Name()) != "player.bin" {
		t.Fatal(g.pfile.Name())
	}
}

func TestTextHelpers(t *testing.T) {
	if got := prettyInt(1234567); got != "1,234,567" {
		t.Fatal(got)
	}
	if got := prettyInt(-123456); got != "-123,456" {
		t.Fatal(got)
	}
	if dispLen("`2Hi `r1there") != 8 {
		t.Fatal(dispLen("`2Hi `r1there"))
	}
	if got := cleanStr("a`2b`cc"); got != "a`2bc" {
		t.Fatal(got)
	}
	if n, ok := jsParseInt(" 42abc"); !ok || n != 42 {
		t.Fatal(n, ok)
	}
	if _, ok := jsParseInt("x"); ok {
		t.Fatal("NaN expected")
	}
}

func TestAssets(t *testing.T) {
	g := New(t.TempDir(), nil)
	g.buildTxtIndex()
	for _, s := range []string{"INTRO1", "HINTS", "ABDUL"} {
		if _, ok := g.txtindex[s]; !ok {
			t.Errorf("lordtxt.lrd has no %s", s)
		}
	}
	if l := g.looks("M", 1); l == "" {
		t.Error("no mlooks")
	}
}

func TestTables(t *testing.T) {
	if len(monsterStats) != 131 || len(trainerStats) != 12 || len(weaponStats) != 16 || len(armourStats) != 16 {
		t.Fatal(len(monsterStats), len(trainerStats), len(weaponStats), len(armourStats))
	}
	if weaponStats[15].Name != "Death Sword" || weaponStats[15].Price != 400000000 || armourStats[1].Name != "Coat" {
		t.Fatal(weaponStats[15], armourStats[1])
	}
}

func TestStaleMutex(t *testing.T) {
	dir := t.TempDir()
	a, b := New(dir, nil), New(dir, nil)
	name := filepath.Join(dir, "state.bin")

	// A lock file left by a killed session (old, no flock) is taken over.
	if err := os.WriteFile(name+".lock", []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	old := time.Now().Add(-time.Minute)
	_ = os.Chtimes(name+".lock", old, old)
	start := time.Now()
	a.fmutex(name)
	if time.Since(start) > 2*time.Second {
		t.Fatal("stale lock was not taken over")
	}

	// A live holder's lock is respected, even once it is old.
	_ = os.Chtimes(name+".lock", old, old)
	if !mutexHeld(name) || staleMutex(name+".lock") {
		t.Fatal("live lock reported stale")
	}
	a.rfmutex(name)
	if mutexHeld(name) {
		t.Fatal("lock not released")
	}
	b.fmutex(name)
	b.rfmutex(name)
}
