package game

import (
	"io"
	"os"
	"syscall"
)

// DIFF: Renegade ends a door with SIGKILL when the caller disconnects, so
// lord.js's exit handler (player.on_now = false) never runs and the warrior
// would stay "online" until someone else happened to play on the same node
// (SafeNode).  Each session therefore holds an OS lock on its player's byte
// of online.lck; the kernel drops it when the process dies, and every new
// session clears on_now for warriors whose byte nobody holds.

// holdOnline takes the liveness lock for the current player until exit.
func (g *Game) holdOnline() {
	if g.player == nil || g.online != nil {
		return
	}
	f, err := os.OpenFile(g.gamedir("online.lck"), os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return
	}
	lk := syscall.Flock_t{Type: syscall.F_WRLCK, Whence: io.SeekStart, Start: int64(g.player.Record), Len: 1}
	if syscall.FcntlFlock(f.Fd(), syscall.F_SETLK, &lk) != nil {
		_ = f.Close()
		return
	}
	g.online = f
}

// sweepOffline clears on_now for warriors with no live session.
func (g *Game) sweepOffline() {
	f, err := os.OpenFile(g.gamedir("online.lck"), os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return
	}
	defer func() {
		// Closing f would drop this process's own liveness lock (POSIX
		// locks belong to the process), so keep it open while one is held.
		if g.online == nil {
			_ = f.Close()
		}
	}()
	for i, p := range g.allPlayers() {
		if !p.OnNow || (g.player != nil && i == g.player.Record) {
			continue
		}
		lk := syscall.Flock_t{Type: syscall.F_WRLCK, Whence: io.SeekStart, Start: int64(i), Len: 1}
		if syscall.FcntlFlock(f.Fd(), syscall.F_GETLK, &lk) != nil || lk.Type != syscall.F_UNLCK {
			continue
		}
		p.Lock()
		p.ReLoad()
		p.OnNow = false
		p.Put()
		p.UnLock()
	}
}
