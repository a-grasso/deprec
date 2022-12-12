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

func statisticAnalysisPYM[T interface {
	model.Commit | model.Issue | model.Release
	GetTimeStamp() time.Time
}](data []T) (avg, last, avgPercentage, lastPercentage float64, yearsSinceLast, monthsSinceLast int) {

	total := len(data)
	if total == 0 {
		return
	}

	sortedKeys, groupedByYM := groupByYM[T](data)

	lastKey := sortedKeys[len(sortedKeys)-1]
	yearsSinceLast, monthsSinceLast = calculateDifferenceOfKeyFromNow(lastKey)

	groupedCounts := funk.Map(groupedByYM, func(k Key, v []T) (Key, int) {
		return k, len(v)
	}).(map[Key]int)

	groupedPercentages := funk.Map(groupedByYM, func(k Key, v []T) (Key, float64) {
		return k, float64(len(v)) / float64(total) * 100
	}).(map[Key]float64)

	perMonth := make([]int, 0, len(sortedKeys))
	percentagesPerMonth := make([]float64, 0, len(sortedKeys))
	for _, key := range sortedKeys {
		perMonth = append(perMonth, groupedCounts[key])
		percentagesPerMonth = append(percentagesPerMonth, groupedPercentages[key])
	}

	avg = float64(total) / float64(len(perMonth))
	avgPercentage = 100.0 / float64(len(percentagesPerMonth))

	last = float64(groupedCounts[sortedKeys[len(sortedKeys)-1]])
	lastPercentage = groupedPercentages[sortedKeys[len(sortedKeys)-1]]

	return
}

func Activity(m *model.DataModel, config configuration.AFConfig) float64 {

	commits := m.Repository.Commits
	avgCommits, lastCommits, avgPercentage, lastPercentage, yearsSinceLastCommit, monthsSinceLastcommit := statisticAnalysisPYM(commits)
	log.Println(avgCommits, lastCommits, avgPercentage, lastPercentage, monthsSinceLastcommit)
	//					13			3				0.39		0.09			11
	if yearsSinceLastCommit > config.CommitThreshold {
		return 0
	}

	issues := m.Repository.Issues
	log.Println(statisticAnalysisPYM(issues))

	//	releases := m.Repository.Releases
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

func groupByYM[T interface {
	model.Commit | model.Issue | model.Release
	GetTimeStamp() time.Time
}](elements []T) ([]Key, map[Key][]T) {
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
