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

	releases := m.Repository.Releases
	//tags := m.Repository.Tags
	issues := m.Repository.Issues
	issueContributions := funk.FlatMap(issues, func(issue model.Issue) []model.IssueContribution {
		return issue.Contributions
	}).([]model.IssueContribution)

	commitAnalysis := statistics.AnalyzeCount(m.Repository.Commits, config.Percentile)
	releaseAnalysis := statistics.AnalyzeCount(releases, config.Percentile)
	//if releaseAnalysis == nil {
	//	releaseAnalysis = statistics.AnalyzeCount(tags, config.Percentile)
	//}
	issueAnalysis := statistics.AnalyzeCount(issues, config.Percentile)
	issueContributionAnalysis := statistics.AnalyzeCount(issueContributions, config.Percentile)

	evalC := evaluate(commitAnalysis, config.CommitThreshold)
	evalR := evaluate(releaseAnalysis, config.ReleaseThreshold)
	evalI := evaluate(issueAnalysis, math.MaxInt)
	evalIC := evaluate(issueContributionAnalysis, math.MaxInt)

	// result := evalC*0.325 + evalR*0.325 + evalI*0.175 + evalIC*0.175
	cr.Intake(evalC, 2)
	cr.Intake(evalR, 2)
	cr.Intake(evalI, 1)
	cr.Intake(evalIC, 1)

	return cr
}

func evaluate(ca *statistics.Result, threshold int) float64 {

	if ca == nil {
		return 0
	}

	percentileAverageDiff := math.Min(1, ca.LastPercentileAverage/ca.FirstPercentileAverage)

	lpaAverageDiff := math.Min(1, ca.LastPercentileAverage/ca.Average)

	_ = 1.0
	if ca.MonthsSinceLast > threshold {
		return 0.0
	}

	return (percentileAverageDiff + lpaAverageDiff) / 2
}
