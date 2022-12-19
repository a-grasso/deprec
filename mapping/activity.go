package mapping

import (
	"deprec/configuration"
	"deprec/model"
	"deprec/statistics"
	"github.com/thoas/go-funk"
	"math"
)

func Activity(m *model.DataModel, config configuration.Activity) model.CoreResult {

	cr := model.CoreResult{Core: model.Activity}

	issueContributions := funk.FlatMap(m.Repository.Issues, func(issue model.Issue) []model.IssueContribution {
		return issue.Contributions
	}).([]model.IssueContribution)

	p := config.Percentile
	handle(m.Repository.Commits, 2, p, &cr)
	handle(m.Repository.Releases, 2, p, &cr)
	handle(m.Repository.Tags, 2, p, &cr)
	handle(m.Repository.Issues, 1, p, &cr)
	handle(issueContributions, 1, p, &cr)

	return cr
}

func handle[T statistics.HasTimestamp](count []T, weight float64, percentile int, cr *model.CoreResult) {

	if count == nil {
		return
	}

	analysis := statistics.AnalyzeCount(count, percentile)

	eval := evaluate(analysis)

	cr.Intake(eval, weight)
}

func evaluate(r statistics.Result) float64 {

	percentileAverageDiff := math.Min(1, r.LastPercentileAverage/r.FirstPercentileAverage)

	lpaAverageDiff := math.Min(1, r.LastPercentileAverage/r.Average)

	return (percentileAverageDiff + lpaAverageDiff) / 2
}
