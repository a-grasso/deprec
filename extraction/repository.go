package extraction

import (
	"context"
	"deprec/configuration"
	"deprec/logging"
	"deprec/model"
	"github.com/google/go-github/v48/github"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2"
	"log"
	"strings"
)

type GitHubExtractor struct {
	RepositoryURL string
	Repository    string
	Owner         string
	Config        *configuration.Configuration
	Client        *GitHubClientWrapper
}

func NewGitHubExtractor(dependency *model.Dependency, config *configuration.Configuration) *GitHubExtractor {

	client := githubClient(config)

	cache := mongoDBClient(config)

	ghWrapper := NewGitHubClientWrapper(client, cache)

	vcs := dependency.MetaData["vcs"]
	owner, repo := parseVCSString(vcs)

	return &GitHubExtractor{RepositoryURL: vcs, Owner: owner, Repository: repo, Config: config, Client: ghWrapper}
}

func checkRateLimits(client *github.Client) {
	limits, _, _ := client.RateLimits(context.TODO())
	logging.SugaredLogger.Infof("rate limit:-> Core: %d Search: %d", limits.Core.Remaining, limits.Search.Remaining)
}

func parseVCSString(vcs string) (string, string) {
	splits := strings.Split(vcs, ".git")
	splits = strings.Split(splits[0], "/")
	return splits[3], splits[4]
}

func mongoDBClient(config *configuration.Configuration) *mongo.Client {
	credentials := options.Credential{
		Username: config.Username,
		Password: config.Password,
	}

	clientOpts := options.Client().ApplyURI(config.URI).SetAuth(credentials)
	cache, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		logging.SugaredLogger.Fatalf("connecting to mongodb database at '%s': %s", config.URI, err)
		log.Fatal(err)
	}
	return cache
}

func githubClient(config *configuration.Configuration) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.APIToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return client
}

func (ghe *GitHubExtractor) Extract(dataModel *model.DataModel) {
	logging.SugaredLogger.Infof("extracting repo '%s'", ghe.RepositoryURL)

	checkRateLimits(ghe.Client.client)

	repositoryData := ghe.extractRepositoryData()

	contributors := ghe.extractContributors()

	commits := ghe.extractCommits()

	*repositoryData.TotalContributors = len(contributors)
	*repositoryData.TotalCommits = len(commits)

	repository := &model.Repository{
		Contributors:   contributors,
		Issues:         nil,
		Commits:        commits,
		Releases:       nil,
		RepositoryData: repositoryData,
	}

	dataModel.Repository = repository

	checkRateLimits(ghe.Client.client)
}

func (ghe *GitHubExtractor) extractRepositoryData() *model.RepositoryData {
	repository, err := ghe.Client.Repositories.Get(context.TODO(), ghe.Owner, ghe.Repository)
	if err != nil {
		logging.SugaredLogger.Errorf("could not extract repository data of '%s' : %s", ghe.RepositoryURL, err)
		return nil
	}

	readme := ghe.extractReadMe()

	var org *string
	if repository.GetOrganization() == nil {
		org = nil
	} else {
		org = repository.GetOrganization().Login
	}

	repositoryData := &model.RepositoryData{
		Owner:              repository.GetOwner().Login,
		Org:                org,
		CreatedAt:          &repository.CreatedAt.Time,
		Size:               repository.Size,
		License:            repository.GetLicense().Name,
		AllowForking:       repository.AllowForking,
		ReadMe:             readme,
		About:              repository.Description,
		Archivation:        repository.Archived,
		Disabled:           repository.Disabled,
		KLOC:               new(int),
		TotalCommits:       new(int),
		TotalIssues:        new(int),
		TotalPRs:           new(int),
		TotalContributors:  new(int),
		Forks:              repository.ForksCount,
		Watchers:           repository.SubscribersCount,
		Stars:              repository.StargazersCount,
		Dependencies:       nil,
		Dependents:         nil,
		CommunityStandards: new(float64),
	}

	return repositoryData
}

