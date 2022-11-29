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
}

func NewGitHubClientWrapper(client *github.Client) *GitHubClientWrapper {

	credential := options.Credential{
		Username: "root",
		Password: "rootpassword",
	}

	clientOpts := options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(credential)
	cache, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Fatal(err)
	}

	return &GitHubClientWrapper{client: client, cache: cache}
}

func (wrapper *GitHubClientWrapper) RepositoriesListContributors(ctx context.Context, owner string, repository string, opts *github.ListContributorsOptions) ([]*github.Contributor, error) {

	coll := wrapper.cache.Database("repositories_list_contributors").Collection(fmt.Sprintf("%s-%s", owner, repository))

	findOptions := options.Find()
	findOptions.SetLimit(5)
	var contributors []*github.Contributor
	cur, err := coll.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var elem *github.Contributor
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		contributors = append(contributors, elem)
	}

	if len(contributors) == 0 {

		opts.ListOptions = github.ListOptions{PerPage: 100}
		for {
			conts, r, err := wrapper.client.Repositories.ListContributors(ctx, owner, repository, opts)
			if err != nil {
				break
			}
			contributors = append(contributors, conts...)
			if r.NextPage == 0 {
				break
			}
			opts.Page = r.NextPage
		}
	}

	for _, c := range contributors {
		_, err := coll.InsertOne(ctx, c, nil)
		if err != nil {
			return nil, err
		}
	}

	return contributors, err
}

func (wrapper *GitHubClientWrapper) RepositoriesList(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error) {

	return wrapper.client.Repositories.List(ctx, user, opts)
}
