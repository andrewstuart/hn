package main

import (
	"flag"
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
