package extraction

import (
	"context"
	"deprec/cache"
	"deprec/githubapi"
	"fmt"
	"github.com/google/go-github/v48/github"
	"github.com/thoas/go-funk"
)

type GitHubClientWrapper struct {
	Cache  cache.Cache
	Client githubapi.Client

	common ServiceWrapper

	Repositories  *RepositoriesServiceWrapper
	Organizations *OrganizationsServiceWrapper
	Issues        *IssuesServiceWrapper
	GraphQL       *GraphQLWrapper
}

type ServiceWrapper struct {
	Cache  cache.Cache
	Client githubapi.Client
}

type RepositoriesServiceWrapper ServiceWrapper
type OrganizationsServiceWrapper ServiceWrapper
type IssuesServiceWrapper ServiceWrapper
type GraphQLWrapper ServiceWrapper

func NewGitHubClientWrapper(client githubapi.Client, cache cache.Cache) *GitHubClientWrapper {

	wrapper := &GitHubClientWrapper{Client: client, Cache: cache}

	wrapper.common.Client = client
	wrapper.common.Cache = cache

	wrapper.Repositories = (*RepositoriesServiceWrapper)(&wrapper.common)
	wrapper.Organizations = (*OrganizationsServiceWrapper)(&wrapper.common)
	wrapper.Issues = (*IssuesServiceWrapper)(&wrapper.common)
	wrapper.GraphQL = (*GraphQLWrapper)(&wrapper.common)

	return wrapper
}

func (s *RepositoriesServiceWrapper) ListContributorStats(ctx context.Context, owner string, repository string) ([]*github.ContributorStats, error) {

	coll := s.Cache.Database("repositories_list_contributor_stats").Collection(fmt.Sprintf("%s-%s", owner, repository))

	f := func() ([]*github.ContributorStats, *github.Response, error) {
		return s.Client.Rest().Repositories.ListContributorsStats(ctx, owner, repository)
	}

	return cache.FetchAsync[*github.ContributorStats](ctx, coll, f)
}

func (s *RepositoriesServiceWrapper) ListContributors(ctx context.Context, owner string, repository string, opts *github.ListContributorsOptions) ([]*github.Contributor, error) {

	coll := s.Cache.Database("repositories_list_contributors").Collection(fmt.Sprintf("%s-%s", owner, repository))

	f := func() ([]*github.Contributor, *github.Response, error) {
		return s.Client.Rest().Repositories.ListContributors(ctx, owner, repository, opts)
	}

	pagination, err := cache.FetchPagination[*github.Contributor](ctx, coll, f, &opts.ListOptions)

	pagination = funk.Filter(pagination, func(contributor *github.Contributor) bool { return contributor.GetLogin() != "gitter-badger" }).([]*github.Contributor)

	return pagination, err
}

func (s *RepositoriesServiceWrapper) List(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, error) {

	coll := s.Cache.Database("repositories_list").Collection(user)

	f := func() ([]*github.Repository, *github.Response, error) {
		return s.Client.Rest().Repositories.List(ctx, user, opts)
	}

	return cache.FetchPagination[*github.Repository](ctx, coll, f, &opts.ListOptions)
}

func (s *RepositoriesServiceWrapper) Get(ctx context.Context, owner string, repo string) (*github.Repository, error) {

	coll := s.Cache.Database("repositories_get").Collection(fmt.Sprintf("%s-%s", owner, repo))

	f := func() (*github.Repository, *github.Response, error) {
		return s.Client.Rest().Repositories.Get(ctx, owner, repo)
	}

	return cache.FetchSingle[github.Repository](ctx, coll, f)
}

func (s *RepositoriesServiceWrapper) GetReadMe(ctx context.Context, owner string, repo string, opts *github.RepositoryContentGetOptions) (*github.RepositoryContent, error) {

	coll := s.Cache.Database("repositories_get_readme").Collection(fmt.Sprintf("%s-%s", owner, repo))

	f := func() (*github.RepositoryContent, *github.Response, error) {
		return s.Client.Rest().Repositories.GetReadme(ctx, owner, repo, opts)
	}

	return cache.FetchSingle[github.RepositoryContent](ctx, coll, f)
}

