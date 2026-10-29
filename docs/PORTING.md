# Porting conventions

The port is a function-by-function translation of `ref/synchronet/lord.js`
(LORD 5.02 JS, Synchronet revision in `ref/synchronet/REVISION`) into the
single package `internal/game`.  Stay faithful: same text, same colours,
same prompts, same order of output, same random rolls, same rules.

## Names and structure

* Globals of lord.js are fields of `*Game` (`g`).  `player` is `g.player`,
  `state` is `g.state`, `settings.forest_fights` is `g.Settings.ForestFights`
  (see `settings.go`), `tournament_settings` is `g.Tournament`.
* Every top-level JS function becomes a method on `*Game` named in
  camelCase: `check_mail` -> `g.checkMail`, `ye_old_bank` -> `g.yeOldBank`,
  `add_hp` -> `g.addHP`.  Nested JS functions become methods named
  `<outer><Inner>` (e.g. forest's `darkhorse_tavern` -> `g.forestDarkhorseTavern`)
  or closures, whichever reads better.  Keep the original order of functions
  in a file and put a `// <js_name>` reference in the doc comment when the
  name changed beyond camel-casing.
* Player_Def properties are `Player` fields (`player.hp_max` -> `p.HPMax`,
  `drag_kills` -> `DragKills`, `pvp_fights` -> `PvPFights`); see `player.go`.
  Opponents (monsters, trainers, the dragon, other players) are `*Player`
  too, exactly like the JS duck typing.  `op.Record` is an `int`.
* All JS numbers that are integers are `int64`.  Where JS division produces a
  fraction that matters (comparisons, later multiplication), use `float64`
  locally; truncate (`int64(x)`) exactly where JS truncates (`parseInt`,
  `random()`, storing into a record field).  `g.random(n int64)` is
  Synchronet's `random()` (0..n-1, 0 when n < 2); `g.randomf(x float64)`
  is `random()` of a fractional number.
* `parseInt(s, 10)` is `jsParseInt(s)` (returns value, ok); NaN cases are
  usually 0 or treated as invalid input, keep the JS behaviour.
* `x.toUpperCase()` on a key: `strings.ToUpper(g.getkey())`.
* `exit(n)` is `g.exit(n)` (unwinds via panic; the on-exit handler saves the
  player).  `throw new Error(msg)` is `g.fatalf(...)`.
* RIP: `g.rip` is always false.  Drop `if (rip) {...}` branches entirely and
  keep the non-RIP `else` branches.  Drop `lrdrip(...)` calls.
* psock (remote game server): drop every `psock !== undefined` branch and
  keep the local-file branch.
* IGMs (in-game modules) are not supported: drop `handle_igm`, the `out*.lrd`
  IGM command lines, `3rdparty.lrd`/`3rdalt.lrd` "other places" lists
  (show the menu as if no IGMs were installed).
* JS bugs: keep the observable behaviour and add a `// BUG (lord.js): ...`
  comment.  Fix only clear crashes/NaN corruption, with a `// DIFF: ...`
  comment.  One known pattern: `fmutex(name)` creates `name+'.lock'`, so a
  JS `file_exists(name)` check for a file only ever created by fmutex never
  matches; check `name+".lock"` and mark it `// DIFF:`.
* JS `// DIFF:`/`// TODO:` comments of the original are worth keeping.

## Runtime (already ported; see `game.go`, `text.go`, `console.go`)

| lord.js | Go |
|---|---|
| `sw(s)`, `sln(s)`, `lw(s, ext)`, `lln(s, ext)` | `g.sw`, `g.sln`, `g.lw(s, ext...)`, `g.lln(s, ext...)` |
| `foreground(n)`, `background(n)`, `sclrscr()` | same names |
| `more()`, `more_nomail()` | `g.more()`, `g.moreNomail()` |
| `getkey()` | `g.getkey() string` (one character) |
| `dk.console.getkey()` / `waitkey(ms)` | `g.con.GetKey()` ("" = undefined), `g.con.WaitKey(ms)` |
| `dk.console.getstr({len:20, integer:true, edit:...})` | `g.consoleGetstr(GetStrOpts{Len: 20, Integer: true, Edit: ...})` |
| `getstr(x, y, len, c, c1, str, mode)` | `g.getstr(x, y, len, c, c1, str, GetStrOpts{})` |
| `dk.console.gotoxy(x, y)` | `g.con.GotoXY(x, y)` |
| `dk.console.print(s)` / `cleareol()` | `g.con.Print(s)` / `g.con.ClearEOL()` |
| `dk.console.attr.value` | `g.con.Attr` |
| `dk.console.rows` / `cols` | `g.con.Rows` / `g.con.Cols` |
| `dk.connection.node` | `g.Ctx.Node` |
| `dk.user.full_name` | `g.Ctx.PlayerKey()` (stable per BBS account) |
| `dk.user.seconds_remaining + seconds_remaining_from - time()` | `g.secondsRemaining()` |
| `curlinenum`, `curcolnum`, `morechk` | `g.curlinenum`, `g.curcolnum`, `g.morechk` |
| `pretty_int`, `disp_len`, `disp_str`, `center`, `space_pad`, `clean_str` | `prettyInt`, `dispLen`, `dispStr`, `center`, `spacePad`, `cleanStr` |
| `looks(sex, charm)` | `g.looks(sex, charm)` |
| `lrdfile(name, more)` | `g.lrdfile(name, more...)` |
| `display_file(f, canskip, more)` | `g.displayFile(f, canskip, more)` (pass the JS defaults: true, true) |
| `show_lord_file`, `show_lord_buffer` | `g.showLordFile`, `g.showLordBuffer` |
| `check_mail()`, `mail_check()` | `g.checkMail()`, `g.mailCheck()` |
| `player_get(n)`, `player_new()`, `player_length()`, `all_players()` | `g.playerGet(int64) *Player`, `g.playerNew(bool)`, `g.playerLength()`, `g.allPlayers()` |
| `player_get(n, true)` (keep locked) | `g.playerGetLocked(n)` |
| `p.put()`, `p.reLoad()`, `p.lock()`, `p.unLock()`, `p.reInit()` | `p.Put()`, `p.ReLoad()`, `p.Lock()`, `p.UnLock()`, `p.ReInit()` |
| `get_state(lock)`, `put_state()`, `state.put()` | `g.getState(lock)`, `g.putState()`, `g.state.Put()` |
| `gamedir(f)` | `g.gamedir(f)` |
| `game_or_exec(f)` | text files: `g.assetLines(f)` / `g.assetExists(f)`; binary tables: `g.gameOrExec(f)` ("" = use the built-in table) |
| `fmutex(f, str)`, `rfmutex(f)` | `g.fmutex(f, str...)`, `g.rfmutex(f)` |
| `file_exists`, `file_remove`, `file_rename`, `file_date` | `fileExists`, `fileRemove`, `fileRename`, `fileDate` |
| `new File(f)` + `readln` loop | `readLines(f) ([]string, ok)` |
| `open('a')` + `writeln` | `appendLines(f, lines...)` |
| `open('w')` + `writeln` | `writeLines(f, lines...)` |
| `time()` | `g.time()` (Unix seconds); `new Date()` -> `g.now()` |
| `mswait(ms)` | `g.mswait(ms)` |
| `strftime`/`format` | `fmt.Sprintf` / `time.Format` |
| Monster_Def/Trainer_Def records | `readRecordAt[monsterRec](path, n)`, `readRecordAt[trainerRec](...)` |
| `monster_stats`, `trainer_stats`, `armour_stats`, `weapon_stats` | `monsterStats`, `trainerStats`, `armourStats`, `weaponStats` (copy before mutating) |
| `castles` / `which_castle` | `castles`, `g.whichCastle` |
| `ver` | `Version` |
| `cond ? a : b` | `pick(cond, a, b)` |
| record number in a file name | `itoa(p.Record)` |

## Ownership while porting

Each group owns the functions in its line ranges of lord.js and writes them
to its own file(s).  Functions another group calls are already declared in
`stubs_<group>.go` with the signature to keep; the owner deletes its stub
file when its real implementations exist.  A function only used inside a
group may get any signature the owner likes.
