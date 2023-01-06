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
	t.Cleanup(CleanDatabase)

	org := ghe.extractOrganization("")

	assert.Nil(t, org)

	CheckNoDatabase(t)
}

func TestExtractRepositoryDataNil(t *testing.T) {
	t.Cleanup(CleanDatabase)

	repoData := ghe.extractRepositoryData("", "")

	assert.Nil(t, repoData)

	CheckNoDatabase(t)
}

func TestExtractReadMeNil(t *testing.T) {
	t.Cleanup(CleanDatabase)

	readme := ghe.extractReadMe("", "")

	assert.Equal(t, "", readme)

	CheckNoDatabase(t)
}

func TestExtractContributorsNil(t *testing.T) {
	t.Cleanup(CleanDatabase)

	contributors := ghe.extractContributors("", "")

	assert.Nil(t, contributors)

	CheckNoDatabase(t)
}

func TestListContributorStatsNil(t *testing.T) {
	t.Cleanup(CleanDatabase)

	contributorStats := ghe.listContributorStats("", "")

	assert.Nil(t, contributorStats)

	CheckNoDatabase(t)
}

func TestListContributorOrganizationsNil(t *testing.T) {
	t.Cleanup(CleanDatabase)

	contributorOrganizations := ghe.listContributorOrganizations("")

	assert.Nil(t, contributorOrganizations)

	CheckNoDatabase(t)
}

func TestListContributorRepositoriesNil(t *testing.T) {

	contributorRepositories := ghe.listContributorRepositories("")

	assert.NotNil(t, contributorRepositories)

	CheckNoDatabase(t)
}

func TestExtractCommitsNil(t *testing.T) {
	t.Cleanup(CleanDatabase)

	contributorStats := ghe.extractCommits("", "")

	assert.Nil(t, contributorStats)

	CheckNoDatabase(t)
}

func TestExtractTagsNil(t *testing.T) {
	t.Cleanup(CleanDatabase)

	tags := ghe.extractTags("", "")

	assert.Nil(t, tags)

	CheckNoDatabase(t)
}
