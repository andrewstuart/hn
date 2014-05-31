package main

import "time"

//A structure created for caching pages for a given amount of time. This avoids heavy traffic to the HN servers.
type PageCache struct {
	Created  time.Time        `json:"created"`
	Pages    map[string]*Page `json:"pages"`
	Articles []*Article       `json:"articles"`
	Next     string           `json:"next"`
}

func NewPageCache() *PageCache {
	pc := PageCache{
		Next:  "/news",
		Pages: make(map[string]*Page),
	}

	pc.GetNext()

	return &pc
}

func (pc *PageCache) GetNext() *Page {
	p := NewPage(pc.Next)
	pc.Pages[p.Url] = p
	pc.Next = p.NextUrl
	pc.Articles = append(pc.Articles, p.Articles...)

	return p
}
