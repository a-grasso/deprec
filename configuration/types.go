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
	Activity   Activity   `json:"Activity"`
	Recentness Recentness `json:"Recentness"`
}

type Activity struct {
	Percentile int `json:"Percentile"`
}

type Recentness struct {
	CommitThreshold  int `json:"CommitThreshold"`
	ReleaseThreshold int `json:"ReleaseThreshold"`
}

type Configuration struct {
	GitHub   `json:"GitHub"`
	MongoDB  `json:"MongoDB"`
	AFConfig `json:"AFConfig"`
}
