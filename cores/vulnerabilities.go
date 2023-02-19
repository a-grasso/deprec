package cores

import (
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/model"
)

func Vulnerabilities(m *model.DataModel, c configuration.Vulnerabilities) model.Core {

	cr := model.NewCore(model.Vulnerabilities)

	if m.VulnerabilityIndex == nil {
		return *cr
	}

	vulnerabilities := m.VulnerabilityIndex.TotalVulnerabilitiesCount

	if vulnerabilities > 0 {
		cr.Intake(model.DM, c.Weights.CVE)
	}

	return *cr
}
