package mapping

import (
	"deprec/configuration"
	"deprec/model"
	"fmt"
	"github.com/thoas/go-funk"
	"log"
	"sort"
	"time"
)

func customNow() time.Time {
	a := time.Now()
	a = time.Date(a.Year(), a.Month(), 1, 0, 0, 0, 0, time.UTC)
	return a
}

type Key struct {
	Year  int
	Month time.Month
}

func toKey(t time.Time) Key {
	return Key{
		Year:  t.Year(),
		Month: t.Month(),
	}
}

func (key *Key) Before(other Key) bool {
	if key.Year == other.Year {
		return key.Month < other.Month
	}

	return key.Year < other.Year
}

func (key *Key) CalculateTimeDifference(to *Key) int {
	var years, months int

	a := customNow()

	if to != nil {
		a = to.ToTime()
	}

	b := key.ToTime()

	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, _ := a.Date()
	y2, M2, _ := b.Date()

	years = y2 - y1
	months = int(M2 - M1)

	if months < 0 {
		months += 12
		years--
	}

	return (years * 12) + months
}

func (key *Key) ToTime() time.Time {
	return time.Date(key.Year, key.Month, 1, 0, 0, 0, 0, time.UTC)
}

func sortKeys(keys []Key) {
	sort.Slice(keys, func(i, j int) bool {
		a := keys[i]
		b := keys[j]

		return a.Before(b)
	})
}

type HasTimeStamp interface {
	GetTimeStamp() time.Time
}

type StatisticAnalysis struct {
	Unit                      string
	Percentile                int
	TotalCount                int
	TotalMonths               int
	MonthsSinceLast           int
	LastCount                 int
	AvgCount                  float64
	AvgFirstPercentileCount   float64
	AvgLastPercentileCount    float64
	FirstPercentileCount      int
	LastPercentileCount       int
	AvgPercentage             float64
	LastPercentage            float64
	FirstPercentilePercentage float64
	LastPercentilePercentage  float64
}

func (sa StatisticAnalysis) String() string {
	return fmt.Sprintf("StatisticAnalysis: %s\nPercentile: %d %%\n\nTotalCount: %d\nTotalMonths: %d\nMonthsSinceLast: %d\nLastCount: %d\nAvgCount: %.3f\nAvgFirstPercentileCount: %.3f\nAvgLastPercentileCount: %.3f\nFirstPercentileCount: %d\nLastPercentileCount: %d\nAvgPercentage: %.3f\nLastPercentage: %.3f\nFirstPercentilePercentage: %.3f\nLastPercentilePercentage: %.3f\n\n", sa.Unit, sa.Percentile, sa.TotalCount, sa.TotalMonths, sa.MonthsSinceLast, sa.LastCount, sa.AvgCount, sa.AvgFirstPercentileCount, sa.AvgLastPercentileCount, sa.FirstPercentileCount, sa.LastPercentileCount, sa.AvgPercentage, sa.LastPercentage, sa.FirstPercentilePercentage, sa.LastPercentilePercentage)
}

func statisticAnalysis[T HasTimeStamp](data []T, percentile int) StatisticAnalysis {

	total := len(data)
	if total == 0 {
		return StatisticAnalysis{}
	}

	sortedKeys, groupedByYM := groupByYM[T](data)

	groupedCounts := funk.Map(groupedByYM, func(k Key, v []T) (Key, int) {
		return k, len(v)
	}).(map[Key]int)

	monthsSinceLast := calcSinceLast(sortedKeys)

	avgCount, lastCount := calcOverall(sortedKeys, groupedCounts)

	firstPercentileCount, lastPercentileCount, avgFirstPercentileCount, avgLastPercentileCount := calcPercentile(percentile, sortedKeys, groupedCounts)
	return StatisticAnalysis{
		Unit:                      "Per Month",
		Percentile:                percentile,
		TotalCount:                total,
		TotalMonths:               len(sortedKeys),
		MonthsSinceLast:           monthsSinceLast,
		LastCount:                 lastCount,
		AvgCount:                  avgCount,
		AvgFirstPercentileCount:   avgFirstPercentileCount,
		AvgLastPercentileCount:    avgLastPercentileCount,
		FirstPercentileCount:      firstPercentileCount,
		LastPercentileCount:       lastPercentileCount,
		AvgPercentage:             toPercentage(avgCount, total),
		LastPercentage:            toPercentage(lastCount, total),
		FirstPercentilePercentage: toPercentage(firstPercentileCount, total),
		LastPercentilePercentage:  toPercentage(lastPercentileCount, total),
	}
}

