package cores

import (
	"github.com/a-grasso/deprec/model"
)

func Rivalry(m *model.DataModel) model.Core {
	cr := model.NewCoreResult(model.Rivalry)

	if m.Distribution == nil {
		return *cr
	}

	if m.Distribution.Artifact == nil || m.Distribution.Library == nil {
		return *cr
	}

	cr.Intake(artifactIsLatest(m), 1)

	return *cr
}

func artifactIsLatest(m *model.DataModel) float64 {

	artifact := m.Distribution.Artifact.Version

	latestVersion := m.Distribution.Library.LatestVersion

	if artifact == latestVersion {
		return 1
	} else {
		return 0
	}
}
