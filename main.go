package main

import (
  "log"
  "github.com/PuerkitoBio/goquery"
  "fmt"
  "code.google.com/p/goncurses"
  "net/http"
  "crypto/tls"
  "strings"
  "strconv"
)

const YCRoot = "https://news.ycombinator.com"
const rowsPerArticle = 3

var scr *goncurses.Window
var doc *goquery.Document
var resp *http.Response
var e error

//Comments structure for HN articles
type Comment struct {
  Text string `json:"text"`
  User string `json:"user"`
  Id int `json:"id"`
  Children []*Comment `json:"children"`
}

//Article structure
type Article struct {
  Title string `json:"title"xml:"`
  Points int `json:"points"`
  Id int `json:"id"`
  Url string `json:"url"`
  SiteLabel string `json:"siteLabel"`
  Comments []*Comment `json:"comments"`
  User string `json:"user"`
  Created string `json:"created"`
}

func main() {
  if scr, e = goncurses.Init(); e != nil {
    log.Fatal(e)
  }
  defer goncurses.End()

  trans := &http.Transport{
    TLSClientConfig : &tls.Config{InsecureSkipVerify: true},
  }

  client := &http.Client{Transport: trans}

  next := YCRoot + "/news"
  exit := false

  ars := make([]*Article, 0)
  page := 0


  for !exit {

    if resp, e = client.Get(next); e != nil {
      log.Print(e)
    }

    if doc, e = goquery.NewDocumentFromResponse(resp); e != nil {
      log.Fatal(e)
    }


    rows := doc.Find("tr").Eq(3).Find("tr")

    nextHref, _ := rows.Last().Find("a").Attr("href")

    if nextHref[0] == '/' {
      next = YCRoot + nextHref
    } else {
      next = YCRoot + "/" + nextHref
    }

    rows.Each(func(i int, row *goquery.Selection) {
      if i % rowsPerArticle > 0 || (len(ars) % 30 == 0 && i != 0){ return }

      ar := Article{}
      up := row.Find("td.title").First().Next()

      if idSt, exists := up.Find("a").First().Attr("id"); exists {
        if id, err := strconv.Atoi(strings.Split(idSt, "_")[1]); err == nil {
          ar.Id = id
        } else {
          fmt.Println(idSt)
          log.Fatal(err)
        }
      }

      title := up.Next()
      link := title.Find("a").First()

      ar.Title = link.Text()
      if url, exists := link.Attr("href"); exists {
        ar.Url = url
      }

      ar.SiteLabel = title.Find("span.comhead").Text()

      row = row.Next()

      row.Find("span").Each(func (i int, s *goquery.Selection) {
        if pts, err := strconv.Atoi(strings.Split(s.Text(), " ")[0]); err == nil {
          ar.Points = pts
        } else {
          fmt.Println(ar)
          fmt.Println(row.Html())
          log.Fatal(err)
        }
      })

      ar.User = row.Find("td.subtext a:first-child").Text()

      ars = append(ars, &ar)

      //scr.Print(ar)
      //scr.Refresh()
      //scr.GetChar()
    })

    scr.Clear()

    start := 30 * page
    end := start + 30

    for i, ar := range ars[start:end] {
      scr.Printf("%d. (%d): %s\n", start + i + 1, ar.Points, ar.Title)
    }

    scr.Print("\n\nPress n to continue or q to quit\n\n")
    scr.Refresh()

    switch goncurses.KeyString(scr.GetChar()) {
    case "q":
      exit = true
    default:
      page += 1
    }
  }
}
