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
	limits, _, err := ghe.Client.client.RateLimits(context.TODO())
	if err != nil {
		logging.SugaredLogger.Debugf("could not check rate limit for github rest api :%s", err)
		return
	}
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
	// TODO: Check connection (Ping)
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

	repositoryData := ghe.extractRepositoryData(ghe.Owner, ghe.Repository)

	contributors := ghe.extractContributors(ghe.Owner, ghe.Repository)

	commits := ghe.extractCommits(ghe.Owner, ghe.Repository)

	releases := ghe.extractReleases(ghe.Owner, ghe.Repository)

	var tags []model.Tag
	if releases == nil {
		tags = ghe.extractTags(ghe.Owner, ghe.Repository)
	}

	issues := ghe.extractIssues(ghe.Owner, ghe.Repository)

	repository := &model.Repository{
		Contributors:   contributors,
		Issues:         issues,
		Commits:        commits,
		Releases:       releases,
		Tags:           tags,
		RepositoryData: repositoryData,
	}

	dataModel.Repository = repository

	ghe.checkRateLimits()
}

func (ghe *GitHubExtractor) extractRepositoryData(owner, repo string) *model.RepositoryData {
	repository, err := ghe.Client.Repositories.Get(context.TODO(), owner, repo)
	if err != nil {
		logging.SugaredLogger.Debugf("could not extract repository data of '%s' : %s", ghe.RepositoryURL, err)
		return nil
	}

	readme := ghe.extractReadMe(owner, repo)

	org := ghe.extractOrganization(repository.GetOrganization().GetLogin())

	repositoryData := &model.RepositoryData{
		Owner:              repository.GetOwner().GetLogin(),
		Org:                org,
		CreatedAt:          repository.GetCreatedAt().Time,
		Size:               repository.GetSize(),
		License:            repository.GetLicense().GetKey(),
		AllowForking:       repository.GetAllowForking(),
		ReadMe:             readme,
		About:              repository.GetDescription(),
		Archivation:        repository.GetArchived(),
		Disabled:           repository.GetDisabled(),
		KLOC:               0,
		TotalPRs:           0,
		Forks:              repository.GetForksCount(),
		Watchers:           repository.GetSubscribersCount(),
		Stars:              repository.GetStargazersCount(),
		Dependencies:       nil,
		Dependents:         nil,
		CommunityStandards: 0,
	}

	return repositoryData
}

func (ghe *GitHubExtractor) extractReadMe(owner, repo string) string {
	readme, err := ghe.Client.Repositories.GetReadMe(context.TODO(), owner, repo, &github.RepositoryContentGetOptions{})
	if err != nil {
		logging.SugaredLogger.Debugf("could not extract readme of '%s' : %s", ghe.RepositoryURL, err)
		return ""
	}

	readmeContent, err := readme.GetContent()
	if err != nil {
		return ""
	}

	return readmeContent
}

func (ghe *GitHubExtractor) extractReleases(owner, repo string) []model.Release {
	releases, err := ghe.Client.Repositories.ListReleases(context.TODO(), owner, repo, &github.ListOptions{})
	if err != nil {
		logging.SugaredLogger.Debugf("could not extract releases of '%s' : %s", ghe.RepositoryURL, err)
		return nil
	}

	var result []model.Release

	for _, release := range releases {

		r := model.Release{
			Author:      release.GetAuthor().GetLogin(),
			Version:     release.GetName(),
			Description: release.GetBody(),
			Changes:     nil,
			Type:        "",
			Date:        release.GetPublishedAt().Time,
		}

		result = append(result, r)
	}

	return result
}

func (ghe *GitHubExtractor) extractTags(owner, repo string) []model.Tag {
	tags, err := ghe.Client.Repositories.ListTags(context.TODO(), owner, repo, &github.ListOptions{})
	if err != nil {
		logging.SugaredLogger.Warnf("could not extract tags of '%s' : %s", ghe.RepositoryURL, err)
		return nil
	}

	var result []model.Tag

	for _, tag := range tags {

		if tag.GetCommit().GetCommitter() == nil {

			tagCommit, err := ghe.Client.Repositories.GetCommit(context.TODO(), owner, repo, tag.GetCommit().GetSHA(), &github.ListOptions{})

			if err != nil {
				continue
			}

			r := model.Tag{
				Author:      tagCommit.GetCommit().GetAuthor().GetEmail(),
				Version:     tag.GetName(),
				Description: tagCommit.GetCommit().GetMessage(),
				Date:        tagCommit.GetCommit().GetCommitter().GetDate(),
			}

			result = append(result, r)

		} else {

			r := model.Tag{
				Author:      tag.GetCommit().GetAuthor().GetLogin(),
				Version:     tag.GetName(),
				Description: tag.GetCommit().GetMessage(),
				Date:        tag.GetCommit().GetCommitter().GetDate(),
			}

			result = append(result, r)
		}
	}
	return result
}

