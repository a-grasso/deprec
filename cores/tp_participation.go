package cores

import (
	"deprec/configuration"
	"deprec/model"
	"github.com/thoas/go-funk"
)

func ThirdPartyParticipation(m *model.DataModel, c configuration.ThirdPartyParticipation) model.CoreResult {

	cr := model.NewCoreResult(model.ThirdPartyParticipation)

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

	ratio := float64(len(users)) / float64(totalUsers)

	cr.IntakeThreshold(ratio, float64(c.ThirdPartyCommitThresholdPercentage), 1)

	return *cr
}
