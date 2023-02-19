package cores

import (
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/model"
	"github.com/thoas/go-funk"
)

func Engagement(m *model.DataModel, c configuration.Engagement) model.Core {

	cr := model.NewCore(model.Engagement)

	if m.Repository == nil {
		return *cr
	}

	totalIssues := len(m.Repository.Issues)

	totalComments := funk.Sum(funk.Map(m.Repository.Issues, func(i model.Issue) int { return len(i.Contributions) }))

	ratio := totalComments / float64(totalIssues)

	ratio *= 100

	cr.IntakeThreshold(ratio, c.IssueCommentsRatioThresholdPercentage, c.Weights.IssueCommentsRatio)

	return *cr
}
