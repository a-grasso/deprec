package model

import (
	"time"
)

type DataModel struct {
	Repository   *Repository
	Distribution *Distribution
}

type Repository struct {
	Contributors []Contributor
	Issues       []Issue
	Commits      []Commit
	Releases     []Release
	Tags         []Tag

	*RepositoryData
}

func (r *Repository) TotalCommits() int {
	return len(r.Commits)
}

func (r *Repository) TotalTags() int {
	return len(r.Tags)
}

func (r *Repository) TotalIssues() int {
	return len(r.Issues)
}

func (r *Repository) TotalReleases() int {
	return len(r.Releases)
}

func (r *Repository) TotalContributors() int {
	return len(r.Contributors)
}

type RepositoryData struct {
	Name      string
	Owner     string
	Org       *Organization
	CreatedAt time.Time
	Size      int

	License      string
	AllowForking bool

	ReadMe      string
	About       string
	Archivation bool
	Disabled    bool

	LOC      int
	TotalPRs int

	Forks    int
	Watchers int
	Stars    int

	Dependencies []string
	Dependents   []string

	CommunityStandards float64
}

type Organization struct {
	Login             string
	PublicRepos       int
	Followers         int
	Following         int
	TotalPrivateRepos int
	OwnedPrivateRepos int
	Collaborators     int
}

type Commit struct {
	Author       string
	Committer    string
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

func (c Commit) GetTimestamp() time.Time {
	return c.Timestamp
}

type Release struct {
	Author      string
	Version     string
	Description string
	Changes     []string
	Type        string
	Date        time.Time
}

func (r Release) GetTimestamp() time.Time {
	return r.Date
}

type Tag struct {
	Author      string
	Version     string
	Description string
	Date        time.Time
}

func (c Tag) GetTimestamp() time.Time {
	return c.Date
}

type IssueState string

const (
	IssueStateOpen   IssueState = "open"
	IssueStateClosed IssueState = "closed"
	IssueStateOther  IssueState = "other"
)

func ToIssueState(s string) IssueState {
	switch s {
	case "open":
		return IssueStateOpen
	case "closed":
		return IssueStateClosed
	default:
		return IssueStateOther
	}
}

type IssueContribution struct {
	Time time.Time
}

func (ic IssueContribution) GetTimestamp() time.Time {
	return ic.Time
}

func (i Issue) GetTimestamp() time.Time {
	return i.CreationTime
}

type Issue struct {
	Number            int
	Author            string
	AuthorAssociation string
	Labels            []string
	State             IssueState
	Title             string
	Content           string
	ClosedBy          string
	Contributions     []IssueContribution
	Contributors      []string
	CreationTime      time.Time
	FirstResponse     time.Time
	LastContribution  time.Time
	ClosingTime       time.Time
}

type Contributor struct {
	Name                    string
	Company                 string
	Sponsors                int
	Organizations           int
	Contributions           int
	Repositories            int
	FirstContribution       *time.Time
	LastContribution        *time.Time
	TotalStatsContributions int
}

type Distribution struct {
	Library  *Library
	Artifact *Artifact
}

type Artifact struct {
	NewVersionAvailable  bool
	ArtifactRepositories []string
	Date                 time.Time
	Vulnerabilities      []string
	Dependents           []string
	Dependencies         []string
	DeprecationWarning   bool
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
