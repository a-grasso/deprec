package configuration

import (
	"encoding/json"
	"fmt"
	"os"
)

func Load(configFilePath string) (*Configuration, error) {
	content, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not read configuration file '%s': %s", configFilePath, err)
	}

	var config Configuration
	err = json.Unmarshal(content, &config)

	if err != nil {
		return nil, fmt.Errorf("could not parse configuration file '%s': %s", configFilePath, err)
	}

	return &config, nil
}
