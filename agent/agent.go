package agent

import (
	"deprec/configuration"
	"deprec/extraction"
	"deprec/mapping"
	"deprec/model"
	"log"
	"strings"
)

type Agent struct {
	Dependency *model.Dependency
	DataModel  *model.DataModel
	Config     *configuration.Configuration
}

func NewAgent(dependency *model.Dependency, configuration *configuration.Configuration) *Agent {
	agent := Agent{Dependency: dependency, DataModel: &model.DataModel{}, Config: configuration}
	return &agent
}

func (agent *Agent) Start() model.AgentResult {
	log.Printf("Starting Extraction...")
	agent.Extraction()
	log.Printf("...Extraction complete")

	log.Printf("Starting Combination & Conclusion...")
	result := agent.CombinationAndConclusion()
	log.Printf("...Combination & Conclusion complete")

	return *result
}

func (agent *Agent) Extraction() {

	if strings.Contains(agent.Dependency.MetaData["vcs"], "github") {
		extraction.NewGitHubExtractor(agent.Dependency, agent.Config.GitHub).Extract(agent.DataModel)
	}
}

func (agent *Agent) CombinationAndConclusion() *model.AgentResult {

	network := mapping.Network(agent.DataModel)

	return &model.AgentResult{
		Dependency: agent.Dependency,
		Result:     network,
	}
}
