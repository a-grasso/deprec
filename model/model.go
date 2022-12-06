package model

import "time"

type Dependency struct {
	Name     string
	Version  string
	MetaData map[string]string
}

type SBOM struct {
	JsonContent string
}

type AgentResult struct {
	Dependency *Dependency
	Result     float64
}

type DataModel struct {
	Repository   *Repository
	Distribution *Distribution
}

type Repository struct {
	Contributors []*Contributor
	Issues       []*Issue
	Commits      []*Commit
	Releases     []*Release

	*RepositoryData
}

type RepositoryData struct {
	Owner     *string
	Org       *string
	CreatedAt *time.Time
	Size      *int

	License      *string
	AllowForking *bool

	ReadMe      *string
	About       *string
	Archivation *bool
	Disabled    *bool

	KLOC              *int
	TotalCommits      *int
	TotalIssues       *int
	TotalPRs          *int
	TotalContributors *int

	Forks    *int
	Watchers *int
	Stars    *int

	Dependencies *[]string
	Dependents   *[]string

	OpenIssues *int

	CommunityStandards *float64
}

type Commit struct {
	Author       string
	Changes      []string
	ChangedFiles []string
	Type         string
	Message      string
	Branch       string
	Timestamp    time.Time
	Additions    int
	Deletions    int
	Total        int
}
type Release struct {
	Author  string
	Version map[string]int
	Changes []string
	Type    string
	Date    time.Time
}
type Issue struct {
	Author           string
	Labels           []string
	Contributions    []string
	Contributors     []*Contributor
	CreationTime     string
	FirstResponse    string
	LastContribution string
	ClosingTime      string
	Content          string
}

type Contributor struct {
	Name               *string
	Sponsors           *[]string
	Organizations      *int
	Contributions      *int
	Repositories       *int
	FirstContribution  *string
	LastContribution   *string
	TotalContributions *int
}

type Distribution struct {
	Library  Library
	Artifact Artifact
}

type Artifact struct {
	NewVersionAvailable  bool
	ArtifactRepositories []string
	Date                 time.Time
	Vulnerabilities      []string
	Dependents           []string
	Dependencies         []string
}

type Library struct {
	Ranking       int
	License       string
	UsedBy        int
	Moved         bool
	Versions      []*map[string]float64
	Description   string
	LastPublished string
}
