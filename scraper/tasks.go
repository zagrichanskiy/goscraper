package scraper

import (
	"fmt"
	"io"
	"net/http"
	"os"
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
