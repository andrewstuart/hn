package main

import (
  "encoding/json"
  "net/http"
)

var p Page

type hns struct {
  p Page
}

func (h hns) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  p.GetNext()
  enc := json.NewEncoder(w)
  enc.Encode(p)
}

func server () (*http.Server) {
  p = Page{
    NextUrl: "news",
  }

  p.Init()

  h := hns{p}

  s := &http.Server{
    Addr: ":80",
    Handler: h,
  }

  return s
}
