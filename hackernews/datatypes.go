package hackernews

import (
	"fmt"
	"time"
)

//A Comment is a structure for storing data about each HN comment as well as nested
//children
type Comment struct {
	Text     string     `json:"text"`
	User     string     `json:"user"`
	Id       int        `json:"id"`
	Created  time.Time  `json:"created,omitempty"`
	Comments []*Comment `json:"comments,omitempty"`
}

//HackerNews article structure
type Article struct {
	Rank        int        `json:"rank"`
	Title       string     `json:"title"`
	Karma       int        `json:"karma"`
	Id          int        `json:"id"`
	Url         string     `json:"url"`
	NumComments int        `json:"numComments"`
	Comments    []*Comment `json:"comments,omitempty"`
	User        string     `json:"user"`
	CreatedAgo  string     `json:"createdAgo,omitempty"`
	Created     time.Time  `json:"created,omitempty"`
}

//A single page for keeping track of where articles originate and where we should
//retreive the next articles
type Page struct {
	NextUrl  string     `json:"next"`
	Url      string     `json:"url"`
	Articles []*Article `json:"articles"`
}

//Get a new page from a url
func NewPage(url string) *Page {
	return &Page{
		Url:      url,
		Articles: make([]*Article, 0, 15),
	}
}

type HnError struct {
	Message string
	Code    HnErrorCode
}

type HnErrorCode int

const (
	TransportError = HnErrorCode(iota)
	SiteError
	CloudFlareError
)

func (e HnError) Error() string {
	return fmt.Sprintf("Error code %d: %s", e.Code, e.Message)
}
