package configuration

type GitHub struct {
	APIToken string `json:"APIToken"`
}

type MongoDB struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	URI      string `json:"URI"`
}

type CoresConfig struct {
	Activity   Activity   `json:"Activity"`
	Recentness Recentness `json:"Recentness"`
	Processing Processing `json:"Processing"`
	Network    Network    `json:"Network"`
	Popularity Popularity `json:"Popularity"`
	CoreTeam   CoreTeam   `json:"CoreTeam"`
	OrgBackup  OrgBackup  `json:"OrgBackup"`
	Engagement Engagement `json:"Engagement"`
}

type Activity struct {
	Percentile float64 `json:"Percentile"`
}
type Network struct {
	Threshold int `json:"Threshold"`
}
type Popularity struct {
	Threshold int `json:"Threshold"`
}

type Recentness struct {
	CommitLimit         int     `json:"CommitLimit"`
	ReleaseLimit        int     `json:"ReleaseLimit"`
	TimeframePercentile float64 `json:"TimeframePercentile"`
}

type Processing struct {
	ClosingTimeLimit int     `json:"ClosingTimeLimit"`
	BurnPercentile   float64 `json:"BurnPercentile"`
}

type Engagement struct {
	IssueCommentsRatioThresholdPercentage float64 `json:"IssueCommentsRatioThresholdPercentage"`
}
type OrgBackup struct {
	CompanyThreshold      int     `json:"CompanyThreshold"`
	SponsorThreshold      float64 `json:"SponsorThreshold"`
	OrganizationThreshold float64 `json:"OrganizationThreshold"`
}

type CoreTeam struct {
	ActiveContributorsPercentile        float64 `json:"ActiveContributorsPercentile"`
	CoreTeamStrengthThresholdPercentage float64 `json:"CoreTeamStrengthThresholdPercentage"`
}

type Extraction struct {
	GitHub GitHub `json:"GitHub"`
}

type Cache struct {
	MongoDB MongoDB `json:"MongoDB"`
}

type Configuration struct {
	Extraction  `json:"Extraction"`
	Cache       `json:"Cache"`
	CoresConfig `json:"CoresConfig"`
}
