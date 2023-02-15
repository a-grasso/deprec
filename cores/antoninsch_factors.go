package cores

import (
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/model"
	"strings"
)

func DeityGiven(m *model.DataModel) model.Core {

	cr := model.NewCore(model.DeityGiven)

	if m.Repository != nil {
		archived := m.Repository.Archivation
		if archived {
			cr.Intake(model.DM, 1)
		}

		readme := strings.ToLower(m.Repository.ReadMe)
		if strings.Contains(readme, "end-of-life") {
			cr.Intake(model.DM, 1)
		}

		about := strings.ToLower(m.Repository.About)
		if strings.Contains(about, "deprecated") || strings.Contains(about, "end-of-life") || strings.Contains(about, "abandoned") {
			cr.Intake(model.DM, 1)
		}
	}

	if distribution := m.Distribution; distribution != nil {

		if artifact := distribution.Artifact; artifact != nil {

			description := strings.ToLower(artifact.Description)
			if strings.Contains(description, "deprecated") || strings.Contains(description, "end-of-life") || strings.Contains(description, "abandoned") {
				cr.Intake(model.DM, 1)
			}
		}
	}

	return *cr
}

func Vulnerabilities(m *model.DataModel) model.Core {

	cr := model.NewCore(model.Vulnerabilities)

	if m.VulnerabilityIndex == nil {
		return *cr
	}

	vulnerabilities := m.VulnerabilityIndex.TotalVulnerabilitiesCount

	if vulnerabilities > 0 {
		cr.Intake(model.DM, 1)
	}

	return *cr
}

func Effort(m *model.DataModel, c configuration.CoresConfig) model.Core {

	cr := model.NewCore(model.Effort)

	activity := Activity(m, c.Activity)

	recentness := Recentness(m, c.Recentness)

	coreTeam := CoreTeam(m, c.CoreTeam)

	cr.Overtake(recentness, 2)
	cr.Overtake(activity, 2)
	cr.Overtake(coreTeam, 1)

	return *cr
}

func Interconnectedness(m *model.DataModel, c configuration.CoresConfig) model.Core {

	cr := model.NewCore(model.Interconnectedness)

	network := Network(m, c.Network)

	popularity := Popularity(m, c.Popularity)

	cr.Overtake(network, 1)

	cr.Overtake(popularity, 1)

	return *cr
}

func Community(m *model.DataModel, c configuration.CoresConfig) model.Core {

	cr := model.NewCore(model.Community)

	thirdPartyParticipation := ThirdPartyParticipation(m, c.ThirdPartyParticipation)

	organizationalBackup := OrganizationalBackup(m, c.OrgBackup)

	cr.Overtake(organizationalBackup, 3)

	cr.Overtake(thirdPartyParticipation, 1)

	return *cr
}

func Support(m *model.DataModel, c configuration.CoresConfig) model.Core {

	cr := model.NewCore(model.Support)

	processing := Processing(m, c.Processing)

	engagement := Engagement(m, c.Engagement)

	cr.Overtake(processing, 2)

	cr.Overtake(engagement, 1)
	return *cr
}

func Circumstances(m *model.DataModel) model.Core {

	cr := model.NewCore(model.Circumstances)

	rivalry := Rivalry(m)

	licensing := Licensing(m)

	cr.Overtake(rivalry, 1)
	cr.Overtake(licensing, 2)

	return *cr
}
