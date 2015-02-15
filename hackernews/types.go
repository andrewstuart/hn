package hackernews

import (
	"fmt"
	"strings"
	"time"
)

//Comments structure for HN articles
type Comment struct {
	Text     string     `json:"text"`
	User     string     `json:"user"`
	Id       int        `json:"id"`
	Created  time.Time  `json:"created,omitempty"`
	Comments []*Comment `json:"comments,omitempty"`
}

//Stringer implementation for sensible logging
func (c *Comment) String() string {
	return strings.Replace(fmt.Sprintf("%s: %s", c.User, c.Text), "\n", " ", -1)
}

//HackerNews article structure
type Article struct {
	Rank        int        `json:"rank"`
	Title       string     `json:"title"xml:"`
	Karma       int        `json:"karma"`
	Id          int        `json:"id"`
	Url         string     `json:"url"`
	NumComments int        `json:"numComments"`
	Comments    []*Comment `json:"comments",omitempty`
	User        string     `json:"user"`
	CreatedAgo  string     `json:"createdAgo,omitempty"`
	Created     time.Time  `json:"created",omitempty`
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
