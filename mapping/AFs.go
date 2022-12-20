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

func Recentness(m *model.DataModel, config configuration.Recentness) model.CoreResult {

	cr := model.CoreResult{Core: model.Recentness}

	commits := m.Repository.Commits
	if commits != nil {
		sort.Slice(commits, func(i, j int) bool {
			return commits[i].Timestamp.Before(commits[j].Timestamp)
		})

		lastCommit := statistics.CalculateTimeDifference(commits[0].GetTimeStamp(), statistics.CustomNow())

		if lastCommit > config.CommitThreshold {
			cr.Intake(0, 1)
		}
	}

	releases := m.Repository.Releases
	if releases != nil {
		sort.Slice(releases, func(i, j int) bool {
			return releases[i].Date.Before(releases[j].Date)
		})

		lastRelease := statistics.CalculateTimeDifference(releases[0].GetTimeStamp(), statistics.CustomNow())

		if lastRelease > config.ReleaseThreshold {
			cr.Intake(0, 1)
		}
	}

	tags := m.Repository.Tags
	if tags != nil {
		sort.Slice(tags, func(i, j int) bool {
			return tags[i].Date.Before(tags[j].Date)
		})

		lastTag := statistics.CalculateTimeDifference(tags[0].GetTimeStamp(), statistics.CustomNow())

		if lastTag > config.ReleaseThreshold {
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

func ProjectSize(m *model.DataModel) float64 {
	return float64(m.Repository.LOC + m.Repository.TotalCommits() + m.Repository.TotalIssues() + m.Repository.TotalContributors() + m.Repository.TotalReleases())
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
		cr.Intake(0, 1)
	}

	readme := strings.ToLower(m.Repository.ReadMe)
	if strings.Contains(readme, "deprecated") || strings.Contains(readme, "end-of-life") {
		cr.Intake(0, 1)
	}

	about := strings.ToLower(m.Repository.About)
	if strings.Contains(about, "deprecated") || strings.Contains(about, "end-of-life") || strings.Contains(about, "abandoned") {
		cr.Intake(0, 1)
	}

	return cr
}

func OrganizationalBackup(m *model.DataModel) model.CoreResult {

	cr := model.CoreResult{Core: model.OrganizationalBackup}

	contributors := m.Repository.Contributors

	sort.Slice(contributors, func(i, j int) bool {
		return contributors[i].Contributions > contributors[j].Contributions
	})

	total := 0.0
	for i, contributor := range contributors {

		weight := float64(i) / float64(len(contributors))

		sponsors := float64(contributor.Sponsors) * weight

		orgs := float64(contributor.Organizations) * weight

		total += sponsors + orgs
	}

	total /= float64(len(contributors))

	cr.Intake(total, 1)

	if m.Repository.Org != nil {
		cr.Intake(1, 3)
	} else {
		cr.Intake(0, 3)

	}

	return cr
}

func ThirdPartyParticipation(m *model.DataModel) {

	_ = m.Repository.Issues

	contributors := m.Repository.Contributors

	noContUser := funk.Filter(contributors, func(c model.Contributor) bool {
		return c.FirstContribution == nil
	}).([]model.Contributor)

	_ = float64(len(noContUser)) / float64(len(contributors))

}

func ContributorPrestige(m *model.DataModel) float64 {

	contributors := m.Repository.Contributors

	commits := m.Repository.Commits

	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Timestamp.Before(commits[j].Timestamp)
	})

	firstCommit := commits[0]
	lastCommit := commits[len(commits)-1]

	repoMonthSpan := statistics.CalculateTimeDifference(firstCommit.Timestamp, lastCommit.Timestamp)

	var prestiges []float64

	for _, c := range contributors {

		var diff float64
		if c.FirstContribution != nil {

			contributionMonthSpan := statistics.CalculateTimeDifference(*c.FirstContribution, *c.LastContribution)

			diff = float64(contributionMonthSpan) / float64(repoMonthSpan)
		}

		prestige := float64(c.Sponsors+c.Organizations+c.Repositories) + diff*10

		prestiges = append(prestiges, prestige)
	}

	result := funk.Sum(prestiges) / float64(len(prestiges))
	return result

}
