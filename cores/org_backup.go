package cores

import (
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/model"
	"github.com/thoas/go-funk"
)

func OrganizationalBackup(m *model.DataModel, c configuration.OrgBackup) model.Core {

	cr := model.NewCoreResult(model.OrganizationalBackup)

	if m.Repository == nil {
		return *cr
	}

	contributors := m.Repository.Contributors

	if len(contributors) == 0 {
		return *cr
	}

	companies := funk.Filter(funk.Map(contributors, func(c model.Contributor) string { return c.Company }), func(company string) bool { return company != "" }).([]string)
	sponsors := funk.Sum(funk.Map(contributors, func(c model.Contributor) int { return c.Sponsors }))
	organizations := funk.Sum(funk.Map(contributors, func(c model.Contributor) int { return c.Organizations }))

	cr.IntakeThreshold(float64(len(companies)), float64(c.CompanyThreshold), 2)
	cr.IntakeThreshold(sponsors, c.SponsorThreshold, 1)
	cr.IntakeThreshold(organizations, c.OrganizationThreshold, 2)

	if m.Repository.Org != nil {
		cr.Intake(1, 3)
	} else {
		cr.Intake(0, 3)
	}

	return *cr
}
