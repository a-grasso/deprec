package agent_test

import (
	"fmt"
	"github.com/a-grasso/deprec/agent"
	"github.com/a-grasso/deprec/cache"
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/logging"
	"github.com/a-grasso/deprec/model"
	"github.com/gocarina/gocsv"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var testCSV = "./../agent.test.csv"
var testConfig = "./../config/config.json"
var testEnv = "./../config/it.env"

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

var config, _ = configuration.Load(testConfig, testEnv)
var csvRows = readCsvFile(testCSV)

func tableEntries() []TableEntry {

	var entries []TableEntry
	for _, row := range csvRows {
		entries = append(entries, Entry(nil, row))
	}
	return entries
}

var mongoCache, _ = cache.NewCache(config.MongoDB)

var _ = Describe("Agent", func() {

	DescribeTable("using evaluation set",
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

			agent := agent.NewAgent(dep, *config)

			agentResult := agent.Run(mongoCache)

			recommendation := agentResult.TopRecommendation()

			actual := recommendation.ToAbbreviation()

			expected := row.Recommendation

			Expect(actual).To(Equal(expected), "Expected: '%s', Was: '%s' | %s", expected, actual, agentResult.Core.ToStringDeep())
		},
		func(row *CSVRow) string {
			return fmt.Sprintf("should result in '%s' when running for dependency '%s:%s', comment: '%s'", row.Recommendation, row.Name, row.Version, row.Comment)
		}, tableEntries())
})
