package main

import (
	"deprec/model"
	"deprec/statistics"
	"errors"
	"fmt"
	cdx "github.com/CycloneDX/cyclonedx-go"
	"github.com/google/go-github/v48/github"
	"github.com/thoas/go-funk"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func getSBOMFromURL(url string) (io.ReadCloser, error) {

	res, err := http.Get("https://github.com/DependencyTrack/dependency-track/releases/download/4.1.0/bom.json")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return res.Body, nil
}

func averageBurnUp(issues []model.Issue, closed []model.Issue) float64 {

	sortedKeys, _ := statistics.GroupByTimestamp(issues)

	_, grouped := statistics.GroupBy(closed, func(i model.Issue) statistics.Key {
		return statistics.TimeToKey(i.ClosingTime)
	})

	mapper := func(k statistics.Key, closed []model.Issue) (statistics.Key, float64) {
		return k, float64(len(closed))
	}

	burnUp := funk.Map(grouped, mapper).(map[statistics.Key]float64)

	analysis := statistics.Analyze(sortedKeys, burnUp, 20)

	eval := analysis.Average

	return eval
}

func averageBurnDown(opened []model.Issue) float64 {

	sortedKeys, grouped := statistics.GroupBy(opened, func(i model.Issue) statistics.Key {
		return statistics.TimeToKey(i.CreationTime)
	})

	statistics.FillInMissingKeys(&sortedKeys)

	mapper := func(k statistics.Key, closed []model.Issue) (statistics.Key, float64) {
		return k, float64(len(closed))
	}

	burnDown := funk.Map(grouped, mapper).(map[statistics.Key]float64)

	analysis := statistics.Analyze(sortedKeys, burnDown, 20)

	eval := analysis.Average

	return eval
}

func errorTooManyContributors(err error) bool {
	if err == nil {
		return false
	}
	var e *github.ErrorResponse
	ok := errors.As(err, &e)
	if !ok {
		return false
	}
	return e.Response.StatusCode == 403 && strings.Contains(e.Message, "list is too large")
}

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
