package model

import (
	"fmt"
	"github.com/thoas/go-funk"
	"strings"
)

type Dependency struct {
	Name     string
	Version  string
	MetaData map[string]string
}

type SBOM struct {
	JsonContent string
}

type AgentResult struct {
	Dependency    *Dependency
	CombConResult CoreResult
	Result        RecommendationResult
}

func (ar *AgentResult) ToString() string {

	header := fmt.Sprintf("Result %s: ", ar.Dependency.Name)
	body := ar.CombConResult.ToStringDeep()

	return header + body
}

func (ar *AgentResult) TopRecommendation() Recommendation {

	result := ar.Result

	var rec Recommendation

	tmp := -1.0
	for recommendation, f := range result {
		if f > tmp {
			rec = recommendation
			tmp = f
		}
	}

	return rec
}

type RecommendationResult map[Recommendation]float64

type Recommendation string

func (r Recommendation) ToAbbreviation() string {
	words := strings.Split(string(r), " ")

	firstChars := funk.Map(words, func(s string) string { return strings.Split(s, "")[0] })

	return funk.Reduce(firstChars, func(acc string, s string) string {
		return acc + s
	}, "").(string)
}

const (
	NoConcerns        Recommendation = "No Concerns"
	NoImmediateAction Recommendation = "No Immediate Action"
	Watchlist         Recommendation = "Watchlist"
	DecisionMaking    Recommendation = "Decision Making"
)
