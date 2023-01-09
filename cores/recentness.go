package cores

import (
	"deprec/configuration"
	"deprec/model"
	"deprec/statistics"
	"github.com/thoas/go-funk"
	"sort"
)

func Recentness(m *model.DataModel, c configuration.Recentness) model.CoreResult {

	cr := model.NewCoreResult(model.Recentness)

	commits := m.Repository.Commits
	if commits != nil {
		sort.Slice(commits, func(i, j int) bool {
			return commits[i].Timestamp.Before(commits[j].Timestamp)
		})

		lastCommit := commits[len(commits)-1]
		lastCommitMonthsSince := statistics.CalculateTimeDifference(lastCommit.Timestamp, statistics.CustomNow())

		averageMonthsLastCommit := averageMonthsSinceLast(commits, c.TimeframePercentileCommits)

		eval := (2*float64(lastCommitMonthsSince) + averageMonthsLastCommit) / 3

		cr.IntakeLimit(eval, float64(c.CommitLimit), 1)
	}

	releases := m.Repository.Releases
	if releases != nil {
		sort.Slice(releases, func(i, j int) bool {
			return releases[i].Date.Before(releases[j].Date)
		})

		lastRelease := releases[len(releases)-1]
		lastReleaseMonthsSince := statistics.CalculateTimeDifference(lastRelease.Date, statistics.CustomNow())

		averageMonthsLastRelease := averageMonthsSinceLast(releases, c.TimeframePercentileReleases)

		eval := (3*float64(lastReleaseMonthsSince) + averageMonthsLastRelease) / 4

		cr.IntakeLimit(eval, float64(c.ReleaseLimit), 1)
	}

	tags := m.Repository.Tags
	if tags != nil {
		sort.Slice(tags, func(i, j int) bool {
			return tags[i].Date.Before(tags[j].Date)
		})

		averageMonthsLastTag := averageMonthsSinceLast(tags, c.TimeframePercentileReleases)

		cr.IntakeLimit(averageMonthsLastTag, float64(c.ReleaseLimit), 1)
	}

	return cr
}

func averageMonthsSinceLast[T statistics.HasTimestamp](elements []T, percentile float64) float64 {
	_, _, timeFrame := statistics.GetPercentilesOf(elements, percentile)

	monthsSince := funk.Map(timeFrame, func(t T) int {
		return statistics.CalculateTimeDifference(t.GetTimestamp(), statistics.CustomNow())
	}).([]int)

	averageMonthsSince := funk.Sum(monthsSince) / float64(len(timeFrame))

	return averageMonthsSince
}
