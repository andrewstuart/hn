package main

import (
	"fmt"

	"github.com/andrewstuart/hn/hackernews"
)

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
