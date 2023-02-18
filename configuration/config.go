package configuration

import (
	"encoding/json"
	"fmt"
	"github.com/a-grasso/deprec/logging"
	"github.com/joho/godotenv"
	"os"
)

func Load(configFilePath, envFilePath string) (*Configuration, error) {
	content, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not read configuration file '%s': %s", configFilePath, err)
	}

	var coresConfig CoresConfig
	err = json.Unmarshal(content, &coresConfig)

	if err != nil {
		return nil, fmt.Errorf("could not parse configuration file '%s': %s", configFilePath, err)
	}

	config := &Configuration{
		Extraction: Extraction{
			GitHub:   GitHub{},
			OSSIndex: OSSIndex{},
		},
		Cache: Cache{
			MongoDB: MongoDB{},
		},
		CoresConfig: coresConfig,
	}

	err = godotenv.Load(envFilePath)
	if err != nil {
		logging.Logger.Warn(fmt.Sprintf("error loading %s file", envFilePath))
	}

	var present bool

	config.Extraction.GitHub.APIToken, present = os.LookupEnv("GITHUB_API_TOKEN")
	if !present {
		logging.Logger.Warn("GITHUB_API_TOKEN environment variable missing!")
	}
	config.Extraction.OSSIndex.Username, present = os.LookupEnv("OSSINDEX_USERNAME")
	if !present {
		logging.Logger.Warn("OSSINDEX_USERNAME environment variable missing!")
	}
	config.Extraction.OSSIndex.Token, present = os.LookupEnv("OSSINDEX_TOKEN")
	if !present {
		logging.Logger.Warn("OSSINDEX_TOKEN environment variable missing!")
	}
	config.Cache.MongoDB.URI, present = os.LookupEnv("CACHE_MONGODB_URI")
	if !present {
		logging.Logger.Warn("CACHE_MONGODB_URI environment variable missing!")
	}
	config.Cache.MongoDB.Username, present = os.LookupEnv("CACHE_MONGODB_USERNAME")
	if !present {
		logging.Logger.Warn("CACHE_MONGODB_USERNAME environment variable missing!")
	}
	config.Cache.MongoDB.Password, present = os.LookupEnv("CACHE_MONGODB_PASSWORD")
	if !present {
		logging.Logger.Warn("CACHE_MONGODB_PASSWORD environment variable missing!")
	}

	return config, nil
}
