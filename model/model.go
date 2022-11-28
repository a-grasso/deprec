package model

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

	Owner   string
	Org     string
	License string

	ReadMe      string
	Archivation bool
	About       string

	KLOC              int
	TotalCommits      int
	TotalIssues       int
	TotalPRs          int
	TotalContributors int

	Forks    int
	Watchers int
	Stars    int

	Dependencies []string
	Dependents   []string

	CommunityStandards float64
}

type Commit struct {
	Author       string
	Changes      []string
	ChangedFiles []string
	Type         string
	Message      string
	Branch       string
	Timestamp    string
}
type Release struct {
	Author  string
	Version map[string]int
	Changes []string
	Type    string
	Date    string
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
	Name              string
	Sponsors          []string
	Organizations     []string
	Repositories      int
	FirstContribution string
	LastContribution  string
}

type Distribution struct {
	Library  Library
	Artifact Artifact
}

type Artifact struct {
	NewVersionAvailable  bool
	ArtifactRepositories []string
	Date                 string
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
