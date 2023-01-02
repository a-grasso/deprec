package cores

import (
	"deprec/configuration"
	"deprec/model"
	"deprec/statistics"
	"github.com/thoas/go-funk"
	"math"
	"sort"
)

func Processing(m *model.DataModel, c configuration.Processing) model.CoreResult {

	cr := model.NewCoreResult(model.Processing)

	issues := m.Repository.Issues

	issues = []model.Issue{}

	if len(issues) == 0 {
		return cr
	}

	sort.Slice(issues, func(i, j int) bool {
		return issues[i].Number < issues[j].Number
	})

	closedIssues := funk.Filter(issues, func(i model.Issue) bool { return i.State == model.IssueStateClosed }).([]model.Issue)

	closingTime := averageClosingTime(closedIssues)
	cr.IntakeLimit(closingTime, float64(c.ClosingTimeLimit), 2)

	burn := averageBurn(issues, closedIssues, c.BurnPercentile)
	cr.Intake(burn, 2)

	return cr
}

func averageBurn(issues []model.Issue, closedIssues []model.Issue, percentile float64) float64 {

	sortedKeysOpen, opened := statistics.GroupBy(issues, func(i model.Issue) statistics.Key {
		return statistics.TimeToKey(i.CreationTime)
	})

	sortedKeysClosed, closed := statistics.GroupBy(closedIssues, func(i model.Issue) statistics.Key {
		return statistics.TimeToKey(i.ClosingTime)
	})

	sortedKeys := append(sortedKeysOpen, sortedKeysClosed...)
	statistics.SortKeys(funk.Uniq(sortedKeys).([]statistics.Key))

	statistics.FillInMissingKeys(&sortedKeys)

	burn := make(map[statistics.Key]float64, 0)

	runningBalance := 0.0
	for _, key := range sortedKeys {

		c := float64(len(closed[key]))
		o := float64(len(opened[key]))

		runningBalance += o
		runningBalance -= c

		if c == 0 && runningBalance == 0 {
			continue
		}

		if o == 0 {
			o = 1
		}

		f := c / o
		burn[key] = f
	}

	analysis := statistics.Analyze(sortedKeys, burn, percentile)

	eval := evaluateBurnAnalysis(analysis)

	return eval
}

func evaluateBurnAnalysis(burn statistics.Result) float64 {

	average := math.Min(1.0, burn.Average)

	percentileAverageDiff := burn.LPAOverSPA()

	lpaAverageDiff := burn.LPAOverAVG()

	// TODO: which is better?
	result1 := (average*2 + percentileAverageDiff + lpaAverageDiff) / 4
	result2 := (average*2 + lpaAverageDiff) / 3
	return (result1 + result2) / 2
}

func averageClosingTime(closedIssues []model.Issue) (months float64) {

	if len(closedIssues) == 0 {
		return math.Inf(1)
	}

	closingTime := funk.Map(closedIssues, func(i model.Issue) (int, int) {
		difference := statistics.CalculateTimeDifference(i.CreationTime, i.ClosingTime)
		return i.Number, difference
	}).(map[int]int)

	months = funk.Sum(funk.Values(closingTime)) / float64(len(closedIssues))

	return
}
