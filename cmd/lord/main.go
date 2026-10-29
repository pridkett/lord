// Command lord is the Go port of LORD 5.02 JS, Synchronet's rewrite of
// Legend of the Red Dragon, as a Renegade door.
//
// Renegade runs it on a PTY with the world as working directory and the
// caller in RENEGADE_CONTEXT (see common/door).  Without it, a local session
// is described by --alias and --minutes.
package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pridkett/doorgames-common/door"
	"github.com/pridkett/lord/internal/assets"
	"github.com/pridkett/lord/internal/game"
)

func main() { os.Exit(run(os.Args[1:])) }

func run(args []string) int {
	o, code, ok := parseArgs(args)
	if !ok {
		return code
	}
	dir, err := filepath.Abs(o.dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "lord:", err)
		return 1
	}
	if o.igm {
		fmt.Fprintln(os.Stderr, "lord: in-game modules (IGMs) are not supported by this port")
		return 1
	}
	switch o.cmd {
	case "init":
		return initWorld(dir)
	case "reset":
		g := game.New(dir, nil)
		g.LoadSettings(o.prefix, o.menu)
		g.Reset()
		fmt.Println("The game has been reset.")
		return 0
	}

	if st, err := os.Stat(dir); err != nil || !st.IsDir() {
		fmt.Fprintf(os.Stderr, "lord: world directory %s does not exist (run 'lord init')\n", dir)
		return 1
	}
	ctx, err := door.LoadContext()
	if err != nil {
		fmt.Fprintln(os.Stderr, "lord:", err)
		return 1
	}
	sess, err := door.Open(ctx, os.Stdin, os.Stdout, door.Options{})
	if err != nil {
		fmt.Fprintln(os.Stderr, "lord:", err)
		return 1
	}
	g := game.New(dir, sess)
	g.LoadSettings(o.prefix, o.menu)
	code, err = g.Run()
	sess.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "lord:", err)
	}
	return code
}

// initWorld creates dir with the default lord.ini.
func initWorld(dir string) int {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "lord:", err)
		return 1
	}
	ini := filepath.Join(dir, "lord.ini")
	if _, err := os.Stat(ini); errors.Is(err, fs.ErrNotExist) {
		b, _ := assets.Read("lord.ini")
		if err := os.WriteFile(ini, b, 0o644); err != nil {
			fmt.Fprintln(os.Stderr, "lord:", err)
			return 1
		}
	}
	fmt.Println("World ready in", dir)
	return 0
}
