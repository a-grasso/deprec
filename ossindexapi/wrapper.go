package ossindexapi

import (
	"context"
	"deprec/cache"
	"github.com/nscuro/ossindex-client"
)

type ClientWrapper struct {
	Cache  *cache.Cache
	Client *Client
}

func NewClientWrapper(client *Client, cache *cache.Cache) *ClientWrapper {
	return &ClientWrapper{
		Cache:  cache,
		Client: client,
	}
}

func (cw *ClientWrapper) GetComponentReport(purl string) ([]ossindex.ComponentReport, error) {

	coll := cw.Cache.Database("ossindex_component_report").Collection(purl)

	f := func() ([]ossindex.ComponentReport, error) {
		reports, err := cw.Client.GetComponentReports(context.TODO(), []string{purl})
		return reports, err
	}

	return cache.FetchMultiple[ossindex.ComponentReport](context.TODO(), coll, f)
}
