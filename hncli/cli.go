package hncli

import (
	"fmt"
	"io"

	"github.com/andrewstuart/hn/hackernews"

	"code.google.com/p/goncurses"
)

type Cli struct {
	CurrentPage int
	Writer      io.Writer
	Screen      goncurses.Screen
	Cache       hackernews.PageCache
}

func (c *Cli) SetContent(content string) error {
	if c.Writer != nil {
		_, err := fmt.Fprint(c.Writer, content)

		if err != nil {
			return fmt.Errorf("Error writing to output: %v", err)
		}
	}

	return nil
}
