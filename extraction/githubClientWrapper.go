package extraction

import (
	"context"
	"fmt"
	"github.com/google/go-github/v48/github"
	"go.mongodb.org/mongo-driver/mongo"
)

type GitHubClientWrapper struct {
	client *github.Client
	cache  *mongo.Client

	common ServiceWrapper

	Repositories  *RepositoriesServiceWrapper
	Organizations *OrganizationsServiceWrapper
}

type ServiceWrapper struct {
	cache  *mongo.Client
	client *github.Client
}

type RepositoriesServiceWrapper ServiceWrapper
type OrganizationsServiceWrapper ServiceWrapper

func NewGitHubClientWrapper(client *github.Client, cache *mongo.Client) *GitHubClientWrapper {

	wrapper := &GitHubClientWrapper{client: client, cache: cache}

	wrapper.common.client = client
	wrapper.common.cache = cache

	wrapper.Repositories = (*RepositoriesServiceWrapper)(&wrapper.common)
	wrapper.Organizations = (*OrganizationsServiceWrapper)(&wrapper.common)

	return wrapper
}

func (s *RepositoriesServiceWrapper) ListContributorStats(ctx context.Context, owner string, repository string) ([]*github.ContributorStats, error) {

	coll := s.cache.Database("repositories_list_contributor_stats").Collection(fmt.Sprintf("%s-%s", owner, repository))

	f := func() ([]*github.ContributorStats, *github.Response, error) {
		return s.client.Repositories.ListContributorsStats(ctx, owner, repository)
	}

	return fetchAsync[*github.ContributorStats](ctx, coll, f)
}

func (s *RepositoriesServiceWrapper) ListContributors(ctx context.Context, owner string, repository string, opts *github.ListContributorsOptions) ([]*github.Contributor, error) {

	coll := s.cache.Database("repositories_list_contributors").Collection(fmt.Sprintf("%s-%s", owner, repository))

	f := func() ([]*github.Contributor, *github.Response, error) {
		return s.client.Repositories.ListContributors(ctx, owner, repository, opts)
	}

	return fetchPagination[*github.Contributor](ctx, coll, f, &opts.ListOptions)
}

func (s *RepositoriesServiceWrapper) List(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, error) {

	coll := s.cache.Database("repositories_list").Collection(user)

	f := func() ([]*github.Repository, *github.Response, error) {
		return s.client.Repositories.List(ctx, user, opts)
	}

	return fetchPagination[*github.Repository](ctx, coll, f, &opts.ListOptions)
}

func (s *RepositoriesServiceWrapper) Get(ctx context.Context, owner string, repo string) (*github.Repository, error) {

	coll := s.cache.Database("repositories_get").Collection(fmt.Sprintf("%s-%s", owner, repo))

	f := func() (*github.Repository, *github.Response, error) {
		return s.client.Repositories.Get(ctx, owner, repo)
	}

	return fetchSingle[github.Repository](ctx, coll, f)
}

func (s *RepositoriesServiceWrapper) GetReadMe(ctx context.Context, owner string, repo string, opts *github.RepositoryContentGetOptions) (*github.RepositoryContent, error) {

	coll := s.cache.Database("repositories_get_readme").Collection(fmt.Sprintf("%s-%s", owner, repo))

	f := func() (*github.RepositoryContent, *github.Response, error) {
		return s.client.Repositories.GetReadme(ctx, owner, repo, opts)
	}

	return fetchSingle[github.RepositoryContent](ctx, coll, f)
}

func (s *OrganizationsServiceWrapper) List(ctx context.Context, user string, opts *github.ListOptions) ([]*github.Organization, error) {

	coll := s.cache.Database("organizations_list").Collection(user)

	f := func() ([]*github.Organization, *github.Response, error) {
		return s.client.Organizations.List(ctx, user, opts)
	}

	return fetchPagination[*github.Organization](ctx, coll, f, opts)
}

func (s *OrganizationsServiceWrapper) Get(ctx context.Context, org string) (*github.Organization, error) {

	coll := s.cache.Database("organizations_get").Collection(org)

	f := func() (*github.Organization, *github.Response, error) {
		return s.client.Organizations.Get(ctx, org)
	}

	return fetchSingle[github.Organization](ctx, coll, f)
}

func (s *RepositoriesServiceWrapper) ListCommits(ctx context.Context, owner string, repository string, opts *github.CommitsListOptions) ([]*github.RepositoryCommit, error) {

	coll := s.cache.Database("repositories_list_commits").Collection(fmt.Sprintf("%s-%s", owner, repository))

	f := func() ([]*github.RepositoryCommit, *github.Response, error) {
		return s.client.Repositories.ListCommits(ctx, owner, repository, opts)
	}

	return fetchPagination[*github.RepositoryCommit](ctx, coll, f, &opts.ListOptions)
}
