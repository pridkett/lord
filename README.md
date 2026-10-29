# Legend of the Red Dragon — Go port

A Go port of **LORD 5.02 JS**, Synchronet's GPL JavaScript rewrite of the
classic BBS door *Legend of the Red Dragon* (`xtrn/lord/lord.js`, by Deuce),
as a native door for the Go port of Renegade BBS.

The original Turbo Pascal source of LORD has never been released (the rights
belong to Metropolis Gameport).  Synchronet's rewrite is the one public
source; it follows LORD 4.0x's text, screens and rules closely.  The port is a
function-by-function translation of it (see `docs/PORTING.md`), with the
JavaScript kept in `ref/synchronet` for reference.  See `NOTICE.md` for origin
and licence.

## Build

```sh
GOPRIVATE=github.com/pridkett/* go build -o lord ./cmd/lord
go test ./...                      # LORD_MONKEY=5m go test -run Monkey ./internal/game for a long soak
```

## Setting up a world

Unlike BRE and SRE, no original archive is needed: the text resources
(`lordtxt.lrd`, the looks and sayings tables, `lord.ini`) are built into the
binary.  A copy placed in the world directory overrides the built-in one, so a
sysop can still customise screens.

```sh
lord --dir /var/lib/doorgames/lord init    # creates the directory and lord.ini
```

The game creates `player.bin`, `state.bin`, the logs and mail files on first
play.  The data files are byte-compatible with Synchronet's LORD, so an
existing Synchronet LORD game directory can be used as a world as is.

## Command line

```
lord [flags] [command]

  play (default)  init  reset

  -d, --dir DIR        world directory ($LORD_DIR, else .)
      --context FILE   Renegade launch context ($RENEGADE_CONTEXT)
  -a, --alias NAME     local caller name ($DOOR_ALIAS, default Sysop)
  -m, --minutes N      local time limit ($DOOR_MINUTES, default 60)
      --version
```

lord.js's own arguments still work: `RESET`, `-P PREFIX`, `-M DIR`.
Settings live in `lord.ini` (see `ref/synchronet/Sysop.doc`): inactivity
days, fights per day, deeds to win, bank transfers, clean mode, tournaments.

## Renegade door definition

```json
{
  "id": "lord-native",
  "name": "Legend of the Red Dragon — native port",
  "command": ["/opt/doorgames/bin/lord"],
  "directory": "/var/lib/doorgames/lord",
  "world": "lord-native",
  "encoding": "cp437",
  "capacity": 8,
  "timeout_seconds": 1800,
  "access": "S10"
}
```

* No adapter or drop file: the game reads `RENEGADE_CONTEXT` and plays on the
  PTY, writing CP437.
* A warrior belongs to the BBS account (`board_id:player_id`), not the alias;
  the warrior's name is chosen in the game.
* LORD is multi-player: several callers share the world at once (record
  locks on `player.bin`, lock files for mail, state and duels), including
  online duels between two callers.

See `docs/RENEGADE_INTEGRATION.md` for wiring it into the Renegade demo.

## Differences from lord.js

* **Not ported:** RIPscrip graphics and the SyncTERM icon cache, the remote
  game server (`lordsrv.js`), and in-game modules (IGMs are JavaScript
  add-ons).  "Other Places" shows no add-ons.
* **Disconnects:** Renegade kills a door when the caller drops, so lord.js's
  exit handler can't run.  Sessions hold OS locks that the kernel releases
  when the process dies; the next session clears the stale "online" flag and
  any lock files the dead session left behind.
* A few lord.js bugs that crash or break play are fixed (blackjack dealing
  the dealer's cards to the player, the duel ending after the first strike,
  `[Tournament]` settings being ignored, `MaxTransferAmount`); the rest are
  kept and marked `// BUG (lord.js)` in the code.  `grep -rn 'DIFF:'
  internal/` lists every deliberate change.
