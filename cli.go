package main

import (
	"log"
	"os/exec"
	"strconv"
	"strings"

	"code.google.com/p/goncurses"
)

const BOTTOM_MARGIN int = 3

func getFitLines(s string) []string {
	_, w := scr.MaxYX()

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

func pagify(t string, n int) string {
	h, _ := scr.MaxYX()
	h -= BOTTOM_MARGIN

	lines := getFitLines(t)

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

func cli() {
	var e error
	scr, e = goncurses.Init()

	if e != nil {
		log.Fatal(e)
	}

	defer goncurses.End()

	exit := false

	pageNum := 0

	p := NewPageCache()

	for !exit {
		scr.Refresh()
		h, _ := scr.MaxYX()

		scr.Clear()

		height := h - BOTTOM_MARGIN

		start := height * pageNum
		end := start + height

		for end > len(p.Articles) {
			p.GetNext()
		}

		for i, ar := range p.Articles[start:end] {
			scr.Printf("%4d.\t(%d)\t%s\n", start+i+1, ar.Karma, ar.Title)
		}

		scr.Print("\n(n: next, p: previous, <num>c: view comments, <num>o: open in browser, q: quit)  ")
		scr.Refresh()

		doneWithInput := false
		input := ""
		for !doneWithInput {
			c := scr.GetChar()

			if c == 127 {
				c = goncurses.Key(goncurses.KEY_BACKSPACE)
			}

			ch := goncurses.KeyString(c)
			switch ch {
			case "c":
				if num, err := strconv.Atoi(input); err == nil {
					if num < 1 {
						doneWithInput = true
						break
					}

					for num-1 > len(p.Articles) {
						p.GetNext()
					}

					text := p.Articles[num-1].PrintComments()
					line := 0

					cont := true
					for cont {
						scr.Clear()
						scr.Print(pagify(text, line))
						scr.Print("\n\n(d/u scroll 30 lines; j/k: scroll 1 line; n/p scroll 1 page; q: quit)")
						scr.Refresh()

						a := scr.GetChar()

						switch goncurses.KeyString(a) {
						case "d":
							line += 30
							break
						case "u":
							line -= 30
							break
						case "j":
							line += 1
							break
						case "k":
							line -= 1
							break
						case "n":
							line += h
							break
						case "p":
							line -= h
							break
						case "q":
							cont = false
							break
						default:
							scr.DelChar()
							scr.DelChar()
							break
						}

						// Verify lines are not negative. Bad mojo
						if line < 0 {
							line = 0
						}
					}

					doneWithInput = true
				} else {
					scr.Clear()
					scr.Print("\n\nPlease enter a number to select a comment\n\n")
					scr.Refresh()
					scr.GetChar()
					doneWithInput = true
				}
			case "o":
				if num, err := strconv.Atoi(input); err == nil {
					for num-1 > len(p.Articles) {
						p.GetNext()
					}

					viewInBrowser := exec.Command("xdg-open", p.Articles[num-1].Url)
					viewInBrowser.Start()
					doneWithInput = true
				} else {
					scr.Clear()
					scr.Print("\n\nPlease enter a number to view an article\n\n")
					scr.Refresh()
					doneWithInput = true
				}
			case "q":
				doneWithInput = true
				exit = true
			case "n":
				pageNum += 1
				doneWithInput = true
			case "p":
				if pageNum > 0 {
					pageNum -= 1
				}
				doneWithInput = true
			case "enter":
				continue
			case "backspace":
				//Not the prettiest but whatever
				cy, cx := scr.CursorYX()
				if len(input) > 0 {
					input = input[:len(input)-1]
					scr.MoveDelChar(cy, cx-3)
					scr.DelChar()
					scr.DelChar()
				} else {
					scr.MoveDelChar(cy, cx-2)
					scr.DelChar()
					scr.DelChar()
				}
			default:
				input += ch
				break
			}
		}
	}
}

var scr *goncurses.Window
