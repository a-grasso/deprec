package cores

import (
	"deprec/model"
	"sort"
)

func OrganizationalBackup(m *model.DataModel) model.CoreResult {

	cr := model.NewCoreResult(model.OrganizationalBackup)

	contributors := m.Repository.Contributors

	sort.Slice(contributors, func(i, j int) bool {
		return contributors[i].Contributions > contributors[j].Contributions
	})

	total := 0.0
	for i, contributor := range contributors {

		weight := float64(i) / float64(len(contributors))

		sponsors := float64(contributor.Sponsors) * weight

		orgs := float64(contributor.Organizations) * weight

		total += sponsors + orgs
	}

	total /= float64(len(contributors))

	cr.Intake(total, 1)

	if m.Repository.Org != nil {
		cr.Intake(1, 3)
	} else {
		cr.Intake(0, 3)
	}

	return cr
}
