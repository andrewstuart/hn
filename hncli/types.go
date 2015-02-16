package hncli

import (
	"io"

	"code.google.com/p/goncurses"
	"github.com/andrewstuart/hn/hackernews"
)

type CharHandler func(string, *hncli)

//A simple cli shell
type hncli struct {
	root, main, help  *goncurses.Window
	Height, offset    int
	handler           CharHandler
	finished          bool
	helpText, content string
	writer            io.Writer
	page              int
	pc                *hackernews.PageCache
}
