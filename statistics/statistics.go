package statistics

import (
	"fmt"
	"github.com/thoas/go-funk"
	"sort"
	"time"
)

type Key struct {
	Year  int
	Month time.Month
}

func CustomNow() time.Time {
	a := time.Now()
	a = time.Date(a.Year(), a.Month(), 1, 0, 0, 0, 0, time.UTC)
	return a
}

func (key *Key) Before(other Key) bool {
	if key.Year == other.Year {
		return key.Month < other.Month
	}

	return key.Year < other.Year
}

func (key *Key) TimeDifferenceTo(to *Key) int {

	a := CustomNow()

	if to != nil {
		a = to.ToTime()
	}

	b := key.ToTime()

	return CalculateTimeDifference(a, b)
}

func SortKeys(keys []Key) {
	sort.Slice(keys, func(i, j int) bool {
		a := keys[i]
		b := keys[j]

		return a.Before(b)
	})
}
func CalculateTimeDifference(from time.Time, to time.Time) int {
	var years, months int

	a := to
	b := from

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

func ToKey(t time.Time) Key {
	return Key{
		Year:  t.Year(),
		Month: t.Month(),
	}
}

type Result struct {
	Unit       string // Basics
	Percentile int    // Basics

	TotalMonths     int // Basics
	MonthsSinceLast int // Basics

	TotalCount *int // Count

	Last    float64 // Basics
	Average float64 // Basics

	LastCountPercentage    *float64 // Count
	AverageCountPercentage *float64 // Count

	FirstPercentileAverage float64 // Basics
	LastPercentileAverage  float64 // Basics

	FirstPercentileCount           *int     // Count
	LastPercentileCount            *int     // Count
	FirstPercentileCountPercentage *float64 // Count
	LastPercentileCountPercentage  *float64 // Count
}

func (sa *Result) String() string {
	return fmt.Sprintf("Statistic Analysis Result: %s\nPercentile: %d %%\n\nTotal: %d\nTotalMonths: %d\nMonthsSinceLast: %d\nLast: %f\nAverage: %.3f\nAvgFirstPercentileCount: %.3f\nAvgLastPercentileCount: %.3f\nFirstPercentileCount: %d\nLastPercentileCount: %d\nAvgPercentage: %.3f\nLastPercentage: %.3f\nFirstPercentilePercentage: %.3f\nLastPercentilePercentage: %.3f\n\n", sa.Unit, sa.Percentile, *sa.TotalCount, sa.TotalMonths, sa.MonthsSinceLast, sa.Last, sa.Average, sa.FirstPercentileAverage, sa.LastPercentileAverage, *sa.FirstPercentileCount, *sa.LastPercentileCount, *sa.AverageCountPercentage, *sa.LastCountPercentage, *sa.FirstPercentileCountPercentage, *sa.LastPercentileCountPercentage)
}

type HasTimestamp interface {
	GetTimeStamp() time.Time
}

func Analyze(sortedKeys []Key, grouped map[Key]float64, percentile int) *Result {
	SortKeys(sortedKeys)

	lastKey, monthsSinceLast := CalcSinceLast(sortedKeys, grouped)
	last := grouped[lastKey]

	avg := CalcOverall(sortedKeys, grouped)

	firstPercentileAverage, lastPercentileAverage := CalcPercentileAverage(percentile, sortedKeys, grouped)

	return &Result{
		Unit:                   "Per Month",
		Percentile:             percentile,
		TotalMonths:            len(sortedKeys),
		MonthsSinceLast:        monthsSinceLast,
		Last:                   last,
		Average:                avg,
		FirstPercentileAverage: firstPercentileAverage,
		LastPercentileAverage:  lastPercentileAverage,
	}
}

func AnalyzeCount[T HasTimestamp](data []T, percentile int) *Result {

	total := len(data)
	if total == 0 {
		return nil
	}

	sortedKeys, grouped := GroupByTimestamp(data)

	mapped := funk.Map(grouped, func(k Key, v []T) (Key, float64) {
		return k, float64(len(v))
	}).(map[Key]float64)

	firstPercentileCount, lastPercentileCount := CalcPercentileCount(percentile, sortedKeys, mapped)

	result := Analyze(sortedKeys, mapped, percentile)

	result.FirstPercentileCount = &firstPercentileCount
	result.LastPercentileCount = &lastPercentileCount

	result.TotalCount = &total

	avgCountP := ToPercentage(result.Average, total)
	lastCountP := ToPercentage(result.Last, total)
	firstPercentileCountP := ToPercentage(firstPercentileCount, total)
	lastPercentileCountP := ToPercentage(lastPercentileCount, total)

	result.AverageCountPercentage = &avgCountP
	result.LastCountPercentage = &lastCountP
	result.FirstPercentileCountPercentage = &firstPercentileCountP
	result.LastPercentileCountPercentage = &lastPercentileCountP

	return result
}

func ToPercentage[T float64 | int](count T, total int) float64 {
	return float64(count) / float64(total) * 100
}

func CalcSinceLast(sortedKeys []Key, grouped map[Key]float64) (lastKey Key, monthsSinceLast int) {
	SortKeys(sortedKeys)

	lastKey = sortedKeys[len(sortedKeys)-1]

	reverseKeys := sortedKeys
	funk.Reverse(reverseKeys)

	for _, key := range reverseKeys {
		if grouped[key] == 0 {
			continue
		}
		lastKey = key
	}

	monthsSinceLast = lastKey.TimeDifferenceTo(nil)
	return
}

func CalcOverall(keys []Key, groupedCounts map[Key]float64) (avg float64) {
	countPerMonth := make([]float64, 0, len(keys))
	for _, key := range keys {
		countPerMonth = append(countPerMonth, groupedCounts[key])
	}

	total := funk.Sum(countPerMonth)

	avg = total / float64(len(keys))
	return
}

func CalcPercentileAverage(p int, sortedKeys []Key, groupedCounts map[Key]float64) (firstPercentileAvgCount, lastPercentileAvgCount float64) {
	SortKeys(sortedKeys)

	totalMonths := len(sortedKeys)

	p = 100 / p
	percentile := float64(totalMonths) / float64(p)

	p20 := int(percentile)
	p80 := int(percentile * 4.0)

	firstPercentile := sortedKeys[p20 : 2*p20+1]
	lastPercentile := sortedKeys[p80:]

	firstPercentileAvgCount = CalcOverall(firstPercentile, groupedCounts)
	lastPercentileAvgCount = CalcOverall(lastPercentile, groupedCounts)

	return
}

func CalcPercentileCount(p int, sortedKeys []Key, groupedCounts map[Key]float64) (firstPercentileCount, lastPercentileCount int) {
	SortKeys(sortedKeys)

	totalMonths := len(sortedKeys)

	p = 100 / p
	percentile := float64(totalMonths) / float64(p)

	p20 := int(percentile)
	p80 := int(percentile * 4.0)

	firstPercentile := sortedKeys[p20 : 2*p20+1]
	lastPercentile := sortedKeys[p80:]

	firstPercentileCount = int(funk.Sum(funk.Map(groupedCounts, func(k Key, v float64) float64 {
		if funk.Contains(firstPercentile, k) {
			return v
		}
		return 0
	})))

	lastPercentileCount = int(funk.Sum(funk.Map(groupedCounts, func(k Key, v float64) float64 {
		if funk.Contains(lastPercentile, k) {
			return v
		}
		return 0
	})))

	return
}

func GroupBy[T any](elements []T, getTime func(T) time.Time) ([]Key, map[Key][]T) {
	grouped := make(map[Key][]T, 300)

	for _, element := range elements {

		key := ToKey(getTime(element))

		if _, exists := grouped[key]; !exists {
			grouped[key] = []T{}
		}

		grouped[key] = append(grouped[key], element)
	}

	keys := make([]Key, 0)
	for key := range grouped {
		keys = append(keys, key)
	}

	FillInMissingKeysAndSort(&keys)

	return keys, grouped
}

func GroupByTimestamp[T HasTimestamp](elements []T) ([]Key, map[Key][]T) {
	return GroupBy(elements, func(hts T) time.Time {
		return hts.GetTimeStamp()
	})
}

func FillInMissingKeysAndSort(keys *[]Key) {
	SortKeys(*keys)

	firstKey := (*keys)[0]

	first := firstKey.ToTime()
	last := CustomNow()

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

	SortKeys(*keys)

	return
}
