package main

import (
	"flag"
	"log"
	"os/exec"
	"strconv"

	"code.google.com/p/goncurses"
)

func main() {
	s := flag.Bool("s", false, "Serves a webpage with rendings of hackernews articles")
	flag.Parse()

	if *s {
		server()
	} else {
		cli()
	}
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

		height := h - 3

		start := height * pageNum
		end := start + height

		for end > len(p.Articles) {
			p.GetNext()
		}

		for i, ar := range p.Articles[start:end] {
			scr.Printf("%d. (%d): %s\n", start+i+1, ar.Points, ar.Title)
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
					for num-1 > len(p.Articles) {
						p.GetNext()
					}

					scr.Clear()
					p.Articles[num-1].PrintComments()
					scr.Refresh()
					scr.GetChar()
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
			}
		}
	}
}

var scr *goncurses.Window
