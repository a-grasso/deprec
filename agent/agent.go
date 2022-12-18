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

	cr := model.CoreResult{Core: model.CombCon}

	activity := mapping.Activity(agent.DataModel, agent.Config.AFConfig.Activity)

	deityGiven := mapping.DeityGiven(agent.DataModel)

	cr.Overtake(deityGiven, 1)
	cr.Overtake(activity, 2)

	return &model.AgentResult{
		Dependency:    agent.Dependency,
		CombConResult: cr,
		Result:        cr.Softmax(),
	}
}
