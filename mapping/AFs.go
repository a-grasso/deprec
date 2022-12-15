package mapping

import (
	"deprec/model"
	"github.com/thoas/go-funk"
	"math"
	"sort"
	"strings"
)

func Network(model *model.DataModel) float64 {
	var result float64

	// Repository
	var repositoryNetwork int

	var contributorRepos int
	var contributorOrgs int
	var contributorSponsors int
	var contributors int

	var orgRepos int
	var orgSponsors int
	var orgCollaborators int
	var orgFollowers int

	if model.Repository == nil {
		repositoryNetwork = 0
	} else {

		contributors += model.Repository.TotalContributors()
		for _, contributor := range model.Repository.Contributors {
			contributorRepos += contributor.Repositories
			contributorOrgs += contributor.Organizations
			contributorSponsors += 0 //TODO
		}

		if model.Repository.Org != nil {
			org := model.Repository.Org
			orgRepos = org.PublicRepos + org.OwnedPrivateRepos // TODO + org.TotalPrivateRepos ???
			orgCollaborators = org.Collaborators
			orgFollowers = org.Followers
		}
	}

	repositoryNetwork += contributorRepos
	repositoryNetwork += contributorOrgs
	repositoryNetwork += contributorSponsors
	repositoryNetwork += contributors

	repositoryNetwork += orgRepos
	repositoryNetwork += orgSponsors
	repositoryNetwork += orgCollaborators
	repositoryNetwork += orgFollowers

	result += float64(repositoryNetwork)
	return result
}

func CoreTeam(m *model.DataModel) float64 {

	percentage := coreTeamPercentage(m.Repository.Contributors)
	coreTeamStrength := math.Min(3*percentage, 1)

	activeContributors := activeContributors(m.Repository.Commits, m.Repository.Contributors)

	return coreTeamStrength*0.4 + activeContributors*0.6
}

func activeContributors(commits []model.Commit, contributors []model.Contributor) float64 {
	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Timestamp.After(commits[j].Timestamp)
	})

	percentile := float64(len(commits)) / float64(20)

	p20 := int(percentile)

	var lastActiveContributors []model.Contributor
	mappedContributors := funk.Map(contributors, func(c model.Contributor) (string, model.Contributor) {
		return c.Name, c
	}).(map[string]model.Contributor)

	i := commits[:p20]
	for _, commit := range i {
		lastActiveContributors = append(lastActiveContributors, mappedContributors[commit.Author])
	}

	lastActiveContributors = funk.UniqBy(lastActiveContributors, func(c model.Contributor) string { return c.Name }).([]model.Contributor)

	totalContributors := len(mappedContributors)
	return float64(len(lastActiveContributors)) / float64(totalContributors)
}

func coreTeamPercentage(contributors []model.Contributor) float64 {

	contributions := funk.Map(contributors, func(c model.Contributor) int {
		return c.Contributions
	}).([]int)

	sort.Slice(contributions, func(i, j int) bool {
		return i > j
	})

	index := findBiggestJump(contributions)

	coreTeam := contributors[:index+1]

	totalContributors := len(contributors)
	return float64(len(coreTeam)) / float64(totalContributors) * 100
}

func findBiggestJump(contributors []int) (index int) {
	var max int
	for i := 1; i < len(contributors); i++ {

		curJump := contributors[i-1] - contributors[i]
		if curJump > max {
			index = i
			max = curJump
		}
	}
	return
}

func Popularity(model *model.DataModel) float64 {

	stars := model.Repository.Stars
	watchers := model.Repository.Watchers
	forks := model.Repository.Forks

	// TODO users := TotalContributors - ContributorsThatContributed :-> does that work?

	return float64(stars + watchers + forks)
}

func Interconnectedness(model *model.DataModel) float64 {
	return Network(model)*0.3 + Popularity(model)*0.7
}

func DeityGiven(model *model.DataModel) float64 {

	archived := model.Repository.Archivation
	if archived {
		return 1
	}

	readme := strings.ToLower(model.Repository.ReadMe)
	if strings.Contains(readme, "deprecated") || strings.Contains(readme, "end-of-life") {
		return 1
	}

	about := strings.ToLower(model.Repository.About)
	if strings.Contains(about, "deprecated") || strings.Contains(about, "end-of-life") || strings.Contains(about, "abandoned") {
		return 1
	}

	return 0
}
