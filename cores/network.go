package cores

import (
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/model"
)

func Network(m model.DataModel, c configuration.Network) model.Core {
	cr := model.NewCore(model.Network)

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

		cr.IntakeThreshold(float64(repositoryNetwork), float64(c.Threshold), c.Weights.RepositoryNetwork)
	}

	return *cr
}
