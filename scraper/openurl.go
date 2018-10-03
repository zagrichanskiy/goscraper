package scraper

import (
	"fmt"
	"net/http"
)

// OpenURL to connect to the server
type OpenURL struct {
	url string
}

// NewOpenURL creates the state
func NewOpenURL(url string) *OpenURL {
	ret := &OpenURL{url}
	return ret
}

// Do of the OpenURL
func (s *OpenURL) Do() State {
	fmt.Println("Openning url: ", s.url)
	resp, err := http.Get(s.url)
	if err != nil {
		fmt.Println("Can't open the url: ", err)
		return exitState
	}

	return NewPrintURLs(resp)
}
