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

	Repositories *RepositoriesServiceWrapper
}

type ServiceWrapper struct {
	cache  *mongo.Client
	client *github.Client
}

type RepositoriesServiceWrapper ServiceWrapper

func NewGitHubClientWrapper(client *github.Client, cache *mongo.Client) *GitHubClientWrapper {

	wrapper := &GitHubClientWrapper{client: client, cache: cache}

	wrapper.common.client = client
	wrapper.common.cache = cache

	wrapper.Repositories = (*RepositoriesServiceWrapper)(&wrapper.common)

	return wrapper
}

func (s *RepositoriesServiceWrapper) ListContributors(ctx context.Context, owner string, repository string, opts *github.ListContributorsOptions) ([]*github.Contributor, error) {

	coll := s.cache.Database("repositories_list_contributors").Collection(fmt.Sprintf("%s-%s", owner, repository))

	f := func() ([]*github.Contributor, *github.Response, error) {
		return s.client.Repositories.ListContributors(ctx, owner, repository, opts)
	}

	return Get[*github.Contributor](ctx, coll, f, &opts.ListOptions)
}

func Get[T any](ctx context.Context, coll *mongo.Collection, f func() ([]T, *github.Response, error), opts *github.ListOptions) ([]T, error) {

	contributors := make([]T, 0)

	checkCache[T](&contributors, coll)

	if len(contributors) == 0 {
		log.Printf("Cache empty, need to go to api for collection '%s' of database '%s'", coll.Name(), coll.Database().Name())

		handlePagination[T](f, &contributors, opts)

		err := updateCache[T](ctx, &contributors, coll)
		if err != nil {
			return nil, err
		}

	} else {
		log.Printf("HIT THE CACHE for collection '%s' of database '%s'", coll.Name(), coll.Database().Name())
	}

	return contributors, nil
}

func (s *RepositoriesServiceWrapper) List(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, error) {

	coll := s.cache.Database("repositories_list").Collection(user)

	f := func() ([]*github.Repository, *github.Response, error) {
		return s.client.Repositories.List(ctx, user, opts)
	}

	return Get[*github.Repository](ctx, coll, f, &opts.ListOptions)
}

func handlePagination[T any](f func() ([]T, *github.Response, error), result *[]T, opts *github.ListOptions) {
	opts.PerPage = 100
	for {
		conts, r, err := f()
		if err != nil {
			break
		}
		*result = append(*result, conts...)
		if r.NextPage == 0 {
			break
		}
		opts.Page = r.NextPage
	}
}

func checkCache[T any](result *[]T, collection *mongo.Collection) {
	cur, err := collection.Find(context.TODO(), bson.D{{}}, options.Find())
	if err != nil {

	}
	for cur.Next(context.TODO()) {
		var elem T
		err := cur.Decode(&elem)
		if err != nil {

		}

		*result = append(*result, elem)
	}
}

func updateCache[T any](ctx context.Context, content *[]T, collection *mongo.Collection) error {
	for _, c := range *content {
		_, err := collection.InsertOne(ctx, c, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
