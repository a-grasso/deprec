package githubapi

import (
	"context"
	"github.com/a-grasso/deprec/configuration"
	"github.com/google/go-github/v48/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Client struct {
	restClient  *github.Client
	graphClient *githubv4.Client
}

func NewClient(config configuration.GitHub) *Client {

	rest, graph := githubClient(config)

	return &Client{
		restClient:  rest,
		graphClient: graph,
	}
}

func githubClient(config configuration.GitHub) (*github.Client, *githubv4.Client) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.APIToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	rest := github.NewClient(tc)
	graph := githubv4.NewClient(tc)

	return rest, graph
}

func (c *Client) Rest() *github.Client {
	return c.restClient
}

func (c *Client) GraphQL() *githubv4.Client {
	return c.graphClient
}
