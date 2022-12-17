package mapping

import (
	"deprec/configuration"
	"deprec/model"
	"deprec/statistics"
	"github.com/thoas/go-funk"
	"log"
	"time"
)

type IssueContributions struct {
	time time.Time
}

func (ic IssueContributions) GetTimeStamp() time.Time {
	return ic.time
}

func Activity(m *model.DataModel, config configuration.Activity) float64 {

	commits := m.Repository.Commits

	commitAnalysis := statistics.AnalyzeCount(commits, config.Percentile)

	if commitAnalysis.MonthsSinceLast > config.CommitThreshold {
		return 0
	}

	issues := m.Repository.Issues
	issueAnalysis := statistics.AnalyzeCount(issues, config.Percentile)
	log.Println(issueAnalysis.String())

	closedIssueAnalysis := statistics.AnalyzeCount(issues, config.Percentile)
	log.Println(closedIssueAnalysis.String())

	issueContributions := funk.FlatMap(issues, func(issue model.Issue) []IssueContributions {
		var result []IssueContributions
		for i := 0; i < issue.Contributions; i++ {
			result = append(result, IssueContributions{time: issue.CreationTime})
		}
		return result
	}).([]IssueContributions)
	issueContributionAnalysis := statistics.AnalyzeCount(issueContributions, config.Percentile)
	log.Println(issueContributionAnalysis.String())

	//releases := m.Repository.Releases
	//log.Println(AnalyzeCount(releases))

	return 0
}
