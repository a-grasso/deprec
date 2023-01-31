package cores

import (
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/model"
	"github.com/a-grasso/deprec/statistics"
	"github.com/thoas/go-funk"
)

func Activity(m *model.DataModel, config configuration.Activity) model.CoreResult {

	cr := model.NewCoreResult(model.Activity)

	if m.Repository == nil {
		return *cr
	}

	commits := m.Repository.Commits
	releases := m.Repository.Releases
	issues := m.Repository.Issues
	issueContributions := funk.FlatMap(issues, func(issue model.Issue) []model.IssueContribution {
		return issue.Contributions
	}).([]model.IssueContribution)

	percentile := config.Percentile
	handle(commits, 3, percentile, cr)
	handle(releases, 3, percentile, cr)
	handle(issues, 2, percentile, cr)
	handle(issueContributions, 1, percentile, cr)

	return *cr
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
