package scraper

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

// PrintURLs state
type PrintURLs struct {
	resp *http.Response
}

// NewPrintURLs creates new state
func NewPrintURLs(resp *http.Response) *PrintURLs {
	ret := &PrintURLs{resp}
	return ret
}

// Do implementation
func (s *PrintURLs) Do() State {
	fmt.Println("Searching links")
	z := html.NewTokenizer(s.resp.Body)

loop:
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			break loop
		case tt == html.StartTagToken:
			t := z.Token()

			if t.Data == "a" {
				printURL(&t)
			}
		}
	}
	return exitState
}

// printURL helper function
func printURL(t *html.Token) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			fmt.Println("Found href:", a.Val)
			break
		}
	}
}
