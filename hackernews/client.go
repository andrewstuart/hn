package hackernews

import (
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const YC_ROOT = "https://news.ycombinator.com/"

//Client for interacting with news.ycombinator.com
type Client struct {
	RootUrl, cfduid string
	http.Client
}

func NewClient(rootUrl string) *Client {
	c := &Client{
		RootUrl: rootUrl,
		Client: http.Client{
			Transport: &http.Transport{
				TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
				DisableCompression: false,
			},
		},
	}

	if head, err := http.NewRequest("HEAD", c.RootUrl, nil); err == nil {
		if resp, err := c.Do(head); err == nil {
			c.cfduid = resp.Cookies()[0].Raw
		} else {
			log.Fatal("Could not determine cfduid (cloudflare id)")
		}
	} else {
		log.Fatal("Error creating HEAD request for cfduid (cloudflare id)")
	}

	return c
}

//Do a request
func (c *Client) doReq(req *http.Request) (*goquery.Document, error) {
	req.Header.Set("cookie", c.cfduid)
	req.Header.Set("referrer", "https://news.ycombinator.com/news")
	req.Header.Set("user-agent", "CLI Scraper (github.com/andrewstuart/hn)")
	req.Header.Set("accept-encoding", "gzip")

	if resp, err := c.Do(req); err == nil {
		if unzipper, err := gzip.NewReader(resp.Body); err == nil {
			if doc, err := goquery.NewDocumentFromReader(unzipper); err == nil {
				return doc, nil
			} else {
				return nil, fmt.Errorf("goquery error: %v", err)
			}
		} else {
			return nil, fmt.Errorf("gzip error: %v", err)
		}
	} else {
		return nil, fmt.Errorf("http error: %v", err)
	}
}
