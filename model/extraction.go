package model

import (
	"time"
)

type DataModel struct {
	Repository         *Repository
	Distribution       *Distribution
	VulnerabilityIndex *VulnerabilityIndex
}

type VulnerabilityIndex struct {
	TotalVulnerabilitiesCount int
}

type Repository struct {
	Contributors []Contributor
	Issues       []Issue
	Commits      []Commit
	Releases     []Release

	*RepositoryData
}

func (r *Repository) TotalCommits() int {
	return len(r.Commits)
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

	LOC int

	Forks    int
	Watchers int
	Stars    int

	Dependencies []string
	Dependents   []string
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
	ChangedFiles []string
	Message      string
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
	Date        time.Time
}

func (r Release) GetTimestamp() time.Time {
	return r.Date
}

type IssueState string

const (
	IssueStateClosed IssueState = "closed"
)

type IssueContribution struct {
	Time time.Time
}

func (ic IssueContribution) GetTimestamp() time.Time {
	return ic.Time
}

type Issue struct {
	Number            int
	Author            string
	AuthorAssociation string
	State             IssueState
	Title             string
	Content           string
	ClosedBy          string
	Contributions     []IssueContribution
	Contributors      []string
	CreationTime      time.Time
	FirstResponse     *time.Time
	LastUpdate        time.Time
	ClosingTime       time.Time
}

func (i Issue) GetTimestamp() time.Time {
	return i.CreationTime
}

type Contributor struct {
	Name              string
	Company           string
	Sponsors          int
	Organizations     int
	Contributions     int
	Repositories      int
	FirstContribution *time.Time
	LastContribution  *time.Time
}

type ContributorInfo struct {
	Repositories struct {
		TotalCount int
	}
	Sponsors struct {
		TotalCount int
	}
	Organizations struct {
		TotalCount int
	}
	Company string
	Login   string
}

type Distribution struct {
	Library  *Library
	Artifact *Artifact
}

type Artifact struct {
	Version              string
	Description          string
	ArtifactRepositories []string
	Date                 time.Time
	Vulnerabilities      []string
	Dependents           []string
	Dependencies         []string
	DeprecationWarning   bool // TODO: Cant interpret false confidently
	Contributors         []string
	Developers           []string
	Organization         string
	Licenses             []string
	MailingLists         []string
}

type Library struct {
	Ranking       *int
	Licenses      []string
	UsedBy        *int
	Versions      []string
	LastUpdated   time.Time
	LatestVersion string
	LatestRelease string
}
