package main

import (
	"deprec/agent"
	"deprec/configuration"
	"deprec/model"
	"errors"
	"flag"
	"fmt"
	cdx "github.com/CycloneDX/cyclonedx-go"
	"log"
	"net/http"
	"os"
)

type input struct {
	sbomPath   string
	configPath string
}

func getInput() (input, error) {
	if len(os.Args) < 3 {
		return input{}, errors.New("SBOM FilePath and ConfigPath Arguments required")
	}

	flag.Parse()

	sbomPath := flag.Arg(0)
	configPath := flag.Arg(1)

	return input{sbomPath, configPath}, nil
}

func exitGracefully(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func decodeSBOM(sbomPath string) *cdx.BOM {

	res, err := http.Get("https://github.com/DependencyTrack/dependency-track/releases/download/4.1.0/bom.json")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	//json, err := os.ReadFile(sbomPath)
	if err != nil {
		log.Fatalf("Could not read sbom file '%s': %s", sbomPath, err)
	}

	bom := new(cdx.BOM)
	decoder := cdx.NewBOMDecoder(res.Body, cdx.BOMFileFormatJSON)
	if err = decoder.Decode(bom); err != nil {
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
			MetaData: parseExternalReference(c.ExternalReferences),
		})
	}

	return result
}

func parseExternalReference(references *[]cdx.ExternalReference) map[string]string {

	if references == nil {
		return nil
	}

	result := make(map[string]string)

	for _, reference := range *references {
		result[string(reference.Type)] = reference.URL
	}
	return result
}

func main() {
	log.Printf("DepRec run started...")

	flag.Usage = func() {
		fmt.Printf("Usage: %s <sbom> <config>\nOptions:\n none", os.Args[0])
	}

	input, err := getInput()

	if err != nil {
		exitGracefully(err)
	}

	config := configuration.Load(input.configPath)

	dxSBOM := decodeSBOM(input.sbomPath)

	customSBOM := parseSBOM(dxSBOM)

	var result []float64
	for i, dep := range customSBOM {

		if i > 20 {
			break
		}
		log.Printf("Running Agent for %s:%s", dep.Name, dep.Version)

		a := agent.NewAgent(dep, config)
		r := a.Start().Result
		result = append(result, r)
	}

	log.Printf("...DepRec run done")
	log.Printf("Results: %f", result)
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
