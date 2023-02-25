package agent

import (
	"fmt"
	"github.com/a-grasso/deprec/cache"
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/cores"
	"github.com/a-grasso/deprec/extraction"
	"github.com/a-grasso/deprec/model"
	"strings"
)

type Result struct {
	Dependency      model.Dependency
	Core            model.Core
	Recommendations model.RecommendationDistribution
	DataSources     []string
}

func (ar *Result) UsedCores() string {

	var usedCores []model.CoreName
	for _, cores := range ar.Core.UnderlyingCores {

		for _, core := range cores {
			coreSum := core.DecisionMaking + core.Watchlist + core.NoImmediateAction + core.NoConcerns
			if coreSum != 0 {
				usedCores = append(usedCores, core.Name)
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

	if ar.Core.IsInconclusive() {
		return model.Inconclusive
	}

	recommendations := ar.Recommendations

	var rec model.Recommendation

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
	DataModel  model.DataModel
}

func NewAgent(dependency model.Dependency, configuration configuration.Configuration) *Agent {
	agent := Agent{Dependency: dependency, DataModel: model.DataModel{}, Config: configuration}
	return &agent
}

func (agent *Agent) Run(cache *cache.Cache) Result {
	dataSources := agent.Extraction(cache)

	result := agent.CombinationAndConclusion()

	return Result{
		Dependency:      agent.Dependency,
		Core:            result,
		Recommendations: result.Recommend(),
		DataSources:     dataSources,
	}
}

func (agent *Agent) Extraction(cache *cache.Cache) []string {

	var dataSources []string

	if vcs, exists := agent.Dependency.ExternalReferences[model.VCS]; exists && strings.Contains(vcs, "github") {
		extractor, err := extraction.NewGitHubExtractor(agent.Dependency, agent.Config.GitHub, cache)
		if err == nil {
			extractor.Extract(&agent.DataModel)
			dataSources = append(dataSources, "github")
		}
	}

	purl := agent.Dependency.PackageURL

	if purl != "" {
		extractor, err := extraction.NewOSSIndexExtractor(agent.Dependency, agent.Config.OSSIndex, cache)
		if err == nil {
			extractor.Extract(&agent.DataModel)
			dataSources = append(dataSources, "ossindex")
		}
	}

	if sha1, exists := agent.Dependency.Hashes[model.SHA1]; exists && sha1 != "" && strings.Contains(purl, "maven") {
		extraction.NewMavenCentralExtractor(agent.Dependency, cache).Extract(&agent.DataModel)
		dataSources = append(dataSources, "mavencentral")
	}

	return dataSources
}

func (agent *Agent) CombinationAndConclusion() model.Core {

	cr := model.NewCore(model.CombCon)

	if agent.DataModel.Repository == nil && agent.DataModel.Distribution == nil {
		return *cr
	}

	config := agent.Config.CoresConfig

	deityGiven := cores.DeityGiven(agent.DataModel, config)

	effort := cores.Effort(agent.DataModel, config)

	interconnectedness := cores.Interconnectedness(agent.DataModel, config)

	community := cores.Community(agent.DataModel, config)

	support := cores.Support(agent.DataModel, config)

	circumstances := cores.Circumstances(agent.DataModel, config)

	cr.Overtake(deityGiven, config.CombCon.Weights.DeityGiven)
	cr.Overtake(effort, config.CombCon.Weights.Effort)
	cr.Overtake(support, config.CombCon.Weights.Support)
	cr.Overtake(circumstances, config.CombCon.Weights.Circumstances)
	cr.Overtake(community, config.CombCon.Weights.Community)
	cr.Overtake(interconnectedness, config.CombCon.Weights.Interconnectedness)

	return *cr
}
