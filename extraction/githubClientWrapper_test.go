package extraction

import (
	"context"
	"deprec/configuration"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
)

var config = configuration.Load("./../test.config.json")

var sut = GitHubExtractor{
	RepositoryURL: "test-Repo-URL",
	Repository:    "test-Repo",
	Owner:         "test-Owner",
	Config:        config,
	Client: &GitHubClientWrapper{
		client:        nil,
		cache:         nil,
		common:        ServiceWrapper{},
		Repositories:  nil,
		Organizations: nil,
	},
}

var cache, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI(config.URI).SetAuth(options.Credential{
	Username: config.Username,
	Password: config.Password,
}))

func init() {
	err := cache.Database("test-collection").Drop(context.TODO())
	if err != nil {
		return
	}
}
func TestCheckCacheSingle(t *testing.T) {

	collection := cache.Database("test-collection").Collection("test-check-cache-single")

	type TestObject struct {
		One   string
		Two   int
		Three []string
	}

	testObject := TestObject{
		One:   "one",
		Two:   2,
		Three: []string{"one", "two", "string"},
	}

	_, err := collection.InsertOne(context.TODO(), testObject)
	if err != nil {
		log.Println("ERROR")
	}

	cachedObject := checkCacheSingle[TestObject](collection)

	assert.Equal(t, testObject, *cachedObject)
}

func TestCheckCacheSingleMultiple(t *testing.T) {

	collection := cache.Database("test-collection").Collection("test-check-cache-single-multiple")

	type TestObject struct {
		One   string
		Two   int
		Three []string
	}

	testObjects := []TestObject{
		{
			One:   "one",
			Two:   2,
			Three: []string{"one", "two", "string"},
		}, {
			One:   "one1",
			Two:   22,
			Three: []string{"one1", "two2", "string3"},
		},
	}

	_, err := collection.InsertMany(context.TODO(), []interface{}{testObjects[0], testObjects[1]})
	if err != nil {
		log.Println("ERROR")
	}

	cachedObjects := checkCacheSingle[TestObject](collection)

	assert.Nil(t, cachedObjects)
}

func TestCheckCacheSingleNil(t *testing.T) {

	collection := cache.Database("test-collection").Collection("test-check-cache-single-nil")

	type TestObject struct {
		One   string
		Two   int
		Three []string
	}

	testObject := TestObject{
		One:   "one",
		Two:   2,
		Three: []string{"one", "two", "string"},
	}

	_, err := collection.InsertOne(context.TODO(), testObject)
	if err != nil {
		log.Println("ERROR")
	}

	_ = collection.Drop(context.TODO())

	cachedObject := checkCacheSingle[TestObject](collection)

	assert.Nil(t, cachedObject)
}

func TestCheckCache(t *testing.T) {

	collection := cache.Database("test-collection").Collection("test-check-cache")

	type TestObject struct {
		One   string
		Two   int
		Three []string
	}

	testObjects := []TestObject{
		{
			One:   "one",
			Two:   2,
			Three: []string{"one", "two", "string"},
		}, {
			One:   "one1",
			Two:   22,
			Three: []string{"one1", "two2", "string3"},
		},
	}

	_, err := collection.InsertMany(context.TODO(), []interface{}{testObjects[0], testObjects[1]})
	if err != nil {
		log.Println("ERROR")
	}

	cachedObjects := checkCache[TestObject](collection)

	assert.Equal(t, testObjects, cachedObjects)
}

func TestUpdateCache(t *testing.T) {

	collection := cache.Database("test-collection").Collection("test-update-cache")

	precheck := checkCache[any](collection)
	assert.Empty(t, precheck)

	type TestObject struct {
		One   string
		Two   int
		Three []string
	}

	testObjects := []TestObject{
		{
			One:   "one",
			Two:   2,
			Three: []string{"one", "two", "string"},
		}, {
			One:   "one1",
			Two:   22,
			Three: []string{"one1", "two2", "string3"},
		},
	}

	updateCache[TestObject](context.TODO(), testObjects, collection)

	aftercheck := checkCache[any](collection)
	assert.NotEmpty(t, aftercheck)
}

func TestUpdateCacheSingleError(t *testing.T) {

	collection := cache.Database("test-collection").Collection("test-update-cache-error")

	precheck := checkCache[any](collection)
	assert.Empty(t, precheck)

	type TestObject struct {
		One   string
		Two   int
		Three []string
	}

	updateCacheSingle[*TestObject](context.TODO(), nil, collection)

	aftercheck := checkCache[any](collection)
	assert.Empty(t, aftercheck)
}

func TestUpdateCacheErrorFirst(t *testing.T) {

	collection := cache.Database("test-collection").Collection("test-update-cache-error-first")

	precheck := checkCache[any](collection)
	assert.Empty(t, precheck)

	type TestObject struct {
		One   string
		Two   int
		Three []string
	}

	updateCache[*TestObject](context.TODO(), []*TestObject{nil}, collection)

	aftercheck := checkCache[any](collection)
	assert.Empty(t, aftercheck)
}

func TestUpdateCacheErrorInbetween(t *testing.T) {

	collection := cache.Database("test-collection").Collection("test-update-cache-error-inbetween")

	precheck := checkCache[any](collection)
	assert.Empty(t, precheck)

	type TestObject struct {
		One   string
		Two   int
		Three []string
	}

	testObjects := []*TestObject{
		{
			One:   "one",
			Two:   2,
			Three: []string{"one", "two", "string"},
		}, nil,
	}

	updateCache[*TestObject](context.TODO(), testObjects, collection)

	aftercheck := checkCache[any](collection)
	assert.Empty(t, aftercheck)
}

func TestUpdateCacheSingle(t *testing.T) {

	collection := cache.Database("test-collection").Collection("test-update-cache-single")

	precheck := checkCache[any](collection)
	assert.Empty(t, precheck)

	type TestObject struct {
		One   string
		Two   int
		Three []string
	}

	testObject := TestObject{
		One:   "one",
		Two:   2,
		Three: []string{"one", "two", "string"},
	}

	updateCacheSingle[TestObject](context.TODO(), testObject, collection)

	aftercheck := checkCache[any](collection)
	assert.NotEmpty(t, aftercheck)
	assert.Equal(t, 1, len(aftercheck))
}

func TestPersistCollectionEmpty(t *testing.T) {
	collection := cache.Database("test-collection").Collection("test-persist-collection-empty")

	assert.False(t, emptyCollectionExists(context.TODO(), collection))

	persistCollection(context.TODO(), collection, 0)

	assert.True(t, emptyCollectionExists(context.TODO(), collection))

	precheck := checkCache[any](collection)
	assert.Empty(t, precheck)
}

func TestPersistCollectionNotEmpty(t *testing.T) {
	collection := cache.Database("test-collection").Collection("test-persist-collection-not-empty")

	assert.False(t, emptyCollectionExists(context.TODO(), collection))

	persistCollection(context.TODO(), collection, 3)

	assert.False(t, emptyCollectionExists(context.TODO(), collection))

	precheck := checkCache[any](collection)
	assert.Empty(t, precheck)
}

func TestEmptyCollectionExistsFalse(t *testing.T) {
	collection := cache.Database("test-collection").Collection("test-empty-collection-exists-false")

	assert.False(t, emptyCollectionExists(context.TODO(), collection))
}

func TestEmptyCollectionExistsTrue(t *testing.T) {
	collection := cache.Database("test-collection").Collection("test-empty-collection-exists-true")

	persistCollection(context.TODO(), collection, 0)

	assert.True(t, emptyCollectionExists(context.TODO(), collection))
}

func TestFetchSingle(t *testing.T) {

}
