package configuration

type GitHub struct {
	APIToken string `json:"APIToken"`
}

type Configuration struct {
	GitHub `json:"GitHub"`
}
