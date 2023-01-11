package cores

import (
	"deprec/configuration"
	"deprec/model"
	"github.com/thoas/go-funk"
	"sort"
)

func OrganizationalBackup(m *model.DataModel, c configuration.OrgBackup) model.CoreResult {

	cr := model.NewCoreResult(model.OrganizationalBackup)

	contributors := m.Repository.Contributors

	if len(contributors) == 0 {
		return *cr
	}

	sort.Slice(contributors, func(i, j int) bool {
		return contributors[i].Contributions > contributors[j].Contributions
	})

	companies := funk.Filter(funk.Map(contributors, func(c model.Contributor) string { return c.Company }), func(c string) bool { return c != "" }).([]string)
	sponsors := funk.Sum(funk.Map(contributors, func(c model.Contributor) int { return c.Sponsors }))
	organizations := funk.Sum(funk.Map(contributors, func(c model.Contributor) int { return c.Organizations }))

	cr.IntakeThreshold(float64(len(companies)), float64(c.CompanyThreshold), 2)
	cr.IntakeThreshold(sponsors, c.SponsorThreshold, 1)
	cr.IntakeThreshold(organizations, c.OrganizationThreshold, 1)

	if m.Repository.Org != nil {
		cr.Intake(1, 3)
	} else {
		cr.Intake(0, 3)
	}

	return *cr
}
