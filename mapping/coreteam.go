package mapping

import (
	"deprec/model"
	"github.com/thoas/go-funk"
	"math"
	"sort"
)

func CoreTeam(m *model.DataModel) float64 {

	percentage := coreTeamPercentage(m.Repository.Contributors)
	coreTeamStrength := math.Min(3*percentage, 100) / 100

	activeContributors := activeContributors(m.Repository.Commits, m.Repository.Contributors)

	// TODO: in relation zu timeline setzen?

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
