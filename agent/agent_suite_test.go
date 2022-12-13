package agent_test

import (
	"deprec/agent"
	"deprec/configuration"
	"deprec/logging"
	"deprec/model"
	"fmt"
	"github.com/gocarina/gocsv"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var testCSV = "./../agent.test.gh.csv"
var testConfig = "./../test.it.config.json"

func TestAgent(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Agent Suite")
}

type CSVRow struct {
	Repository string  `csv:"repository"` // .csv column headers
	Name       string  `csv:"name"`
	Result     float64 `csv:"result"`
	Comment    string  `csv:"comment"`
	Version    string  `csv:"version"`
	Ignore     bool    `csv:"ignore"`
}

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

var config = configuration.Load(testConfig)
var csvRows = readCsvFile(testCSV)

func tableEntries() []TableEntry {

	var entries []TableEntry
	for _, row := range csvRows {
		entries = append(entries, Entry(nil, row))
	}
	return entries
}

var _ = Describe("Agent", func() {

	DescribeTable("using only repository data",
		func(row *CSVRow) {

			if row.Ignore {
				return
			}

			config = &configuration.Configuration{
				GitHub:  config.GitHub,
				MongoDB: config.MongoDB,
			}

			dep := model.Dependency{
				Name:     row.Name,
				Version:  row.Version,
				MetaData: map[string]string{"vcs": row.Repository},
			}

			agent := agent.NewAgent(&dep, config)

			actual := agent.Start().Result
			expected := row.Result

			Expect(actual).To(Equal(expected))
		},
		func(row *CSVRow) string {
			return fmt.Sprintf("should result in '%.2f' when running for dependency '%s', comment: '%s'", row.Result, row.Name, row.Comment)
		}, tableEntries())
})
