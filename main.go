package main

import (
  "log"
  "os/exec"
  "code.google.com/p/goncurses"
  "strconv"
)

var scr *goncurses.Window

func main() {
  var e error
  scr, e = goncurses.Init()
  if e != nil {
    log.Fatal(e)
  }

  defer goncurses.End()

  exit := false

  pageNum := 0

  p := Page{
    NextUrl: "news",
  }

  p.GetCFDUid()

  for !exit {
    scr.Refresh()
    h, _ := scr.MaxYX()

    scr.Clear()

    height := h - 5

    start := height * pageNum
    end := start + height

    for end > len(p.Articles) {
      p.GetNext()
    }

    for i, ar := range p.Articles[start:end] {
      scr.Printf("%d. (%d): %s\n", start + i + 1, ar.Points, ar.Title)
    }

    scr.Print("\n(n: next, p: previous, <num>c: view comments, <num>o: open in browser, q: quit)  ")
    scr.Refresh()

    doneWithInput := false
    input := ""
    for !doneWithInput {
      c := scr.GetChar()
      chr := goncurses.KeyString(c)
      switch chr {
      case "c":
        if num, err := strconv.Atoi(input); err == nil {
          for num - 1 > len(p.Articles) {
            p.GetNext()
          }

          scr.Clear()
          p.Articles[num - 1].PrintComments()
          scr.Refresh()
          scr.GetChar()
          doneWithInput = true
        } else {
          scr.Clear()
          scr.Print("\n\nPlease enter a number to select a comment\n\n")
          scr.Refresh()
          scr.GetChar()
          doneWithInput = true
        }
      case "o":
        if num, err := strconv.Atoi(input); err == nil {
          for num - 1 > len(p.Articles) {
            p.GetNext()
          }

          viewInBrowser := exec.Command("xdg-open", p.Articles[num - 1].Url)
          viewInBrowser.Start()
          doneWithInput = true
        } else {
          scr.Clear()
          scr.Print("\n\nPlease enter a number to view an article\n\n")
          scr.Refresh()
          doneWithInput = true
        }
      case "q":
        doneWithInput = true
        exit = true
      case "n":
        pageNum += 1
        doneWithInput = true
      case "p":
        if pageNum > 0 {
          pageNum -= 1
        }
        doneWithInput = true
      default:
        input += chr
      }
    }
  }
}
