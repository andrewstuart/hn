package hncli

import (
	"fmt"

	"github.com/andrewstuart/hn/hackernews"
)

func (cli *hncli) showComments(text string) {
	cli.SetContent(text)
	cli.ScrollTop()
	cli.SetHelp("(d/u scroll 30 lines; j/k: scroll 1 line; n/p scroll 1 page; q: quit to story view)")
	cli.handler = commentHandler
}

func commentHandler(input string, cli *hncli) {
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
		cli.ScrollTop()
		break
	case "n":
		cli.Scroll(cli.Height)
		break
	case "p":
		cli.Scroll(-cli.Height)
		break
	case "q":
		cli.showStories()
		break
	}
}

func getComments(a *hackernews.Article) string {
	return fmt.Sprintf("")
}

//Recursively get comments
func commentString(cs []*hackernews.Comment, off string) string {
	s := ""

	for _, c := range cs {
		s += off + fmt.Sprintf(off+"%s\n\n", c)

		if len(c.Comments) > 0 {
			s += commentString(c.Comments, off+"  ")
		}
	}

	return s
}
