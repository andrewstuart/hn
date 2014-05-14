package main

import (
  "log"
  "github.com/PuerkitoBio/goquery"
  "os/exec"
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

var trans *http.Transport = &http.Transport{
  TLSClientConfig : &tls.Config{InsecureSkipVerify: true},
}

var client *http.Client = &http.Client{Transport: trans}

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

func (a *Article) GetComments() (comments []*Comment) {
  comments = make([]*Comment, 0)

  articleUrl := YCRoot + "/item?id=" + strconv.Itoa(a.Id)

  resp, e := client.Get(articleUrl)

  if e != nil {
    log.Fatal(e)
  }

  if doc, e := goquery.NewDocumentFromResponse(resp); e != nil {
    log.Fatal(e)
  } else {

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
  }

  a.Comments = comments
  return comments;
}

func (a *Article) GoString() string {
  return fmt.Sprintf("(%d) %s: %s\n\n", a.Points, a.User, a.Title)
}

func (a *Article) PrintComments() {
  a.GetComments()

  scr.Print(a)

  for i, comment := range a.Comments {
    scr.Printf("%d. %s: %s\n", i, comment.User, comment.Text)
  }
}

type Page struct {
  NextUrl string `json:"next"`
  Articles []*Article `json:"articles"`
}

func (p *Page) GetNext() {
  url := YCRoot

  if p.NextUrl[0] != '/' {
    url += "/"
  }

  url += p.NextUrl

  if resp, e := client.Get(url); e != nil {
    log.Fatal(e)
  } else {

    if doc, e := goquery.NewDocumentFromResponse(resp); e != nil {
      log.Fatal(e)
    } else {

      //Get all the trs with subtext for children then go back one (for the first row)
      rows := doc.Find(".subtext").ParentsFilteredUntil("tr", "tbody").Prev()

      p.NextUrl, _ = doc.Find("td.title").Last().Find("a").Attr("href")

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

        p.Articles = append(p.Articles, &ar)
      })
    }
  }
}

func main() {
  scr, e := goncurses.Init()
  if e != nil {
    log.Fatal(e)
  }

  defer goncurses.End()

  exit := false

  pageNum := 0

  p := Page{
    NextUrl: "news",
  }

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
