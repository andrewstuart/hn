package shell

import (
	"os"

	"code.google.com/p/goncurses"
)

func (c *Cli) Quit() {
	goncurses.End()
	os.Exit(0)
}
