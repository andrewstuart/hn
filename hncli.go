package main

import (
	"log"
	"strings"

	"code.google.com/p/goncurses"
)

type hncli struct {
	root, main, help  *goncurses.Window
	Height, pos       int
	handler           CharHandler
	finished          bool
	helpText, content string
}

func (h *hncli) Refresh() {
	h.root.Refresh()
	h.root.Move(h.Height-1, 0)
}

const MENU_HEIGHT int = 3

var hc hncli

//Ha.
var singleDone = false

//Returns an instance of the CLI. Call once
func GetCli() hncli {
	if !singleDone {
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

		singleDone = true
	}

	return hc
}

func (h *hncli) SetContent(text string) {
	h.content = text
	h.main.Clear()

	h.main.Printf(text)
	h.main.Refresh()
}

func (h *hncli) SetHelp(text string) {
	h.helpText = text
	h.help.Clear()
	h.help.Print(text)
	h.help.Refresh()
}

//Scroll the content that was added with SetContent
func (h *hncli) Scroll(amount int) {
	newPos := h.pos + amount

	if newPos < 0 {
		h.pos = 0
	} else {
		h.pos = newPos
	}

	h.SetContent(h.paginate())
}

func (h *hncli) ResetScroll() {
	h.pos = 0
}

type CharHandler func(string)

func (h *hncli) SetKeyHandler(hand CharHandler) {
	h.handler = hand
}

func (h *hncli) Quit() {
	goncurses.End()
	h.finished = true
}

func (h *hncli) Run() {
	for !h.finished {
		c := h.root.GetChar()

		if c == 127 {
			c = goncurses.Key(goncurses.KEY_BACKSPACE)
		}

		ch := goncurses.KeyString(c)

		h.handler(ch)
	}
}

//Show an alert, wait for a character, then reset
func (h *hncli) Alert(text string) {
	hText := h.helpText
	h.SetHelp(text)
	h.root.GetChar()
	h.SetHelp(hText)
}

func (h *hncli) DelChar() {
	cy, cx := h.root.CursorYX()
	h.root.MoveDelChar(cy, cx-3)
	h.root.DelChar()
	h.root.DelChar()
}

func (h *hncli) getFitLines() []string {
	_, w := h.main.MaxYX()

	a := strings.Split(h.content, "\n")

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

func (hc *hncli) paginate() string {
	h, _ := hc.main.MaxYX()

	lines := hc.getFitLines()

	nLines := len(lines)

	//If text won't fit
	if nLines > h {
		//Determine start point
		if hc.pos < nLines-h {
			//If n is not in last h elements, reslice starting at n
			lines = lines[hc.pos : hc.pos+h]
		} else {
			//Else, Get last h (at most) elements
			lines = lines[nLines-h:]
		}
	}

	return strings.Join(lines, "\n")
}
