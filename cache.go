package main

import (
	"math/rand"
	"time"

	"github.com/andrewstuart/hn/hackernews"
)

//A structure created for caching pages for a given amount of time. This avoids heavy traffic to the HN servers.
type PageCache struct {
	Created  time.Time                   `json:"created"`
	Pages    map[string]*hackernews.Page `json:"pages"`
	Articles map[int]*hackernews.Article `json:"articles"`
	Next     string                      `json:"next"`
}

func RandomString() string {
	rand.Seed(time.Now().Unix())
	b := make([]byte, 80)

	return string(b)
}

func NewPageCache() *PageCache {
	pc := PageCache{
		Next:  "news",
		Pages: make(map[string]*hackernews.Page),
	}

	pc.GetNext()

	return &pc
}

func (pc *PageCache) GetNext() *hackernews.Page {
	p := hackernews.NewPage(pc.Next)
	pc.Pages[p.Url] = p
	pc.Next = p.NextUrl

	for _, art := range p.Articles {
		pc.Articles[art.Id] = art
	}

	return p
}
