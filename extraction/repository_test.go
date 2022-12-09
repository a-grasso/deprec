package extraction

import (
	"context"
	"deprec/configuration"
	"deprec/model"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

var repositoryConfig = configuration.Load("./../test.ut.config.json")

var testDependency = &model.Dependency{
	Name:     "test-dependency",
	Version:  "stable",
	MetaData: map[string]string{"vcs": "https://github.com//.git"},
}

var ghe = NewGitHubExtractor(testDependency, repositoryConfig)

func init() {
	databases, err := cache.ListDatabases(context.TODO(), bson.D{})
	if err != nil {
		return
	}

	for _, database := range databases.Databases {
		err := cache.Database(database.Name).Drop(context.TODO())
		if err != nil {
			continue
		}
	}
}

func checkNoDatabase(t *testing.T) {
	databases, err := cache.ListDatabases(context.TODO(), bson.D{})
	if err != nil {
		assert.FailNow(t, "Could not fetch databases")
	}
	assert.True(t, len(databases.Databases) == 3) // 3 databases can not be dropped (admin, local, config)
}

func TestExtractOrganizationNil(t *testing.T) {

	org := ghe.extractOrganization("")

	assert.Nil(t, org)

	checkNoDatabase(t)
}

func TestExtractRepositoryDataNil(t *testing.T) {

	ghe.Owner = ""
	ghe.Repository = ""

	repoData := ghe.extractRepositoryData()

	assert.Nil(t, repoData)

	checkNoDatabase(t)
}

func TestExtractReadMeNil(t *testing.T) {

	ghe.Owner = ""
	ghe.Repository = ""

	readme := ghe.extractReadMe()

	assert.Equal(t, "", readme)

	checkNoDatabase(t)
}

func TestExtractContributorsNil(t *testing.T) {

	ghe.Owner = ""
	ghe.Repository = ""

	contributors := ghe.extractContributors()

	assert.Nil(t, contributors)

	checkNoDatabase(t)
}

func TestListContributorStatsNil(t *testing.T) {

	ghe.Owner = ""
	ghe.Repository = ""

	contributorStats := ghe.listContributorStats()

	assert.Nil(t, contributorStats)

	checkNoDatabase(t)
}

func TestListContributorOrganizationsNil(t *testing.T) {

	contributorOrganizations := ghe.listContributorOrganizations("")

	assert.Nil(t, contributorOrganizations)

	checkNoDatabase(t)
}

func TestListContributorRepositoriesNil(t *testing.T) {

	contributorRepositories := ghe.listContributorRepositories("")

	assert.NotNil(t, contributorRepositories)

	checkNoDatabase(t)
}

func TestExtractCommitsNil(t *testing.T) {

	ghe.Owner = ""
	ghe.Repository = ""

	contributorStats := ghe.extractCommits()

	assert.Nil(t, contributorStats)

	checkNoDatabase(t)
}
