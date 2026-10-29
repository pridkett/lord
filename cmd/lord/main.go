// Command lord is the Go port of LORD 5.02 JS, Synchronet's rewrite of
// Legend of the Red Dragon.
package main

import (
	"fmt"
	"os"

	"github.com/pridkett/doorgames-common/door"
	"github.com/pridkett/lord/internal/game"
)

func main() {
	ctx, err := door.LoadContext()
	if err != nil {
		fmt.Fprintln(os.Stderr, "lord:", err)
		os.Exit(1)
	}
	sess, err := door.Open(ctx, os.Stdin, os.Stdout, door.Options{})
	if err != nil {
		fmt.Fprintln(os.Stderr, "lord:", err)
		os.Exit(1)
	}
	dir, _ := os.Getwd()
	g := game.New(dir, sess)
	code, err := g.Run()
	sess.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "lord:", err)
	}
	os.Exit(code)
}
