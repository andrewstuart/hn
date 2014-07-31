package main

import (
	"fmt"
	"os/exec"
	"strconv"
)

func storyHandler(ch string) {
	input := ""
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
				a := scr.GetChar()

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

func commentHandler(input string) {
	switch input {
	case "d":
		cli.Scroll(30)
		break
	case "u":
		cli.Scroll(-30)
		break
	case "j":
		cli.Scroll(1)
		break
	case "k":
		cli.Scroll(1)
		break
	case "n":
		cli.Scroll(cli.Height)
		break
	case "p":
		cli.Scroll(cli.Height)
		break
	case "q":
		cli.SetKeyHandler(storyHandler)
		cli.SetContent(stories)
		break
	}
}

var cli hncli

func runCli() {
	cli = GetCli()

	exit := false
	pageNum := 0
	p := NewPageCache()

	for !exit {
		h := cli.Height

		start := h * pageNum
		end := start + h

		for end > len(p.Articles) {
			p.GetNext()
		}

		content := ""

		for i, ar := range p.Articles[start:end] {
			content += fmt.Sprintf("%d. (%d): %s\n", start+i+1, ar.Karma, ar.Title)
		}

		cli.PrintHelp("\n(n: next, p: previous, <num>c: view comments, <num>o: open in browser, q: quit)  ")

		cli.Refresh()

	}
}
