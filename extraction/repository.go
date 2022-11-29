package extraction

import (
	"context"
	"deprec/configuration"
	"deprec/model"
	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
	"log"
	"strings"
)

type GitHubExtractor struct {
	Repository string
	Config     configuration.GitHub
	Client     *GitHubClientWrapper
}

func NewGitHubExtractor(dependency *model.Dependency, gh configuration.GitHub) *GitHubExtractor {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gh.APIToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	ghWrapper := NewGitHubClientWrapper(client)
	return &GitHubExtractor{Repository: dependency.MetaData["vcs"], Config: gh, Client: ghWrapper}
}

func (re *GitHubExtractor) Extract(dataModel *model.DataModel) {

	contributors := re.extractContributors()

	repository := &model.Repository{Contributors: contributors}

	dataModel.Repository = repository
}

func (re *GitHubExtractor) parseVCSString(vcs string) (string, string) {
	splits := strings.Split(vcs, ".git")
	splits = strings.Split(splits[0], "/")
	return splits[3], splits[4]
}

func (re *GitHubExtractor) extractContributors() []*model.Contributor {
	log.Printf("Extracting Contributor of repo %s", re.Repository)

	owner, repo := re.parseVCSString(re.Repository)

	contributors, err := re.Client.RepositoriesListContributors(context.TODO(), owner, repo, nil)

	if err != nil {
		return nil
	}

	var result []*model.Contributor
	for _, c := range contributors {
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
	log.Printf("Listing Repos of Contributor %s", user)

	var allRepos []*github.Repository
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	for {
		projects, r, err := re.Client.RepositoriesList(context.TODO(), user, opt)

		if err != nil {
			break
		}
		allRepos = append(allRepos, projects...)
		if r.NextPage == 0 {
			break
		}
		opt.Page = r.NextPage
	}

	return len(allRepos)
}
