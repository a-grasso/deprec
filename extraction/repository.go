package extraction

import (
	"context"
	"deprec/configuration"
	"deprec/model"
	"github.com/google/go-github/v48/github"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2"
	"log"
	"strings"
)

type GitHubExtractor struct {
	Repository string
	Config     *configuration.Configuration
	Client     *GitHubClientWrapper
}

func NewGitHubExtractor(dependency *model.Dependency, config *configuration.Configuration) *GitHubExtractor {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.APIToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	credentials := options.Credential{
		Username: config.Username,
		Password: config.Password,
	}

	clientOpts := options.Client().ApplyURI(config.URI).SetAuth(credentials)
	cache, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Fatal(err)
	}

	ghWrapper := NewGitHubClientWrapper(client, cache)
	return &GitHubExtractor{Repository: dependency.MetaData["vcs"], Config: config, Client: ghWrapper}
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

	contributors, err := re.Client.Repositories.ListContributors(context.TODO(), owner, repo, &github.ListContributorsOptions{})

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

	repos, err := re.Client.Repositories.List(context.TODO(), user, &github.RepositoryListOptions{})
	if err != nil {
		return 0
	}

	return len(repos)
}
