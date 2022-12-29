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
	ContributorPrestige  Core = "ContributorPrestige"
	Processing           Core = "Processing"
	Effort               Core = "Effort"
	Interconnectedness   Core = "Interconnectedness"
	Network              Core = "Network"
	Popularity           Core = "Popularity"
	Community            Core = "Community"
	Support              Core = "Support"
	Circumstances        Core = "Circumstances"
)

type CoreResult struct {
	Core              Core
	NoConcerns        float64
	NoImmediateAction float64

	Watchlist      float64
	DecisionMaking float64

	UnderlyingCores map[float64][]CoreResult
}

func NewCoreResult(core Core) CoreResult {
	return CoreResult{Core: core, UnderlyingCores: make(map[float64][]CoreResult)}
}

const Separator string = " <---> "

func (cr *CoreResult) ToString() string {

	rec := cr.Softmax()
	topCore := fmt.Sprintf("Top Core: %v", cr.Core)
	softmaxResult := fmt.Sprintf("%s -> %.3f | %s -> %.3f | %s -> %.3f | %s -> %.3f", NoConcerns, rec[NoConcerns], NoImmediateAction, rec[NoImmediateAction], Watchlist, rec[Watchlist], DecisionMaking, rec[DecisionMaking])
	underlyingCores := fmt.Sprintf("Underlying Cores: %v", funk.Map(cr.UnderlyingCores, func(weight float64, cr []CoreResult) (float64, []Core) {
		ads := funk.Map(cr, func(cr CoreResult) Core {
			return cr.Core
		}).([]Core)
		return weight, ads
	}))

	return topCore + Separator + softmaxResult + Separator + underlyingCores
}

func (cr *CoreResult) ToStringDeep() string {

	rec := cr.Softmax()
	topCore := fmt.Sprintf("Top Core: %v", cr.Core)
	softmaxResult := fmt.Sprintf("%s -> %.3f | %s -> %.3f | %s -> %.3f | %s -> %.3f", NoConcerns, rec[NoConcerns], NoImmediateAction, rec[NoImmediateAction], Watchlist, rec[Watchlist], DecisionMaking, rec[DecisionMaking])
	underlyingCores := fmt.Sprintf("Underlying Cores: %v", funk.Map(cr.UnderlyingCores, func(cr CoreResult) string { return fmt.Sprintf("\n{\n%v\n}\n", cr.ToString()) }))

	return topCore + Separator + softmaxResult + Separator + underlyingCores
}

func (cr *CoreResult) Normalized() CoreResult {

	var total float64
	total += cr.NoConcerns
	total += cr.NoImmediateAction
	total += cr.Watchlist
	total += cr.DecisionMaking

	if total == 0 {
		total = 1
	}

	return CoreResult{
		Core:              cr.Core,
		NoConcerns:        cr.NoConcerns / total,
		NoImmediateAction: cr.NoImmediateAction / total,
		Watchlist:         cr.Watchlist / total,
		DecisionMaking:    cr.DecisionMaking / total,
		UnderlyingCores:   cr.UnderlyingCores,
	}
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

func (cr *CoreResult) IntakeThreshold(value, threshold, weight float64) {

	v := math.Min(1, value/threshold)

	cr.Intake(v, weight)
}

func (cr *CoreResult) IntakeLimit(value, limit, weight float64) {

	v := math.Max(0, 1-value/limit)

	cr.Intake(v, weight)
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

	normalized := from
	//normalized := from.Normalized()

	cr.NoConcerns += normalized.NoConcerns * weight
	cr.NoImmediateAction += normalized.NoImmediateAction * weight
	cr.Watchlist += normalized.Watchlist * weight
	cr.DecisionMaking += normalized.DecisionMaking * weight

	cr.UnderlyingCores[weight] = append(cr.UnderlyingCores[weight], normalized)
}