func (ghe *GitHubExtractor) extractReadMe() *string {
	readme, err := ghe.Client.Repositories.GetReadMe(context.TODO(), ghe.Owner, ghe.Repository, &github.RepositoryContentGetOptions{})
	if err != nil {
		return nil
	}

	readmeContent, err := readme.GetContent()
	if err != nil {
		return nil
	}

	return &readmeContent
}

func (ghe *GitHubExtractor) extractCommits() []*model.Commit {
	commits, err := ghe.Client.Repositories.ListCommits(context.TODO(), ghe.Owner, ghe.Repository, &github.CommitsListOptions{})
	if err != nil {
		return nil
	}

	var result []*model.Commit

	for _, c := range commits {

		var changedFiles []string
		files := c.Files
		for _, f := range files {
			changedFiles = append(changedFiles, f.GetFilename())
		}

		commit := &model.Commit{
			Author:       c.GetAuthor().GetLogin(),
			Changes:      nil,
			ChangedFiles: changedFiles,
			Type:         "",
			Message:      c.GetCommit().GetMessage(),
			Branch:       "",
			Timestamp:    c.GetCommit().GetAuthor().GetDate(),
			Additions:    c.GetCommit().GetStats().GetAdditions(),
			Deletions:    c.GetCommit().GetStats().GetDeletions(),
			Total:        c.GetCommit().GetStats().GetTotal(),
		}

		result = append(result, commit)
	}

	return result
}

func (ghe *GitHubExtractor) extractContributors() []*model.Contributor {

	contributors, err := ghe.Client.Repositories.ListContributors(context.TODO(), ghe.Owner, ghe.Repository, &github.ListContributorsOptions{})
	if err != nil {
		return []*model.Contributor{}
	}

	var result []*model.Contributor
	for _, c := range contributors {

		user := c.GetLogin()
		projects := len(ghe.listContributorRepositories(user))
		stats := ghe.listContributorStats(user)
		firstContribution, lastContribution, total := siftContributorStats(stats)

		orgs := len(ghe.listContributorOrganizations(user))
		contributor := model.Contributor{
			Name:               c.Login,
			Sponsors:           nil,
			Organizations:      &orgs,
			Contributions:      c.Contributions,
			Repositories:       &projects,
			FirstContribution:  firstContribution,
			LastContribution:   lastContribution,
			TotalContributions: total,
		}

		result = append(result, &contributor)
	}

	return result
}

func siftContributorStats(stats *github.ContributorStats) (*string, *string, *int) {

	if stats == nil {
		return nil, nil, nil
	}

	var first, last *string
	for i, week := range stats.Weeks {
		if i == 0 {
			tmp := week.Week.String()
			first = &tmp
		}

		if i == len(stats.Weeks)-1 {
			tmp := week.Week.String()
			last = &tmp
		}
	}

	return first, last, stats.Total
}

func (ghe *GitHubExtractor) listContributorStats(user string) *github.ContributorStats {
	contributorStats, err := ghe.Client.Repositories.ListContributorStats(context.TODO(), ghe.Owner, ghe.Repository)

	if err != nil {
		return nil
	}

	var result *github.ContributorStats
	for _, stat := range contributorStats {
		if user == stat.GetAuthor().GetLogin() {
			result = stat
		}
	}

	return result
}

func (ghe *GitHubExtractor) listContributorRepositories(user string) []*github.Repository {

	repos, err := ghe.Client.Repositories.List(context.TODO(), user, &github.RepositoryListOptions{})
	if err != nil {
		return nil
	}

	return repos
}

func (ghe *GitHubExtractor) listContributorOrganizations(user string) []*github.Organization {

	orgs, err := ghe.Client.Organizations.List(context.TODO(), user, &github.ListOptions{})
	if err != nil {
		return nil
	}

	return orgs
}
