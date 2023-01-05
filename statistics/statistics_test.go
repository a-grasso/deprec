package statistics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPercentileOf_10elements_50percentile(t *testing.T) {
	elements := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	percentile := 50

	expectedFirst := []int{1, 2, 3, 4, 5}
	expectedSecond := []int{6, 7, 8, 9, 10}
	expectedLast := []int{6, 7, 8, 9, 10}

	actualFirst, actualSecond, actualLast := GetPercentilesOf(elements, float64(percentile))

	assert.EqualValues(t, expectedFirst, actualFirst)
	assert.EqualValues(t, expectedSecond, actualSecond)
	assert.EqualValues(t, expectedLast, actualLast)
}

func TestGetPercentileOf_3elements_20percentile(t *testing.T) {
	elements := []int{1, 2, 3}

	percentile := 20

	expectedFirst := []int{1}
	expectedSecond := []int{2}
	expectedLast := []int{3}

	actualFirst, actualSecond, actualLast := GetPercentilesOf(elements, float64(percentile))

	assert.EqualValues(t, expectedFirst, actualFirst)
	assert.EqualValues(t, expectedSecond, actualSecond)
	assert.EqualValues(t, expectedLast, actualLast)
}

func TestGetPercentileOf_2elements_20percentile(t *testing.T) {
	elements := []int{1, 2}

	percentile := 20

	expectedFirst := []int{1}
	expectedSecond := []int{2}
	expectedLast := []int{2}

	actualFirst, actualSecond, actualLast := GetPercentilesOf(elements, float64(percentile))

	assert.EqualValues(t, expectedFirst, actualFirst)
	assert.EqualValues(t, expectedSecond, actualSecond)
	assert.EqualValues(t, expectedLast, actualLast)
}

func TestGetPercentileOf_1elements_20percentile(t *testing.T) {
	elements := []int{1}

	percentile := 20

	expectedFirst := []int{1}
	expectedLast := []int{1}

	actualFirst, actualSecond, actualLast := GetPercentilesOf(elements, float64(percentile))

	assert.EqualValues(t, expectedFirst, actualFirst)
	assert.Nil(t, nil, actualSecond)
	assert.EqualValues(t, expectedLast, actualLast)
}

func TestGetPercentileOf_1elements_50percentile(t *testing.T) {
	elements := []int{1}

	percentile := 50

	expectedFirst := []int{1}
	expectedLast := []int{1}

	actualFirst, actualSecond, actualLast := GetPercentilesOf(elements, float64(percentile))

	assert.EqualValues(t, expectedFirst, actualFirst)
	assert.Nil(t, nil, actualSecond)
	assert.EqualValues(t, expectedLast, actualLast)
}

func TestGetPercentileOf_10elements_2percentile(t *testing.T) {
	elements := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	percentile := 2

	expectedFirst := []int{1}
	expectedLast := []int{10}

	actualFirst, actualSecond, actualLast := GetPercentilesOf(elements, float64(percentile))

	assert.EqualValues(t, expectedFirst, actualFirst)
	assert.Nil(t, nil, actualSecond)
	assert.EqualValues(t, expectedLast, actualLast)
}

func TestGetPercentileOf_100elements_2percentile(t *testing.T) {

	var elements []int
	for i := 0; i < 100; i++ {
		elements = append(elements, i)
	}

	percentile := 2

	expectedFirst := []int{0, 1}
	expectedSecond := []int{2, 3}
	expectedLast := []int{98, 99}

	actualFirst, actualSecond, actualLast := GetPercentilesOf(elements, float64(percentile))

	assert.EqualValues(t, expectedFirst, actualFirst)
	assert.EqualValues(t, expectedSecond, actualSecond)
	assert.EqualValues(t, expectedLast, actualLast)
}

func TestGetPercentileOf_10elements_30percentile(t *testing.T) {

	var elements []int
	for i := 0; i < 10; i++ {
		elements = append(elements, i)
	}

	percentile := 30

	expectedFirst := []int{0, 1, 2}
	expectedSecond := []int{3, 4, 5}
	expectedLast := []int{6, 7, 8, 9}

	actualFirst, actualSecond, actualLast := GetPercentilesOf(elements, float64(percentile))

	assert.EqualValues(t, expectedFirst, actualFirst)
	assert.EqualValues(t, expectedSecond, actualSecond)
	assert.EqualValues(t, expectedLast, actualLast)
}

