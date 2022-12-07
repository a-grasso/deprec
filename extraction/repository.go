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

	repositoryData.TotalContributors = len(contributors)
	repositoryData.TotalCommits = len(commits)

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

	org := repository.GetOrganization().GetLogin()

	repositoryData := &model.RepositoryData{
		Owner:              repository.GetOwner().GetLogin(),
		Org:                org,
		CreatedAt:          repository.GetCreatedAt().Time,
		Size:               repository.GetSize(),
		License:            repository.GetLicense().GetName(),
		AllowForking:       repository.GetAllowForking(),
		ReadMe:             readme,
		About:              repository.GetDescription(),
		Archivation:        repository.GetArchived(),
		Disabled:           repository.GetDisabled(),
		KLOC:               0,
		TotalCommits:       0,
		TotalIssues:        0,
		TotalPRs:           0,
		TotalContributors:  0,
		Forks:              repository.GetForksCount(),
		Watchers:           repository.GetSubscribersCount(),
		Stars:              repository.GetStargazersCount(),
		Dependencies:       nil,
		Dependents:         nil,
		CommunityStandards: 0,
	}

	return repositoryData
}

func (ghe *GitHubExtractor) extractReadMe() string {
	readme, err := ghe.Client.Repositories.GetReadMe(context.TODO(), ghe.Owner, ghe.Repository, &github.RepositoryContentGetOptions{})
	if err != nil {
		return ""
	}

	readmeContent, err := readme.GetContent()
	if err != nil {
		return ""
	}

	return readmeContent
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
			Committer:    c.GetCommitter().GetLogin(),
			Changes:      nil,
			ChangedFiles: changedFiles,
			Type:         "",
			Message:      c.GetCommit().GetMessage(),
			Branch:       "",
			Timestamp:    c.GetCommit().GetCommitter().GetDate(),
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
			Name:               c.GetLogin(),
			Sponsors:           nil,
			Organizations:      orgs,
			Contributions:      c.GetContributions(),
			Repositories:       projects,
			FirstContribution:  firstContribution,
			LastContribution:   lastContribution,
			TotalContributions: total,
		}

		result = append(result, &contributor)
	}

	return result
}

func siftContributorStats(stats *github.ContributorStats) (string, string, int) {

	if stats == nil {
		return "", "", 0
	}

	var first, last string
	for i, week := range stats.Weeks {
		if i == 0 {
			tmp := week.GetWeek().String()
			first = tmp
		}

		if i == len(stats.Weeks)-1 {
			tmp := week.Week.String()
			last = tmp
		}
	}

	return first, last, stats.GetTotal()
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