func (s *OrganizationsServiceWrapper) List(ctx context.Context, user string, opts *github.ListOptions) ([]*github.Organization, error) {

	coll := s.Cache.Database("organizations_list").Collection(user)

	f := func() ([]*github.Organization, *github.Response, error) {
		return s.Client.Rest().Organizations.List(ctx, user, opts)
	}

	return cache.FetchPagination[*github.Organization](ctx, coll, f, opts)
}

func (s *OrganizationsServiceWrapper) Get(ctx context.Context, org string) (*github.Organization, error) {

	coll := s.Cache.Database("organizations_get").Collection(org)

	f := func() (*github.Organization, *github.Response, error) {
		return s.Client.Rest().Organizations.Get(ctx, org)
	}

	return cache.FetchSingle[github.Organization](ctx, coll, f)
}

func (s *RepositoriesServiceWrapper) ListCommits(ctx context.Context, owner string, repository string, opts *github.CommitsListOptions) ([]*github.RepositoryCommit, error) {

	coll := s.Cache.Database("repositories_list_commits").Collection(fmt.Sprintf("%s-%s", owner, repository))

	f := func() ([]*github.RepositoryCommit, *github.Response, error) {
		return s.Client.Rest().Repositories.ListCommits(ctx, owner, repository, opts)
	}

	return cache.FetchPagination[*github.RepositoryCommit](ctx, coll, f, &opts.ListOptions)
}

func (s *RepositoriesServiceWrapper) ListReleases(ctx context.Context, owner string, repository string, opts *github.ListOptions) ([]*github.RepositoryRelease, error) {

	coll := s.Cache.Database("repositories_list_releases").Collection(fmt.Sprintf("%s-%s", owner, repository))

	f := func() ([]*github.RepositoryRelease, *github.Response, error) {
		return s.Client.Rest().Repositories.ListReleases(ctx, owner, repository, opts)
	}

	return cache.FetchPagination[*github.RepositoryRelease](ctx, coll, f, opts)
}

func (s *IssuesServiceWrapper) ListByRepo(ctx context.Context, owner string, repository string, opts *github.IssueListByRepoOptions) ([]*github.Issue, error) {

	coll := s.Cache.Database("issues_list_by_repo").Collection(fmt.Sprintf("%s-%s", owner, repository))

	f := func() ([]*github.Issue, *github.Response, error) {
		return s.Client.Rest().Issues.ListByRepo(ctx, owner, repository, opts)
	}

	return cache.FetchPagination[*github.Issue](ctx, coll, f, &opts.ListOptions)
}

func (s *IssuesServiceWrapper) ListComments(ctx context.Context, owner string, repository string, number int, opts *github.IssueListCommentsOptions) ([]*github.IssueComment, error) {

	coll := s.Cache.Database("issues_list_comments").Collection(fmt.Sprintf("%s-%s-%d", owner, repository, number))

	f := func() ([]*github.IssueComment, *github.Response, error) {
		return s.Client.Rest().Issues.ListComments(ctx, owner, repository, number, opts)
	}

	return cache.FetchPagination[*github.IssueComment](ctx, coll, f, &opts.ListOptions)
}

func (s *RepositoriesServiceWrapper) ListTags(ctx context.Context, owner string, repository string, opts *github.ListOptions) ([]*github.RepositoryTag, error) {

	coll := s.Cache.Database("repositories_list_tags").Collection(fmt.Sprintf("%s-%s", owner, repository))

	f := func() ([]*github.RepositoryTag, *github.Response, error) {
		return s.Client.Rest().Repositories.ListTags(ctx, owner, repository, opts)
	}

	return cache.FetchPagination[*github.RepositoryTag](ctx, coll, f, opts)
}

func (s *RepositoriesServiceWrapper) GetCommit(ctx context.Context, owner string, repository string, sha string, opts *github.ListOptions) (*github.RepositoryCommit, error) {

	coll := s.Cache.Database("repositories_get_commit").Collection(fmt.Sprintf("%s-%s-%s", owner, repository, sha))

	f := func() (*github.RepositoryCommit, *github.Response, error) {
		return s.Client.Rest().Repositories.GetCommit(ctx, owner, repository, sha, opts)
	}

	return cache.FetchSingle(ctx, coll, f)
}
