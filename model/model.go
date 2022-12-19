package model

import (
	"fmt"
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

type RecommendationResult map[Recommendation]float64

type Recommendation string

const (
	NoConcerns        Recommendation = "No Concerns"
	NoImmediateAction Recommendation = "No Immediate Action"
	Watchlist         Recommendation = "Watchlist"
	DecisionMaking    Recommendation = "Decision Making"
)
