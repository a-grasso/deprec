package configuration

type GitHub struct {
	APIToken string `json:"APIToken"`
}

type OSSIndex struct {
	Username string `json:"Username"`
	Token    string `json:"Token"`
}

type MongoDB struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	URI      string `json:"URI"`
}

type CoresConfig struct {
	Activity                Activity                `json:"Activity"`
	Recentness              Recentness              `json:"Recentness"`
	Processing              Processing              `json:"Processing"`
	Network                 Network                 `json:"Network"`
	Popularity              Popularity              `json:"Popularity"`
	CoreTeam                CoreTeam                `json:"CoreTeam"`
	OrgBackup               OrgBackup               `json:"OrgBackup"`
	Engagement              Engagement              `json:"Engagement"`
	ThirdPartyParticipation ThirdPartyParticipation `json:"ThirdPartyParticipation"`
}

type Activity struct {
	Percentile float64 `json:"Percentile"`
}
type ThirdPartyParticipation struct {
	CommitLimit                         int `json:"CommitLimit"`
	ThirdPartyCommitThresholdPercentage int `json:"ThirdPartyCommitThresholdPercentage"`
}

type Network struct {
	Threshold int `json:"Threshold"`
}
type Popularity struct {
	Threshold int `json:"Threshold"`
}

type Recentness struct {
	CommitLimit                 int     `json:"CommitLimit"`
	ReleaseLimit                int     `json:"ReleaseLimit"`
	TimeframePercentileCommits  float64 `json:"TimeframePercentileCommits"`
	TimeframePercentileReleases float64 `json:"TimeframePercentileReleases"`
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
	ActiveContributorsPercentile          float64 `json:"ActiveContributorsPercentile"`
	ActiveContributorsThresholdPercentage float64 `json:"ActiveContributorsThresholdPercentage"`
	CoreTeamStrengthThresholdPercentage   float64 `json:"CoreTeamStrengthThresholdPercentage"`
}

type Extraction struct {
	GitHub   GitHub   `json:"GitHub"`
	OSSIndex OSSIndex `json:"OSSIndex"`
}

type Cache struct {
	MongoDB MongoDB `json:"MongoDB"`
}

type Configuration struct {
	Extraction  `json:"Extraction"`
	Cache       `json:"Cache"`
	CoresConfig `json:"CoresConfig"`
}
