package scraper

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"golang.org/x/net/html"
)

// Status of the task.
type Status struct {
	Message string
	Ok      bool
}

// Task interface to run tasks.
type Task interface {
	Do(ch chan Status)
}

// BladeTask to download blades.
type BladeTask struct {
	Dir  string
	Link string
	File string
}

// NewBladeTask creates new task to download blade.
func NewBladeTask(rootURL string, rootDir string, file string) *BladeTask {
	rootURL = addSlash(rootURL)
	rootDir = addSlash(rootDir)
	return &BladeTask{
		rootDir,
		rootURL + file,
		rootDir + file}
}

// Do downloads blade.
func (t *BladeTask) Do(ch chan Status) {
	if err := os.Mkdir(t.Dir, 0775); err != nil && !os.IsExist(err) {
		fmt.Println("Can't create download directory:", t.Dir)
		panic(err)
	}
	fmt.Println("Downloading", t.Link, "to", t.File)
	download(ch, t.Link, t.File)
}

// SdkTask to download sdk.
type SdkTask struct {
	Dir  string
	Link string
	File string
}

// NewSdkTask creates new task for downloading sdk.
func NewSdkTask(rootURL string, rootDir string, expr string) *SdkTask {
	rootURL = addSlash(rootURL)
	rootDir = addSlash(rootDir)

	re, err := regexp.Compile(expr)
	if err != nil {
		fmt.Printf("Can't compile regular expression %s: %v\n", expr, err)
		return nil
	}

	resp, err := http.Get(rootURL)
	if err != nil {
		fmt.Printf("Can't open %s: %v\n", rootURL, err)
		return nil
	}
	defer resp.Body.Close()

	files := make([]string, 0)
	fillUrls(resp, &files)
	for _, file := range files {
		if re.MatchString(file) {
			return &SdkTask{
				rootDir,
				rootURL + file,
				rootDir + file}
		}
	}

	fmt.Printf("Can't match pattern %s in sdk folder\n", expr)
	return nil
}

func fillUrls(r *http.Response, urls *[]string) {
	z := html.NewTokenizer(r.Body)

	for tt := z.Next(); tt != html.ErrorToken; tt = z.Next() {
		if tt == html.StartTagToken {
			t := z.Token()

			if t.Data == "a" {
				for _, a := range t.Attr {
					if a.Key == "href" {
						*urls = append(*urls, a.Val)
						break
					}
				}
			}
		}
	}
}

func addSlash(s string) string {
	if s[len(s)-1] != '/' {
		return s + "/"
	}
	return s
}

func download(ch chan Status, link string, file string) {
	resp, err := http.Get(link)
	if err != nil {
		fmt.Printf("Can't open link %s: %v\n", link, err)
		ch <- Status{file, false}
		return
	}
	defer resp.Body.Close()

	f, err := os.Create(file)
	if err != nil {
		fmt.Printf("Can't create %s: %v\n", file, err)
		ch <- Status{file, false}
		return
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		fmt.Printf("Can't download %s, %v\n", link, err)
		ch <- Status{file, false}
		return
	}

	ch <- Status{file, true}
}
