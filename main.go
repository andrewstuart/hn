package main

import (
  "net/http"
  "encoding/json"
  "log"
  "code.google.com/p/goncurses"
)

type Article struct {
  Title string `json:"title"`
  Url string `json:"url"`
  Id int `json:"id"`
  Comments int `json:"commentCount"`
  Points int `json:"points"`
  Poster string `json:"postedBy"`
}

type apiDoc struct {
  NextId int `json:"nextId"`
  Items []Article `json:"items"`
}

var scr goncurses.Screen

func printArticles (a []Article) (m map[int]Article) {
  for i, art := range a {
    m[i] = art
    scr.Printf("%d. (%d) %s\n", i, art.Points, art.Title)
  }

  scr.Print("\n\nEnter a number to display the article, or n to get next.\n\n")
  scr.Refresh()

  return m
}

func main() {
  var err error
  scr, err = goncurses.Init()
  if err != nil {
    log.Fatal(err)
  }

  resp , err2 := http.Get("http://api.ihackernews.com/page")

  if err2 != nil {
    log.Fatal(err2)
  }

  if resp.StatusCode != http.StatusOK {
    log.Fatal(resp.Status)
  }

  r := new(apiDoc)

  err2 = json.NewDecoder(resp.Body).Decode(r)

  if err2 != nil {
    log.Fatal(err2)
  }

  arts := printArticles(r.Items)

  exit := false
  for st := ""; !exit {
    k := scr.GetChar()

    s := KeyString(k)

    if s != "enter" && s != "n" {
      st += KeyString(k)
    } else if s == "enter" {
      end := fmt.Sprintf("http://api.ihackernews.com/post/%d", arts[int(st)].Id)
      resp, err2 = http.Get(end)

      if err2 != nil {
        log.Fatal(err2)
      }

    } else if s == "^[" {
      exit = true
    }
  }

  defer goncurses.End()
}
