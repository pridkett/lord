package main

import (
	"flag"
	"os"
	"strings"

	"github.com/pridkett/doorgames-common/cli"
	"github.com/pridkett/doorgames-common/door"
	"github.com/pridkett/lord/internal/game"
)

// options are the command-line settings.
type options struct {
	cmd    string // the command word: play, reset, init
	dir    string // world directory (--dir, LORD_DIR, ".")
	prefix string // -P: lord.js's game data prefix
	menu   string // -M: lord.js's menu directory
	igm    bool   // +IGM / -IGM given (not supported)
}

func newApp(o *options) *cli.App {
	fs := flag.NewFlagSet("lord", flag.ContinueOnError)
	fs.StringVar(&o.dir, "dir", "", "world `dir`: player.bin, state.bin, lord.ini, logs, mail\n(default $LORD_DIR, else the current directory)")
	door.AddFlags(fs)
	return &cli.App{
		Name:    "lord",
		Title:   "Legend of the Red Dragon " + game.Version + " — Go port of Synchronet's LORD",
		Version: game.Version,
		Description: `Run with no command, lord plays one door session for the caller Renegade
names in RENEGADE_CONTEXT, or a local session from a shell.  Several
callers may play the same world at once.  The other commands are the
sysop's, run from a shell.`,
		Flags:   fs,
		Short:   map[string]string{"d": "dir", "a": "alias", "m": "minutes"},
		Default: "play",
		Commands: []*cli.Command{
			{Name: "play", Summary: "Play a door session",
				Help: `Plays one session.  Under Renegade the caller comes from RENEGADE_CONTEXT;
from a shell it is a local session for --alias (default Sysop) with
--minutes of time (default 60).  A player is tied to the BBS account
(board and user number), not the alias; the warrior's name is chosen in
the game.`},
			{Name: "init", Summary: "Create a world directory with the default lord.ini",
				Help: `Creates the world directory if needed and writes the default lord.ini
(only when there is none).  A world needs nothing else: the game creates
its data files on first play.  Never touches an existing game.`},
			{Name: "reset", Summary: "Start a new game (deletes all players)",
				Help: `Deletes state.bin and player.bin, like lord.js's RESET: every warrior is
gone and the next caller starts day one of a new game.  lord.ini, the
logs and the flower garden are kept.`},
		},
		Env: append([]cli.Env{
			{Name: "LORD_DIR", Help: "world directory; see --dir"},
		}, door.Env...),
		Examples: []string{
			"lord -d /var/lib/doorgames/lord init   set up a new world",
			"lord --alias Tester --minutes 30       play locally",
			"lord reset                             start over in the world in .",
		},
		Notes: `lord.js's own arguments still work: RESET (any case), -P PREFIX (game
data prefix) and -M DIR (menu directory).  +IGM and -IGM are accepted
but in-game modules (JavaScript add-ons) are not supported.`,
		Passthrough: func(arg string) bool {
			switch strings.ToUpper(arg) {
			case "+IGM", "-IGM", "-P", "-M":
				return true
			}
			return false
		},
	}
}

// parseArgs parses the command line.  ok is false when the program should
// exit with code (after help, --version or a usage error).
func parseArgs(args []string) (o options, code int, ok bool) {
	// -P and -M take a value; pull them out before the flag parser sees
	// the value as a command word.
	var rest []string
	for i := 0; i < len(args); i++ {
		switch strings.ToUpper(args[i]) {
		case "-P", "-M", "+IGM", "-IGM":
			if i+1 < len(args) {
				switch strings.ToUpper(args[i]) {
				case "-P":
					o.prefix = args[i+1]
				case "-M":
					o.menu = args[i+1]
				default:
					o.igm = true
				}
				i++
				continue
			}
		}
		rest = append(rest, args[i])
	}
	app := newApp(&o)
	res, err := app.Parse(rest)
	if err != nil {
		return o, cli.ExitCode(err), false
	}
	o.cmd = res.Command.Name
	if o.dir == "" {
		o.dir = os.Getenv("LORD_DIR")
	}
	if o.dir == "" {
		o.dir = "."
	}
	return o, 0, true
}
