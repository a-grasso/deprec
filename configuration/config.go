package configuration

import (
	"encoding/json"
	"log"
	"os"
)

func Load(configFilePath string) *Configuration {
	content, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("Could not read configuration file '%s': %s", configFilePath, err)
	}

	var config Configuration
	err = json.Unmarshal(content, &config)

	if err != nil {
		log.Fatalf("Could not parse configuration file '%s': %s", configFilePath, err)
	}

	return &config
}
