package services

import (
	"encoding/json"
	"fmt"
	"os"
)

//Config configurations of program
type Config struct {
	To      string `json:to`
	From    string `json:form`
	DoBkp   bool   `json:doBkp`
	PathBkp string `json:pathBkp`
}

//LoadConfig Load configs from settings.json
func (c *Config) LoadConfig(path string) {
	configFile, err := os.Open(path)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&c)
	// return config
}
