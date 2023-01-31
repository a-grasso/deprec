package deprec

import (
	"github.com/CycloneDX/cyclonedx-go"
	cdx "github.com/CycloneDX/cyclonedx-go"
	"github.com/a-grasso/deprec/agent"
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/logging"
	"github.com/a-grasso/deprec/model"
	"sync"
)

type Result struct {
	Results map[string]model.CoreResult
}

type Client struct {
	Configuration *configuration.Configuration
}

func NewClient(config *configuration.Configuration) *Client {
	return &Client{
		Configuration: config,
	}
}

type RunConfig struct {
	Mode       RunMode
	NumWorkers int
}

type RunMode string

const (
	Linear   RunMode = "linear"
	Parallel RunMode = "parallel"
)

func (c *Client) Run(sbom *cyclonedx.BOM, runConfig RunConfig) *Result {
	logging.Logger.Info("DepRec run started...")

	dependencies := parseSBOM(sbom)

	var agentResults []agent.Result
	if runConfig.Mode == Linear {
		agentResults = linear(*c.Configuration, dependencies)
	} else if runConfig.Mode == Parallel {
		agentResults = parallel(dependencies, runConfig.NumWorkers, *c.Configuration)
	}

	return convertAgentResults(agentResults)

}

func convertAgentResults(agentResults []agent.Result) *Result {

	resultMap := make(map[string]model.CoreResult, 0)

	for _, agentResult := range agentResults {
		resultMap[agentResult.Dependency.Name] = agentResult.Core
	}

	return &Result{Results: resultMap}
}

func linear(config configuration.Configuration, dependencies []model.Dependency) []agent.Result {
	var agentResults []agent.Result
	totalDependencies := len(dependencies)
	for i, dep := range dependencies {

		if i > 50 {
			break
		}

		logging.SugaredLogger.Infof("running agent for dependency '%s:%s' %d/%d", dep.Name, dep.Version, i, totalDependencies)

		a := agent.NewAgent(dep, config)
		agentResult := a.Run()
		agentResults = append(agentResults, agentResult)
	}

	logging.Logger.Info("...DepRec run done")
	for _, ar := range agentResults {
		logging.SugaredLogger.Infof("%s --->> %s", ar.Dependency.Name, ar.TopRecommendation())
		logging.SugaredLogger.Infof("{\n%s\n}", ar.Core.ToStringDeep())
	}

	return agentResults
}

func parallel(deps []model.Dependency, numWorkers int, config configuration.Configuration) []agent.Result {
	agentResults := make(chan agent.Result, len(deps))
	dependencies := make(chan model.Dependency, len(deps))

	var wg sync.WaitGroup

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)

		w := w

		go func() {
			defer wg.Done()
			worker(config, dependencies, agentResults, w)
		}()
	}

	for i, dep := range deps {

		if i > 50 {
			break
		}

		dependencies <- dep
	}
	close(dependencies)

	wg.Wait()

	close(agentResults)

	logging.Logger.Info("...DepRec run done")

	var result []agent.Result
	for ar := range agentResults {
		result = append(result, ar)
		logging.SugaredLogger.Infof("%s --->> %s", ar.Dependency.Name, ar.TopRecommendation())
	}

	return result
}

func worker(configuration configuration.Configuration, dependencies <-chan model.Dependency, results chan<- agent.Result, worker int) {

	for dep := range dependencies {
		logging.SugaredLogger.Infof("worker %d running agent for dependency '%s:%s' %d/%d", worker, dep.Name, dep.Version, 0, 0)

		a := agent.NewAgent(dep, configuration)
		results <- a.Run()
	}
}

func parseSBOM(sbom *cdx.BOM) []model.Dependency {
	var result []model.Dependency

	for _, c := range *sbom.Components {
		result = append(result, model.Dependency{
			Name:               c.Name,
			Version:            c.Version,
			PackageURL:         c.PackageURL,
			Hashes:             parseHashes(c),
			ExternalReferences: parseExternalReference(c),
		})
	}

	return result
}

func parseExternalReference(component cdx.Component) map[model.ExternalReference]string {

	references := component.ExternalReferences

	if references == nil {
		logging.SugaredLogger.Infof("SBOM component '%s' has no external references", component.Name)
		return nil
	}

	result := make(map[model.ExternalReference]string)

	for _, reference := range *references {
		result[model.ExternalReference(reference.Type)] = reference.URL
	}

	return result
}

func parseHashes(component cdx.Component) map[model.HashAlgorithm]string {

	hashes := component.Hashes

	if hashes == nil {
		logging.SugaredLogger.Infof("SBOM component '%s' has no external references", component.Name)
		return nil
	}

	result := make(map[model.HashAlgorithm]string)

	for _, hash := range *hashes {
		result[model.HashAlgorithm(hash.Algorithm)] = hash.Value
	}

	return result
}
