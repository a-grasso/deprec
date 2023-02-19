package cores

import (
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/model"
	"strings"
)

func Licensing(m *model.DataModel, c configuration.Licensing) model.Core {
	cr := model.NewCore(model.Licensing)

	licenses := map[string]float64{
		"mit":    model.NC,
		"apache": model.NC,
		"isc":    model.NC,
		"wtfpl":  model.NC,
		"bsd":    model.NC,
		"gpl":    model.W,
	}

	if m.Repository != nil {
		license := m.Repository.License
		license = strings.ToLower(license)

		for l, f := range licenses {

			if strings.Contains(license, l) {
				cr.Intake(f, c.Weights.Repository)
			}
		}
	}

	if m.Distribution != nil {

		if m.Distribution.Artifact != nil {
			artifactLicenses := m.Distribution.Artifact.Licenses
			for _, license := range artifactLicenses {
				license = strings.ToLower(license)
				for l, f := range licenses {
					if strings.Contains(license, l) {
						cr.Intake(f, c.Weights.Artifact)
					}
				}
			}
		}

		if m.Distribution.Library != nil {
			libraryLicenses := m.Distribution.Library.Licenses
			for _, license := range libraryLicenses {
				for l, f := range licenses {
					if strings.Contains(license, l) {
						cr.Intake(f, c.Weights.Library)
					}
				}
			}
		}
	}

	return *cr
}
