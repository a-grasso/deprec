package mapping

import (
	"deprec/model"
	"deprec/statistics"
	"github.com/thoas/go-funk"
	"log"
	"math"
	"sort"
	"strings"
	"time"
)

func Processing(m *model.DataModel) float64 {

	issues := m.Repository.Issues

	nonBotIssues := funk.Filter(issues, func(i model.Issue) bool { return !strings.Contains(i.Author, "[bot]") }).([]model.Issue)

	if len(nonBotIssues) == 0 {
		return 0
	}

	sort.Slice(nonBotIssues, func(i, j int) bool {
		return nonBotIssues[i].Number < nonBotIssues[j].Number
	})

	openedIssues := funk.Filter(nonBotIssues, func(i model.Issue) bool { return i.State == model.IssueStateOpen }).([]model.Issue)
	closedIssues := funk.Filter(nonBotIssues, func(i model.Issue) bool { return i.State == model.IssueStateClosed }).([]model.Issue)

	avgClosingTime := calcAvgClosingTime(closedIssues)

	bu := calcAvgBurnUp(closedIssues)
	bd := calcAvgBurnDown(openedIssues)
	avgB := calcAvgBurn(nonBotIssues, closedIssues)

	log.Println(avgClosingTime, bu, bd, avgB)

	return 0
}

func calcAvgBurnUp(closed []model.Issue) statistics.Result {

	sortedKeys, grouped := statistics.GroupBy(closed, func(i model.Issue) time.Time {
		return i.ClosingTime
	})

	mapper := func(k statistics.Key, closed []model.Issue) (statistics.Key, float64) {
		return k, float64(len(closed))
	}

	burnUp := funk.Map(grouped, mapper).(map[statistics.Key]float64)

	return statistics.Analyze(sortedKeys, burnUp, 20)

}

func calcAvgBurnDown(opened []model.Issue) statistics.Result {

	sortedKeys, grouped := statistics.GroupBy(opened, func(i model.Issue) time.Time {
		return i.CreationTime
	})

	mapper := func(k statistics.Key, closed []model.Issue) (statistics.Key, float64) {
		return k, float64(len(closed))
	}

	burnDown := funk.Map(grouped, mapper).(map[statistics.Key]float64)

	return statistics.Analyze(sortedKeys, burnDown, 20)
}

func calcAvgBurn(issues []model.Issue, closedIssues []model.Issue) statistics.Result {

	sortedKeysOpen, opened := statistics.GroupBy(issues, func(i model.Issue) time.Time {
		return i.CreationTime
	})

	sortedKeysClosed, closed := statistics.GroupBy(closedIssues, func(i model.Issue) time.Time {
		return i.ClosingTime
	})

	sortedKeys := append(sortedKeysOpen, sortedKeysClosed...)
	statistics.SortKeys(funk.Uniq(sortedKeys).([]statistics.Key))

	burn := make(map[statistics.Key]float64, 0)

	for _, key := range sortedKeys {

		c := float64(len(closed[key]))
		o := math.Max(1, float64(len(opened[key])))

		burn[key] = c / o
	}

	return statistics.Analyze(sortedKeys, burn, 20)
}

func calcAvgClosingTime(closedIssues []model.Issue) float64 {

	closingTime := funk.Map(closedIssues, func(i model.Issue) (int, int) {
		difference := statistics.CalculateTimeDifference(i.CreationTime, i.ClosingTime)
		return i.Number, difference
	}).(map[int]int)

	return funk.Sum(funk.Values(closingTime)) / float64(len(closedIssues))
}
