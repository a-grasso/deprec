package deprec_test

import (
	"context"
	"fmt"
	"github.com/a-grasso/deprec/agent"
	"github.com/a-grasso/deprec/cache"
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/logging"
	"github.com/a-grasso/deprec/model"
	"github.com/gocarina/gocsv"
	"github.com/thoas/go-funk"
	"log"
	"os"
	"sync"
	"testing"
)

var testCSV = "./agent.test.gh.csv"
var testConfig = "./config/config.json"
var testEnv = "./config/it.env"

type CSVRow struct {
	Repository     string `csv:"repository"` // .csv column headers
	SHA1           string `csv:"sha1"`       // .csv column headers
	PackageURL     string `csv:"purl"`       // .csv column headers
	Name           string `csv:"name"`
	Recommendation string `csv:"recommendation"`
	Comment        string `csv:"comment"`
	Version        string `csv:"version"`
	Ignore         bool   `csv:"ignore"`
}

var config, _ = configuration.Load(testConfig, testEnv)
var csvRows = readCsvFile(testCSV)

func readCsvFile(filePath string) []*CSVRow {
	in, err := os.Open(filePath)
	if err != nil {
		logging.SugaredLogger.Fatalf("Unable to Read Input File %s : %s", filePath, err)
	}
	defer in.Close()

	var rows []*CSVRow

	if err := gocsv.UnmarshalFile(in, &rows); err != nil {
		logging.SugaredLogger.Fatalf("Unable to parse file as CSV for %s : %s", filePath, err)
	}

	return rows
}

func TestEvaluation(t *testing.T) {

	errors := make(map[string]float64, 0)

	var confidence = 0.75

	dependencies := depsFromRows()

	agentResults := parallel(dependencies, 5, *config)

	correct, correctConfident, total := evaluation(agentResults, confidence, errors)

	sum := funk.Sum(funk.Map(funk.Values(errors), func(f float64) float64 { return f * f }))

	mse := sum / float64(len(errors))

	log.Println(fmt.Sprintf("Mean Squared Error: %f", mse))
	log.Println(fmt.Sprintf("Correct Classified (Highest Softmax Value): %2.2f %% (%d/%d))", float64(correct)/float64(total)*100, correct, total))
	log.Println(fmt.Sprintf("Confident Correct Classified (Highest Softmax Value > 0.75): %2.2f %% (%d/%d)", float64(correctConfident)/float64(total)*100, correctConfident, total))
}

func evaluation(agentResults []agent.Result, confidence float64, errors map[string]float64) (int, int, int) {
	var correct int
	var correctConfident int
	var total int

	for _, agentResult := range agentResults {

		dep := agentResult.Dependency

		rows := funk.Filter(csvRows, func(row *CSVRow) bool {

			if row.Name != dep.Name {
				return false
			}

			if row.Version != dep.Version {
				return false
			}

			if row.Repository != dep.ExternalReferences[model.VCS] {
				return false
			}

			if row.SHA1 != dep.Hashes[model.SHA1] {
				return false
			}

			if row.PackageURL != dep.PackageURL {
				return false
			}

			return true
		}).([]*CSVRow)

		row := rows[0]

		var expectedRec model.Recommendation
		switch row.Recommendation {
		case "DM":
			expectedRec = model.DecisionMaking
		case "W":
			expectedRec = model.Watchlist
		case "NIA":
			expectedRec = model.NoImmediateAction
		case "NC":
			expectedRec = model.NoConcerns
		}

		expectedRecVal := agentResult.Recommendations[expectedRec]

		if agentResult.TopRecommendation() == expectedRec {
			correct++
			if expectedRecVal > confidence {
				correctConfident++
			}
		}
		total++

		diff := 1 - expectedRecVal

		errors[row.Name] = diff
	}

	return correct, correctConfident, total
}

func depsFromRows() []model.Dependency {
	var deps []model.Dependency

	for _, row := range csvRows {

		if row.Name == "" {
			continue
		}

		if row.Ignore {
			continue
		}

		dep := model.Dependency{
			Name:               row.Name,
			Version:            row.Version,
			PackageURL:         row.PackageURL,
			Hashes:             map[model.HashAlgorithm]string{model.SHA1: row.SHA1},
			ExternalReferences: map[model.ExternalReference]string{model.VCS: row.Repository},
		}

		deps = append(deps, dep)
	}
	return deps
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

	for _, dep := range deps {

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
