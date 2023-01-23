package cache

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckCacheSingle(t *testing.T) {
	t.Cleanup(CleanDatabase)

	db := cache.Database("test-collection")
	collection := &Collection{
		name:       "test-check-cache-single",
		db:         &Database{name: "TestD", Database: db},
		Collection: db.Collection("test-check-cache-single"),
	}

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
		t.Fatalf("Error Inserting Into Cache")
	}

	cachedObject := checkCacheSingle[TestObject](collection)

	assert.Equal(t, testObject, *cachedObject)
}

func TestCheckCacheSingleMultiple(t *testing.T) {
	t.Cleanup(CleanDatabase)

	db := cache.Database("test-collection")
	collection := &Collection{
		name:       "test-check-cache-single-multiple",
		db:         &Database{name: "TestD", Database: db},
		Collection: db.Collection("test-check-cache-single-multiple"),
	}

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
		t.Fatalf("Error Inserting Into Cache")
	}

	cachedObjects := checkCacheSingle[TestObject](collection)

	assert.Nil(t, cachedObjects)
}

func TestCheckCacheSingleNil(t *testing.T) {
	t.Cleanup(CleanDatabase)

	db := cache.Database("test-collection")
	collection := &Collection{
		name:       "test-check-cache-single-nil",
		db:         &Database{name: "TestD", Database: db},
		Collection: db.Collection("test-check-cache-single-nil"),
	}

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
		t.Fatalf("Error Inserting Into Cache")
	}

	_ = collection.Drop(context.TODO())

	cachedObject := checkCacheSingle[TestObject](collection)

	assert.Nil(t, cachedObject)
}

func TestCheckCache(t *testing.T) {
	t.Cleanup(CleanDatabase)

	db := cache.Database("test-collection")
	collection := &Collection{
		name:       "test-check-cache",
		db:         &Database{name: "TestD", Database: db},
		Collection: db.Collection("test-check-cache"),
	}

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
		t.Fatalf("Error Inserting Into Cache")
	}

	cachedObjects := checkCache[TestObject](collection)

	assert.Equal(t, testObjects, cachedObjects)
}

func TestUpdateCache(t *testing.T) {
	t.Cleanup(CleanDatabase)

	db := cache.Database("test-collection")
	collection := &Collection{
		name:       "test-update-cache",
		db:         &Database{name: "TestD", Database: db},
		Collection: db.Collection("test-update-cache"),
	}

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
	t.Cleanup(CleanDatabase)

	db := cache.Database("test-collection")
	collection := &Collection{
		name:       "test-update-cache-error",
		db:         &Database{name: "TestD", Database: db},
		Collection: db.Collection("test-update-cache-error"),
	}

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
	t.Cleanup(CleanDatabase)

	db := cache.Database("test-collection")
	collection := &Collection{
		name:       "test-update-cache-error-first",
		db:         &Database{name: "TestD", Database: db},
		Collection: db.Collection("test-update-cache-error-first"),
	}

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
	t.Cleanup(CleanDatabase)

	db := cache.Database("test-collection")
	collection := &Collection{
		name:       "test-update-cache-error-inbetween",
		db:         &Database{name: "TestD", Database: db},
		Collection: db.Collection("test-update-cache-error-inbetween"),
	}

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
	t.Cleanup(CleanDatabase)

	db := cache.Database("test-collection")
	collection := &Collection{
		name:       "test-update-cache-single",
		db:         &Database{name: "TestD", Database: db},
		Collection: db.Collection("test-update-cache-single"),
	}

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
	t.Cleanup(CleanDatabase)

	db := cache.Database("test-collection")
	collection := &Collection{
		name:       "test-persist-collection-empty",
		db:         &Database{name: "TestD", Database: db},
		Collection: db.Collection("test-persist-collection-empty"),
	}

	assert.False(t, emptyCollectionExists(context.TODO(), collection))

	persistCollection(context.TODO(), collection, 0)

	assert.True(t, emptyCollectionExists(context.TODO(), collection))

	precheck := checkCache[any](collection)
	assert.Empty(t, precheck)
}

func TestPersistCollectionNotEmpty(t *testing.T) {
	t.Cleanup(CleanDatabase)

	db := cache.Database("test-collection")
	collection := &Collection{
		name:       "test-persist-collection-not-empty",
		db:         &Database{name: "TestD", Database: db},
		Collection: db.Collection("test-persist-collection-not-empty"),
	}

	assert.False(t, emptyCollectionExists(context.TODO(), collection))

	persistCollection(context.TODO(), collection, 3)

	assert.False(t, emptyCollectionExists(context.TODO(), collection))

	precheck := checkCache[any](collection)
	assert.Empty(t, precheck)
}

func TestEmptyCollectionExistsFalse(t *testing.T) {
	t.Cleanup(CleanDatabase)

	db := cache.Database("test-collection")
	collection := &Collection{
		name:       "test-empty-collection-exists-false",
		db:         &Database{name: "TestD", Database: db},
		Collection: db.Collection("test-empty-collection-exists-false"),
	}

	assert.False(t, emptyCollectionExists(context.TODO(), collection))
}

func TestEmptyCollectionExistsTrue(t *testing.T) {
	t.Cleanup(CleanDatabase)

	db := cache.Database("test-collection")
	collection := &Collection{
		name:       "test-empty-collection-exists-true",
		db:         &Database{name: "TestD", Database: db},
		Collection: db.Collection("test-empty-collection-exists-true"),
	}

	persistCollection(context.TODO(), collection, 0)

	assert.True(t, emptyCollectionExists(context.TODO(), collection))
}
