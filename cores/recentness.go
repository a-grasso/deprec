package cores

import (
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/model"
	"github.com/a-grasso/deprec/statistics"
	"github.com/thoas/go-funk"
	"sort"
)

func Recentness(m model.DataModel, c configuration.Recentness) model.Core {

	cr := model.NewCore(model.Recentness)

	repositoryPart(cr, c, m.Repository)

	return *cr
}

func repositoryPart(cr *model.Core, c configuration.Recentness, repository *model.Repository) {
	if repository == nil {
		return
	}

	commits := repository.Commits
	if commits != nil {
		sort.Slice(commits, func(i, j int) bool {
			return commits[i].Timestamp.Before(commits[j].Timestamp)
		})

		lastCommit := commits[len(commits)-1]
		monthsSinceLastCommit := statistics.CalculateTimeDifference(lastCommit.Timestamp, statistics.CustomNow())

		averageMonthsSinceLastCommits := averageMonthsSinceLast(commits, c.TimeframePercentileCommits)

		cr.IntakeLimit(float64(monthsSinceLastCommit), float64(c.CommitLimit), c.Weights.MonthsSinceLastCommit)

		cr.IntakeLimit(averageMonthsSinceLastCommits, float64(c.CommitLimit), c.Weights.AverageMonthsSinceLastCommits)
	}

	releases := repository.Releases
	if releases != nil {
		sort.Slice(releases, func(i, j int) bool {
			return releases[i].Date.Before(releases[j].Date)
		})

		lastRelease := releases[len(releases)-1]
		monthsSinceLastRelease := statistics.CalculateTimeDifference(lastRelease.Date, statistics.CustomNow())

		cr.IntakeLimit(float64(monthsSinceLastRelease), float64(c.ReleaseLimit), c.Weights.MonthsSinceLastRelease)
	}
}

func averageMonthsSinceLast[T statistics.HasTimestamp](elements []T, percentile float64) float64 {
	_, _, timeFrame, _ := statistics.GetPercentilesOf(elements, percentile)

	monthsSince := funk.Map(timeFrame, func(t T) int {
		return statistics.CalculateTimeDifference(t.GetTimestamp(), statistics.CustomNow())
	}).([]int)

	averageMonthsSince := funk.Sum(monthsSince) / float64(len(timeFrame))

	return averageMonthsSince
}
