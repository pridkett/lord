// Package assets holds the distribution files of Synchronet's LORD
// (xtrn/lord at the revision in ref/synchronet/REVISION): the text screens
// in lordtxt.lrd, the looks tables and the NPC sayings.  A copy in the
// world directory takes precedence, as the exec directory did in lord.js.
package assets

import (
	"embed"
	"io/fs"
	"strings"
)

//go:embed data
var files embed.FS

// Read returns a built-in file by name (case-insensitive).
func Read(name string) ([]byte, error) {
	ents, err := files.ReadDir("data")
	if err != nil {
		return nil, err
	}
	for _, e := range ents {
		if strings.EqualFold(e.Name(), name) {
			return files.ReadFile("data/" + e.Name())
		}
	}
	return nil, fs.ErrNotExist
}

// Names lists the built-in files.
func Names() []string {
	ents, _ := files.ReadDir("data")
	var out []string
	for _, e := range ents {
		out = append(out, e.Name())
	}
	return out
}
