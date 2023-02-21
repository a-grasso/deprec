package cores

import (
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/model"
)

func Rivalry(m model.DataModel, c configuration.Rivalry) model.Core {
	cr := model.NewCore(model.Rivalry)

	if m.Distribution == nil {
		return *cr
	}

	if m.Distribution.Artifact == nil || m.Distribution.Library == nil {
		return *cr
	}

	cr.Intake(artifactIsLatest(m), c.Weights.IsLatest)

	return *cr
}

func artifactIsLatest(m model.DataModel) float64 {

	artifact := m.Distribution.Artifact.Version

	latestVersion := m.Distribution.Library.LatestVersion

	if artifact == latestVersion {
		return 1
	} else {
		return 0
	}
}
