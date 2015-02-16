package hackernews

import (
	"fmt"
	"math/rand"
	"time"
)

//A structure created for caching pages for a given amount of time. This avoids heavy traffic to the HN servers.
type PageCache struct {
	Created  time.Time        `json:"created"`
	Pages    map[string]*Page `json:"pages"`
	Articles []*Article       `json:"articles"`
	Next     string           `json:"next"`
	*Client
}

//Make a random string.
func RandomString(length int) string {
	if length == 0 {
		length = 80
	}

	rand.Seed(time.Now().Unix())
	b := make([]byte, 0, length)

	return string(b)
}

//Function to return a new page cache.
//As a note, it will be empty. GetNext() will need to be called for anything to show up.
func NewPageCache() *PageCache {
	pc := PageCache{
		Next:     "news",
		Pages:    make(map[string]*Page),
		Articles: make([]*Article, 0, 100),
		Client:   NewClient(YcRoot),
	}

	return &pc
}

//Gets the next page from a cache.
func (pc *PageCache) GetNext() (*Page, error) {
	p, err := pc.RetrievePage(pc.Next)

	if err != nil {
		return nil, fmt.Errorf("Error retrieving next page:\n\t %v", err)
	}

	pc.Pages[p.Url] = p
	pc.Next = p.NextUrl

	pc.Articles = append(pc.Articles, p.Articles...)

	return p, nil
}

//Get a given number of articles
func (pc *PageCache) GetArticles(articleCount int) error {
	articlesNeeded := len(pc.Articles) + articleCount
	for len(pc.Articles) < articlesNeeded {
		_, err := pc.GetNext()

		if err != nil {
			return err
		}
	}

	return nil
}
