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
	"time"
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

	gitHubClientWrapper := NewGitHubClientWrapper(client, cache)

	vcs := dependency.MetaData["vcs"]
	owner, repo := parseVCSString(vcs)

	return &GitHubExtractor{RepositoryURL: vcs, Owner: owner, Repository: repo, Config: config, Client: gitHubClientWrapper}
}

func (ghe *GitHubExtractor) checkRateLimits() {
	limits, _, _ := ghe.Client.client.RateLimits(context.TODO())
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

	ghe.checkRateLimits()

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

	ghe.checkRateLimits()
}

func (ghe *GitHubExtractor) extractRepositoryData() *model.RepositoryData {
	repository, err := ghe.Client.Repositories.Get(context.TODO(), ghe.Owner, ghe.Repository)
	if err != nil {
		logging.SugaredLogger.Errorf("could not extract repository data of '%s' : %s", ghe.RepositoryURL, err)
		return nil
	}

	readme := ghe.extractReadMe()

	org := ghe.extractOrganization(repository.GetOrganization().GetLogin())

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
		logging.SugaredLogger.Errorf("could not extract readme of '%s' : %s", ghe.RepositoryURL, err)
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
		logging.SugaredLogger.Errorf("could not extract commits of '%s' : %s", ghe.RepositoryURL, err)
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
		logging.SugaredLogger.Errorf("could not extract contributors of '%s' : %s", ghe.RepositoryURL, err)
		return nil
	}

	var result []*model.Contributor
	contributorStats := ghe.listContributorStats()
	for _, c := range contributors {

		user := c.GetLogin()
		projects := len(ghe.listContributorRepositories(user))
		firstContribution, lastContribution, total := ghe.siftContributorStats(contributorStats, user)

		orgs := len(ghe.listContributorOrganizations(user))
		contributor := model.Contributor{
			Name:                    c.GetLogin(),
			Sponsors:                nil,
			Organizations:           orgs,
			Contributions:           c.GetContributions(),
			Repositories:            projects,
			FirstContribution:       firstContribution,
			LastContribution:        lastContribution,
			TotalStatsContributions: total,
		}

		result = append(result, &contributor)
	}

	return result
}

func (ghe *GitHubExtractor) siftContributorStats(contributorStats []*github.ContributorStats, user string) (time.Time, time.Time, int) {
	var stats *github.ContributorStats
	for _, cs := range contributorStats {
		if user == cs.GetAuthor().GetLogin() {
			stats = cs
		}
	}

	if stats == nil {
		logging.SugaredLogger.Errorf("could not find stats of contributor '%s' from repo '%s'", user, ghe.RepositoryURL)
		return time.Time{}, time.Time{}, 0
	}

	var first, last time.Time
	for i, week := range stats.Weeks {
		if i == 0 {
			tmp := week.GetWeek().Time
			first = tmp
		}

		if i == len(stats.Weeks)-1 {
			tmp := week.Week.Time
			last = tmp
		}
	}

	return first, last, stats.GetTotal()
}

func (ghe *GitHubExtractor) listContributorStats() []*github.ContributorStats {
	contributorStats, err := ghe.Client.Repositories.ListContributorStats(context.TODO(), ghe.Owner, ghe.Repository)

	if err != nil {
		logging.SugaredLogger.Errorf("could not extract stats of contributors from repo '%s' : %s", ghe.RepositoryURL, err)
		return nil
	}

	return contributorStats
}

func (ghe *GitHubExtractor) listContributorRepositories(user string) []*github.Repository {

	repos, err := ghe.Client.Repositories.List(context.TODO(), user, &github.RepositoryListOptions{})
	if err != nil {
		logging.SugaredLogger.Errorf("could not list repositories of contributor '%s' : %s", user, err)
		return nil
	}

	return repos
}

func (ghe *GitHubExtractor) listContributorOrganizations(user string) []*github.Organization {

	orgs, err := ghe.Client.Organizations.List(context.TODO(), user, &github.ListOptions{})
	if err != nil {
		logging.SugaredLogger.Errorf("could not list organizations of contributor '%s' : %s", user, err)
		return nil
	}

	return orgs
}

func (ghe *GitHubExtractor) extractOrganization(o string) *model.Organization {

	org, err := ghe.Client.Organizations.Get(context.TODO(), o)

	if err != nil {
		logging.SugaredLogger.Errorf("could not extract organization data of '%s' : %s", o, err)
		return nil
	}

	organization := &model.Organization{
		Login:             org.GetLogin(),
		PublicRepos:       org.GetPublicRepos(),
		Followers:         org.GetFollowers(),
		Following:         org.GetFollowing(),
		TotalPrivateRepos: org.GetTotalPrivateRepos(),
		OwnedPrivateRepos: org.GetOwnedPrivateRepos(),
		Collaborators:     org.GetCollaborators(),
	}

	return organization
}
