package hncli

import (
	"os"
	"strings"

	"code.google.com/p/goncurses"
)

const MenuHeight int = 3

func (h *hncli) Refresh() {
	h.root.Refresh()
	h.root.Move(h.Height-1, 0)
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

//Scroll back to beginning
func (h *hncli) ScrollTop() {
	h.offset = 0
	h.drawPage()
}

func (h *hncli) Quit() {
	goncurses.End()
	os.Exit(0)
}

func (h *hncli) Run() {
	for {
		c := h.root.GetChar()

		if c == 127 {
			c = goncurses.Key(goncurses.KEY_BACKSPACE)
		}

		ch := goncurses.KeyString(c)
		h.handler(ch, h)
	}
}

//Show an alert, wait for a character, then reset
func (h *hncli) Alert(text string) {
	hText := h.helpText
	h.SetHelp(text + "   (Press any key to continue)")
	h.root.GetChar()
	h.SetHelp(hText)
}

//Delete a character
func (h *hncli) DelChar() {
	cy, cx := h.root.CursorYX()
	h.root.MoveDelChar(cy, cx-3)
	h.root.DelChar()
	h.root.DelChar()
}

const CommentPadding = "  "

func (h *hncli) getFitLines() []string {
	_, w := h.main.MaxYX()

	lineArray := strings.Split(h.content, "\n")

	p := make([]string, 0, len(lineArray))

	for _, line := range lineArray {

		//Remember padding for each line
		var pad string
		for len(line) > len(CommentPadding)+len(pad) && line[len(pad):len(pad)+len(CommentPadding)] == CommentPadding {
			pad += CommentPadding
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
