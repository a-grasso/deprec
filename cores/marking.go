package cores

import (
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/model"
	"strings"
)

func Marking(m *model.DataModel, c configuration.Marking) model.Core {

	cr := model.NewCore(model.Marking)

	if m.Repository != nil {
		archived := m.Repository.Archivation
		if archived {
			cr.Intake(model.DM, c.Weights.Archivation)
		}

		readme := strings.ToLower(m.Repository.ReadMe)
		for _, keyword := range c.ReadMeKeywords {
			if strings.Contains(readme, keyword) {
				cr.Intake(model.DM, c.Weights.ReadMe)
			}
		}

		about := strings.ToLower(m.Repository.About)
		for _, keyword := range c.AboutKeywords {
			if strings.Contains(about, keyword) {
				cr.Intake(model.DM, c.Weights.About)
			}
		}
	}

	if distribution := m.Distribution; distribution != nil {

		if artifact := distribution.Artifact; artifact != nil {

			description := strings.ToLower(artifact.Description)

			for _, keyword := range c.ArtifactDescriptionKeywords {
				if strings.Contains(description, keyword) {
					cr.Intake(model.DM, c.Weights.ReadMe)
				}
			}
		}
	}

	return *cr
}
