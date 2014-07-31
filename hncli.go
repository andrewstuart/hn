package main

import (
	"log"
	"strings"

	"code.google.com/p/goncurses"
)

type hncli struct {
	root, main, help *goncurses.Window
	Height           int
	handler          CharHandler
	finished         bool
}

func (h hncli) Refresh() {
	h.root.Refresh()
}

const MENU_HEIGHT int = 3

var hc hncli
var instantiated = false

//Returns an instance of the CLI. Call once
func GetCli() hncli {
	if !instantiated {
		hc = hncli{}

		root, e := goncurses.Init()

		if e != nil {
			goncurses.End()
			log.Fatal(e)
		}

		h, w := root.MaxYX()

		hc.root = root
		hc.main = root.Sub(h-MENU_HEIGHT, w, 0, 0)
		hc.help = root.Sub(MENU_HEIGHT, w, h-MENU_HEIGHT, 0)

		hc.Height = h - MENU_HEIGHT

		instantiated = true
	}

	return hc
}

func (h hncli) SetContent(format string, args ...interface{}) {
	h.main.Printf(format, args...)
}

func (h hncli) SetHelp(text string) {
	h.help.Print(text)
}

//Scroll the content that was added with SetContent
func (h hncli) Scroll(amount int) {

}

type CharHandler func(string)

func (h hncli) SetKeyHandler(CharHandler) {
	h.handler = CharHandler
}

func (h hncli) Quit() {
	goncurses.End()
	h.finished = true
}

func (h hncli) Run() {
	for !finished {
		c := scr.GetChar()

		if c == 127 {
			c = goncurses.Key(goncurses.KEY_BACKSPACE)
		}

		ch := goncurses.KeyString(c)

		h.handler(ch)
	}
}

func (h hncli) getFitLines(s string) []string {
	_, w := h.main.MaxYX()

	a := strings.Split(s, "\n")

	p := make([]string, 0, len(a))

	//Newlines for stuff
	for _, line := range a {
		for len(line) > w {
			//Current line length
			l := w
			if l > len(line) {
				l = len(line)
			}

			//Find last space
			for line[l] != ' ' {
				l--
			}

			//Add substring to slice
			p = append(p, line[:l])

			line = line[l:]
		}

		p = append(p, line)
	}

	return p
}

func (hc hncli) paginate(t string, n int) string {
	h, _ := hc.main.MaxYX()

	lines := hc.getFitLines(t)

	nLines := len(lines)

	//If text won't fit
	if nLines > h {
		//Determine start point
		if n < nLines-h {
			//If n is not in last h elements, reslice starting at n
			lines = lines[n : n+h]
		} else {
			//Else, Get last h (at most) elements
			lines = lines[nLines-h:]
		}
	}

	return strings.Join(lines, "\n")
}
