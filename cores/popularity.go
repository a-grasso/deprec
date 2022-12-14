package cores

import (
	"deprec/configuration"
	"deprec/model"
)

func Popularity(m *model.DataModel, c configuration.Popularity) model.CoreResult {
	cr := model.NewCoreResult(model.Popularity)

	if m.Repository != nil {
		var repositoryPopularity int

		repositoryPopularity += m.Repository.Stars
		repositoryPopularity += m.Repository.Watchers
		repositoryPopularity += m.Repository.Forks

		cr.IntakeThreshold(float64(repositoryPopularity), float64(c.Threshold), 1)
	}

	return *cr
}
