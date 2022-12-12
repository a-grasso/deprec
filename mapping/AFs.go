package mapping

import (
	"deprec/model"
	"strings"
)

func Network(model *model.DataModel) float64 {
	var result float64

	// Repository
	var repositoryNetwork float64

	var contributorRepos int
	var contributorOrgs int
	var contributorSponsors int
	var contributors int

	if model.Repository == nil {
		repositoryNetwork = 0
	} else {

		contributors += model.Repository.TotalContributors
		for _, contributor := range model.Repository.Contributors {
			contributorRepos += contributor.Repositories
			contributorOrgs += contributor.Organizations
			contributorSponsors += 0 //TODO
		}
	}

	result += repositoryNetwork

	return result
}

func Popularity(model *model.DataModel) float64 {

	stars := model.Repository.Stars
	watchers := model.Repository.Watchers
	forks := model.Repository.Forks

	contributors := model.Repository.TotalContributors

	return float64(stars + watchers + forks + contributors)
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
