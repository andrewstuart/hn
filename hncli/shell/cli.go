package shell

import (
	"io"

	"github.com/andrewstuart/hn/hackernews"

	"code.google.com/p/goncurses"
)

type Cli struct {
	CurrentPage int
	Writer      io.Writer
	Window      *goncurses.Window
	Cache       hackernews.PageCache
}

func NewCli(w io.Writer) (*Cli, error) {
	c := &Cli{
		Writer: w,
	}

	if w == nil {
		if newWindow, err := goncurses.Init(); err == nil {
			c.Window = newWindow
		} else {
			return nil, err
		}
	}

	return c, nil
}
