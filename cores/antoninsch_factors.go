package cores

import (
	"deprec/configuration"
	"deprec/model"
	"strings"
)

func ProjectSize(m *model.DataModel) float64 {
	return float64(m.Repository.LOC + m.Repository.TotalCommits() + m.Repository.TotalIssues() + m.Repository.TotalContributors() + m.Repository.TotalReleases())
}

func DeityGiven(m *model.DataModel) model.CoreResult {

	cr := model.NewCoreResult(model.DeityGiven)

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

func Effort(m *model.DataModel, c configuration.CoresConfig) model.CoreResult {
	cr := model.NewCoreResult(model.Effort)

	activity := Activity(m, c.Activity)

	recentness := Recentness(m, c.Recentness)

	coreTeam := CoreTeam(m, c.CoreTeam)

	cr.Overtake(recentness, 5)
	cr.Overtake(activity, 2)
	cr.Overtake(coreTeam, 1)

	return cr
}

func Interconnectedness(m *model.DataModel, c configuration.CoresConfig) model.CoreResult {
	cr := model.NewCoreResult(model.Interconnectedness)

	network := Network(m, c.Network)

	popularity := Popularity(m, c.Popularity)

	cr.Overtake(network, 1)

	cr.Overtake(popularity, 1)

	return cr
}

func Community(m *model.DataModel, c configuration.CoresConfig) model.CoreResult {
	cr := model.NewCoreResult(model.Community)

	//contributorPrestige := ContributorPrestige(m)

	//thirdPartyParticipation := ThirdPartyParticipation(m)

	organizationalBackup := OrganizationalBackup(m, c.OrgBackup)

	//cr.Overtake(contributorPrestige, 0)

	cr.Overtake(organizationalBackup, 1)

	return cr
}

func Support(m *model.DataModel, c configuration.CoresConfig) model.CoreResult {
	cr := model.NewCoreResult(model.Support)

	processing := Processing(m, c.Processing)

	engagement := Engagement(m, c.Engagement)

	cr.Overtake(processing, 2)

	cr.Overtake(engagement, 1)
	return cr
}

func Ecosystem(m *model.DataModel) {

}

func Circumstances(m *model.DataModel, c configuration.CoresConfig) model.CoreResult {

	cr := model.NewCoreResult(model.Circumstances)

	return cr
}
