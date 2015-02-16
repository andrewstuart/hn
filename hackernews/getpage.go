package hackernews

import (
	"fmt"
	"html"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

//Get a new page by passing a url
func (c *Client) RetrievePage(url string) (*Page, error) {
	//Trim leading slash if necessary
	if url[0] == '/' {
		url = url[1:]
	}

	//All urls must start with YC root (or test)
	urlForReq := fmt.Sprintf("%s/%s", c.RootUrl, url)
	req, err := http.NewRequest("GET", urlForReq, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request for url %s: %v", url, err)
	}

	doc, err := c.doReq(req)

	if err != nil {
		return nil, fmt.Errorf("Error doing request:\n\t %v", err)
	}

	//Get all the trs with subtext for children then go back one (for the first row)
	rows := doc.Find(".subtext").ParentsFilteredUntil("tr", "tbody").Prev()

	p := NewPage(url)

	//Get the next url
	if nextUrl, found := doc.Find("td.title").Last().Find("a").Attr("href"); found {
		p.NextUrl = nextUrl
	} else {
		return nil, fmt.Errorf("Could not retreive next hackernews page. Time to go outside?")
	}

	//Make sure NextUrl doesn't start with forward slash
	for len(p.NextUrl) > 0 && p.NextUrl[0] == '/' {
		p.NextUrl = p.NextUrl[1:]
	}

	//Parse articles
	rows.Each(func(i int, row *goquery.Selection) {
		ar := Article{
			Rank: len(p.Articles) + i,
		}

		title := row.Find(".title").Eq(1)
		link := title.Find("a").First()

		ar.Title = html.UnescapeString(link.Text())

		if url, exists := link.Attr("href"); exists {
			ar.Url = url
		}

		//Rows are used in pairs currently
		row = row.Next()

		row.Find("span").Each(func(i int, s *goquery.Selection) {
			if karma, err := strconv.Atoi(strings.Split(s.Text(), " ")[0]); err == nil {
				ar.Karma = karma
			}

			if idSt, exists := s.Attr("id"); exists {
				if id, err := strconv.Atoi(strings.Split(idSt, "_")[1]); err == nil {
					ar.Id = id
				}
			}
		})

		sub := row.Find("td.subtext")
		t := html.UnescapeString(sub.Text())

		//We can ignore the error safely here
		ar.Created, _ = parseCreated(t)

		//Get the username
		ar.User = html.UnescapeString(sub.Find("a").First().Text())

		//Get number of comments
		comStr := strings.Split(sub.Find("a").Last().Text(), " ")[0]

		if comNum, err := strconv.Atoi(comStr); err == nil {
			ar.NumComments = comNum
		}

		p.Articles = append(p.Articles, &ar)
	})

	return p, nil
}

//Parse out from a string the create time
var timeName = map[string]time.Duration{
	"second": time.Second,
	"minute": time.Minute,
	"hour":   time.Hour,
	"day":    24 * time.Hour,
}

var agoRegexp = regexp.MustCompile(`((?:\w*\W){2})(?:ago)`)

func parseCreated(s string) (time.Time, error) {
	agoStr := agoRegexp.FindStringSubmatch(s)

	if len(agoStr) < 2 {
		return time.Time{}, fmt.Errorf(`No "ago" string found in string %s`, s)
	}

	words := strings.Split(agoStr[1], " ")

	if count, err := strconv.Atoi(words[0]); err == nil {
		durText := words[1]

		if durText[len(durText)-1] == 's' {
			durText = durText[:len(durText)-1]
		}

		dur := timeName[durText]
		diff := -int64(count) * int64(dur)
		return time.Now().Add(time.Duration(diff)).Round(dur), nil
	} else {
		return time.Time{}, err
	}
}
