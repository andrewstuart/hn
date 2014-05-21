package main

import (
  "encoding/json"
  "net/http"
  "log"
  "os/exec"
)

const port string = "8000"

var p Page

type hns struct {
  p Page
}

func (h hns) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  headers := w.Header()

  headers["Access-Control-Allow-Origin"] = []string{"*"}

  p.GetNext()
  enc := json.NewEncoder(w)
  enc.Encode(p)
}

func server () {
  // log.Fatal("huh?")
  p = Page{
    NextUrl: "news",
  }

  p.Init()

  h := hns{p}

  s := &http.Server{
    Addr: ":" + port,
    Handler: h,
  }

  view := exec.Command("xdg-open", "http://localhost:" + port)

  view.Start()

  err := s.ListenAndServe()

  if err != nil {
    log.Fatal(err)
  }
}

