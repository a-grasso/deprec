package mapping

import (
	"deprec/model"
	"deprec/statistics"
	"github.com/thoas/go-funk"
	"math"
	"sort"
	"strings"
)

func Processing(m *model.DataModel) model.CoreResult {

	cr := model.CoreResult{Core: model.Processing}

	issues := funk.Filter(m.Repository.Issues, func(i model.Issue) bool { return !strings.Contains(i.Author, "[bot]") }).([]model.Issue)

	if len(issues) == 0 {
		return cr
	}

	sort.Slice(issues, func(i, j int) bool {
		return issues[i].Number < issues[j].Number
	})

	openedIssues := funk.Filter(issues, func(i model.Issue) bool { return i.State == model.IssueStateOpen }).([]model.Issue)
	closedIssues := funk.Filter(issues, func(i model.Issue) bool { return i.State == model.IssueStateClosed }).([]model.Issue)

	avgClosingTime := calcAvgClosingTime(closedIssues)

	cr.Intake(1-avgClosingTime, 2)

	if len(closedIssues) != 0 {
		cr.Intake(averageBurnUp(closedIssues), 1)

		cr.Intake(averageBurn(issues, closedIssues), 1)
	}

	if len(openedIssues) != 0 {
		cr.Intake(averageBurnDown(openedIssues), 1)
	}

	return cr
}

func evaluate2(r statistics.Result) float64 {
	return 0
}

func averageBurnUp(closed []model.Issue) float64 {

	sortedKeys, grouped := statistics.GroupBy(closed, func(i model.Issue) statistics.Key {
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

	burn := make(map[statistics.Key]float64, 0)

	for _, key := range sortedKeys {

		c := float64(len(closed[key]))
		o := math.Max(1, float64(len(opened[key])))

		burn[key] = c / o
	}

	analysis := statistics.Analyze(sortedKeys, burn, 20)

	eval := evaluate2(analysis)

	return eval
}

func calcAvgClosingTime(closedIssues []model.Issue) (months float64) {

	closingTime := funk.Map(closedIssues, func(i model.Issue) (int, int) {
		difference := statistics.CalculateTimeDifference(i.CreationTime, i.ClosingTime)
		return i.Number, difference
	}).(map[int]int)

	months = funk.Sum(funk.Values(closingTime)) / float64(len(closedIssues))

	return
}
