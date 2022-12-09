package agent

import (
	"deprec/configuration"
	"deprec/logging"
	"deprec/model"
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var config = configuration.Load("./../test.it.config.json")

var testCSV = "./../agent.test.gh.csv"

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

type CSVRow struct {
	Repository string  `csv:"repository"` // .csv column headers
	Name       string  `csv:"name"`
	Result     float64 `csv:"result"`
	Comment    string  `csv:"comment"`
	Version    string  `csv:"version"`
}

func TestAgentGitHubRepository(t *testing.T) {

	config := configuration.Configuration{
		GitHub:  config.GitHub,
		MongoDB: config.MongoDB,
	}

	csvRows := readCsvFile(testCSV)

	for _, row := range csvRows {

		dep := model.Dependency{
			Name:     row.Name,
			Version:  row.Version,
			MetaData: map[string]string{"vcs": row.Repository},
		}

		logging.SugaredLogger.Infof("Running Agent Test for Dependency '%s:%s' | Comment: %s", row.Name, row.Version, row.Comment)
		agent := NewAgent(&dep, &config)

		actual := agent.Start().Result

		logging.SugaredLogger.Infof("Result: '%s' :-> expected: '%f' with comment: '%s' got: '%f'", row.Repository, row.Result, row.Comment, actual)
		assert.Equal(t, row.Result, actual, fmt.Sprintf("Result: '%s' :-> expected: '%f' with comment: '%s' got: '%f'", row.Repository, row.Result, row.Comment, actual))
	}
}