func TestGetPercentileOf_10elements_33percentile(t *testing.T) {

	var elements []int
	for i := 0; i < 10; i++ {
		elements = append(elements, i)
	}

	percentile := 33.33333333

	expectedFirst := []int{0, 1, 2}
	expectedSecond := []int{3, 4, 5}
	expectedLast := []int{6, 7, 8, 9}

	actualFirst, actualSecond, actualLast := GetPercentilesOf(elements, float64(percentile))

	assert.EqualValues(t, expectedFirst, actualFirst)
	assert.EqualValues(t, expectedSecond, actualSecond)
	assert.EqualValues(t, expectedLast, actualLast)
}

func TestGetPercentileOf_10elements_34percentile(t *testing.T) {

	var elements []int
	for i := 0; i < 10; i++ {
		elements = append(elements, i)
	}

	percentile := 34

	expectedFirst := []int{0, 1, 2}
	expectedSecond := []int{3, 4, 5}
	expectedLast := []int{6, 7, 8, 9}

	actualFirst, actualSecond, actualLast := GetPercentilesOf(elements, float64(percentile))

	assert.EqualValues(t, expectedFirst, actualFirst)
	assert.EqualValues(t, expectedSecond, actualSecond)
	assert.EqualValues(t, expectedLast, actualLast)
}

func TestGetPercentileOf_10elements_20percentile(t *testing.T) {

	var elements []int
	for i := 0; i < 10; i++ {
		elements = append(elements, i)
	}

	percentile := 20

	expectedFirst := []int{0, 1}
	expectedSecond := []int{2, 3}
	expectedLast := []int{8, 9}

	actualFirst, actualSecond, actualLast := GetPercentilesOf(elements, float64(percentile))

	assert.EqualValues(t, expectedFirst, actualFirst)
	assert.EqualValues(t, expectedSecond, actualSecond)
	assert.EqualValues(t, expectedLast, actualLast)
}

func TestGetPercentileOf_10elements_60percentile(t *testing.T) {

	var elements []int
	for i := 0; i < 10; i++ {
		elements = append(elements, i)
	}

	percentile := 60

	expectedFirst := []int{0, 1, 2, 3, 4}
	expectedSecond := []int{5, 6, 7, 8, 9}
	expectedLast := []int{5, 6, 7, 8, 9}

	actualFirst, actualSecond, actualLast := GetPercentilesOf(elements, float64(percentile))

	assert.EqualValues(t, expectedFirst, actualFirst)
	assert.EqualValues(t, expectedSecond, actualSecond)
	assert.EqualValues(t, expectedLast, actualLast)
}

func TestGetPercentileOf_10elements_40percentile(t *testing.T) {

	var elements []int
	for i := 0; i < 10; i++ {
		elements = append(elements, i)
	}

	percentile := 40

	expectedFirst := []int{0, 1, 2}
	expectedSecond := []int{3, 4, 5}
	expectedLast := []int{6, 7, 8, 9}

	actualFirst, actualSecond, actualLast := GetPercentilesOf(elements, float64(percentile))

	assert.EqualValues(t, expectedFirst, actualFirst)
	assert.EqualValues(t, expectedSecond, actualSecond)
	assert.EqualValues(t, expectedLast, actualLast)
}

func TestGetPercentileOf_10elements_25percentile(t *testing.T) {

	var elements []int
	for i := 0; i < 10; i++ {
		elements = append(elements, i)
	}

	percentile := 25

	expectedFirst := []int{0, 1}
	expectedSecond := []int{2, 3, 4}
	expectedLast := []int{7, 8, 9}

	actualFirst, actualSecond, actualLast := GetPercentilesOf(elements, float64(percentile))

	assert.EqualValues(t, expectedFirst, actualFirst)
	assert.EqualValues(t, expectedSecond, actualSecond)
	assert.EqualValues(t, expectedLast, actualLast)
}
