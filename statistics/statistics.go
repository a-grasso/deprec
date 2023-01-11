package statistics

import (
	"fmt"
	"github.com/thoas/go-funk"
	"math"
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

func CustomYear(year int) time.Time {
	a := time.Now()
	a = time.Date(year, 12, 1, 0, 0, 0, 0, time.UTC)
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

func TimeToKey(t time.Time) Key {
	return Key{
		Year:  t.Year(),
		Month: t.Month(),
	}
}

type Result struct {
	Unit       string  // Basics
	Percentile float64 // Basics

	TotalMonths     int // Basics
	MonthsSinceLast int // Basics

	TotalCount *int // Count

	Last    float64 // Basics
	Average float64 // Basics

	LastCountPercentage    *float64 // Count
	AverageCountPercentage *float64 // Count

	SecondPercentileAverage float64 // Basics
	LastPercentileAverage   float64 // Basics

	SecondPercentileCount           *int     // Count
	LastPercentileCount             *int     // Count
	SecondPercentileCountPercentage *float64 // Count
	LastPercentileCountPercentage   *float64 // Count
}

func (sa *Result) Ratio(numerator, denominator float64) float64 {
	if denominator == 0 {
		denominator = 1
	}

	f := numerator / denominator

	return math.Min(1, f)
}

func (sa *Result) LPAOverAVG() float64 {

	lpa := sa.LastPercentileAverage
	avg := sa.Average

	return sa.Ratio(lpa, avg)
}

func (sa *Result) LPAOverSPA() float64 {

	lpa := sa.LastPercentileAverage
	spa := sa.SecondPercentileAverage

	return sa.Ratio(lpa, spa)
}

func (sa *Result) String() string {
	return fmt.Sprintf("Statistic Analysis Result: %s\nPercentile: %f %%\n\nTotal: %d\nTotalMonths: %d\nMonthsSinceLast: %d\nLast: %f\nAverage: %.3f\nAvgSecondPercentileCount: %.3f\nAvgLastPercentileCount: %.3f\nSecondPercentileCount: %d\nLastPercentileCount: %d\nAvgPercentage: %.3f\nLastPercentage: %.3f\nSecondPercentilePercentage: %.3f\nLastPercentilePercentage: %.3f\n\n", sa.Unit, sa.Percentile, *sa.TotalCount, sa.TotalMonths, sa.MonthsSinceLast, sa.Last, sa.Average, sa.SecondPercentileAverage, sa.LastPercentileAverage, *sa.SecondPercentileCount, *sa.LastPercentileCount, *sa.AverageCountPercentage, *sa.LastCountPercentage, *sa.SecondPercentileCountPercentage, *sa.LastPercentileCountPercentage)
}

type HasTimestamp interface {
	GetTimestamp() time.Time
}

func Analyze(sortedKeys []Key, grouped map[Key]float64, percentile float64) Result {

	lastKey, monthsSinceLast := CalcSinceLast(sortedKeys, grouped)
	last := grouped[lastKey]

	avg := CalcOver(sortedKeys, grouped)

	secondPercentileAverage, lastPercentileAverage := CalcPercentileAverage(percentile, sortedKeys, grouped)

	return Result{
		Unit:                    "Per Month",
		Percentile:              percentile,
		TotalMonths:             len(sortedKeys),
		MonthsSinceLast:         monthsSinceLast,
		Last:                    last,
		Average:                 avg,
		SecondPercentileAverage: secondPercentileAverage,
		LastPercentileAverage:   lastPercentileAverage,
	}
}

func AnalyzeForActivity[T HasTimestamp](data []T, percentile float64) Result {

	total := len(data)

	sortedKeys, grouped := GroupByTimestamp(data)

	FillInMissingKeys(&sortedKeys)

	mapped := funk.Map(grouped, func(k Key, v []T) (Key, float64) {
		return k, float64(len(v))
	}).(map[Key]float64)

	secondPercentileCount, lastPercentileCount := CalcPercentileCount(percentile, sortedKeys, mapped)

	result := Analyze(sortedKeys, mapped, percentile)

	result.SecondPercentileCount = &secondPercentileCount
	result.LastPercentileCount = &lastPercentileCount

	result.TotalCount = &total

	avgCountP := ToPercentage(result.Average, total)
	lastCountP := ToPercentage(result.Last, total)
	secondPercentileCountP := ToPercentage(secondPercentileCount, total)
	lastPercentileCountP := ToPercentage(lastPercentileCount, total)

	result.AverageCountPercentage = &avgCountP
	result.LastCountPercentage = &lastCountP
	result.SecondPercentileCountPercentage = &secondPercentileCountP
	result.LastPercentileCountPercentage = &lastPercentileCountP

	return result
}

func ToPercentage[T float64 | int](value T, total int) float64 {
	return float64(value) / float64(total) * 100
}

func CalcSinceLast(sortedKeys []Key, grouped map[Key]float64) (lastKey Key, monthsSinceLast int) {

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

func CalcOver(keys []Key, grouped map[Key]float64) (avg float64) {

	perMonth := make([]float64, 0, len(keys))
	for _, key := range keys {
		perMonth = append(perMonth, grouped[key])
	}

	total := funk.Sum(perMonth)

	avg = total / float64(len(keys))
	return
}

func CalcPercentileAverage(p float64, sortedKeys []Key, grouped map[Key]float64) (secondPercentileAvg, lastPercentileAvg float64) {

	_, secondPercentile, lastPercentile := GetPercentilesOf(sortedKeys, p)

	secondPercentileAvg = CalcOver(secondPercentile, grouped)
	lastPercentileAvg = CalcOver(lastPercentile, grouped)

	return
}

func CalcPercentileCount(p float64, sortedKeys []Key, groupedCounts map[Key]float64) (secondPercentileCount, lastPercentileCount int) {

	_, secondPercentile, lastPercentile := GetPercentilesOf(sortedKeys, p)

	secondPercentileCount = int(funk.Sum(funk.Map(groupedCounts, func(k Key, v float64) float64 {
		if funk.Contains(secondPercentile, k) {
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

func GetPercentilesOf[T any](elements []T, p float64) (first, second, last []T) {

	// e.g. total = 1000
	total := len(elements)

	// e.g. p = 12,5
	slices := math.Round(100 / p)       // slices = 8
	perSlice := float64(total) / slices // perSlice = 125

	elementsPerSlice := make([][]T, int(slices))
	scope := math.Max(1, perSlice) // = 125

	for i := 0; i < int(slices); i++ {

		if i >= total {
			break
		}

		j := int(scope * float64(i))
		k := int(scope * float64(i+1))
		elementsPerSlice[i] = elements[j:k]
	}

	sizeDiff := int(slices) - total
	offset := int(math.Max(0, float64(sizeDiff)))

	first = elementsPerSlice[0]                             // [0   : 125]
	second = elementsPerSlice[1]                            // [125 : 250]
	last = elementsPerSlice[len(elementsPerSlice)-offset-1] // [875 : max]

	return
}

func GroupBy[T any](elements []T, toKey func(T) Key) ([]Key, map[Key][]T) {
	grouped := make(map[Key][]T)

	for _, element := range elements {

		key := toKey(element)

		if _, exists := grouped[key]; !exists {
			grouped[key] = []T{}
		}

		grouped[key] = append(grouped[key], element)
	}

	keys := make([]Key, 0)
	for key := range grouped {
		keys = append(keys, key)
	}

	SortKeys(keys)

	return keys, grouped
}

func GroupByTimestamp[T HasTimestamp](elements []T) ([]Key, map[Key][]T) {
	return GroupBy(elements, func(hts T) Key {
		return TimeToKey(hts.GetTimestamp())
	})
}

func FillInMissingKeys(keys *[]Key) {
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
