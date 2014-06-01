package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

const commentRoute string = "/comments/"

var articles map[string]*Article

func getComments(w http.ResponseWriter, r *http.Request) {
	w.Header()["Access-Control-Allow-Origin"] = []string{"*"}

	idSt := r.URL.Path[len(commentRoute):]

	if id, err := strconv.Atoi(idSt); err == nil {
		enc := json.NewEncoder(w)
		if ar, cached := arsCache[id]; cached {
			enc.Encode(ar)
		} else {
			log.Print("Not cached")
		}
	} else {
		log.Print(err)
	}
}

func getPage(w http.ResponseWriter, r *http.Request) {
	reqUrl := r.URL.Path[len("/page/"):]

	w.Header()["Access-Control-Allow-Origin"] = []string{"*"}
	enc := json.NewEncoder(w)

	if page, cacheExists := pc.Pages[reqUrl]; cacheExists {
		enc.Encode(page)
	} else if reqUrl == pc.Next {
		page = pc.GetNext()
		enc.Encode(page)
	}
}

func send(w http.ResponseWriter, r *http.Request) {
	w.Header()["Access-Control-Allow-Origin"] = []string{"*"}

	enc := json.NewEncoder(w)
	enc.Encode(pc.Pages)
}

func reCache() {
	<-time.After(15 * time.Minute)

	pc = NewPageCache()
	go reCache()
}

var pc *PageCache

func server(addr string) {
	articles = make(map[string]*Article)
	pc = NewPageCache()

	for _, art := range pc.Articles {
		articles[strconv.Itoa(art.Id)] = art
	}

	http.HandleFunc("/page/", getPage)
	http.HandleFunc("/", send)
	http.HandleFunc(commentRoute, getComments)

	err := http.ListenAndServe(addr, nil)

	if err != nil {
		log.Fatal(err)
	}

	go reCache()
}
