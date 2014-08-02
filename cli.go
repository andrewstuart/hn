package main

import (
	"fmt"
	"os/exec"
	"strconv"
)

var input string = ""

func storyHandler(ch string) {
	switch ch {
	case "c":
		if num, err := strconv.Atoi(input); err == nil {
			if num < 1 {
				break
			}

			for num-1 > len(p.Articles) {
				p.GetNext()
			}

			text := p.Articles[num-1].PrintComments()

			cli.SetContent(text)
			cli.SetHelp("(d/u scroll 30 lines; j/k: scroll 1 line; n/p scroll 1 page; q: quit to story view)")
			cli.SetKeyHandler(commentHandler)
			input = ""
		} else {
			cli.Alert("Please enter a number to select a comment")
		}
		break
	case "o":
		if num, err := strconv.Atoi(input); err == nil {
			for num-1 > len(p.Articles) {
				p.GetNext()
			}

			viewInBrowser := exec.Command("xdg-open", p.Articles[num-1].Url)
			viewInBrowser.Start()
		} else {
			cli.Alert("Please enter a number to view an article")
		}
		input = ""
		break
	case "q":
		cli.Quit()
		break
	case "n":
		//Go forward 1 page
		pageNum += 1
		cli.SetContent(getStories(pageNum))
		break
	case "p":
		//Go back 1 page, unless page < 0
		if pageNum > 0 {
			pageNum -= 1
		}
		cli.SetContent(getStories(pageNum))
		break
	case "enter":
		break
	case "backspace":
		if len(input) > 0 {
			input = input[:len(input)-1]
			cli.DelChar()
		} else {
			cli.DelChar()
		}
		break
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
		cli.Scroll(-1)
		break
	case "g":
		cli.ResetScroll()
		break
	case "n":
		cli.Scroll(cli.Height)
		break
	case "p":
		cli.Scroll(-cli.Height)
		break
	case "q":
		cli.ResetScroll()
		cli.SetContent(stories)
		cli.SetKeyHandler(storyHandler)
		break
	}
}

var cli hncli
var p *PageCache

var stories string

func getStories(pageNum int) string {
	h := cli.Height

	start = h * pageNum
	end = start + h

	for end > len(p.Articles) {
		p.GetNext()
	}

	str := ""
	for i, ar := range p.Articles[start:end] {
		str += fmt.Sprintf("%4d.\t(%d)\t%s\n", start+i+1, ar.Karma, ar.Title)
	}

	return str
}

var pageNum, start, end int

func runCli() {
	pageNum = 0
	cli = GetCli()

	p = NewPageCache()

	stories = getStories(pageNum)

	cli.SetContent(stories)
	cli.SetKeyHandler(storyHandler)
	cli.SetHelp("\n(n: next, p: previous, <num>c: view comments, <num>o: open in browser, q: quit)  ")

	cli.Refresh()
	cli.Run()
}
