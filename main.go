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
		err = os.Mkdir(path, 0666)
		if err != nil {
			fmt.Println("Can't create program directory")
			panic(err)
		}
	} else {
		fmt.Println("Program directory is ", path)
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

// Status of downloading.
type Status struct {
	file string
	ok   bool
}

func download(path string, latest string, c scraper.Config) {
	fmt.Println("Downloading")

	downloadDir := filepath.Join(path, latest)
	if err := os.Mkdir(downloadDir, 0666); !os.IsExist(err) {
		fmt.Println("Can't create download directory:", downloadDir)
		panic(err)
	}

	ch := make(chan Status)
	tasks := 0
	runTasks := func(ch chan Status, t *int, toDownload bool, file string) {
		if toDownload {
			*t++
			filePath := filepath.Join(downloadDir, file)
			link := c.LatestURL + file
			fmt.Println("Downloading of", file)
			go downloadLink(ch, filePath, link)
		}
	}

	runTasks(ch, &tasks, c.Download.Blade1, c.Blade1)
	runTasks(ch, &tasks, c.Download.Blade2, c.Blade2)
	runTasks(ch, &tasks, c.Download.Blade3, c.Blade3)

	for tasks > 0 {
		status := <-ch
		fmt.Println(status.file, "downloading status:", status.ok)
		tasks--
	}

	fmt.Println("Downloading is done")
}

func downloadLink(ch chan Status, file string, link string) {
	fmt.Println("Downloading", file, "from", link)
	ch <- Status{file, true}
}

func main() {
	var (
		// ProgramPath to the program direcotry.
		ProgramPath = os.Getenv("HOME") + "/" + ProgramDir
		// ConfigPath to the configuration file.
		ConfigPath = ProgramPath + "/" + ConfigName
	)

	checkProgramDir(ProgramPath)
	c := scraper.InitConfig(ConfigPath)

	latest := getLatest(c.RootURL)
	fmt.Println("Latest builds on server are of:", latest)
	fmt.Println("Latest downloaded builds are of: ", c.Updated)

	if c.Updated != latest {
		download(ProgramPath, latest, c)
	}
}
