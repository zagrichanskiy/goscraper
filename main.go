package main

import (
	"fmt"
	"os"

	"github.com/zagrichanskiy/goscraper/scraper"
)

const (
	// ProgramDir specifies the name of the directory to store data.
	ProgramDir = ".goscraper"
	// ConfigName is the name of the config file.
	ConfigName = "config.json"
)

var (
	// Path to the program direcotry.
	Path = os.Getenv("HOME") + "/" + ProgramDir
	// ConfigPath to the configuration file.
	ConfigPath = Path + "/" + ConfigName
)

func main() {
	checkProgramDir(Path)
	c := scraper.InitConfig(ConfigPath)
	fmt.Printf("%+v\n", c)
}

func checkProgramDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Program directory doesn't exist, creating")
		err = os.Mkdir(path, 0666)
		if err != nil {
			fmt.Println("Can't create program directory")
			panic(err)
		}
	} else {
		fmt.Println("Program directory is ", Path)
	}
}
