package extraction

import (
	"context"
	"deprec/logging"
	"fmt"
	"github.com/google/go-github/v48/github"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
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

	single, err := fetchSingle[*github.Repository](ctx, coll, f)
	return *single, err
}

func (s *RepositoriesServiceWrapper) GetReadMe(ctx context.Context, owner string, repo string, opts *github.RepositoryContentGetOptions) (*github.RepositoryContent, error) {

	coll := s.cache.Database("repositories_get_readme").Collection(fmt.Sprintf("%s-%s", owner, repo))

	f := func() (*github.RepositoryContent, *github.Response, error) {
		return s.client.Repositories.GetReadme(ctx, owner, repo, opts)
	}

	single, err := fetchSingle[*github.RepositoryContent](ctx, coll, f)
	return *single, err
}

func (s *OrganizationsServiceWrapper) List(ctx context.Context, user string, opts *github.ListOptions) ([]*github.Organization, error) {

	coll := s.cache.Database("organizations_list").Collection(user)

	f := func() ([]*github.Organization, *github.Response, error) {
		return s.client.Organizations.List(ctx, user, opts)
	}

	return fetchPagination[*github.Organization](ctx, coll, f, opts)
}

func (s *RepositoriesServiceWrapper) ListCommits(ctx context.Context, owner string, repository string, opts *github.CommitsListOptions) ([]*github.RepositoryCommit, error) {

	coll := s.cache.Database("repositories_list_commits").Collection(fmt.Sprintf("%s-%s", owner, repository))

	f := func() ([]*github.RepositoryCommit, *github.Response, error) {
		return s.client.Repositories.ListCommits(ctx, owner, repository, opts)
	}

	return fetchPagination[*github.RepositoryCommit](ctx, coll, f, &opts.ListOptions)
}

func fetchSingle[T any](ctx context.Context, coll *mongo.Collection, f func() (T, *github.Response, error)) (*T, error) {

	cachedObject := checkCacheSingle[T](coll)
	if cachedObject != nil {
		logging.SugaredLogger.Debugf("HIT THE CACHE | collection '%s' of database '%s'", coll.Name(), coll.Database().Name())
		return cachedObject, nil
	}

	logging.SugaredLogger.Debugf("CACHE EMPTY | consuming API for collection '%s' of database '%s'", coll.Name(), coll.Database().Name())

	object, _, err := f()
	if err != nil {
		return nil, err
	}

	updateCacheSingle[T](ctx, object, coll)

	return &object, nil
}

func fetchPagination[T any](ctx context.Context, coll *mongo.Collection, f func() ([]T, *github.Response, error), opts *github.ListOptions) ([]T, error) {

	cachedObjects := checkCache[T](coll)
	if cachedObjects != nil {
		logging.SugaredLogger.Debugf("HIT THE CACHE | collection '%s' of database '%s'", coll.Name(), coll.Database().Name())
		return cachedObjects, nil
	}

	logging.SugaredLogger.Debugf("CACHE EMPTY | consuming API for collection '%s' of database '%s'", coll.Name(), coll.Database().Name())

	objects, err := handlePagination[T](f, opts)
	if err != nil {
		return nil, err
	}

	updateCache[T](ctx, objects, coll)

	return objects, nil
}

func fetchAsync[T any](ctx context.Context, coll *mongo.Collection, f func() ([]T, *github.Response, error)) ([]T, error) {

	cachedObjects := checkCache[T](coll)
	if cachedObjects != nil {
		logging.SugaredLogger.Debugf("HIT THE CACHE | collection '%s' of database '%s'", coll.Name(), coll.Database().Name())
		return cachedObjects, nil
	}

	logging.SugaredLogger.Debugf("CACHE EMPTY | consuming API for collection '%s' of database '%s'", coll.Name(), coll.Database().Name())

	objects, err := handleAsync[[]T](f)
	if err != nil {
		return nil, err
	}

	updateCache[T](ctx, objects, coll)

	return objects, nil
}

func handleAsync[T any](f func() (T, *github.Response, error)) (T, error) {
	var object T
	var err error

	for {
		tmp, _, tmpErr := f()
		object = tmp
		err = tmpErr

		_, isAcceptedError := err.(*github.AcceptedError)

		if isAcceptedError {
			logging.Logger.Info("waiting for async request of GitHub...")
		} else {
			break
		}

		time.Sleep(10000)
	}
	return object, err
}

func handlePagination[T any](f func() ([]T, *github.Response, error), opts *github.ListOptions) ([]T, error) {
	objects := make([]T, 0)

	opts.PerPage = 100
	for {
		content, r, err := f()
		if err != nil {
			return nil, err
		}
		objects = append(objects, content...)
		if r.NextPage == 0 {
			break
		}
		opts.Page = r.NextPage
	}
	return objects, nil
}

func checkCacheSingle[T any](collection *mongo.Collection) *T {
	cachedObjects := checkCache[T](collection)
	if len(cachedObjects) == 1 {
		return &cachedObjects[0]
	}
	return nil
}

func checkCache[T any](collection *mongo.Collection) []T {

	if !emptyCollectionExists(context.TODO(), collection) {
		return nil
	}

	cur, err := collection.Find(context.TODO(), bson.D{{}}, options.Find())
	if err != nil {
		logging.SugaredLogger.Errorf("checking cache for collection '%s' of database '%s': %s", collection.Name(), collection.Database().Name(), err)
		return nil
	}

	result := make([]T, 0)
	for cur.Next(context.TODO()) {
		var elem T
		err = cur.Decode(&elem)
		if err != nil {
			logging.SugaredLogger.Errorf("decoding element of collection '%s' from database '%s': %s", collection.Name(), collection.Database().Name(), err)
			return nil
		}

		result = append(result, elem)
	}

	return result
}

func emptyCollectionExists(ctx context.Context, coll *mongo.Collection) bool {
	names, err := coll.Database().ListCollectionNames(ctx, bson.D{}, nil)
	if err != nil {
		logging.SugaredLogger.Errorf("listing collection names of database '%s': %s", coll.Database().Name(), err)
		return false
	}
	for _, name := range names {
		if name == coll.Name() {
			return true
		}
	}
	return false
}

func updateCacheSingle[T any](ctx context.Context, content T, collection *mongo.Collection) {
	updateCache[T](ctx, []T{content}, collection)
}

func updateCache[T any](ctx context.Context, content []T, collection *mongo.Collection) {

	persistCollection(ctx, collection, len(content))

	for _, c := range content {
		_, err := collection.InsertOne(ctx, c, nil)
		if err != nil {
			logging.SugaredLogger.Errorf("updating cache for collection '%s' of database '%s': %s", collection.Name(), collection.Database().Name(), err)
			logging.SugaredLogger.Infof("cleaning cache where updating was throwing error for collection '%s' of database '%s'", collection.Name(), collection.Database().Name())
			err = collection.Drop(ctx)
			if err != nil {
				logging.SugaredLogger.Errorf("dropping cache for collection '%s' of database '%s': %s", collection.Name(), collection.Database().Name(), err)
			}
			break
		}
	}
}

func persistCollection(ctx context.Context, collection *mongo.Collection, length int) {
	if length == 0 {
		_, _ = collection.InsertOne(ctx, bson.D{})
		_, _ = collection.DeleteOne(ctx, bson.D{})
	}
}
