package model

import (
	"fmt"
	"github.com/thoas/go-funk"
	"gonum.org/v1/gonum/mat"
	"math"
)

type Core string

const (
	CombCon              Core = "Combination And Conclusion"
	Activity             Core = "Activity"
	CoreTeam             Core = "Core Team"
	DeityGiven           Core = "Deity-Given"
	Recentness           Core = "Recentness"
	OrganizationalBackup Core = "Organizational Backup"
	Processing           Core = "Processing"
)

type CoreResult struct {
	Core              Core
	NoConcerns        float64
	NoImmediateAction float64

	Watchlist      float64
	DecisionMaking float64

	UnderlyingCores []CoreResult
}

const Separator string = " <---> "

func (cr *CoreResult) ToString() string {

	rec := cr.Softmax()
	topCore := fmt.Sprintf("Top Core: %v", cr.Core)
	softmaxResult := fmt.Sprintf("%s -> %.3f | %s -> %.3f | %s -> %.3f | %s -> %.3f", NoConcerns, rec[NoConcerns], NoImmediateAction, rec[NoImmediateAction], Watchlist, rec[Watchlist], DecisionMaking, rec[DecisionMaking])
	underlyingCores := fmt.Sprintf("Underlying Cores: %v", funk.Map(cr.UnderlyingCores, func(cr CoreResult) Core { return cr.Core }))

	return topCore + Separator + softmaxResult + Separator + underlyingCores
}

func (cr *CoreResult) ToStringDeep() string {

	rec := cr.Softmax()
	topCore := fmt.Sprintf("Top Core: %v", cr.Core)
	softmaxResult := fmt.Sprintf("%s -> %.3f | %s -> %.3f | %s -> %.3f | %s -> %.3f", NoConcerns, rec[NoConcerns], NoImmediateAction, rec[NoImmediateAction], Watchlist, rec[Watchlist], DecisionMaking, rec[DecisionMaking])
	underlyingCores := fmt.Sprintf("Underlying Cores: %v", funk.Map(cr.UnderlyingCores, func(cr CoreResult) string { return fmt.Sprintf("\n{\n%v\n}\n", cr.ToString()) }))

	return topCore + Separator + softmaxResult + Separator + underlyingCores
}

func (cr *CoreResult) Softmax() RecommendationResult {

	matrix := mat.NewDense(4, 1, []float64{cr.NoConcerns, cr.NoImmediateAction, cr.Watchlist, cr.DecisionMaking})

	var sum float64
	// Calculate the sum
	for _, v := range matrix.RawMatrix().Data {
		sum += math.Exp(v)
	}

	resultMatrix := mat.NewDense(matrix.RawMatrix().Rows, matrix.RawMatrix().Cols, nil)
	// Calculate softmax value for each element
	resultMatrix.Apply(func(i int, j int, v float64) float64 {
		return math.Exp(v) / sum
	}, matrix)

	result := make(map[Recommendation]float64)

	col := resultMatrix.ColView(0)

	result[NoConcerns] = col.At(0, 0)
	result[NoImmediateAction] = col.At(1, 0)
	result[Watchlist] = col.At(2, 0)
	result[DecisionMaking] = col.At(3, 0)

	return result
}

func (cr *CoreResult) Intake(value float64, weight float64) {

	if value >= 0.75 {
		cr.NoConcerns += weight
		return
	}

	if value >= 0.5 {
		cr.NoImmediateAction += weight
		return
	}

	if value >= 0.25 {
		cr.Watchlist += weight
		return
	}

	cr.DecisionMaking += weight

	return
}

func (cr *CoreResult) Overtake(from CoreResult, weight float64) {
	cr.NoConcerns += from.NoConcerns * weight
	cr.NoImmediateAction += from.NoImmediateAction * weight
	cr.Watchlist += from.Watchlist * weight
	cr.DecisionMaking += from.DecisionMaking * weight

	cr.UnderlyingCores = append(cr.UnderlyingCores, from)
}
