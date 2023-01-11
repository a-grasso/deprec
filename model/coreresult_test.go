package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func NewCR() CoreResult {
	return CoreResult{}
}

func TestIntakeLimit_0_4(t *testing.T) {

	cr := NewCR()

	expected := []float64{1, 0, 0, 0}

	cr.IntakeLimit(0, 1, 1)

	assert.Equal(t, cr.NoConcerns, expected[0])
	assert.Equal(t, cr.NoImmediateAction, expected[1])
	assert.Equal(t, cr.Watchlist, expected[2])
	assert.Equal(t, cr.DecisionMaking, expected[3])
}

func TestIntakeLimit_1_4(t *testing.T) {

	cr := NewCR()

	expected := []float64{1, 0, 0, 0}

	cr.IntakeLimit(0.25, 1, 1)

	assert.Equal(t, cr.NoConcerns, expected[0])
	assert.Equal(t, cr.NoImmediateAction, expected[1])
	assert.Equal(t, cr.Watchlist, expected[2])
	assert.Equal(t, cr.DecisionMaking, expected[3])
}

func TestIntakeLimit_2_4(t *testing.T) {

	cr := NewCR()

	expected := []float64{0, 1, 0, 0}

	cr.IntakeLimit(0.5, 1, 1)

	assert.Equal(t, cr.NoConcerns, expected[0])
	assert.Equal(t, cr.NoImmediateAction, expected[1])
	assert.Equal(t, cr.Watchlist, expected[2])
	assert.Equal(t, cr.DecisionMaking, expected[3])
}
func TestIntakeLimit_3_4(t *testing.T) {

	cr := NewCR()

	expected := []float64{0, 0, 1, 0}

	cr.IntakeLimit(0.75, 1, 1)

	assert.Equal(t, cr.NoConcerns, expected[0])
	assert.Equal(t, cr.NoImmediateAction, expected[1])
	assert.Equal(t, cr.Watchlist, expected[2])
	assert.Equal(t, cr.DecisionMaking, expected[3])
}
func TestIntakeLimit_4_4(t *testing.T) {

	cr := NewCR()

	expected := []float64{0, 0, 0, 1}

	cr.IntakeLimit(1, 1, 1)

	assert.Equal(t, cr.NoConcerns, expected[0])
	assert.Equal(t, cr.NoImmediateAction, expected[1])
	assert.Equal(t, cr.Watchlist, expected[2])
	assert.Equal(t, cr.DecisionMaking, expected[3])
}

func TestIntakeLimit_99_100(t *testing.T) {

	cr := NewCR()

	expected := []float64{0, 0, 1, 0}

	cr.IntakeLimit(99, 100, 1)

	assert.Equal(t, cr.NoConcerns, expected[0])
	assert.Equal(t, cr.NoImmediateAction, expected[1])
	assert.Equal(t, cr.Watchlist, expected[2])
	assert.Equal(t, cr.DecisionMaking, expected[3])
}

func TestIntakeThreshold_99_100(t *testing.T) {

	cr := NewCR()

	expected := []float64{1, 0, 0, 0}

	cr.IntakeThreshold(99, 100, 1)

	assert.Equal(t, cr.NoConcerns, expected[0])
	assert.Equal(t, cr.NoImmediateAction, expected[1])
	assert.Equal(t, cr.Watchlist, expected[2])
	assert.Equal(t, cr.DecisionMaking, expected[3])
}

func TestIntakeThreshold_4_4(t *testing.T) {

	cr := NewCR()

	expected := []float64{1, 0, 0, 0}

	cr.IntakeThreshold(1, 1, 1)

	assert.Equal(t, cr.NoConcerns, expected[0])
	assert.Equal(t, cr.NoImmediateAction, expected[1])
	assert.Equal(t, cr.Watchlist, expected[2])
	assert.Equal(t, cr.DecisionMaking, expected[3])
}

func TestIntakeThreshold_3_4(t *testing.T) {

	cr := NewCR()

	expected := []float64{1, 0, 0, 0}

	cr.IntakeThreshold(0.75, 1, 1)

	assert.Equal(t, cr.NoConcerns, expected[0])
	assert.Equal(t, cr.NoImmediateAction, expected[1])
	assert.Equal(t, cr.Watchlist, expected[2])
	assert.Equal(t, cr.DecisionMaking, expected[3])
}

func TestIntakeThreshold_2_4(t *testing.T) {

	cr := NewCR()

	expected := []float64{0, 1, 0, 0}

	cr.IntakeThreshold(0.5, 1, 1)

	assert.Equal(t, cr.NoConcerns, expected[0])
	assert.Equal(t, cr.NoImmediateAction, expected[1])
	assert.Equal(t, cr.Watchlist, expected[2])
	assert.Equal(t, cr.DecisionMaking, expected[3])
}

func TestIntakeThreshold_1_4(t *testing.T) {

	cr := NewCR()

	expected := []float64{0, 0, 1, 0}

	cr.IntakeThreshold(0.25, 1, 1)

	assert.Equal(t, cr.NoConcerns, expected[0])
	assert.Equal(t, cr.NoImmediateAction, expected[1])
	assert.Equal(t, cr.Watchlist, expected[2])
	assert.Equal(t, cr.DecisionMaking, expected[3])
}

func TestIntakeThreshold_0_4(t *testing.T) {

	cr := NewCR()

	expected := []float64{0, 0, 0, 1}

	cr.IntakeThreshold(0, 1, 1)

	assert.Equal(t, cr.NoConcerns, expected[0])
	assert.Equal(t, cr.NoImmediateAction, expected[1])
	assert.Equal(t, cr.Watchlist, expected[2])
	assert.Equal(t, cr.DecisionMaking, expected[3])
}
