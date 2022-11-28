package extraction

import (
	"context"
	"deprec/configuration"
	"deprec/model"
	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
	"strings"
)

type GitHubExtractor struct {
	Repository string
	DataModel  *model.DataModel
	Config     configuration.GitHub
	ghClient   *github.Client
}

func NewGitHubExtractor(dependency *model.Dependency, dataModel *model.DataModel, gh configuration.GitHub) *GitHubExtractor {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gh.APIToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return &GitHubExtractor{Repository: dependency.MetaData["vcs"], Config: gh, DataModel: dataModel, ghClient: client}
}

func (re *GitHubExtractor) Extract(dataModel *model.DataModel) {

	contributors := re.extractContributors()

	repository := &model.Repository{Contributors: contributors}

	dataModel.Repository = repository
}

func (re *GitHubExtractor) parseVCSString(vcs string) (string, string) {
	splits := strings.Split(vcs, "/")
	return splits[3], splits[4]
}

func (re *GitHubExtractor) extractContributors() []*model.Contributor {

	owner, repo := re.parseVCSString(re.Repository)

	var allContributors []*github.Contributor
	opt := &github.ListContributorsOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	for {
		contributors, r, _ := re.ghClient.Repositories.ListContributors(context.TODO(), owner, repo, opt)
		allContributors = append(allContributors, contributors...)
		if r.NextPage == 0 {
			break
		}
		opt.Page = r.NextPage
	}

	var result []*model.Contributor
	for _, c := range allContributors {
		projects := re.listContributorRepositories(c.GetLogin())
		result = append(result, &model.Contributor{
			Name:              c.GetLogin(),
			Sponsors:          nil,
			Organizations:     nil,
			Repositories:      projects,
			FirstContribution: "",
			LastContribution:  "",
		})
	}

	return result
}

func (re *GitHubExtractor) listContributorRepositories(user string) int {
	projects, _, _ := re.ghClient.Repositories.List(context.TODO(), user, nil)
	return len(projects)
}
