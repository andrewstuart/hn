package main

import (
  "net/http"
  "encoding/json"
  "fmt"
  "log"
  "strconv"
  "code.google.com/p/goncurses"
)

type Article struct {
  Title string `json:"title"`
  Url string `json:"url"`
  Id int `json:"id"`
  CommentCount int `json:"commentCount"`
  Points int `json:"points"`
  Poster string `json:"postedBy"`
}

func handleFatal(err error) {
  if(err != nil) {
    log.Fatal(err)
  }
}

type CommentDoc struct {
  Comments []Comment `json:"items"`
}

type Comment struct {
  Username string `json:"username"`
  Text string `json:"comment"`
  Id string `json:"id"`
  Children []Comment `json:"children"`
}

type apiDoc struct {
  NextId int `json:"nextId"`
  Articles []Article `json:"items"`
}

var scr *goncurses.Window

var artMap map[int]*Article

func getComments (id int) (c []Comment) {
  end := fmt.Sprintf("http://hndroidapi.appspot.com/nestedcomments/format/json/id/%d", artMap[id].Id)
  resp, err2 := http.Get(end)

  handleFatal(err2)

  cd := new(CommentDoc)
  handleFatal(json.NewDecoder(resp.Body).Decode(cd))

  return cd.Comments
}

func printComments (c []Comment) {
  handleFatal(scr.Clear())
  for i, comm := range c {
    scr.Printf("%d. (%d) %s: %s\n\n", i + 1, len(comm.Children), comm.Username, comm.Text)
  }

  scr.Print("\n\nPress any key to return\n\n")
  scr.GetChar()
}

func getArticles (start int) (a []Article, next int) {
  url := "http://api.ihackernews.com/page"

  if(next > 0) {
    url = fmt.Sprintf("%s/%d", url, start)
  }

  resp , err2 := http.Get(url)

  handleFatal(err2)

  if resp.StatusCode != http.StatusOK {
    log.Fatal(resp.Status)
  }

  r := new(apiDoc)
  handleFatal(json.NewDecoder(resp.Body).Decode(r))

  next = r.Articles[len(r.Articles) - 1].Id

  scr.Println(next)
  scr.Refresh()
  scr.GetChar()

  a = r.Articles
  return a, next
}

func printArticles (a []Article) {
  artMap = make(map[int]*Article)

  for i, art := range a {
    artMap[i] = &art
    scr.Printf("%d. (%d) %s\n", i, art.Points, art.Title)
  }

  scr.Print("\n\nEnter a number to display the article, or n to get next.\n\n")
  scr.Refresh()
}

func main() {
  var err error
  scr, err = goncurses.Init()
  defer goncurses.End()

  if err != nil {
    log.Fatal(err)
  }

  artList, next := getArticles(0)

  scr.Clear()

  printArticles(artList)

  exit := false
  st := ""
  for !exit {

    k := scr.GetChar()

    s := goncurses.KeyString(k)

    switch s {
    case "n":
      artList, next = getArticles(next)
      scr.Clear()
      printArticles(artList)
    case "^[":
      exit = true
    case "enter":
      id, err := strconv.Atoi(st)

      if(err != nil) {
        scr.Printf("Bad input.\n\n")
        scr.Refresh()
        st = ""
        continue
      }

      c := getComments(id)

      printComments(c)

      scr.Clear()
      printArticles(artList)

      st = ""
    default:
      st += s
    }
  }
}
