package model

import (
	"fmt"
	"github.com/thoas/go-funk"
	"gonum.org/v1/gonum/mat"
	"math"
)

type CoreName string

const (
	CombCon            CoreName = "Combination And Conclusion"
	Activity           CoreName = "Activity"
	CoreTeam           CoreName = "Core Team"
	Rivalry            CoreName = "Rivalry"
	DeityGiven         CoreName = "Deity-Given"
	Recentness         CoreName = "Recentness"
	Backup             CoreName = "Backup"
	Participation      CoreName = "Participation"
	Prestige           CoreName = "Prestige"
	Processing         CoreName = "Processing"
	Effort             CoreName = "Effort"
	Interconnectedness CoreName = "Interconnectedness"
	Network            CoreName = "Network"
	Popularity         CoreName = "Popularity"
	Vulnerabilities    CoreName = "Vulnerabilities"
	Community          CoreName = "Community"
	Support            CoreName = "Support"
	Circumstances      CoreName = "Circumstances"
	ProjectQuality     CoreName = "ProjectQuality"
	Engagement         CoreName = "Engagement"
	Licensing          CoreName = "Licensing"
	Marking            CoreName = "Marking"
)

const (
	NC  float64 = 0.875
	NIA float64 = 0.625
	W   float64 = 0.375
	DM  float64 = 0.125
)

type Core struct {
	Name              CoreName
	NoConcerns        float64
	NoImmediateAction float64

	Watchlist      float64
	DecisionMaking float64

	UnderlyingCores map[float64][]Core
}

func NewCore(core CoreName) *Core {
	return &Core{Name: core, UnderlyingCores: make(map[float64][]Core)}
}

const Separator string = " <---> "

func (cr *Core) ToString() string {

	rec := cr.Softmax()
	topCore := fmt.Sprintf("Top Core: %v", cr.Name)
	softmaxResult := fmt.Sprintf("%s -> %.3f | %s -> %.3f | %s -> %.3f | %s -> %.3f", NoConcerns, rec[NoConcerns], NoImmediateAction, rec[NoImmediateAction], Watchlist, rec[Watchlist], DecisionMaking, rec[DecisionMaking])
	underlyingCores := fmt.Sprintf("Underlying Cores: %v", funk.Map(cr.UnderlyingCores, func(weight float64, cr []Core) (float64, []CoreName) {
		ads := funk.Map(cr, func(cr Core) CoreName {
			return cr.Name
		}).([]CoreName)
		return weight, ads
	}))

	return topCore + Separator + softmaxResult + Separator + underlyingCores
}

func (cr *Core) ToStringDeep() string {

	rec := cr.Softmax()
	topCore := fmt.Sprintf("Top Core: %v", cr.Name)
	softmaxResult := fmt.Sprintf("%s -> %.3f | %s -> %.3f | %s -> %.3f | %s -> %.3f", NoConcerns, rec[NoConcerns], NoImmediateAction, rec[NoImmediateAction], Watchlist, rec[Watchlist], DecisionMaking, rec[DecisionMaking])
	underlyingCores := fmt.Sprintf("Underlying Cores: %v", funk.Map(cr.UnderlyingCores, func(weight float64, cr []Core) string {
		return fmt.Sprintf("\n{ Weight: %f\n%v\n}\n", weight, funk.Map(cr, func(c Core) string { return fmt.Sprintf("\n{\n%v\n}\n", c.ToStringDeep()) }))
	}))

	return topCore + Separator + softmaxResult + Separator + underlyingCores
}

func (cr *Core) GetAllCores() []Core {

	var result []Core

	result = append(result, *cr)

	for _, factors := range cr.UnderlyingCores {

		for _, factor := range factors {

			result = append(result, factor)

			for _, statements := range factor.UnderlyingCores {

				for _, statement := range statements {
					result = append(result, statement)
				}
			}
		}
	}
	return result
}
func (cr *Core) Sum() float64 {
	return cr.DecisionMaking + cr.Watchlist + cr.NoImmediateAction + cr.NoConcerns
}

func (cr *Core) IsInconclusive() bool {
	values := []float64{cr.NoConcerns, cr.NoImmediateAction, cr.Watchlist, cr.DecisionMaking}
	unique := funk.Uniq(values).([]float64)
	if len(unique) == 1 {
		return true
	}
	return false
}

func (cr *Core) Normalized() Core {

	var total float64
	total += cr.NoConcerns
	total += cr.NoImmediateAction
	total += cr.Watchlist
	total += cr.DecisionMaking

	if total == 0 {
		total = 1
	}

	return Core{
		Name:              cr.Name,
		NoConcerns:        cr.NoConcerns / total,
		NoImmediateAction: cr.NoImmediateAction / total,
		Watchlist:         cr.Watchlist / total,
		DecisionMaking:    cr.DecisionMaking / total,
		UnderlyingCores:   cr.UnderlyingCores,
	}
}

type RecommendationDistribution map[Recommendation]float64

func (cr *Core) Softmax() RecommendationDistribution {

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

	result := make(RecommendationDistribution)

	col := resultMatrix.ColView(0)

	result[NoConcerns] = col.At(0, 0)
	result[NoImmediateAction] = col.At(1, 0)
	result[Watchlist] = col.At(2, 0)
	result[DecisionMaking] = col.At(3, 0)

	return result
}

func (cr *Core) IntakeThreshold(value, threshold, weight float64) {

	v := math.Min(1, value/threshold)

	cr.Intake(v, weight)
}

func (cr *Core) IntakeLimit(value, limit, weight float64) {

	r := value / limit
	v := math.Max(0, 1-r)

	if v < 0.25 && v > 0 {
		v += 0.25
	}

	v = math.Max(0, v)
	v = math.Min(1, v)

	cr.Intake(v, weight)
}

func (cr *Core) Intake(value float64, weight float64) {

	if value > 1 {
		return
	}

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

	if value >= 0 {
		cr.DecisionMaking += weight
		return
	}

	return
}

func (cr *Core) Overtake(from Core, weight float64) {

	normalized := from.Normalized()

	cr.NoConcerns += normalized.NoConcerns * weight
	cr.NoImmediateAction += normalized.NoImmediateAction * weight
	cr.Watchlist += normalized.Watchlist * weight
	cr.DecisionMaking += normalized.DecisionMaking * weight

	cr.UnderlyingCores[weight] = append(cr.UnderlyingCores[weight], from)
}
