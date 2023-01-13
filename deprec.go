package main

import (
	"bytes"
	"deprec/agent"
	"deprec/configuration"
	"deprec/logging"
	"deprec/model"
	"errors"
	"flag"
	"fmt"
	cdx "github.com/CycloneDX/cyclonedx-go"
	"os"
	"strconv"
	"strings"
	"sync"
)

type input struct {
	sbomPath   string
	configPath string
	numWorkers int
}

func main() {
	logging.Logger.Info("DepRec run started...")

	flag.Usage = func() {
		fmt.Printf("Usage: %s <sbom> <config> <workers>\nOptions:\n none", os.Args[0])
	}

	input, err := getInput()

	if err != nil {
		exitGracefully(err)
	}

	config, err := configuration.Load(input.configPath)
	if err != nil {
		exitGracefully(err)
	}

	cdxBom, err := decodeSBOM(input.sbomPath)
	if err != nil {
		exitGracefully(err)
	}

	deps := parseSBOM(cdxBom)

	linear(config, deps)

	//parallel(deps, input, config)
}

func linear(config *configuration.Configuration, dependencies []*model.Dependency) {
	var agentResults []agent.Result
	totalDependencies := len(dependencies)
	for i, dep := range dependencies {

		if i > 50 {
			break
		}

		logging.SugaredLogger.Infof("running agent for dependency '%s:%s' %d/%d", dep.Name, dep.Version, i, totalDependencies)

		a := agent.NewAgent(dep, config)
		agentResult := a.Start()
		agentResults = append(agentResults, agentResult)
	}

	logging.Logger.Info("...DepRec run done")
	for _, ar := range agentResults {
		logging.SugaredLogger.Infof("%s --->> %s", ar.Dependency.Name, ar.TopRecommendation())
		logging.SugaredLogger.Infof("{\n%s\n}", ar.Core.ToStringDeep())
	}
}

func parallel(deps []*model.Dependency, input input, config *configuration.Configuration) {
	agentResults := make(chan *agent.Result, len(deps))
	dependencies := make(chan *model.Dependency, len(deps))

	var wg sync.WaitGroup

	for w := 0; w < input.numWorkers; w++ {
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

	for ar := range agentResults {
		logging.SugaredLogger.Infof("%s --->> %s", ar.Dependency.Name, ar.TopRecommendation())
	}
}

func worker(configuration *configuration.Configuration, dependencies <-chan *model.Dependency, results chan<- *agent.Result, worker int) {

	for dep := range dependencies {
		logging.SugaredLogger.Infof("worker %d running agent for dependency '%s:%s' %d/%d", worker, dep.Name, dep.Version, 0, 0)

		a := agent.NewAgent(dep, configuration)
		asd := a.Start()
		results <- &asd
	}
}

func getInput() (input, error) {
	if len(os.Args) < 3 {
		return input{}, errors.New("cli argument error: SBOM file and config path arguments required")
	}

	flag.Parse()

	sbom := flag.Arg(0)
	config := flag.Arg(1)
	workers, err := strconv.Atoi(flag.Arg(2))
	if err != nil {
		workers = 5
	}

	return input{sbom, config, workers}, nil
}

func exitGracefully(err error) {
	logging.SugaredLogger.Fatalf("exited gracefully : %v\n", err)
}

func decodeSBOM(sbomPath string) (*cdx.BOM, error) {

	json, err := os.ReadFile(sbomPath)
	if err != nil {
		return nil, fmt.Errorf("could not read sbom file '%s': %s", sbomPath, err)
	}
	reader := bytes.NewReader(json)

	bom := new(cdx.BOM)
	decoder := cdx.NewBOMDecoder(reader, cdx.BOMFileFormatJSON)
	if err = decoder.Decode(bom); err != nil {
		return nil, fmt.Errorf("could not decode SBOM: %s", err)
	}

	calcSBOMStats(bom)

	return bom, nil
}

func calcSBOMStats(bom *cdx.BOM) {
	noVCS := 0
	vcsGitHub := 0
	for _, component := range *bom.Components {
		if component.ExternalReferences == nil {
			noVCS += 1
			continue
		}

		externalReference := parseExternalReference(component)
		vcs, exists := externalReference["vcs"]

		if !exists {
			noVCS += 1
			continue
		}

		if strings.Contains(vcs, "github.com") {
			vcsGitHub += 1
		}
	}

	logging.SugaredLogger.Infof("%d/%d/%d github/vcs/total", vcsGitHub, len(*bom.Components)-noVCS, len(*bom.Components))
}

func parseSBOM(sbom *cdx.BOM) []*model.Dependency {
	var result []*model.Dependency

	for _, c := range *sbom.Components {
		result = append(result, &model.Dependency{
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
