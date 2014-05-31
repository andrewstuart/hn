package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"code.google.com/p/goncurses"
	"github.com/PuerkitoBio/goquery"
)

const YC_ROOT = "https://news.ycombinator.com/"
const ROWS_PER_ARTICLE = 3

var agoRegexp = regexp.MustCompile(`((?:\w*\W){2})(?:ago)`)

var trans *http.Transport = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

var client *http.Client = &http.Client{Transport: trans}

//Comments structure for HN articles
type Comment struct {
	Text     string     `json:"text"`
	User     string     `json:"user"`
	Id       int        `json:"id"`
	Created  time.Time  `json:"created,omitempty"`
	Comments []*Comment `json:"comments"`
}

func (c *Comment) String() string {
	return fmt.Sprintf("%s: %s\n", c.User, c.Text)
}

//Article structure
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

var arsCache = make(map[int]*Article)

var timeName = map[string]time.Duration{
	"second": time.Second,
	"minute": time.Minute,
	"hour":   time.Hour,
	"day":    24 * time.Hour,
}

func parseCreated(s string) time.Time {
	agoStr := agoRegexp.FindStringSubmatch(s)

	words := strings.Split(agoStr[1], " ")

	if count, err := strconv.Atoi(words[0]); err == nil {
		durText := words[1]

		if durText[len(durText)-1] == 's' {
			durText = durText[:len(durText)-1]
		}

		dur := timeName[durText]
		diff := -int64(count) * int64(dur)
		return time.Now().Add(time.Duration(diff)).Round(dur)
	} else {
		return time.Time{}
	}
}

//Retreives comments for a given article
func (a *Article) GetComments() {
	if _, exists := arsCache[a.Id]; exists {
		return
	}

	a.Comments = make([]*Comment, 0)

	articleUrl := YC_ROOT + "/item?id=" + strconv.Itoa(a.Id)

	resp, e := client.Get(articleUrl)

	if e != nil {
		log.Println(e)
	}

	if doc, e := goquery.NewDocumentFromResponse(resp); e != nil {
		log.Println(e)
	} else {

		commentStack := make([]*Comment, 1, 10)

		doc.Find("span.comment").Each(func(i int, comment *goquery.Selection) {
			text := ""
			user := comment.Parent().Find("a").First().Text()

			text += comment.Text()

			//Get around HN's little weird "reply" nesting randomness
			//Is it part of the comment, or isn't it?
			if last5 := len(text) - 5; last5 > 0 && text[last5:] == "reply" {
				text = text[:last5]
			}

			c := &Comment{
				User:     user,
				Text:     text,
				Comments: make([]*Comment, 0),
			}

			//Get comment create time
			t := comment.Prev().Text()

			c.Created = parseCreated(t)

			//Get id
			if idAttr, exists := comment.Prev().Find("a").Last().Attr("href"); exists {
				idSt := strings.Split(idAttr, "=")[1]

				if id, err := strconv.Atoi(idSt); err == nil {
					c.Id = id
				}
			}

			//Track the comment offset for nesting.
			if width, exists := comment.Parent().Prev().Prev().Find("img").Attr("width"); exists {
				offset, _ := strconv.Atoi(width)
				offset = offset / 40

				lastEle := len(commentStack) - 1 //Index of the last element in the stack

				if offset > lastEle {
					commentStack = append(commentStack, c)
					commentStack[lastEle].Comments = append(commentStack[lastEle].Comments, c)
				} else {

					if offset < lastEle {
						commentStack = commentStack[:offset+1] //Trim the stack
					}

					commentStack[offset] = c

					//Add the comment to its parents
					if offset == 0 {
						a.Comments = append(a.Comments, c)
					} else {
						commentStack[offset-1].Comments = append(commentStack[offset-1].Comments, c)
					}
				}
			}
		})
	}

	arsCache[a.Id] = a

	<-time.After(5 * time.Minute)

	delete(arsCache, a.Id)
}

func (a *Article) String() string {
	return fmt.Sprintf("(%d) %s: %s\n\n", a.Karma, a.User, a.Title)
}

func commentString(cs []*Comment, off string) string {
	s := ""

	for i, c := range cs {
		s += off + fmt.Sprintf("%d. %s\n", i+1, c)

		if len(c.Comments) > 0 {
			s += commentString(c.Comments, off+strconv.Itoa(i+1)+".")
		}
	}

	return s
}

//TODO Use stringer interface
func (a *Article) PrintComments() {
	a.GetComments()

	scr.Print(a)

	cs := commentString(a.Comments, "")

	scr.Print(cs)
}

type Page struct {
	NextUrl  string     `json:"next"`
	Url      string     `json:"url"`
	Articles []*Article `json:"articles"`
	cfduid   string
}

//Get a new page by passing a url
func NewPage(url string) *Page {
	p := Page{
		Url: url,
	}

	url = YC_ROOT + url

	if resp, err := client.Get(url); err == nil {
		c := resp.Cookies()
		p.cfduid = c[0].Raw
	} else {
		goncurses.End()
		log.Println(resp)
		log.Println(err)
	}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		goncurses.End()
		log.Println(err)
	}

	req.Header.Set("cookie", p.cfduid)
	req.Header.Set("referrer", "https://news.ycombinator.com/news")
	//TODO Find out what to do for this
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/35.0.1916.114 Safari/537.36")

	req.Write(os.Stdout)

	if resp, e := client.Do(req); e != nil {
		log.Println(e)
	} else {

		if doc, e := goquery.NewDocumentFromResponse(resp); e != nil {
			log.Println(e)
		} else {

			//Get all the trs with subtext for children then go back one (for the first row)
			rows := doc.Find(".subtext").ParentsFilteredUntil("tr", "tbody").Prev()

			var a bool

			p.NextUrl, a = doc.Find("td.title").Last().Find("a").Attr("href")

			for p.NextUrl[0] == '/' {
				p.NextUrl = p.NextUrl[1:]
			}

			if !a {
				goncurses.End()
				log.Println("Could not retreive next hackernews page. Time to go outside?")
			}

			rows.Each(func(i int, row *goquery.Selection) {
				ar := Article{
					Rank: len(p.Articles) + i,
				}

				title := row.Find(".title").Eq(1)
				link := title.Find("a").First()

				ar.Title = link.Text()

				if url, exists := link.Attr("href"); exists {
					ar.Url = url
				}

				row = row.Next()

				row.Find("span").Each(func(i int, s *goquery.Selection) {
					if karma, err := strconv.Atoi(strings.Split(s.Text(), " ")[0]); err == nil {
						ar.Karma = karma
					} else {
						log.Println(err)
					}

					if idSt, exists := s.Attr("id"); exists {
						if id, err := strconv.Atoi(strings.Split(idSt, "_")[1]); err == nil {
							ar.Id = id
						} else {
							log.Println(err)
						}
					}
				})

				sub := row.Find("td.subtext")
				t := sub.Text()

				ar.Created = parseCreated(t)

				ar.User = sub.Find("a").First().Text()

				comStr := strings.Split(sub.Find("a").Last().Text(), " ")[0]

				if comNum, err := strconv.Atoi(comStr); err == nil {
					ar.NumComments = comNum
				}

				p.Articles = append(p.Articles, &ar)

			})
		}
	}

	return &p
}
