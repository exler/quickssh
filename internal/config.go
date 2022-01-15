package internal

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
)

type Profile struct {
	Hostname string
	Port     int
	Username string
	Keyfile  string
}

const (
	appDir           = "quickssh"
	configFilename   = "config.json"
	profilesFilename = "profiles.json"
)

func loadJSON(path string, data interface{}) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		tmp := make(map[string]interface{})
		saveJSON(path, &tmp)
	}

	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln("Cannot read configuration file")
	}

	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Fatalf("Cannot decode JSON data")
	}
}

func saveJSON(filename string, data interface{}) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Fatalln("Cannot encode data as JSON")
	}

	err = os.WriteFile(filename, file, 0600)
	if err != nil {
		log.Fatalln("Cannot save configuration file")
	}
}

func getAppPath(config string) string {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalln("Cannot get user configuration directory")
	}

	configPath := filepath.Join(userConfigDir, appDir)
	if err = os.Mkdir(configPath, 0700); os.IsNotExist(err) {
		log.Fatalln("Cannot create app configuration directory")
	}

	configPath = filepath.Join(configPath, config)
	return configPath
}

func GetProfiles() (profiles map[string]Profile) {
	profilesPath := getAppPath(profilesFilename)
	loadJSON(profilesPath, &profiles)
	return
}

func SetProfiles(profiles map[string]Profile) {
	profilesPath := getAppPath(profilesFilename)
	saveJSON(profilesPath, &profiles)
}

func GetConfig() (config map[string]interface{}) {
	configPath := getAppPath(configFilename)
	loadJSON(configPath, &config)
	return
}

func SetConfig(config map[string]interface{}) {
	configPath := getAppPath(configFilename)
	saveJSON(configPath, &config)
}
