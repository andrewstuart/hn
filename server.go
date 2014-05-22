package main

import (
  "encoding/json"
  "net/http"
  "log"
  "os/exec"
  "strconv"
)

const port string = "8000"

const commentRoute string = "/comments/"

func getComments (w http.ResponseWriter, r *http.Request) {
  w.Header()["Access-Control-Allow-Origin"] = []string{"*"}

  idSt := r.URL.Path[len(commentRoute):]

  if id, err := strconv.Atoi(idSt); err == nil {
    ar := p.GetComments(id)
    enc := json.NewEncoder(w)
    enc.Encode(ar)
  } else {
    log.Print(err)
  }
}

func next (w http.ResponseWriter, r *http.Request) {
  w.Header()["Access-Control-Allow-Origin"] = []string{"*"}

  p.GetNext()

  enc := json.NewEncoder(w)

  enc.Encode(p)
}

var p Page

func send(w http.ResponseWriter, r *http.Request) {
  w.Header()["Access-Control-Allow-Origin"] = []string{"*"}

  enc := json.NewEncoder(w)
  enc.Encode(p)
}

func server () {
  // log.Fatal("huh?")
  p = Page{
    NextUrl: "news",
  }

  p.Init()
  p.GetNext()

  view := exec.Command("xdg-open", "http://localhost:" + port)

  view.Start()

  http.HandleFunc("/next/", next)
  http.HandleFunc("/", send)
  http.HandleFunc(commentRoute, getComments)

  err := http.ListenAndServe(":" + port, nil)

  if err != nil {
    log.Fatal(err)
  }
}

