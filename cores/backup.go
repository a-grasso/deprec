package cores

import (
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/model"
	"github.com/thoas/go-funk"
)

func Backup(m model.DataModel, c configuration.Backup) model.Core {

	cr := model.NewCore(model.Backup)

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

	cr.IntakeThreshold(float64(len(companies)), float64(c.CompanyThreshold), c.Weights.Companies)
	cr.IntakeThreshold(sponsors, c.SponsorThreshold, c.Weights.Sponsors)
	cr.IntakeThreshold(organizations, c.OrganizationThreshold, c.Weights.Organizations)

	if m.Repository.Org != nil {
		cr.Intake(model.NC, c.Weights.RepositoryOrganization)
	} else {
		cr.Intake(model.W, c.Weights.RepositoryOrganization)
	}

	return *cr
}
