package main

import (
	"deprec/model"
	"deprec/statistics"
	"github.com/thoas/go-funk"
)

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

	eval := analysis.Average

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

	eval := analysis.Average

	return eval
}
