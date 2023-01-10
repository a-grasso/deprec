package cores

import (
	"deprec/configuration"
	"deprec/model"
	"deprec/statistics"
	"github.com/thoas/go-funk"
)

func Activity(m *model.DataModel, config configuration.Activity) model.CoreResult {

	cr := model.NewCoreResult(model.Activity)

	issueContributions := funk.FlatMap(m.Repository.Issues, func(issue model.Issue) []model.IssueContribution {
		return issue.Contributions
	}).([]model.IssueContribution)

	percentile := config.Percentile
	handle(m.Repository.Commits, 3, percentile, &cr)
	handle(m.Repository.Releases, 3, percentile, &cr)
	handle(m.Repository.Issues, 2, percentile, &cr)
	handle(issueContributions, 1, percentile, &cr)

	return cr
}

func handle[T statistics.HasTimestamp](count []T, weight float64, percentile float64, cr *model.CoreResult) {

	if len(count) == 0 {
		return
	}

	analysis := statistics.AnalyzeForActivity(count, percentile)

	eval := evaluateActivityAnalysis(analysis)

	cr.Intake(eval, weight)
}

func evaluateActivityAnalysis(r statistics.Result) float64 {

	percentileAverageDiff := r.LPAOverSPA()

	lpaAverageDiff := r.LPAOverAVG()

	return (percentileAverageDiff + lpaAverageDiff) / 2
}
