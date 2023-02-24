package cores

import (
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/model"
	"github.com/a-grasso/deprec/statistics"
)

func Activity(m model.DataModel, config configuration.Activity) model.Core {

	cr := model.NewCore(model.Activity)

	if m.Repository == nil {
		return *cr
	}

	commits := m.Repository.Commits
	releases := m.Repository.Releases
	issues := m.Repository.Issues

	percentile := config.Percentile
	handle(commits, config.Weights.Commits, percentile, cr)
	handle(releases, config.Weights.Releases, percentile, cr)
	handle(issues, config.Weights.Issues, percentile, cr)

	return *cr
}

func handle[T statistics.HasTimestamp](count []T, weight float64, percentile float64, cr *model.Core) {

	if len(count) == 0 {
		return
	}

	analysis := statistics.AnalyzeForActivity(count, percentile)

	percentileAverageDiff := analysis.LPAOverSPA()

	lpaAverageDiff := analysis.LPAOverAVG()

	cr.Intake(percentileAverageDiff, weight)

	cr.Intake(lpaAverageDiff, weight)
}
