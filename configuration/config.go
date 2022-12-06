package configuration

import (
	"deprec/logging"
	"encoding/json"
	"os"
)

func Load(configFilePath string) *Configuration {
	content, err := os.ReadFile(configFilePath)
	if err != nil {
		logging.SugaredLogger.Fatalf("could not read configuration file '%s': %s", configFilePath, err)
	}

	var config Configuration
	err = json.Unmarshal(content, &config)

	if err != nil {
		logging.SugaredLogger.Fatalf("could not parse configuration file '%s': %s", configFilePath, err)
	}

	return &config
}
