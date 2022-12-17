package extraction

import (
	"deprec/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testDependency = &model.Dependency{
	Name:     "test-dependency",
	Version:  "stable",
	MetaData: map[string]string{"vcs": "https://github.com//.git"},
}

var ghe = NewGitHubExtractor(testDependency, config)

func TestExtractOrganizationNil(t *testing.T) {
	t.Cleanup(cleanDatabase)

	org := ghe.extractOrganization("")

	assert.Nil(t, org)

	checkNoDatabase(t)
}

func TestExtractRepositoryDataNil(t *testing.T) {
	t.Cleanup(cleanDatabase)

	repoData := ghe.extractRepositoryData("", "")

	assert.Nil(t, repoData)

	checkNoDatabase(t)
}

func TestExtractReadMeNil(t *testing.T) {
	t.Cleanup(cleanDatabase)

	readme := ghe.extractReadMe("", "")

	assert.Equal(t, "", readme)

	checkNoDatabase(t)
}

func TestExtractContributorsNil(t *testing.T) {
	t.Cleanup(cleanDatabase)

	contributors := ghe.extractContributors("", "")

	assert.Nil(t, contributors)

	checkNoDatabase(t)
}

func TestListContributorStatsNil(t *testing.T) {
	t.Cleanup(cleanDatabase)

	contributorStats := ghe.listContributorStats("", "")

	assert.Nil(t, contributorStats)

	checkNoDatabase(t)
}

func TestListContributorOrganizationsNil(t *testing.T) {
	t.Cleanup(cleanDatabase)

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
	t.Cleanup(cleanDatabase)

	contributorStats := ghe.extractCommits("", "")

	assert.Nil(t, contributorStats)

	checkNoDatabase(t)
}

func TestExtractTagsNil(t *testing.T) {
	t.Cleanup(cleanDatabase)

	tags := ghe.extractTags("", "")

	assert.Nil(t, tags)

	checkNoDatabase(t)
}
