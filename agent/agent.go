package agent

import (
	"fmt"
	"github.com/a-grasso/deprec/cache"
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/cores"
	"github.com/a-grasso/deprec/extraction"
	"github.com/a-grasso/deprec/model"
	"github.com/thoas/go-funk"
	"strings"
)

type Result struct {
	Dependency      model.Dependency
	Core            model.CoreResult
	Recommendations model.RecommendationDistribution
	DataSources     []string
}

func (ar *Result) UsedCores() string {

	var usedCores []model.Core
	for _, cores := range ar.Core.UnderlyingCores {

		for _, core := range cores {
			coreSum := core.DecisionMaking + core.Watchlist + core.NoImmediateAction + core.NoConcerns
			if coreSum != 0 {
				usedCores = append(usedCores, core.Core)
			}
		}
	}

	return fmt.Sprint(usedCores)
}

func (ar *Result) RecommendationsInsights() string {

	dm := ar.Recommendations[model.DecisionMaking]
	w := ar.Recommendations[model.Watchlist]
	nia := ar.Recommendations[model.NoImmediateAction]
	nc := ar.Recommendations[model.NoConcerns]

	return fmt.Sprintf("DM: %.3f | W: %.3f | NIA: %.3f | NC: %.3f", dm, w, nia, nc)
}

func (ar *Result) ToString() string {

	header := fmt.Sprintf("Result %s: ", ar.Dependency.Name)
	body := ar.Core.ToStringDeep()

	return header + body
}

func (ar *Result) TopRecommendation() model.Recommendation {

	recommendations := ar.Recommendations

	var rec model.Recommendation

	unique := funk.Uniq(funk.Values(recommendations)).([]float64)
	if len(unique) == 1 {
		return model.Inconclusive
	}

	tmp := -1.0
	for recommendation, f := range recommendations {
		if f > tmp {
			rec = recommendation
			tmp = f
		}
	}

	return rec
}

type Agent struct {
	Dependency model.Dependency
	Config     configuration.Configuration
	DataModel  *model.DataModel
}

func NewAgent(dependency model.Dependency, configuration configuration.Configuration) *Agent {
	agent := Agent{Dependency: dependency, DataModel: &model.DataModel{}, Config: configuration}
	return &agent
}

func (agent *Agent) Run() Result {
	dataSources := agent.Extraction()

	result := agent.CombinationAndConclusion()

	return Result{
		Dependency:      agent.Dependency,
		Core:            result,
		Recommendations: result.Softmax(),
		DataSources:     dataSources,
	}
}

func (agent *Agent) Extraction() []string {

	var dataSources []string
	cache := cache.NewCache(agent.Config.MongoDB)

	if vcs, exists := agent.Dependency.ExternalReferences[model.VCS]; exists && strings.Contains(vcs, "github") {
		extraction.NewGitHubExtractor(agent.Dependency, agent.Config.GitHub, cache).Extract(agent.DataModel)
		dataSources = append(dataSources, "github")
	}

	if agent.Dependency.PackageURL != "" {
		extraction.NewOSSIndexExtractor(agent.Dependency, agent.Config.OSSIndex, cache).Extract(agent.DataModel)
		dataSources = append(dataSources, "ossindex")
	}

	if sha1, exists := agent.Dependency.Hashes[model.SHA1]; exists && sha1 != "" {
		extraction.NewMavenCentralExtractor(agent.Dependency, cache).Extract(agent.DataModel)
		dataSources = append(dataSources, "mavencentral")
	}

	return dataSources
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
	cr.Overtake(vulnerabilities, 5)
	cr.Overtake(effort, 1)
	cr.Overtake(support, 1)
	cr.Overtake(circumstances, 1)
	cr.Overtake(community, 1)
	cr.Overtake(interconnectedness, 1)

	return *cr
}
