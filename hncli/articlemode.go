package hncli

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/andrewstuart/hn/hackernews"
)

func (cli *hncli) showStories() {
	cli.SetContent(cli.getStoriesString(cli.page))
	cli.handler = storyHandler
	cli.SetHelp("(n: next, p: previous, <num>c: view comments, <num>o: open in browser, q: quit)  ")
}

func getArticle(a *hackernews.Article) string {
	return fmt.Sprintf("(%d)\t%s - %s", a.Karma, a.User, a.Title)
}

func (cli *hncli) getStoriesString(pageNum int) string {
	h := cli.Height

	start := h * pageNum
	end := start + h

	for end > len(cli.pc.Articles) {
		cli.pc.GetNext()
	}

	str := ""
	for i, ar := range cli.pc.Articles[start:end] {
		str += fmt.Sprintf("%4d.\t(%d)\t%s\n", start+i+1, ar.Karma, ar.Title)
	}

	return str
}

func storyHandler(ch string, cli *hncli) {
	if cli == nil {
		return
	}

	switch ch {
	case "c":
		if num, err := strconv.Atoi(input); err == nil {
			if num < 1 {
				break
			}

			for num-1 > len(cli.pc.Articles) {
				cli.pc.GetNext()
			}

			//TODO comment-parsing functionality
			text := getComments(cli.pc.Articles[num-1])

			cli.showComments(text)
			input = ""
		} else {
			cli.Alert("Please enter a number to select a comment")
		}
		input = ""
		break
	case "o":
		if num, err := strconv.Atoi(input); err == nil {
			for num-1 > len(cli.pc.Articles) {
				cli.pc.GetNext()
			}

			viewInBrowser := exec.Command("xdg-open", cli.pc.Articles[num-1].Url)
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
		cli.page += 1
		cli.SetContent(cli.getStoriesString(cli.page))
		input = ""
		break
	case "p":
		//Go back 1 page, unless page < 0
		if cli.page > 0 {
			cli.page -= 1
		}
		cli.SetContent(cli.getStoriesString(cli.page))
		break
	case "enter":
		cli.Refresh()
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
