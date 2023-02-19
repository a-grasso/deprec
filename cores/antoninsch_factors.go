package cores

import (
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/model"
)

func DeityGiven(m *model.DataModel, c configuration.CoresConfig) model.Core {

	cr := model.NewCore(model.DeityGiven)

	marking := Marking(m, c.Marking)

	vulnerabilities := Vulnerabilities(m, c.Vulnerabilities)

	cr.Overtake(marking, c.DeityGiven.Weights.Marking)

	cr.Overtake(vulnerabilities, c.DeityGiven.Weights.Vulnerabilities)

	return *cr
}

func Effort(m *model.DataModel, c configuration.CoresConfig) model.Core {

	cr := model.NewCore(model.Effort)

	activity := Activity(m, c.Activity)

	recentness := Recentness(m, c.Recentness)

	coreTeam := CoreTeam(m, c.CoreTeam)

	cr.Overtake(recentness, c.Effort.Weights.Recentness)
	cr.Overtake(activity, c.Effort.Weights.Activity)
	cr.Overtake(coreTeam, c.Effort.Weights.CoreTeam)

	return *cr
}

func Interconnectedness(m *model.DataModel, c configuration.CoresConfig) model.Core {

	cr := model.NewCore(model.Interconnectedness)

	network := Network(m, c.Network)

	popularity := Popularity(m, c.Popularity)

	cr.Overtake(network, c.Interconnectedness.Weights.Network)

	cr.Overtake(popularity, c.Interconnectedness.Weights.Popularity)

	return *cr
}

func Community(m *model.DataModel, c configuration.CoresConfig) model.Core {

	cr := model.NewCore(model.Community)

	participation := Participation(m, c.Participation)

	backup := Backup(m, c.Backup)

	prestige := Prestige(m, c.Prestige)

	cr.Overtake(prestige, c.Community.Weights.Prestige)

	cr.Overtake(backup, c.Community.Weights.Backup)

	cr.Overtake(participation, c.Community.Weights.Participation)

	return *cr
}

func Support(m *model.DataModel, c configuration.CoresConfig) model.Core {

	cr := model.NewCore(model.Support)

	processing := Processing(m, c.Processing)

	engagement := Engagement(m, c.Engagement)

	cr.Overtake(processing, c.Support.Weights.Processing)

	cr.Overtake(engagement, c.Support.Weights.Engagement)
	return *cr
}

func Circumstances(m *model.DataModel, c configuration.CoresConfig) model.Core {

	cr := model.NewCore(model.Circumstances)

	rivalry := Rivalry(m, c.Rivalry)

	quality := ProjectQuality(m, c.ProjectQuality)

	licensing := Licensing(m, c.Licensing)

	cr.Overtake(rivalry, 1)
	cr.Overtake(licensing, 2)
	cr.Overtake(quality, 1)

	return *cr
}
