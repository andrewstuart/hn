package main

import (
  "log"
  "fmt"
  "github.com/PuerkitoBio/goquery"
  "code.google.com/p/goncurses"
  "net/http"
  "crypto/tls"
  "strings"
  "strconv"
)

const YCRoot = "https://news.ycombinator.com"
const rowsPerArticle = 3

var trans *http.Transport = &http.Transport{
  TLSClientConfig : &tls.Config{InsecureSkipVerify: true},
}

var client *http.Client = &http.Client{Transport: trans}


//Comments structure for HN articles
type Comment struct {
  Text string `json:"text"`
  User string `json:"user"`
  Id int `json:"id"`
  Comments []*Comment `json:"comments"`
}

func (c *Comment) String() string {
  return fmt.Sprintf("%s: %s\n", c.User, c.Text)
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

func (a *Article) GetComments() {
  a.Comments = make([]*Comment, 0)

  articleUrl := YCRoot + "/item?id=" + strconv.Itoa(a.Id)

  resp, e := client.Get(articleUrl)

  if e != nil {
    log.Fatal(e)
  }

  if doc, e := goquery.NewDocumentFromResponse(resp); e != nil {
    log.Fatal(e)
  } else {

    commentStack := make([]*Comment, 1, 10)

    doc.Find("span.comment").Each(func (i int, comment *goquery.Selection) {
      text := ""
      user := comment.Parent().Find("a").First().Text()

      text += comment.Text()

      //Get around HN's little weird "reply" nesting randomness
      //Is it part of the comment, or isn't it?
      if last5 := len(text) - 5; last5 > 0 && text[last5:] == "reply" {
        text = text[:last5]
      }

      c := &Comment{
        User: user,
        Text: text,
        Comments : make([]*Comment, 0),
      }

      //Get id
      if idAttr, exists := comment.Prev().Find("a").Last().Attr("href"); exists {
        idSt := strings.Split(idAttr, "=")[1]

        if id, err := strconv.Atoi(idSt); err == nil {
          c.Id = id
        }
      }


      //Track the comment offset for nesting.
      //TODO: Better selectors
      if width, exists := comment.Parent().Prev().Prev().Find("img").Attr("width"); exists {
        offset, _ := strconv.Atoi(width)
        offset = offset / 40

        lastEle := len(commentStack) - 1 //Index of the last element in the stack

        if offset > lastEle {
          commentStack = append(commentStack, c)
          commentStack[lastEle].Comments = append(commentStack[lastEle].Comments, c)
        } else {

          if offset < lastEle {
            commentStack = commentStack[:offset + 1] //Trim the stack
          }

          commentStack[offset] = c

          //Add the comment to its parents
          if offset == 0 {
            a.Comments = append(a.Comments, c)
          } else {
            commentStack[offset - 1].Comments = append(commentStack[offset - 1].Comments, c)
          }
        }
      }
    })
  }
}

func (a *Article) String() string {
  return fmt.Sprintf("(%d) %s: %s\n\n", a.Points, a.User, a.Title)
}

func commentString (cs []*Comment, off string) string {
  s := ""

  for i, c := range cs {
    s += off + fmt.Sprintf("%d. %s\n", i + 1, c)

    if len(c.Comments) > 0 {
      s += commentString(c.Comments, off + strconv.Itoa(i + 1) + ".")
    }
  }

  return s
}

//TODO Use stringer interface
func (a *Article) PrintComments() {
  a.GetComments()

  scr.Print(a)

  cs := commentString(a.Comments, "")

  scr.Print(cs)
}

type Page struct {
  NextUrl string `json:"next"`
  Articles []*Article `json:"articles"`
  CFDUid string
}

func (p *Page) GetCFDUid() {
  url := YCRoot + "/news"

  if resp, err := client.Head(url); err == nil {
    c := resp.Cookies()
    p.CFDUid = c[0].Raw
  } else {
    goncurses.End()
    log.Fatal(err)
  }
}

func (p *Page) GetNext() {
  url := YCRoot

  if p.NextUrl[0] != '/' {
    url += "/"
  }

  url += p.NextUrl

  req, err := http.NewRequest("GET", url, nil)

  if err != nil {
    goncurses.End()
    log.Fatal(err)
  }

  req.Header.Set("cookie", p.CFDUid)

  if resp, e := client.Do(req); e != nil {
    log.Fatal(e)
  } else {

    if doc, e := goquery.NewDocumentFromResponse(resp); e != nil {
      log.Fatal(e)
    } else {

      //Get all the trs with subtext for children then go back one (for the first row)
      rows := doc.Find(".subtext").ParentsFilteredUntil("tr", "tbody").Prev()

      var a bool

      p.NextUrl, a = doc.Find("td.title").Last().Find("a").Attr("href")

      if !a {
        goncurses.End()
        log.Fatal("Could not retreive next hackernews page. Time to go outside?")
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

        p.Articles = append(p.Articles, &ar)
      })
    }
  }
}

