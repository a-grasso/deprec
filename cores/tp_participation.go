package cores

import (
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/model"
	"github.com/thoas/go-funk"
)

func Participation(m *model.DataModel, c configuration.Participation) model.Core {

	cr := model.NewCore(model.Participation)

	if m.Repository == nil {
		return *cr
	}

	contributors := m.Repository.Contributors

	noContUser := funk.Filter(contributors, func(c model.Contributor) bool {
		return c.FirstContribution == nil
	}).([]model.Contributor)

	lowContUSer := funk.Filter(contributors, func(contributor model.Contributor) bool {
		return contributor.Contributions <= c.CommitLimit
	}).([]model.Contributor)

	noContUser = append(noContUser, lowContUSer...)

	users := funk.UniqBy(noContUser, func(c model.Contributor) string { return c.Name }).([]model.Contributor)

	totalUsers := len(contributors)

	ratio := float64(len(users)) / float64(totalUsers) * 100

	cr.IntakeThreshold(ratio, float64(c.ThirdPartyCommitThresholdPercentage), c.Weights.ThirdPartyCommits)

	return *cr
}
