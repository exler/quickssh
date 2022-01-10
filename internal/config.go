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
	Password string
}

const (
	appDir           = "quickssh"
	configFilename   = "config.json"
	profilesFilename = "profiles.json"
)

func getConfigPath(config string) string {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalln("Cannot get user configuration directory")
	}

	configPath := filepath.Join(userConfigDir, appDir)
	if err = os.Mkdir(configPath, 0600); os.IsNotExist(err) {
		log.Fatalln("Cannot create app configuration directory")
	}

	configPath = filepath.Join(configPath, config)
	return configPath
}

func GetProfiles() (config map[string]Profile) {
	configPath := getConfigPath(profilesFilename)

	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		SetProfiles(make(map[string]Profile))
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalln("Cannot read configuration file")
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Cannot decode JSON data")
	}
	return
}

func SetProfiles(config map[string]Profile) {
	configPath := getConfigPath(profilesFilename)
	file, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		log.Fatalln("Cannot encode data as JSON")
	}

	err = os.WriteFile(configPath, file, 0600)
	if err != nil {
		log.Fatalln("Cannot save configuration file")
	}
}
