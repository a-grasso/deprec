package cores

import (
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/model"
	"github.com/a-grasso/deprec/statistics"
	"github.com/thoas/go-funk"
	"sort"
)

func Recentness(m *model.DataModel, c configuration.Recentness) model.Core {

	cr := model.NewCoreResult(model.Recentness)

	repositoryPart(cr, c, m.Repository)

	distributionPart(cr, c, m.Distribution)

	return *cr
}

func distributionPart(cr *model.Core, c configuration.Recentness, distribution *model.Distribution) {
	if distribution == nil {
		return
	}

	artifact := distribution.Artifact
	if artifact != nil {
		date := artifact.Date

		monthsSince := statistics.CalculateTimeDifference(date, statistics.CustomNow())

		cr.IntakeLimit(float64(monthsSince), float64(c.ArtifactLimit), 1)

	}

	library := distribution.Library
	if library != nil {
		lastUpdated := library.LastUpdated

		monthsSince := statistics.CalculateTimeDifference(lastUpdated, statistics.CustomNow())

		cr.IntakeLimit(float64(monthsSince), float64(c.LibraryLimit), 1)
	}
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
		lastCommitMonthsSince := statistics.CalculateTimeDifference(lastCommit.Timestamp, statistics.CustomNow())

		averageMonthsLastCommit := averageMonthsSinceLast(commits, c.TimeframePercentileCommits)

		eval := (2*float64(lastCommitMonthsSince) + averageMonthsLastCommit) / 3

		cr.IntakeLimit(eval, float64(c.CommitLimit), 2)
	}

	releases := repository.Releases
	if releases != nil {
		sort.Slice(releases, func(i, j int) bool {
			return releases[i].Date.Before(releases[j].Date)
		})

		lastRelease := releases[len(releases)-1]
		lastReleaseMonthsSince := statistics.CalculateTimeDifference(lastRelease.Date, statistics.CustomNow())

		averageMonthsLastRelease := averageMonthsSinceLast(releases, c.TimeframePercentileReleases)

		eval := (3*float64(lastReleaseMonthsSince) + averageMonthsLastRelease) / 4

		cr.IntakeLimit(eval, float64(c.ReleaseLimit), 2)
	}
}

func averageMonthsSinceLast[T statistics.HasTimestamp](elements []T, percentile float64) float64 {
	_, _, timeFrame := statistics.GetPercentilesOf(elements, percentile)

	monthsSince := funk.Map(timeFrame, func(t T) int {
		return statistics.CalculateTimeDifference(t.GetTimestamp(), statistics.CustomNow())
	}).([]int)

	averageMonthsSince := funk.Sum(monthsSince) / float64(len(timeFrame))

	return averageMonthsSince
}
