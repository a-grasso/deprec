package model

import (
	"github.com/thoas/go-funk"
	"strings"
)

type Recommendation string

const (
	NoConcerns        Recommendation = "No Concerns"
	NoImmediateAction Recommendation = "No Immediate Action"
	Watchlist         Recommendation = "Watchlist"
	DecisionMaking    Recommendation = "Decision Making"
)

func (r Recommendation) ToAbbreviation() string {
	words := strings.Split(string(r), " ")

	firstChars := funk.Map(words, func(s string) string { return strings.Split(s, "")[0] })

	return funk.Reduce(firstChars, func(acc string, s string) string {
		return acc + s
	}, "").(string)
}
