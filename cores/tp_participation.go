package cores

import (
	"deprec/model"
	"github.com/thoas/go-funk"
)

func ThirdPartyParticipation(m *model.DataModel) {

	_ = m.Repository.Issues

	contributors := m.Repository.Contributors

	noContUser := funk.Filter(contributors, func(c model.Contributor) bool {
		return c.FirstContribution == nil
	}).([]model.Contributor)

	_ = float64(len(noContUser)) / float64(len(contributors))
}
