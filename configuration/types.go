package configuration

type GitHub struct {
	APIToken string `json:"APIToken"`
}

type MongoDB struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	URI      string `json:"URI"`
}

type AFConfig struct {
	CommitThreshold int `json:"CommitThreshold"`
}

type Configuration struct {
	GitHub   `json:"GitHub"`
	MongoDB  `json:"MongoDB"`
	AFConfig `json:"AFConfig"`
}
