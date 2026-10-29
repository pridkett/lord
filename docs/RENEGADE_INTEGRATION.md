# Wiring native LORD into Renegade

What a Renegade agent needs to add `lord-native` next to `bre-native` and
`sre`.  Nothing here has been applied to the Renegade repository yet.

## How LORD differs from BRE/SRE

| | BRE / SRE | LORD |
|---|---|---|
| Original archive | required (commercial files) | **none**: text resources are built into the binary |
| World setup | unpack + reset | `lord --dir WORLD init` (idempotent; only writes `lord.ini` if missing) |
| World "ready" check | `DATA/GAME.DAT`, `DATA/GALAXY.II`, ... | the directory exists (data files are created on first play) |
| Capacity | 1 | multi-user: `capacity` 8, all ids on the same world share `world` |
| Normal exit codes | `[1]` / `[1,2]` | none needed: exits 0 after quit, death, time-up and hang-up; 1 only for real errors |
| Identity | alias | `board_id:player_id` (alias changes keep the warrior) |
| Encoding | cp437 | cp437 |

## Changes to make in Renegade

1. **Installer** (`games/native-realms/install.py`): add a `lord` game.
   Build with `go install github.com/pridkett/lord/cmd/lord@<tag>` (same
   `GOPRIVATE` handling as bre/sre) or take a prebuilt binary from
   `--bin-dir`; setup is `lord --dir "$ROOT/lord-native" init`.  No archive
   option or hash check.  Never run `reset` on an existing world.
2. **Demo discovery** (`cmd/renegade/native_realms.go`): add an entry for
   binary `lord`, id `lord-native`, ready when the world directory exists,
   with `Capacity: 8`, no `NormalExitCodes`.  The current loop hard-codes
   `Capacity: 1`; make it a per-game field.
3. **play.sh / play images**: register `lord-native` when the binary and
   world exist, like the other native doors.  The existing `lord` id
   (talisto/LORD.ts, Node) can stay as it is or be replaced by the native
   port; that is the owner's call.  If replaced, keep the old world: LORD.ts
   uses SQLite and its saves don't carry over.
4. **Menus**: a `door` item targeting `lord-native` in
   `examples/play/menu.json` and the default games menu.
5. **Private image**: the binary is built from a private repository, so it
   belongs only in the private `play-native` image, like bre/sre.

## Verification

1. `lord --dir WORLD init` twice: the second run changes nothing.
2. Browser caller: title screen with correct CP437 art, create a warrior,
   fight in the forest, quit → back at the BBS menu with no error.
3. Same account from `renegade local` → the same warrior.
4. Two callers at once: both appear under (P)eople Online; one can attack
   or duel the other.
5. Disconnect mid-session, call again: the warrior is not stuck "online" and
   nothing waits on a stale lock.
6. Let the time limit run out: "GOD STRIKES YOU UNCONSCIOUS", back to the BBS.
