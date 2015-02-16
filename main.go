package main

import (
	"flag"
	"os"

	"github.com/andrewstuart/hn/hncli"
)

func main() {
	// s := flag.Bool("s", false, "Serves a webpage with rendings of hackernews articles")
	// p := flag.String("p", "8000", "Sets the port for the server")
	flag.Parse()

	cli := hncli.NewCli()

	for _, flag := range flag.Args() {
		if flag == "-" {
			cli.UseWriter(os.Stdout)
		}
	}

	cli.Run()
}
