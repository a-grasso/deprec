package cores

import (
	"deprec/configuration"
	"deprec/model"
)

func Network(m *model.DataModel, c configuration.Network) model.CoreResult {
	cr := model.NewCoreResult(model.Network)

	if m.Repository != nil {
		var repositoryNetwork int

		repositoryNetwork += m.Repository.TotalContributors()
		for _, contributor := range m.Repository.Contributors {
			repositoryNetwork += contributor.Repositories
			repositoryNetwork += contributor.Organizations
		}

		if m.Repository.Org != nil {
			org := m.Repository.Org
			repositoryNetwork += org.PublicRepos
			repositoryNetwork += org.OwnedPrivateRepos
			repositoryNetwork += org.Collaborators
			repositoryNetwork += org.Followers
		}

		cr.IntakeThreshold(float64(repositoryNetwork), float64(c.Threshold), 1)
	}

	return *cr
}
