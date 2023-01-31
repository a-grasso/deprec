package mavencentralapi

import (
	"context"
	"fmt"
	"github.com/a-grasso/deprec/cache"
	"github.com/vifraa/gopom"
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

func (cw *ClientWrapper) SearchMavenCentralSHA1(sha1 string) (*MavenCentralSearch, error) {

	coll := cw.Cache.Database("mavencentral_search_sha").Collection(sha1)

	f := func() (*MavenCentralSearch, error) {
		reports, err := cw.Client.SearchMavenCentralSHA1(sha1)
		return reports, err
	}

	return cache.FetchSingle[MavenCentralSearch](context.TODO(), coll, f)
}

func (cw *ClientWrapper) GetArtifactPom(groupId, artifactId, version string) (*gopom.Project, error) {

	coll := cw.Cache.Database("mavencentral_browse_pom").Collection(fmt.Sprintf("%s-%s-%s", groupId, artifactId, version))

	f := func() (*gopom.Project, error) {
		reports, err := cw.Client.GetArtifactPom(groupId, artifactId, version)
		return reports, err
	}

	return cache.FetchSingle[gopom.Project](context.TODO(), coll, f)
}

func (cw *ClientWrapper) GetLibraryMetadata(groupId, artifactId string) (*Metadata, error) {

	coll := cw.Cache.Database("mavencentral_browse_metadata").Collection(fmt.Sprintf("%s-%s", groupId, artifactId))

	f := func() (*Metadata, error) {
		reports, err := cw.Client.GetLibraryMetadata(groupId, artifactId)
		return reports, err
	}

	return cache.FetchSingle[Metadata](context.TODO(), coll, f)
}
