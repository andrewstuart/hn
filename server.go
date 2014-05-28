package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

const port string = "8000"

const commentRoute string = "/comments/"

func getComments(w http.ResponseWriter, r *http.Request) {
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

var pages = make(map[string]*Page)

func next(w http.ResponseWriter, r *http.Request) {
	reqUrl := r.URL.Path[len("/next/"):]

	w.Header()["Access-Control-Allow-Origin"] = []string{"*"}
	enc := json.NewEncoder(w)

	if pages[reqUrl] != nil {
		enc.Encode(pages[reqUrl])
	} else {
		pages[p.GetNext()] = &p
		enc.Encode(p)
	}
}

var p Page

func send(w http.ResponseWriter, r *http.Request) {
	w.Header()["Access-Control-Allow-Origin"] = []string{"*"}

	enc := json.NewEncoder(w)
	enc.Encode(p)
}

func newPage() {
	p = Page{
		NextUrl: "news",
	}

	p.Init()
	p.GetNext()

	<-time.After(1 * time.Minute)

	newPage()
}

func server() {
	newPage()

	http.HandleFunc("/next/", next)
	http.HandleFunc("/", send)
	http.HandleFunc(commentRoute, getComments)

	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Fatal(err)
	}
}
