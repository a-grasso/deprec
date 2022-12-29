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

	p := c.TimeframePercentile

	commits := m.Repository.Commits
	if commits != nil {
		sort.Slice(commits, func(i, j int) bool {
			return commits[i].Timestamp.Before(commits[j].Timestamp)
		})

		averageMonthsLastCommit := averageMonthsSinceLast(commits, p)

		cr.IntakeLimit(averageMonthsLastCommit, float64(c.CommitLimit), 1)
	}

	releases := m.Repository.Releases
	if releases != nil {
		sort.Slice(releases, func(i, j int) bool {
			return releases[i].Date.Before(releases[j].Date)
		})

		averageMonthsLastRelease := averageMonthsSinceLast(releases, p)

		cr.IntakeLimit(averageMonthsLastRelease, float64(c.ReleaseLimit), 1)
	}

	tags := m.Repository.Tags
	if tags != nil {
		sort.Slice(tags, func(i, j int) bool {
			return tags[i].Date.Before(tags[j].Date)
		})

		averageMonthsLastTag := averageMonthsSinceLast(tags, p)

		cr.IntakeLimit(averageMonthsLastTag, float64(c.ReleaseLimit), 1)
	}

	return cr
}

func averageMonthsSinceLast[T statistics.HasTimestamp](elements []T, p int) float64 {
	_, _, timeFrame := statistics.GetPercentilesOf(elements, p)

	monthsSince := funk.Map(timeFrame, func(t T) int {
		return statistics.CalculateTimeDifference(t.GetTimestamp(), statistics.CustomNow())
	}).([]int)

	averageMonthsSince := funk.Sum(monthsSince) / float64(len(timeFrame))

	return averageMonthsSince
}
