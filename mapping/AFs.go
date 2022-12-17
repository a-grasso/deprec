package mapping

import (
	"deprec/model"
	"github.com/thoas/go-funk"
	"strings"
)

func Network(model *model.DataModel) float64 {
	var result float64

	var repositoryNetwork int

	var contributorRepos int
	var contributorOrgs int
	var contributorSponsors int
	var contributors int

	var orgRepos int
	var orgSponsors int
	var orgCollaborators int
	var orgFollowers int

	if model.Repository == nil {
		repositoryNetwork = 0
	} else {

		contributors += model.Repository.TotalContributors()
		for _, contributor := range model.Repository.Contributors {
			contributorRepos += contributor.Repositories
			contributorOrgs += contributor.Organizations
			contributorSponsors += 0 //TODO
		}

		if model.Repository.Org != nil {
			org := model.Repository.Org
			orgRepos = org.PublicRepos + org.OwnedPrivateRepos // TODO + org.TotalPrivateRepos ???
			orgCollaborators = org.Collaborators
			orgFollowers = org.Followers
		}
	}

	repositoryNetwork += contributorRepos
	repositoryNetwork += contributorOrgs
	repositoryNetwork += contributorSponsors
	repositoryNetwork += contributors

	repositoryNetwork += orgRepos
	repositoryNetwork += orgSponsors
	repositoryNetwork += orgCollaborators
	repositoryNetwork += orgFollowers

	result += float64(repositoryNetwork)
	return result
}

func Popularity(model *model.DataModel) float64 {

	stars := model.Repository.Stars
	watchers := model.Repository.Watchers
	forks := model.Repository.Forks

	// TODO users := TotalContributors - ContributorsThatContributed :-> does that work?

	return float64(stars + watchers + forks)
}

func Interconnectedness(model *model.DataModel) float64 {
	return Network(model)*0.3 + Popularity(model)*0.7
}

func DeityGiven(model *model.DataModel) float64 {

	archived := model.Repository.Archivation
	if archived {
		return 1
	}

	readme := strings.ToLower(model.Repository.ReadMe)
	if strings.Contains(readme, "deprecated") || strings.Contains(readme, "end-of-life") {
		return 1
	}

	about := strings.ToLower(model.Repository.About)
	if strings.Contains(about, "deprecated") || strings.Contains(about, "end-of-life") || strings.Contains(about, "abandoned") {
		return 1
	}

	return 0
}

func OrganizationalBackup(m *model.DataModel) float64 {

	contOrgs := funk.Sum(funk.Map(m.Repository.Contributors, func(c model.Contributor) int { return c.Organizations }))

	contSpons := funk.Sum(funk.Map(m.Repository.Contributors, func(c model.Contributor) int { return c.Sponsors }))

	organization := 0.0
	if m.Repository.Org != nil {
		organization = 1.0
	}

	return organization*0.5 + contOrgs*0.25 + contSpons*0.25
}
