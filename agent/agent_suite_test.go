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
	Repository     string `csv:"repository"` // .csv column headers
	SHA1           string `csv:"sha1"`       // .csv column headers
	PackageURL     string `csv:"purl"`       // .csv column headers
	Name           string `csv:"name"`
	Recommendation string `csv:"recommendation"`
	Comment        string `csv:"comment"`
	Version        string `csv:"version"`
	Ignore         bool   `csv:"ignore"`
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

var config, _ = configuration.Load(testConfig)
var csvRows = readCsvFile(testCSV)

func tableEntries() []TableEntry {

	var entries []TableEntry
	for _, row := range csvRows {
		entries = append(entries, Entry(nil, row))
	}
	return entries
}

var _ = Describe("Agent", func() {

	DescribeTable("using only all data",
		func(row *CSVRow) {

			if row.Name == "" {
				return
			}

			if row.Ignore {
				return
			}

			config = &configuration.Configuration{
				Extraction: configuration.Extraction{
					GitHub:   config.GitHub,
					OSSIndex: config.OSSIndex,
				},
				Cache: configuration.Cache{
					MongoDB: config.MongoDB,
				},
				CoresConfig: config.CoresConfig,
			}

			dep := model.Dependency{
				Name:               row.Name,
				Version:            row.Version,
				PackageURL:         row.PackageURL,
				Hashes:             map[model.HashAlgorithm]string{model.SHA1: row.SHA1},
				ExternalReferences: map[model.ExternalReference]string{model.VCS: row.Repository},
			}

			agent := agent.NewAgent(&dep, config)

			agentResult := agent.Start()

			recommendation := agentResult.TopRecommendation()

			actual := recommendation.ToAbbreviation()

			expected := row.Recommendation

			Expect(actual).To(Equal(expected), "Expected: '%s', Was: '%s' | %s", expected, actual, agentResult.Core.ToStringDeep())
		},
		func(row *CSVRow) string {
			return fmt.Sprintf("should result in '%s' when running for dependency '%s', comment: '%s'", row.Recommendation, row.Name, row.Comment)
		}, tableEntries())
})
