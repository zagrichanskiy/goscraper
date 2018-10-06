package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/zagrichanskiy/goscraper/scraper"
	"golang.org/x/net/html"
)

const (
	// ProgramDir specifies the name of the directory to store data.
	ProgramDir = ".goscraper"
	// ConfigName is the name of the config file.
	ConfigName = "config.json"
)

func checkProgramDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Program directory doesn't exist, creating")
		err = os.Mkdir(path, 0775)
		if err != nil {
			fmt.Println("Can't create program directory")
			panic(err)
		}
	} else {
		fmt.Println("Program directory is", path)
	}
}

func getLatest(url string) string {
	fmt.Println("Checking latests builds on server")
	// Connect to the host
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Can't connect to the host")
		panic(err)
		// TODO: Try connection after some timeout.
	}
	defer resp.Body.Close()

	urls := make([]string, 0)
	fillUrls(resp, &urls)

	return urls[len(urls)-2]
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

func download(dir string, c scraper.Config) {
	fmt.Println("Downloading")

	// Initialization.
	ch := make(chan scraper.Status)
	tasks := make([]scraper.Task, 0)
	addTask := func(tasks *[]scraper.Task, download bool, url string, file string) {
		if download {
			t := scraper.NewBladeTask(url, dir, file)
			*tasks = append(*tasks, t)
		}
	}

	// Creating downloading tasks.
	addTask(&tasks, c.Blade1.Download, c.LatestURL, c.Blade1.File)
	addTask(&tasks, c.Blade2.Download, c.LatestURL, c.Blade2.File)
	addTask(&tasks, c.Blade3.Download, c.LatestURL, c.Blade3.File)

	// Invoking tasks.
	for _, task := range tasks {
		go task.Do(ch)
	}

	// Waiting for tasks to finish.
	for range tasks {
		s := <-ch
		switch {
		case s.Ok:
			fmt.Println("Downloaded: ", s.Message)
		case !s.Ok:
			fmt.Println("Not downloaded: ", s.Message)
			// TODO: Handle download errors.
		}
	}

	fmt.Println("Downloading is done")
}

func main() {
	var (
		// ProgramPath to the program direcotry.
		ProgramPath = os.Getenv("HOME") + "/" + ProgramDir
		// ConfigPath to the configuration file.
		ConfigPath = ProgramPath + "/" + ConfigName
	)

	// Program directory and configuration file.
	checkProgramDir(ProgramPath)
	c := scraper.InitConfig(ConfigPath)

	// Getting the date of the latest images on server.
	latest := getLatest(c.RootURL)
	fmt.Println("Latest builds on server are of:", latest)
	fmt.Println("Latest downloaded builds are of: ", c.Updated)

	// Downloading images.
	if c.Updated != latest {
		download(filepath.Join(ProgramPath, latest), c)
	}

	// TODO: Update 'Updated' field in config with latest info.
}
