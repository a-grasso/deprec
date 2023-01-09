package cores

import (
	"deprec/configuration"
	"deprec/model"
	"github.com/thoas/go-funk"
)

func Engagement(m *model.DataModel, c configuration.Engagement) model.CoreResult {

	cr := model.NewCoreResult(model.Engagement)

	totalIssues := len(m.Repository.Issues)

	totalComments := funk.Sum(funk.Map(m.Repository.Issues, func(i model.Issue) int { return len(i.Contributions) }))

	ratio := totalComments / float64(totalIssues)

	ratio *= 100

	cr.IntakeThreshold(ratio, c.IssueCommentsRatioThresholdPercentage, 1)

	return cr

}
