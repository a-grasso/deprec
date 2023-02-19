package cores

import (
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/model"
)

func ProjectQuality(m *model.DataModel, c configuration.ProjectQuality) model.Core {
	cr := model.NewCore(model.ProjectQuality)

	if m.Repository != nil {

		if m.Repository.ReadMe != "" {
			cr.Intake(model.NIA, c.Weights.ReadMe)
		}

		if m.Repository.License != "" {
			cr.Intake(model.NIA, c.Weights.License)
		}

		if m.Repository.About != "" {
			cr.Intake(model.NIA, c.Weights.About)
		}

		if m.Repository.AllowForking {
			cr.Intake(model.NIA, c.Weights.AllowForking)
		}
	}

	return *cr
}
