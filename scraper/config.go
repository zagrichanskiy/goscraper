package scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

// DownloadInfo what to download.
type DownloadInfo struct {
	Download bool
	File     string
}

// Config represetns configuration in json.
type Config struct {
	RootURL   string
	LatestURL string
	SdkURL    string
	Blade1    DownloadInfo
	Blade2    DownloadInfo
	Blade3    DownloadInfo
	Sdk       DownloadInfo
	Updated   string
}

var defaultConfig = Config{
	RootURL:   "http://dehil-ae03.debads.europe.delphiauto.net:82/platforms/csp/csp/",
	LatestURL: "http://dehil-ae03.debads.europe.delphiauto.net:82/platforms/csp/csp/latest-release/",
	SdkURL:    "http://dehil-ae03.debads.europe.delphiauto.net:82/platforms/csp/csp/latest-release/sdk",
	Blade1: DownloadInfo{
		Download: true,
		File:     "csp-image-blade-i-intel-corei7-64-dom0.wic.md5"},
	Blade2: DownloadInfo{
		Download: true,
		File:     "csp-image-blade-ii-intel-corei7-64-dom0.wic.md5"},
	Blade3: DownloadInfo{
		Download: true,
		File:     "csp-image-blade-iii-intel-corei7-64-dom0.wic.md5"},
	Sdk: DownloadInfo{
		Download: true,
		// File:     "\\.sh$"}}
		File: "\\.host\\.manifest$"}}

// InitConfig opens configuration file or creates new.
func InitConfig(path string) Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Configuration doesn't exist, creating")
		err = write(path, defaultConfig)
		if err != nil {
			panic(err)
		}
		return defaultConfig
	}

	fmt.Println("Open configuration ", path)
	var config Config
	err := read(path, &config)
	if err != nil {
		panic(err)
	}
	return config
}

func write(path string, config Config) error {
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return errors.New("Can't marshal config")
	}
	data = append(data, '\n')

	err = ioutil.WriteFile(path, data, 0666)
	if err != nil {
		return errors.New("Can't write config to file")
	}

	return nil
}

func read(path string, config *Config) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.New("Can't read file " + path)
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		return errors.New("Can't unmarshal data from file " + path)
	}

	return nil
}
