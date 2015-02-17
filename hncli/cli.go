package hncli

import (
	"log"

	"github.com/andrewstuart/hn/hackernews"
	"github.com/andrewstuart/hn/hncli/shell"
)

var pc *hackernews.PageCache
var cli *shell.Cli

func Run() {
	pc = hackernews.NewPageCache()
	cli, err := shell.NewCli(nil)

	if err != nil {
		cli.Quit()
		log.Fatal(err)
	}
}
