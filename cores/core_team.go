package cores

import (
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/model"
	"github.com/a-grasso/deprec/statistics"
	"github.com/thoas/go-funk"
	"math"
	"sort"
)

func CoreTeam(m *model.DataModel, c configuration.CoreTeam) model.CoreResult {

	cr := model.NewCoreResult(model.CoreTeam)

	if m.Repository == nil {
		return *cr
	}

	contributors := m.Repository.Contributors
	commits := m.Repository.Commits

	if contributors == nil {
		return *cr
	}

	percentage := coreTeamPercentage(contributors)
	//TODO: Needs overhaul, as too punishing for big projects 50+ contributors (all those with ~2 commits)

	cr.IntakeThreshold(percentage, c.CoreTeamStrengthThresholdPercentage, 1)

	if commits == nil {
		return *cr
	}

	active := activeContributors(commits, contributors, c.ActiveContributorsPercentile)

	cr.Intake(active, 2)

	// TODO: in relation zu timeline setzen?

	return *cr
}

func activeContributors(commits []model.Commit, contributors []model.Contributor, percentile float64) float64 {
	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Timestamp.After(commits[j].Timestamp)
	})

	var lastActiveContributors []model.Contributor
	mappedContributors := funk.Map(contributors, func(c model.Contributor) (string, model.Contributor) {
		return c.Name, c
	}).(map[string]model.Contributor)

	lastCommits, _, _ := statistics.GetPercentilesOf(commits, percentile)

	for _, commit := range lastCommits {
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

	index := findBiggestJump(contributions)

	coreTeam := contributors[:index]

	totalContributors := len(contributors)
	return float64(len(coreTeam)) / float64(totalContributors) * 100
}

func findBiggestJump(contributions []int) (index int) {

	var maxAbs int
	var maxRel float64

	var indexAbs int
	var indexRel int
	for i := 1; i < len(contributions); i++ {

		curJumpRel := 1 - (float64(contributions[i]) / float64(contributions[i-1]))
		if curJumpRel > maxRel {
			indexRel = i
			maxRel = curJumpRel
		}

		curJumpAbs := contributions[i-1] - contributions[i]
		if curJumpAbs > maxAbs {
			indexAbs = i
			maxAbs = curJumpAbs
		}
	}

	return int(math.Round((float64(indexAbs) + float64(indexRel)) / 2))
}
