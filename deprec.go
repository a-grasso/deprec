package main

import (
	"deprec/agent"
	"deprec/configuration"
	"deprec/logging"
	"deprec/model"
	"errors"
	"flag"
	"fmt"
	cdx "github.com/CycloneDX/cyclonedx-go"
	"net/http"
	"os"
)

type input struct {
	sbomPath   string
	configPath string
}

func main() {
	logging.Logger.Info("DepRec run started...")

	flag.Usage = func() {
		fmt.Printf("Usage: %s <sbom> <config>\nOptions:\n none", os.Args[0])
	}

	input, err := getInput()

	if err != nil {
		exitGracefully(err)
	}

	config := configuration.Load(input.configPath)

	cdxBom := decodeSBOM(input.sbomPath)

	dependencies := parseSBOM(cdxBom)

	var result []float64
	for i, dep := range dependencies {

		if i > 10 {
			break
		}

		logging.SugaredLogger.Infof("running agent for dependency '%s:%s'", dep.Name, dep.Version)

		a := agent.NewAgent(dep, config)
		agentResult := a.Start().Result
		result = append(result, agentResult)
	}

	logging.Logger.Info("...DepRec run done")
	logging.SugaredLogger.Infof("results: %f", result)
}

func getInput() (input, error) {
	if len(os.Args) < 3 {
		return input{}, errors.New("cli argument error: SBOM file and config path arguments required")
	}

	flag.Parse()

	sbom := flag.Arg(0)
	config := flag.Arg(1)

	return input{sbom, config}, nil
}

func exitGracefully(err error) {
	logging.SugaredLogger.Fatalf("exited gracefully : %v\n", err)
}

func decodeSBOM(sbomPath string) *cdx.BOM {

	res, err := http.Get("https://github.com/DependencyTrack/dependency-track/releases/download/4.1.0/bom.json")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	//json, err := os.ReadFile(sbomPath)
	//if err != nil {
	//	log.Fatalf("Could not read sbom file '%s': %s", sbomPath, err)
	//}

	bom := new(cdx.BOM)
	decoder := cdx.NewBOMDecoder(res.Body, cdx.BOMFileFormatJSON)
	if err = decoder.Decode(bom); err != nil {
		logging.SugaredLogger.Fatalf("could not decode SBOM: %s", err)
		panic(err)
	}
	return bom
}

func parseSBOM(sbom *cdx.BOM) []*model.Dependency {
	var result []*model.Dependency

	for _, c := range *sbom.Components {
		result = append(result, &model.Dependency{
			Name:     c.Name,
			Version:  c.Version,
			MetaData: parseExternalReference(c),
		})
	}

	return result
}

func parseExternalReference(component cdx.Component) map[string]string {

	references := component.ExternalReferences

	if references == nil {
		logging.SugaredLogger.Infof("SBOM component '%s' has no external references", component.Name)
		return nil
	}

	result := make(map[string]string)

	for _, reference := range *references {
		result[string(reference.Type)] = reference.URL
	}

	return result
}

/*
func RequireValidSBOM(bom *cdx.BOM, fileFormat cdx.BOMFileFormat) {
	var inputFormat string
	switch fileFormat {
	case cdx.BOMFileFormatJSON:
		inputFormat = "json"
	case cdx.BOMFileFormatXML:
		inputFormat = "xml"
	}

	bomFile, err := os.Create(fmt.Sprintf("bom.%s", inputFormat))
	defer func() {
		if err := bomFile.Close(); err != nil && err.Error() != "file already closed" {
			fmt.Printf("failed to close bom file: %v\n", err)
		}
	}()

	encoder := cdx.NewBOMEncoder(bomFile, fileFormat)
	encoder.SetPretty(true)
	err = encoder.Encode(bom)

	valCmd := exec.Command("cyclonedx", "validate", "--input-file", bomFile.Name(), "--input-format", inputFormat, "--input-version", "v1_4", "--fail-on-errors") // #nosec G204
	valOut, err := valCmd.CombinedOutput()
	if err != nil {
		// Provide some context when test is failing
		fmt.Printf("validation error: %s\n", string(valOut))
	}
}
*/
