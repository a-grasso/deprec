package model

type Dependency struct {
	Name     string
	Version  string
	MetaData map[string]string
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
	Contributors *Contributors
}
type Contributors struct {
	Sponsors          []string
	Organizations     []string
	Projects          []string
	FirstContribution string
	LastContribution  string
}

type Distribution struct {
}
