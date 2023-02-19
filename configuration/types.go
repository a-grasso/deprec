package configuration

type GitHub struct {
	APIToken string `json:"APIToken,omitempty"`
}

type OSSIndex struct {
	Username string `json:"Username,omitempty"`
	Token    string `json:"Token,omitempty"`
}

type MongoDB struct {
	Username string `json:"Username,omitempty"`
	Password string `json:"Password,omitempty"`
	URI      string `json:"URI,omitempty"`
}

type CoresConfig struct {
	CombCon CombCon `json:"CombCon"`

	DeityGiven         DeityGiven         `json:"DeityGiven"`
	Circumstances      Circumstances      `json:"Circumstances"`
	Effort             Effort             `json:"Effort"`
	Support            Support            `json:"Support"`
	Community          Community          `json:"Community"`
	Interconnectedness Interconnectedness `json:"Interconnectedness"`

	Vulnerabilities Vulnerabilities `json:"Vulnerabilities"`
	Activity        Activity        `json:"Activity"`
	Recentness      Recentness      `json:"Recentness"`
	Processing      Processing      `json:"Processing"`
	Network         Network         `json:"Network"`
	Popularity      Popularity      `json:"Popularity"`
	CoreTeam        CoreTeam        `json:"CoreTeam"`
	Backup          Backup          `json:"Backup"`
	Engagement      Engagement      `json:"Engagement"`
	Participation   Participation   `json:"Participation"`
	Prestige        Prestige        `json:"Prestige"`
	Licensing       Licensing       `json:"Licensing"`
	Rivalry         Rivalry         `json:"Rivalry"`
	ProjectQuality  ProjectQuality  `json:"ProjectQuality"`
	Marking         Marking         `json:"Marking"`
}

type CombCon struct {
	Weights struct {
		DeityGiven         float64 `json:"DeityGiven,omitempty"`
		Circumstances      float64 `json:"Circumstances,omitempty"`
		Effort             float64 `json:"Effort,omitempty"`
		Support            float64 `json:"Support,omitempty"`
		Community          float64 `json:"Community,omitempty"`
		Interconnectedness float64 `json:"Interconnectedness,omitempty"`
	} `json:"Weights"`
}

type DeityGiven struct {
	Weights struct {
		Marking         float64 `json:"Marking,omitempty"`
		Vulnerabilities float64 `json:"Vulnerabilities,omitempty"`
	} `json:"Weights"`
}
type Circumstances struct {
	Weights struct {
		Rivalry        float64 `json:"Rivalry,omitempty"`
		Licensing      float64 `json:"Licensing,omitempty"`
		ProjectQuality float64 `json:"ProjectQuality,omitempty"`
	} `json:"Weights"`
}
type Effort struct {
	Weights struct {
		Activity   float64 `json:"Activity,omitempty"`
		Recentness float64 `json:"Recentness,omitempty"`
		CoreTeam   float64 `json:"CoreTeam,omitempty"`
	} `json:"Weights"`
}

type Support struct {
	Weights struct {
		Processing float64 `json:"Processing,omitempty"`
		Engagement float64 `json:"Engagement,omitempty"`
	} `json:"Weights"`
}

type Community struct {
	Weights struct {
		Prestige      float64 `json:"Prestige,omitempty"`
		Backup        float64 `json:"Backup,omitempty"`
		Participation float64 `json:"Participation,omitempty"`
	} `json:"Weights"`
}

type Interconnectedness struct {
	Weights struct {
		Popularity float64 `json:"Popularity,omitempty"`
		Network    float64 `json:"Network,omitempty"`
	} `json:"Weights"`
}

type Marking struct {
	ReadMeKeywords              []string `json:"ReadMeKeywords,omitempty"`
	AboutKeywords               []string `json:"AboutKeywords,omitempty"`
	ArtifactDescriptionKeywords []string `json:"ArtifactDescriptionKeywords,omitempty"`

	Weights struct {
		ReadMe      float64 `json:"ReadMe,omitempty"`
		About       float64 `json:"About,omitempty"`
		Archivation float64 `json:"Archivation,omitempty"`
		Artifact    float64 `json:"Artifact,omitempty"`
	} `json:"Weights"`
}

type ProjectQuality struct {
	Weights struct {
		ReadMe       float64 `json:"ReadMe,omitempty"`
		License      float64 `json:"License,omitempty"`
		About        float64 `json:"About,omitempty"`
		AllowForking float64 `json:"AllowForking,omitempty"`
	} `json:"Weights"`
}

