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

	//issues := funk.Filter(m.Repository.Issues, func(i model.Issue) bool { return !strings.Contains(i.Author, "[bot]") }).([]model.Issue)
	issues := m.Repository.Issues

	if len(issues) == 0 {
		return cr
	}

	sort.Slice(issues, func(i, j int) bool {
		return issues[i].Number < issues[j].Number
	})

	openedIssues := funk.Filter(issues, func(i model.Issue) bool { return i.State == model.IssueStateOpen }).([]model.Issue)
	closedIssues := funk.Filter(issues, func(i model.Issue) bool { return i.State == model.IssueStateClosed }).([]model.Issue)

	avgClosingTime := calcAvgClosingTime(closedIssues)

	cr.IntakeLimit(avgClosingTime, float64(c.ClosingTimeLimit), 2)

	dep := m.Repository.Name
	dep = dep + ""

	if len(closedIssues) != 0 {
		//burnUp := averageBurnUp(issues, closedIssues)
		//cr.Intake(burnUp, 1)

		burn := averageBurn(issues, closedIssues)
		cr.Intake(burn, 2)
	}

	if len(openedIssues) != 0 {
		//burnDown := averageBurnDown(openedIssues)
		//cr.Intake(burnDown, 1)
	}

	return cr
}

func averageBurnUp(issues []model.Issue, closed []model.Issue) float64 {

	sortedKeys, _ := statistics.GroupByTimestamp(issues)

	_, grouped := statistics.GroupBy(closed, func(i model.Issue) statistics.Key {
		return statistics.TimeToKey(i.ClosingTime)
	})

	mapper := func(k statistics.Key, closed []model.Issue) (statistics.Key, float64) {
		return k, float64(len(closed))
	}

	burnUp := funk.Map(grouped, mapper).(map[statistics.Key]float64)

	analysis := statistics.Analyze(sortedKeys, burnUp, 20)

	eval := evaluate(analysis)

	return eval
}

func averageBurnDown(opened []model.Issue) float64 {

	sortedKeys, grouped := statistics.GroupBy(opened, func(i model.Issue) statistics.Key {
		return statistics.TimeToKey(i.CreationTime)
	})

	statistics.FillInMissingKeys(&sortedKeys)

	mapper := func(k statistics.Key, closed []model.Issue) (statistics.Key, float64) {
		return k, float64(len(closed))
	}

	burnDown := funk.Map(grouped, mapper).(map[statistics.Key]float64)

	analysis := statistics.Analyze(sortedKeys, burnDown, 20)

	eval := evaluate(analysis)

	return eval
}

func averageBurn(issues []model.Issue, closedIssues []model.Issue) float64 {

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

	analysis := statistics.Analyze(sortedKeys, burn, 20)

	eval := evaluateBurnAnalysis(analysis)

	return eval
}

func evaluateBurnAnalysis(burn statistics.Result) float64 {

	average := math.Min(1.0, burn.Average)

	percentileAverageDiff := math.Min(1, burn.LastPercentileAverage/burn.SecondPercentileAverage)

	lpaAverageDiff := math.Min(1, burn.LastPercentileAverage/burn.Average)

	result := (average*2 + percentileAverageDiff + lpaAverageDiff) / 4
	result = (average*2 + lpaAverageDiff) / 3
	return result
}

func calcAvgClosingTime(closedIssues []model.Issue) (months float64) {

	closingTime := funk.Map(closedIssues, func(i model.Issue) (int, int) {
		difference := statistics.CalculateTimeDifference(i.CreationTime, i.ClosingTime)
		return i.Number, difference
	}).(map[int]int)

	months = funk.Sum(funk.Values(closingTime)) / float64(len(closedIssues))

	return
}
