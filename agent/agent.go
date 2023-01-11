package agent

import (
	"deprec/configuration"
	"deprec/cores"
	"deprec/extraction"
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

	if vcs, exists := agent.Dependency.ExternalReferences[model.VCS]; exists && strings.Contains(vcs, "github") {
		extraction.NewGitHubExtractor(agent.Dependency, agent.Config).Extract(agent.DataModel)
	}

	if agent.Dependency.PackageURL != "" {
		extraction.NewOSSIndexExtractor(agent.Dependency, agent.Config).Extract(agent.DataModel)
	}

	if sha1, exists := agent.Dependency.Hashes[model.SHA1]; exists && sha1 != "" {
		extraction.NewMavenCentralExtractor(agent.Dependency, agent.Config).Extract(agent.DataModel)
	}
}

func (agent *Agent) CombinationAndConclusion() model.CoreResult {

	cr := model.NewCoreResult(model.CombCon)

	if agent.DataModel.Repository == nil && agent.DataModel.Distribution == nil {
		return *cr
	}

	deityGiven := cores.DeityGiven(agent.DataModel)

	effort := cores.Effort(agent.DataModel, agent.Config.CoresConfig)

	interconnectedness := cores.Interconnectedness(agent.DataModel, agent.Config.CoresConfig)

	community := cores.Community(agent.DataModel, agent.Config.CoresConfig)

	support := cores.Support(agent.DataModel, agent.Config.CoresConfig)

	circumstances := cores.Circumstances(agent.DataModel, agent.Config.CoresConfig)

	cr.Overtake(deityGiven, 100)
	cr.Overtake(effort, 2)
	cr.Overtake(support, 1)
	cr.Overtake(circumstances, 0)
	cr.Overtake(community, 0)
	cr.Overtake(interconnectedness, 0)

	return *cr
}
