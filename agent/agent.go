package agent

import (
	"deprec/configuration"
	"deprec/cores"
	"deprec/extraction"
	"deprec/model"
	"fmt"
	"strings"
)

type Result struct {
	Dependency      *model.Dependency
	Core            model.CoreResult
	Recommendations model.RecommendationDistribution
}

func (ar *Result) ToString() string {

	header := fmt.Sprintf("Result %s: ", ar.Dependency.Name)
	body := ar.Core.ToStringDeep()

	return header + body
}

func (ar *Result) TopRecommendation() model.Recommendation {

	result := ar.Recommendations

	var rec model.Recommendation

	tmp := -1.0
	for recommendation, f := range result {
		if f > tmp {
			rec = recommendation
			tmp = f
		}
	}

	return rec
}

type Agent struct {
	Dependency *model.Dependency
	DataModel  *model.DataModel
	Config     *configuration.Configuration
}

func NewAgent(dependency *model.Dependency, configuration *configuration.Configuration) *Agent {
	agent := Agent{Dependency: dependency, DataModel: &model.DataModel{}, Config: configuration}
	return &agent
}

func (agent *Agent) Start() Result {
	agent.Extraction()

	result := agent.CombinationAndConclusion()

	return Result{
		Dependency:      agent.Dependency,
		Core:            result,
		Recommendations: result.Softmax(),
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
	vulnerabilities := cores.Vulnerabilities(agent.DataModel)

	effort := cores.Effort(agent.DataModel, agent.Config.CoresConfig)

	interconnectedness := cores.Interconnectedness(agent.DataModel, agent.Config.CoresConfig)

	community := cores.Community(agent.DataModel, agent.Config.CoresConfig)

	support := cores.Support(agent.DataModel, agent.Config.CoresConfig)

	circumstances := cores.Circumstances(agent.DataModel)

	cr.Overtake(deityGiven, 100)
	cr.Overtake(vulnerabilities, 25)
	cr.Overtake(effort, 1)
	cr.Overtake(support, 1)
	cr.Overtake(circumstances, 1)
	cr.Overtake(community, 1)
	cr.Overtake(interconnectedness, 1)

	return *cr
}
