package main

import (
  "log"
  "github.com/PuerkitoBio/goquery"
  "os/exec"
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

var trans *http.Transport = &http.Transport{
  TLSClientConfig : &tls.Config{InsecureSkipVerify: true},
}

var client *http.Client = &http.Client{Transport: trans}

func (a *Article) GetComments() (comments []*Comment) {
  comments = make([]*Comment, 0)

  articleUrl := YCRoot + "/item?id=" + strconv.Itoa(a.Id)

  resp, e := client.Get(articleUrl)

  if e != nil {
    log.Fatal(e)
  }

  if doc, e = goquery.NewDocumentFromResponse(resp); e != nil {
    log.Fatal(e)
  }

  doc.Find(".comment").Each(func (i int, comment *goquery.Selection) {
    text := ""
    user := comment.Parent().Find("a").First().Text()

    comment.Find("font").Each(func (j int, paragraph *goquery.Selection) {
      text += paragraph.Text() + "\n"
    })

    comments = append(comments, &Comment{
      User: user,
      Text: text,
    })
  })

  a.Comments = comments
  return comments;
}

func (a *Article) PrintHead() {
  scr.Printf("(%d) %s: %s\n\n", a.Points, a.User, a.Title)
}

func (a *Article) PrintComments() {
  a.GetComments()

  a.PrintHead()

  for i, comment := range a.Comments {
    scr.Printf("%d. %s: %s\n", i, comment.User, comment.Text)
  }
}

var next string = YCRoot + "/news"
var ars []*Article = make([]*Article, 0, 30)

func getArticles() {
  if resp, e = client.Get(next); e != nil {
    log.Print(e)
  }

  if doc, e = goquery.NewDocumentFromResponse(resp); e != nil {
    log.Fatal(e)
  }

  rows := doc.Find(".subtext").ParentsFilteredUntil("tr", "tbody").Prev()

  nextHref, _ := doc.Find("td.title").Last().Find("a").Attr("href")

  if nextHref[0] == '/' {
    next = YCRoot + nextHref
  } else {
    next = YCRoot + "/" + nextHref
  }

  rows.Each(func(i int, row *goquery.Selection) {
    ar := Article{}

    title := row.Find(".title").Eq(1)
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
        log.Fatal(err)
      }

      if idSt, exists := s.Attr("id"); exists {
        if id, err := strconv.Atoi(strings.Split(idSt, "_")[1]); err == nil {
          ar.Id = id
        } else {
          log.Fatal(err)
        }
      }
    })

    ar.User = row.Find("td.subtext a").First().Text()

    ars = append(ars, &ar)
  })
}

func main() {
  if scr, e = goncurses.Init(); e != nil {
    log.Fatal(e)
  }
  defer goncurses.End()

  exit := false

  page := 0

  for !exit {

    scr.Clear()

    start := 30 * page
    end := start + 30

    for end > len(ars) {
      getArticles()
    }

    for i, ar := range ars[start:end] {
      scr.Printf("%d. (%d): %s\n", start + i + 1, ar.Points, ar.Title)
    }

    scr.Print("\n\nPress n to go forward or p to go back\n\nEnter a number and press c to view comments or v to view the article\n\nPress q to quit\n\n")
    scr.Refresh()

    doneWithInput := false
    input := ""
    for !doneWithInput {
      c := scr.GetChar()
      chr := goncurses.KeyString(c)
      switch chr {
      case "c":
        if num, err := strconv.Atoi(input); err == nil {
          scr.Clear()
          ars[num - 1].PrintComments()
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
      case "v":
        if num, err := strconv.Atoi(input); err == nil {
          viewInBrowser := exec.Command("xdg-open", ars[num - 1].Url)
          viewInBrowser.Run()
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
        page += 1
        doneWithInput = true
      case "p":
        if page > 0 {
          page -= 1
        }
        doneWithInput = true
      default:
        input += chr
      }
    }
  }
}