func toPercentage[T float64 | int](count T, total int) float64 {
	return float64(count) / float64(total) * 100
}

func calcSinceLast(sortedKeys []Key) (monthsSinceLast int) {
	lastKey := sortedKeys[len(sortedKeys)-1]
	monthsSinceLast = lastKey.CalculateTimeDifference(nil)
	return
}

func calcPercentile(p int, sortedKeys []Key, groupedCounts map[Key]int) (firstPercentileCount, lastPercentileCount int, firstPercentileAvgCount, lastPercentileAvgCount float64) {

	totalMonths := len(sortedKeys)

	p = 100 / p
	percentile := float64(totalMonths) / float64(p)

	p20 := int(percentile)
	p80 := int(percentile * 4.0)

	firstPercentile := sortedKeys[:p20]
	lastPercentile := sortedKeys[p80:]

	firstPercentileCount = int(funk.Sum(funk.Map(groupedCounts, func(k Key, v int) int {
		if funk.Contains(firstPercentile, k) {
			return v
		}
		return 0
	})))

	lastPercentileCount = int(funk.Sum(funk.Map(groupedCounts, func(k Key, v int) int {
		if funk.Contains(lastPercentile, k) {
			return v
		}
		return 0
	})))

	firstPercentileAvgCount, _ = calcOverall(firstPercentile, groupedCounts)
	lastPercentileAvgCount, _ = calcOverall(lastPercentile, groupedCounts)

	return
}

func calcOverall(sortedKeys []Key, groupedCounts map[Key]int) (avgCount float64, lastCount int) {

	countPerMonth := make([]int, 0, len(sortedKeys))
	for _, key := range sortedKeys {
		countPerMonth = append(countPerMonth, groupedCounts[key])
	}

	total := funk.Sum(countPerMonth)

	avgCount = total / float64(len(sortedKeys))
	lastCount = groupedCounts[sortedKeys[len(sortedKeys)-1]]

	return
}

type IssueContributions struct {
	time time.Time
}

func (ic IssueContributions) GetTimeStamp() time.Time {
	return ic.time
}

func Activity(m *model.DataModel, config configuration.Activity) float64 {

	commits := m.Repository.Commits
	commitAnalysis := statisticAnalysis(commits, config.Percentile)

	log.Println(commitAnalysis.String())
	if commitAnalysis.MonthsSinceLast > config.CommitThreshold {
		return 0
	}

	issues := m.Repository.Issues
	log.Println(statisticAnalysis(issues, config.Percentile).String())

	issueContributions := funk.FlatMap(issues, func(issue model.Issue) []IssueContributions {
		var result []IssueContributions
		for i := 0; i < issue.Contributions; i++ {
			result = append(result, IssueContributions{time: issue.CreationTime})
		}
		return result
	}).([]IssueContributions)

	log.Println(statisticAnalysis(issueContributions, config.Percentile).String())

	//releases := m.Repository.Releases
	//log.Println(statisticAnalysis(releases))

	return 0
}

func groupByYM[T HasTimeStamp](elements []T) ([]Key, map[Key][]T) {
	grouped := make(map[Key][]T)

	for _, element := range elements {

		key := toKey(element.GetTimeStamp())

		if _, exists := grouped[key]; !exists {
			grouped[key] = []T{}
		}

		grouped[key] = append(grouped[key], element)
	}

	sortedKeys := make([]Key, 0)
	for key, _ := range grouped {
		sortedKeys = append(sortedKeys, key)
	}

	fillInMissingMonths(&sortedKeys)

	sortKeys(sortedKeys)

	return sortedKeys, grouped
}

func fillInMissingMonths(keys *[]Key) {

	sortKeys(*keys)

	firstKey := (*keys)[0]
	lastKey := (*keys)[len(*keys)-1]

	first := firstKey.ToTime()
	last := lastKey.ToTime()

	tmp := first
	for {

		missingKey := Key{
			Year:  tmp.Year(),
			Month: tmp.Month(),
		}
		if !funk.Contains(*keys, missingKey) {
			*keys = append(*keys, missingKey)
		}
		tmp = tmp.AddDate(0, 1, 0)
		if tmp.After(last) {
			break
		}
	}

	return
}