func (ghe *GitHubExtractor) extractIssues(owner, repo string) []model.Issue {
	issues, err := ghe.Client.Issues.ListByRepo(context.TODO(), owner, repo, &github.IssueListByRepoOptions{
		State: "all",
	})
	if err != nil {
		logging.SugaredLogger.Debugf("could not extract issues of '%s' : %s", ghe.RepositoryURL, err)
		return nil
	}

	var result []model.Issue

	for _, issue := range issues {

		//var issueContributors []string
		//var firstResponse time.Time
		// if issue.GetComments() != 0 {
		// 	sort := "created"
		// 	comments, err := ghe.Client.Issues.ListComments(context.TODO(), owner, repo, issue.GetNumber(), &github.IssueListCommentsOptions{
		// 		Sort: &sort,
		// 	})
		//
		// 	if err != nil {
		// 		logging.SugaredLogger.Debugf("could not extract comments of issue '%d' for repo '%s'", issue.Number, repo)
		// 	}
		//
		// 	firstResponse = comments[0].GetCreatedAt()
		//
		// 	commentators := funk.Map(comments, func(comment *github.IssueComment) string { return comment.GetUser().GetLogin() })
		//
		// 	issueContributors = funk.Uniq(commentators).([]string)
		// }

		var contributions []model.IssueContribution
		for i := 0; i < issue.GetComments(); i++ {
			contributions = append(contributions, model.IssueContribution{Time: issue.GetCreatedAt()})
		}

		i := model.Issue{
			Number:            issue.GetNumber(),
			Author:            issue.GetUser().GetLogin(),
			AuthorAssociation: issue.GetAuthorAssociation(),
			Labels:            nil,
			State:             model.ToIssueState(issue.GetState()),
			Title:             issue.GetTitle(),
			Content:           issue.GetBody(),
			ClosedBy:          issue.GetClosedBy().GetLogin(),
			Contributions:     contributions,
			Contributors:      nil,
			CreationTime:      issue.GetCreatedAt(),
			FirstResponse:     time.Time{},
			LastContribution:  issue.GetUpdatedAt(),
			ClosingTime:       issue.GetClosedAt(),
		}

		result = append(result, i)
	}

	return result
}

func (ghe *GitHubExtractor) extractCommits(owner, repo string) []model.Commit {
	commits, err := ghe.Client.Repositories.ListCommits(context.TODO(), owner, repo, &github.CommitsListOptions{})
	if err != nil {
		logging.SugaredLogger.Debugf("could not extract commits of '%s' : %s", ghe.RepositoryURL, err)
		return nil
	}

	var result []model.Commit

	for _, c := range commits {

		var changedFiles []string
		files := c.Files
		for _, f := range files {
			changedFiles = append(changedFiles, f.GetFilename())
		}

		commit := model.Commit{
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

func (ghe *GitHubExtractor) extractContributors(owner, repo string) []model.Contributor {

	contributors, err := ghe.Client.Repositories.ListContributors(context.TODO(), owner, repo, &github.ListContributorsOptions{})

	if err != nil {
		logging.SugaredLogger.Debugf("could not extract contributors of '%s' : %s", ghe.RepositoryURL, err)
		return nil
	}

	var result []model.Contributor
	contributorStats := ghe.listContributorStats(owner, repo)
	for _, c := range contributors {

		user := c.GetLogin()
		repos := len(ghe.listContributorRepositories(user))
		firstContribution, lastContribution, total := ghe.siftContributorStats(contributorStats, user)

		orgs := len(ghe.listContributorOrganizations(user))
		contributor := model.Contributor{
			Name:                    user,
			Sponsors:                0,
			Organizations:           orgs,
			Contributions:           c.GetContributions(),
			Repositories:            repos,
			FirstContribution:       firstContribution,
			LastContribution:        lastContribution,
			TotalStatsContributions: total,
		}

		result = append(result, contributor)
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
		logging.SugaredLogger.Debugf("could not find stats of contributor '%s' from repo '%s'", user, ghe.RepositoryURL)
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

func (ghe *GitHubExtractor) listContributorStats(owner, repo string) []*github.ContributorStats {
	contributorStats, err := ghe.Client.Repositories.ListContributorStats(context.TODO(), owner, repo)

	if err != nil {
		logging.SugaredLogger.Debugf("could not extract stats of contributors from repo '%s' : %s", ghe.RepositoryURL, err)
		return nil
	}

	return contributorStats
}

func (ghe *GitHubExtractor) listContributorRepositories(user string) []*github.Repository {

	repos, err := ghe.Client.Repositories.List(context.TODO(), user, &github.RepositoryListOptions{})
	if err != nil {
		logging.SugaredLogger.Debugf("could not list repositories of contributor '%s' : %s", user, err)
		return nil
	}

	return repos
}

func (ghe *GitHubExtractor) listContributorOrganizations(user string) []*github.Organization {

	organizations, err := ghe.Client.Organizations.List(context.TODO(), user, &github.ListOptions{})
	if err != nil {
		logging.SugaredLogger.Debugf("could not list organizations of contributor '%s' : %s", user, err)
		return nil
	}

	return organizations
}

func (ghe *GitHubExtractor) extractOrganization(o string) *model.Organization {

	if o == "" {
		logging.SugaredLogger.Debug("could not extract organization data of '' : does not exist")
		return nil
	}

	org, err := ghe.Client.Organizations.Get(context.TODO(), o)

	if err != nil {
		logging.SugaredLogger.Debugf("could not extract organization data of '%s' : %s", o, err)
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
