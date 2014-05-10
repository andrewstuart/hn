package main

import (
  "log"
  "github.com/PuerkitoBio/goquery"
  "code.google.com/p/goncurses"
  "net/http"
  "crypto/tls"
)
var scr *goncurses.Window
var doc *goquery.Document
var resp *http.Response
var e error


func main() {
  if scr, e = goncurses.Init(); e != nil {
    log.Fatal(e)
  }
  defer goncurses.End()

  trans := &http.Transport{
    TLSClientConfig : &tls.Config{InsecureSkipVerify: true},
  }

  client := &http.Client{Transport: trans}

  if resp, e = client.Get("https://news.ycombinator.com/news"); e != nil {
    log.Print(e)
  }

  if doc, e = goquery.NewDocumentFromResponse(resp); e != nil {
    log.Fatal(e)
  }

  doc.Find("tbody td.title a:first-child").Each(func(i int, s *goquery.Selection) {
    if t := s.Text(); t != "More" {
      scr.Printf("%d. %s\n", i + 1, t)
    }
  })
  scr.Refresh()
  scr.GetChar()
}
