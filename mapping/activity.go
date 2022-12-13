package mapping

import (
	"deprec/configuration"
	"deprec/model"
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
	Total             int
	YearsSinceLast    int
	MonthsSinceLast   int
	LastCount         int
	AvgCount          float64
	AvgFirst20Count   float64
	AvgLast20Count    float64
	First20Count      int
	Last20Count       int
	AvgPercentage     float64
	LastPercentage    float64
	First20Percentage float64
	Last20Percentage  float64
}

func statisticAnalysisPYM[T HasTimeStamp](data []T) StatisticAnalysis {

	total := len(data)
	if total == 0 {
		return StatisticAnalysis{}
	}

	sortedKeys, groupedByYM := groupByYM[T](data)

	groupedCounts := funk.Map(groupedByYM, func(k Key, v []T) (Key, int) {
		return k, len(v)
	}).(map[Key]int)

	yearsSinceLast, monthsSinceLast := calcSinceLast(sortedKeys)

	avgCount, lastCount := calcOverall(sortedKeys, groupedCounts)

	first20Count, last20Count, avgFirst20Count, avgLast20Count := calc20Percentile(sortedKeys, groupedCounts)
	return StatisticAnalysis{
		Total:             total,
		YearsSinceLast:    yearsSinceLast,
		MonthsSinceLast:   monthsSinceLast,
		LastCount:         lastCount,
		AvgCount:          avgCount,
		AvgFirst20Count:   avgFirst20Count,
		AvgLast20Count:    avgLast20Count,
		First20Count:      first20Count,
		Last20Count:       last20Count,
		AvgPercentage:     toPercentage(avgCount, total),
		LastPercentage:    toPercentage(lastCount, total),
		First20Percentage: toPercentage(first20Count, total),
		Last20Percentage:  toPercentage(last20Count, total),
	}
}

func toPercentage[T float64 | int](count T, total int) float64 {
	return float64(count) / float64(total) * 100
}

func calcSinceLast(sortedKeys []Key) (yearsSinceLast, monthsSinceLast int) {
	lastKey := sortedKeys[len(sortedKeys)-1]
	yearsSinceLast, monthsSinceLast = calculateDifferenceOfKeyFromNow(lastKey)
	return
}

func calc20Percentile(sortedKeys []Key, groupedCounts map[Key]int) (first20Count, last20Count int, first20AvgCount, last20AvgCount float64) {

	totalMonths := len(sortedKeys)

	percentile := float64(totalMonths) / 5.0

	p20 := int(percentile)
	p80 := int(percentile * 4.0)

	first20 := sortedKeys[:p20]
	last20 := sortedKeys[p80:]

	first20Count = int(funk.Sum(funk.Map(groupedCounts, func(k Key, v int) int {
		if funk.Contains(first20, k) {
			return v
		}
		return 0
	})))

	last20Count = int(funk.Sum(funk.Map(groupedCounts, func(k Key, v int) int {
		if funk.Contains(last20, k) {
			return v
		}
		return 0
	})))

	first20AvgCount, _ = calcOverall(first20, groupedCounts)
	last20AvgCount, _ = calcOverall(last20, groupedCounts)

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

func Activity(m *model.DataModel, config configuration.AFConfig) float64 {

	commits := m.Repository.Commits
	commitAnalysis := statisticAnalysisPYM(commits)

	log.Println(commitAnalysis)
	if commitAnalysis.YearsSinceLast > config.CommitThreshold {
		return 0
	}

	issues := m.Repository.Issues
	log.Println(statisticAnalysisPYM(issues))

	issueContributions := funk.FlatMap(issues, func(issue model.Issue) []IssueContributions {
		var result []IssueContributions
		for i := 0; i < issue.Contributions; i++ {
			result = append(result, IssueContributions{time: issue.CreationTime})
		}
		return result
	}).([]IssueContributions)

	log.Println(statisticAnalysisPYM(issueContributions))

	//releases := m.Repository.Releases
	//log.Println(statisticAnalysisPYM(releases))

	return 0
}

func calculateDifferenceOfKeyFromNow(key Key) (year, month int) {
	a := customNow()

	b := key.ToTime()

	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, _ := a.Date()
	y2, M2, _ := b.Date()

	year = y2 - y1
	month = int(M2 - M1)

	if month < 0 {
		month += 12
		year--
	}

	return
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
