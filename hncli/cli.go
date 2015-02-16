package hncli

import (
	"io"
	"log"

	"code.google.com/p/goncurses"
)

var input string = ""

func NewCli() *hncli {
	return &hncli{}
}

//UseWriter allows users to override ncurses interface and use stdout
func (hc *hncli) UseWriter(w io.Writer) {
	hc.writer = w
}

//Set up ncurses if necessary
func (hc *hncli) Init() {
	if hc.writer != nil {
		return
	}

	root, e := goncurses.Init()

	if e != nil {
		goncurses.End()
		log.Fatal(e)
	}

	h, w := root.MaxYX()

	hc.root = root
	hc.main = root.Sub(h-MenuHeight, w, 0, 0)
	hc.help = root.Sub(MenuHeight, w, h-MenuHeight, 0)

	hc.Height = h - MenuHeight
}
