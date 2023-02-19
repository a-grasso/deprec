package cores

import (
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/model"
	"github.com/a-grasso/deprec/statistics"
	"github.com/thoas/go-funk"
	"sort"
)

func Prestige(m *model.DataModel, c configuration.Prestige) model.Core {

	cr := model.NewCore(model.Prestige)

	if m.Repository == nil {
		return *cr
	}

	contributors := m.Repository.Contributors

	sort.Slice(contributors, func(i, j int) bool {
		return contributors[i].Contributions > contributors[j].Contributions
	})

	commits := m.Repository.Commits

	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Timestamp.Before(commits[j].Timestamp)
	})

	firstCommit := commits[0]
	lastCommit := commits[len(commits)-1]

	repoMonthSpan := statistics.CalculateTimeDifference(firstCommit.Timestamp, lastCommit.Timestamp)

	var prestiges []float64

	for i, c := range contributors {

		var diff float64
		if c.FirstContribution != nil {

			contributionMonthSpan := statistics.CalculateTimeDifference(*c.FirstContribution, *c.LastContribution)

			diff = float64(contributionMonthSpan) / float64(repoMonthSpan)
		}

		prestige := float64(c.Sponsors+c.Organizations+c.Repositories) + diff*10

		prestige *= float64(len(contributors)-i) / float64(len(contributors))

		prestiges = append(prestiges, prestige)
	}

	result := funk.Sum(prestiges) / float64(len(prestiges))
	cr.Intake(result, c.Weights.Contributors)

	return *cr
}
