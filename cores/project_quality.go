package cores

import "github.com/a-grasso/deprec/model"

func ProjectQuality(m *model.DataModel) model.Core {
	cr := model.NewCore(model.ProjectQuality)

	if m.Repository != nil {

		if m.Repository.ReadMe != "" {
			cr.Intake(model.NIA, 1)
		}

		if m.Repository.License != "" {
			cr.Intake(model.NIA, 1)
		}

		if m.Repository.About != "" {
			cr.Intake(model.NIA, 1)
		}

		if m.Repository.AllowForking {
			cr.Intake(model.NIA, 1)
		}
	}

	return *cr
}
