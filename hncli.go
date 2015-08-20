package main

import (
	"log"
	"strings"
	"github.com/rthornton128/goncurses"
	//"github.com/tncardoso/gocurses"
	//"code.google.com/p/goncurses"
)

type hncli struct {
	root, main, help  *goncurses.Window
	Height, offset    int
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
	h.help.Move(1, 0)
	h.help.Print(text)
	h.help.Refresh()
}

//Scroll the content that was added with SetContent
func (h *hncli) Scroll(amount int) {
	newOffset := h.offset + amount

	if newOffset < 0 {
		h.offset = 0
	} else {
		h.offset = newOffset
	}

	h.drawPage()
}

func (h *hncli) ResetScroll() {
	h.offset = 0
	h.drawPage()
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
	h.SetHelp(text + "   (Press any key to continue)")
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

	lineArray := strings.Split(h.content, "\n")

	p := make([]string, 0, len(lineArray))

	for _, line := range lineArray {

		//Remember padding for each line
		var pad string
		for len(line) > len(COMMENT_PAD)+len(pad) && line[len(pad):len(pad)+len(COMMENT_PAD)] == COMMENT_PAD {
			pad += COMMENT_PAD
		}

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

			//Split lines
			p = append(p, line[:l])
			line = line[l:]

			//Strip leading spaces
			for len(line) > 0 && line[0] == ' ' {
				line = line[1:]
			}

			//Pad line if need be
			if len(line) > 0 {
				line = pad + line
			}
		}

		p = append(p, line)
	}

	return p
}

func (h *hncli) drawPage() {
	lines := h.getFitLines()

	nLines := len(lines)

	//If text won't fit
	if nLines > h.Height {
		//Determine start point
		if h.offset < nLines-h.Height {
			//If n is not in last h elements, reslice starting at n
			lines = lines[h.offset : h.offset+h.Height]
		} else {
			//Else, Get last h (at most) elements
			lines = lines[nLines-h.Height:]
		}
	}

	h.main.Clear()
	h.main.Print(strings.Join(lines, "\n"))
	h.main.Refresh()
}
