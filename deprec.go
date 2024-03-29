package deprec

import (
	"context"
	"github.com/CycloneDX/cyclonedx-go"
	cdx "github.com/CycloneDX/cyclonedx-go"
	"github.com/a-grasso/deprec/agent"
	"github.com/a-grasso/deprec/cache"
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/logging"
	"github.com/a-grasso/deprec/model"
	"sync"
)

type Result struct {
	Results map[string]agent.Result
}

type Client struct {
	Configuration configuration.Configuration
}

func NewClient(config configuration.Configuration) *Client {
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
	logging.Logger.Info("deprec run started...")
	defer logging.Logger.Info("...deprec run done")

	dependencies := parseSBOM(sbom)

	var agentResults []agent.Result
	if runConfig.Mode == Linear {
		agentResults = linear(c.Configuration, dependencies)
	} else if runConfig.Mode == Parallel {
		agentResults = parallel(dependencies, runConfig.NumWorkers, c.Configuration)
	}

	return convertAgentResults(agentResults)
}

func convertAgentResults(agentResults []agent.Result) *Result {

	resultMap := make(map[string]agent.Result, 0)

	for _, agentResult := range agentResults {
		resultMap[agentResult.Dependency.Name] = agentResult
	}

	return &Result{Results: resultMap}
}

func linear(config configuration.Configuration, dependencies []model.Dependency) []agent.Result {
	var agentResults []agent.Result
	totalDependencies := len(dependencies)

	cache, err := cache.NewCache(config.MongoDB)
	if err == nil {
		defer cache.Client.Disconnect(context.TODO())
	}

	for i, dep := range dependencies {

		if i > 50 {
			break
		}

		logging.SugaredLogger.Infof("running agent for dependency '%s:%s' %d/%d", dep.Name, dep.Version, i, totalDependencies)

		a := agent.NewAgent(dep, config)
		agentResult := a.Run(cache)
		agentResults = append(agentResults, agentResult)
	}

	return agentResults
}

func parallel(deps []model.Dependency, numWorkers int, config configuration.Configuration) []agent.Result {
	agentResults := make(chan agent.Result, len(deps))
	dependencies := make(chan model.Dependency, len(deps))

	var wg sync.WaitGroup

	cache, err := cache.NewCache(config.MongoDB)
	if err == nil {
		defer cache.Client.Disconnect(context.TODO())
	}

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)

		w := w

		go func() {
			defer wg.Done()
			worker(config, cache, dependencies, agentResults, w)
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

	var result []agent.Result
	for ar := range agentResults {
		result = append(result, ar)
	}

	return result
}

func worker(configuration configuration.Configuration, cache *cache.Cache, dependencies <-chan model.Dependency, results chan<- agent.Result, worker int) {

	for dep := range dependencies {
		logging.SugaredLogger.Infof("worker %d running agent for dependency '%s:%s' %d/%d", worker, dep.Name, dep.Version, 0, 0)

		a := agent.NewAgent(dep, configuration)
		results <- a.Run(cache)
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
		logging.SugaredLogger.Infof("SBOM component '%s' has no hashes", component.Name)
		return nil
	}

	result := make(map[model.HashAlgorithm]string)

	for _, hash := range *hashes {
		result[model.HashAlgorithm(hash.Algorithm)] = hash.Value
	}

	return result
}
