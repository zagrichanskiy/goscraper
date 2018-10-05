package scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

// DownloadObj in the Config struct for json.
type DownloadObj struct {
	Blade1 bool
	Blade2 bool
	Blade3 bool
	Sdk    bool
}

// Config represetns configuration in json.
type Config struct {
	RootURL   string
	LatestURL string
	Blade1    string
	Blade2    string
	Blade3    string
	Sdk       string
	SdkReg    string
	Download  DownloadObj
	Updated   string
}

var defaultConfig = Config{
	RootURL:   "http://dehil-ae03.debads.europe.delphiauto.net:82/platforms/csp/csp/",
	LatestURL: "http://dehil-ae03.debads.europe.delphiauto.net:82/platforms/csp/csp/latest-release/",
	Blade1:    "csp-image-blade-i-intel-corei7-64-dom0.wic",
	Blade2:    "csp-image-blade-ii-intel-corei7-64-dom0.wic",
	Blade3:    "csp-image-blade-iii-intel-corei7-64-dom0.wic",
	Sdk:       "sdk",
	SdkReg:    "\\.sh$",
	Download: DownloadObj{
		Blade1: true,
		Blade2: true,
		Blade3: true,
		Sdk:    false}}

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
