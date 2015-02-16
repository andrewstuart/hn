package hackernews

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPageCache(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gw := gzip.NewWriter(w)
		defer gw.Close()

		switch r.Method {
		case "HEAD":
			w.Header().Set("Set-Cookie", "__cfudid=123")
			break
		case "GET":
			fmt.Fprintln(gw, TestResponse)
			break
		}

	}))
	defer ts.Close()

	pc := NewPageCache()

	pc.Client = NewClient(ts.URL)

	p, err := pc.GetNext()

	if err != nil {
		t.Errorf("Error getting next page:\v\t%v", err)
	}

	if len(p.Articles) != 30 {
		t.Fatalf("Not enough pages returned: %d", len(p.Articles))
	}

	if len(pc.Articles) != 30 {
		t.Fatalf("Wrong number of pages in cache: %d", len(pc.Articles))
	}
}
