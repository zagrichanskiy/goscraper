package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	// DirName folder in the $HOME
	DirName = ".goscraper"
	// ConfigName name of config
	ConfigName = "config.json"
	// DefaultURL to look up builds
	DefaultURL = "http://dehil-ae03.debads.europe.delphiauto.net:82/platforms/csp/csp/"
)

// JSONFormat data to store in config file
type JSONFormat struct {
	URL     string
	Sdk     string
	Blade1  string
	Blade2  string
	Blade3  string
	Updated string
}

func main() {
	var (
		dirPath    = filepath.Join(os.Getenv("HOME"), DirName)
		configPath = filepath.Join(dirPath, ConfigName)
	)

	if ok := checkFolder(dirPath); !ok {
		os.Exit(1)
	}

	if ok := checkConfig(configPath); !ok {
		os.Exit(1)
	}
}

func checkFolder(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Directory ", path, " doesn't exist, creating...")
		err := os.Mkdir(path, 0777)
		if err != nil {
			fmt.Println("Can't create directory ", path, ": ", err)
			return false
		}
	}

	return true
}

func checkConfig(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Config ", path, " doesn't exist, creating...")

		jFile, err := os.Create(path)
		if err != nil {
			fmt.Println("Can't create ", path, ": ", err)
			return false
		}
		defer jFile.Close()

		jData, err := json.MarshalIndent(&JSONFormat{URL: DefaultURL}, "", "    ")
		if err != nil {
			fmt.Println("Can't format json: ", err)
			return false
		}

		fmt.Fprintf(jFile, "%s\n", jData)
	}
	return true
}
