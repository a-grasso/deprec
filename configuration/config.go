package configuration

import (
	"encoding/json"
	"log"
	"os"
)

func Load() *Configuration {
	content, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Could not read configuration file 'config.json': %s", err)
	}

	var config Configuration
	err = json.Unmarshal(content, &config)

	if err != nil {
		log.Fatalf("Could not parse configuration file 'config.json': %s", err)
	}

	return &config
}
