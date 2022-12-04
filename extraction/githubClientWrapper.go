package extraction

import (
	"context"
	"fmt"
	"github.com/google/go-github/v48/github"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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

func (s *RepositoriesServiceWrapper) ListContributors(ctx context.Context, owner string, repository string, opts *github.ListContributorsOptions) ([]*github.Contributor, error) {

	coll := s.cache.Database("repositories_list_contributors").Collection(fmt.Sprintf("%s-%s", owner, repository))

	f := func() ([]*github.Contributor, *github.Response, error) {
		return s.client.Repositories.ListContributors(ctx, owner, repository, opts)
	}

	return fetch[*github.Contributor](ctx, coll, f, &opts.ListOptions)
}

func (s *RepositoriesServiceWrapper) List(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, error) {

	coll := s.cache.Database("repositories_list").Collection(user)

	f := func() ([]*github.Repository, *github.Response, error) {
		return s.client.Repositories.List(ctx, user, opts)
	}

	return fetch[*github.Repository](ctx, coll, f, &opts.ListOptions)
}

func (s *OrganizationsServiceWrapper) List(ctx context.Context, user string, opts *github.ListOptions) ([]*github.Organization, error) {

	coll := s.cache.Database("organizations_list").Collection(user)

	f := func() ([]*github.Organization, *github.Response, error) {
		return s.client.Organizations.List(ctx, user, opts)
	}

	return fetch[*github.Organization](ctx, coll, f, opts)
}

func fetch[T any](ctx context.Context, coll *mongo.Collection, f func() ([]T, *github.Response, error), opts *github.ListOptions) ([]T, error) {

	cachedObjects := checkCache[T](coll)
	if cachedObjects != nil {
		log.Printf("HIT THE CACHE for collection '%s' of database '%s'", coll.Name(), coll.Database().Name())
		return cachedObjects, nil
	}

	log.Printf("CACHE EMPTY, need to go to api for collection '%s' of database '%s'", coll.Name(), coll.Database().Name())

	objects, err := handlePagination[T](f, opts)
	if err != nil {
		return nil, err
	}

	updateCache[T](ctx, objects, coll)

	return objects, nil
}

func checkCache[T any](collection *mongo.Collection) []T {

	if !emptyCollectionExists(context.TODO(), collection) {
		return nil
	}

	cur, err := collection.Find(context.TODO(), bson.D{{}}, options.Find())
	if err != nil {
		log.Printf("ERROR checking cache for collection '%s' of database '%s': %s", collection.Name(), collection.Database().Name(), err)
		return nil
	}

	result := make([]T, 0)
	for cur.Next(context.TODO()) {
		var elem T
		err = cur.Decode(&elem)
		if err != nil {
			log.Printf("ERROR decoding element of collection '%s' from database '%s': %s", collection.Name(), collection.Database().Name(), err)
			return nil
		}

		result = append(result, elem)
	}

	return result
}

func emptyCollectionExists(ctx context.Context, coll *mongo.Collection) bool {
	names, err := coll.Database().ListCollectionNames(ctx, bson.D{}, nil)
	if err != nil {
		log.Printf("ERROR listing collection names of database '%s': %s", coll.Database().Name(), err)
		return false
	}
	for _, name := range names {
		if name == coll.Name() {
			return true
		}
	}
	return false
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

func updateCache[T any](ctx context.Context, content []T, collection *mongo.Collection) {

	persistCollection(ctx, collection, len(content))

	for _, c := range content {
		_, err := collection.InsertOne(ctx, c, nil)
		if err != nil {
			log.Printf("ERROR updating cache for collection '%s' of database '%s': %s", collection.Name(), collection.Database().Name(), err)
			log.Printf("Cleaning cache where updating was throwing error for collection '%s' of database '%s'", collection.Name(), collection.Database().Name())
			err = collection.Drop(ctx)
			if err != nil {
				log.Printf("ERROR dropping cache for collection '%s' of database '%s': %s", collection.Name(), collection.Database().Name(), err)
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