type Vulnerabilities struct {
	Weights struct {
		CVE float64 `json:"CVE,omitempty"`
	} `json:"Weights"`
}
type Rivalry struct {
	Weights struct {
		IsLatest float64 `json:"IsLatest,omitempty"`
	} `json:"Weights"`
}

type Licensing struct {
	Weights struct {
		Repository float64 `json:"Repository,omitempty"`
		Artifact   float64 `json:"Artifact,omitempty"`
		Library    float64 `json:"Library,omitempty"`
	} `json:"Weights"`
}

type Prestige struct {
	Weights struct {
		Contributors float64 `json:"Contributors,omitempty"`
	} `json:"Weights"`
}
type Activity struct {
	Percentile float64 `json:"Percentile,omitempty"`
	Weights    struct {
		Commits            float64 `json:"Commits,omitempty"`
		Releases           float64 `json:"Releases,omitempty"`
		Issues             float64 `json:"Issues,omitempty"`
		IssueContributions float64 `json:"IssueContributions,omitempty"`
	} `json:"Weights"`
}

type Participation struct {
	CommitLimit               int `json:"CommitLimit,omitempty"`
	ThirdPartyCommitThreshold int `json:"ThirdPartyCommitThreshold,omitempty"`
	Weights                   struct {
		ThirdPartyCommits float64 `json:"ThirdPartyCommits,omitempty"`
	} `json:"Weights"`
}

type Network struct {
	Threshold int `json:"Threshold,omitempty"`
	Weights   struct {
		RepositoryNetwork float64 `json:"RepositoryNetwork,omitempty"`
	} `json:"Weights"`
}
type Popularity struct {
	Threshold int `json:"Threshold,omitempty"`
	Weights   struct {
		RepositoryPopularity float64 `json:"RepositoryPopularity,omitempty"`
	} `json:"Weights"`
}

type Recentness struct {
	CommitLimit                int     `json:"CommitLimit,omitempty"`
	ReleaseLimit               int     `json:"ReleaseLimit,omitempty"`
	TimeframePercentileCommits float64 `json:"TimeframePercentileCommits,omitempty"`
	Weights                    struct {
		MonthsSinceLastCommit         float64 `json:"MonthsSinceLastCommit,omitempty"`
		AverageMonthsSinceLastCommits float64 `json:"AverageMonthsSinceLastCommits,omitempty"`
		MonthsSinceLastRelease        float64 `json:"MonthsSinceLastRelease,omitempty"`
	} `json:"Weights"`
}

type Processing struct {
	ClosingTimeLimit int     `json:"ClosingTimeLimit,omitempty"`
	BurnPercentile   float64 `json:"BurnPercentile,omitempty"`
	Weights          struct {
		AverageClosingTime float64 `json:"AverageClosingTime,omitempty"`
		Burn               float64 `json:"Burn,omitempty"`
	} `json:"Weights"`
}

type Engagement struct {
	IssueCommentsRatioThreshold float64 `json:"IssueCommentsRatioThreshold,omitempty"`
	Weights                     struct {
		IssueCommentsRatio float64 `json:"IssueCommentsRatio,omitempty"`
	} `json:"Weights"`
}
type Backup struct {
	CompanyThreshold      int     `json:"CompanyThreshold,omitempty"`
	SponsorThreshold      float64 `json:"SponsorThreshold,omitempty"`
	OrganizationThreshold float64 `json:"OrganizationThreshold,omitempty"`
	Weights               struct {
		Companies              float64 `json:"Companies,omitempty"`
		Sponsors               float64 `json:"Sponsors,omitempty"`
		Organizations          float64 `json:"Organizations,omitempty"`
		RepositoryOrganization float64 `json:"RepositoryOrganization,omitempty"`
	} `json:"Weights"`
}

type CoreTeam struct {
	ActiveContributorsPercentile          float64 `json:"ActiveContributorsPercentile,omitempty"`
	ActiveContributorsThresholdPercentage float64 `json:"ActiveContributorsThresholdPercentage,omitempty"`
	CoreTeamStrengthThresholdPercentage   float64 `json:"CoreTeamStrengthThresholdPercentage,omitempty"`
	Weights                               struct {
		ActiveContributors float64 `json:"ActiveContributors,omitempty"`
		CoreTeamStrength   float64 `json:"CoreTeamStrength,omitempty"`
	} `json:"Weights"`
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
