package agent

import (
	"deprec/configuration"
	"deprec/extraction"
	"deprec/mapping"
	"deprec/model"
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
	agent.Extraction()

	result := agent.CombinationAndConclusion()

	return *result
}

func (agent *Agent) Extraction() {

	if strings.Contains(agent.Dependency.MetaData["vcs"], "github") {
		extraction.NewGitHubExtractor(agent.Dependency, agent.Config).Extract(agent.DataModel)
	}
}

func (agent *Agent) CombinationAndConclusion() *model.AgentResult {

	network := mapping.Network(agent.DataModel)
	popularity := mapping.Popularity(agent.DataModel)
	activity := mapping.Activity(agent.DataModel, agent.Config.AFConfig.Activity)

	deityGiven := mapping.DeityGiven(agent.DataModel)

	coreTeam := mapping.CoreTeam(agent.DataModel)

	result := activity*0.55 + network*0.12 + popularity*0.28 + coreTeam*0.05

	if deityGiven == 1 {
		result = 1
	} else {
		result = network
	}

	return &model.AgentResult{
		Dependency: agent.Dependency,
		Result:     result,
	}
}
