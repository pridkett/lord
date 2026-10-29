package game

import "testing"

// Two games on one world exchange a battle response through war<rec>.bin.
func TestDuelResponseFiles(t *testing.T) {
	dir := t.TempDir()
	a, b := New(dir, nil), New(dir, nil)
	a.player = &Player{Record: 0}
	b.player = &Player{Record: 1}

	a.obSendResp(b.player, "42")
	// A second response before the first is read is dropped (file_mutex).
	a.obSendResp(b.player, "7")
	if !fileExists(b.gamedir("war1.bin")) {
		t.Fatal("war1.bin not written")
	}
	if got := b.obReadMsg(b.gamedir("war1.bin")); got != "42" {
		t.Fatalf("obReadMsg = %q, want 42", got)
	}
	if fileExists(b.gamedir("war1.bin")) {
		t.Fatal("war1.bin not removed after reading")
	}

	b.fmutex(b.gamedir("fight1.rec"))
	if !fileMutexExists(a.gamedir("fight1.rec")) {
		t.Fatal("fight lock not visible to the other game")
	}
	b.obCleanup()
	if fileMutexExists(a.gamedir("fight1.rec")) {
		t.Fatal("obCleanup left the fight lock")
	}
}
