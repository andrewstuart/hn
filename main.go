package main

import (
	"flag"
)

func main() {
	s := flag.Bool("s", false, "Serves a webpage with rendings of hackernews articles")
	p := flag.String("p", "8000", "Sets the port for the server")
	flag.Parse()

	if *s {
		port := ":" + *p
		server(port)
	} else {
		runCli()
	}
}
