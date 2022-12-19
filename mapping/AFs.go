package mapping

import (
	"deprec/configuration"
	"deprec/model"
	"deprec/statistics"
	"github.com/thoas/go-funk"
	"sort"
	"strings"
)

func LicenseRestrictiveness(model *model.DataModel) {

	//license := model.Repository.License

}

func Recentness(m *model.DataModel, config configuration.Configuration) model.CoreResult {

	cr := model.CoreResult{Core: model.Recentness}

	commits := m.Repository.Commits
	if commits != nil {
		sort.Slice(commits, func(i, j int) bool {
			return commits[i].Timestamp.Before(commits[j].Timestamp)
		})

		lastCommit := statistics.CalculateTimeDifference(commits[0].GetTimeStamp(), statistics.CustomNow())

		if lastCommit > config.Recentness.CommitThreshold {
			cr.Intake(0, 1)
		}
	}

	releases := m.Repository.Releases
	if releases != nil {
		sort.Slice(releases, func(i, j int) bool {
			return releases[i].Date.Before(releases[j].Date)
		})

		lastRelease := statistics.CalculateTimeDifference(releases[0].GetTimeStamp(), statistics.CustomNow())

		if lastRelease > config.Recentness.ReleaseThreshold {
			cr.Intake(0, 1)
		}
	}

	tags := m.Repository.Tags
	if tags != nil {
		sort.Slice(tags, func(i, j int) bool {
			return tags[i].Date.Before(tags[j].Date)
		})

		lastTag := statistics.CalculateTimeDifference(tags[0].GetTimeStamp(), statistics.CustomNow())

		if lastTag > config.Recentness.ReleaseThreshold {
			cr.Intake(0, 1)
		}
	}

	return cr
}

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

func DeityGiven(m *model.DataModel) model.CoreResult {

	cr := model.CoreResult{Core: model.DeityGiven}

	archived := m.Repository.Archivation
	if archived {
		cr.Intake(0, 100)
	}

	readme := strings.ToLower(m.Repository.ReadMe)
	if strings.Contains(readme, "deprecated") || strings.Contains(readme, "end-of-life") {
		cr.Intake(0, 100)
	}

	about := strings.ToLower(m.Repository.About)
	if strings.Contains(about, "deprecated") || strings.Contains(about, "end-of-life") || strings.Contains(about, "abandoned") {
		cr.Intake(0, 100)
	}

	return cr
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
