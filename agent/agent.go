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

	return model.AgentResult{
		Dependency:    agent.Dependency,
		CombConResult: result,
		Result:        result.Softmax(),
	}
}

func (agent *Agent) Extraction() {

	if strings.Contains(agent.Dependency.MetaData["vcs"], "github") {
		extraction.NewGitHubExtractor(agent.Dependency, agent.Config).Extract(agent.DataModel)
	}
}

func (agent *Agent) CombinationAndConclusion() model.CoreResult {

	cr := model.CoreResult{Core: model.CombCon}

	deityGiven := mapping.DeityGiven(agent.DataModel)

	activity := mapping.Activity(agent.DataModel, agent.Config.AFConfig.Activity)

	coreTeam := mapping.CoreTeam(agent.DataModel)

	cr.Overtake(deityGiven, 1)
	cr.Overtake(coreTeam, 1)
	cr.Overtake(activity, 2)

	return cr
}
