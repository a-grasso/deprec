package cache

import (
	"context"
	"errors"
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/logging"
	"github.com/google/go-github/v48/github"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Cache struct {
	*mongo.Client
}

type Database struct {
	name string
	*mongo.Database
}

type Collection struct {
	name string
	db   *Database
	*mongo.Collection
}

func (c *Cache) Database(db string) *Database {
	if c.Client == nil {
		return &Database{db, nil}
	}

	return &Database{db, c.Client.Database(db)}
}

func (db *Database) Name() string {
	return db.name
}

func (db *Database) Collection(coll string) *Collection {
	if db.Database == nil {
		return &Collection{
			name:       coll,
			db:         db,
			Collection: nil,
		}
	}

	return &Collection{coll, db, db.Database.Collection(coll)}
}

func (c *Collection) Name() string {
	return c.name
}

func (c *Collection) IsBroken() bool {
	return c.Collection == nil
}

func (c *Collection) Database() *Database {
	return c.db
}

func NewCache(config configuration.MongoDB) (*Cache, error) {

	if config.URI == "" || config.Password == "" || config.Username == "" {
		return &Cache{
			nil,
		}, errors.New("could not create cache, config invalid")
	}

	client := mongoDBClient(config)

	return &Cache{
		client,
	}, nil
}

func mongoDBClient(config configuration.MongoDB) *mongo.Client {
	credentials := options.Credential{
		Username: config.Username,
		Password: config.Password,
	}

	clientOpts := options.Client().ApplyURI(config.URI).SetAuth(credentials)
	cache, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		logging.SugaredLogger.Warnf("connecting to mongodb database at '%s' failed: %s", config.URI, err)
		return nil
	}

	timeout, _ := context.WithTimeout(context.TODO(), 1*time.Second)
	err = cache.Ping(timeout, nil)
	if err != nil {
		logging.SugaredLogger.Warnf("pinging mongodb database at '%s' failed: %s", config.URI, err)
		return nil
	}

	return cache
}

func FetchSingle[T any](ctx context.Context, coll *Collection, f func() (*T, error)) (*T, error) {

	cachedObject := checkCacheSingle[T](coll)
	if cachedObject != nil {
		logging.SugaredLogger.Debugf("CACHE HIT | collection '%s' of database '%s'", coll.Name(), coll.Database().Name())
		return cachedObject, nil
	}

	logging.SugaredLogger.Debugf("EMPTY CACHE | consuming API for collection '%s' of database '%s'", coll.Name(), coll.Database().Name())

	object, err := f()
	if err != nil {
		return nil, err
	}

	updateCacheSingle[*T](ctx, object, coll)

	return object, nil
}

func FetchMultiple[T any](ctx context.Context, coll *Collection, f func() ([]T, error)) ([]T, error) {

	cachedObjects := checkCache[T](coll)
	if cachedObjects != nil {
		logging.SugaredLogger.Debugf("CACHE HIT | collection '%s' of database '%s'", coll.Name(), coll.Database().Name())
		return cachedObjects, nil
	}

	logging.SugaredLogger.Debugf("EMPTY CACHE | consuming API for collection '%s' of database '%s'", coll.Name(), coll.Database().Name())

	objects, err := f()
	if err != nil {
		return nil, err
	}

	updateCache[T](ctx, objects, coll)

	return objects, nil
}

func FetchPagination[T any](ctx context.Context, coll *Collection, f func() ([]T, *github.Response, error), opts *github.ListOptions) ([]T, error) {

	cachedObjects := checkCache[T](coll)
	if cachedObjects != nil {
		logging.SugaredLogger.Debugf("CACHE HIT | collection '%s' of database '%s'", coll.Name(), coll.Database().Name())
		return cachedObjects, nil
	}

	logging.SugaredLogger.Debugf("EMPTY CACHE | consuming API for collection '%s' of database '%s'", coll.Name(), coll.Database().Name())

	objects, err := handlePagination[T](f, opts)
	if err != nil {
		return nil, err
	}

	updateCache[T](ctx, objects, coll)

	return objects, nil
}

func FetchBatchQuery[T any](ctx context.Context, coll *Collection, f func() (map[string]T, error)) ([]T, error) {

	cached := checkCache[T](coll)

	if cached != nil {
		logging.SugaredLogger.Debugf("CACHE HIT | collection '%s' of database '%s'", coll.Name(), coll.Database().Name())
		return cached, nil
	}

	logging.SugaredLogger.Debugf("EMPTY CACHE | consuming API for collection '%s' of database '%s'", coll.Name(), coll.Database().Name())

	queryResponse, err := f()
	if err != nil {
		return nil, err
	}

	values := funk.Values(queryResponse).([]T)
	updateCache[T](ctx, values, coll)

	return values, nil
}

func FetchAsync[T any](ctx context.Context, coll *Collection, f func() ([]T, *github.Response, error)) ([]T, error) {

	cachedObjects := checkCache[T](coll)
	if cachedObjects != nil {
		logging.SugaredLogger.Debugf("CACHE HIT | collection '%s' of database '%s'", coll.Name(), coll.Database().Name())
		return cachedObjects, nil
	}

	logging.SugaredLogger.Debugf("EMPTY CACHE | consuming API for collection '%s' of database '%s'", coll.Name(), coll.Database().Name())

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
			logging.Logger.Debug("waiting for async request of GitHub...")
		} else {
			break
		}

		time.Sleep(5 * time.Second)
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

func checkCacheSingle[T any](collection *Collection) *T {
	cachedObjects := checkCache[T](collection)
	if len(cachedObjects) == 1 {
		return &cachedObjects[0]
	}
	return nil
}

func checkCache[T any](collection *Collection) []T {

	if collection.IsBroken() {
		return nil
	}

	if !emptyCollectionExists(context.TODO(), collection) {
		return nil
	}

	cur, err := collection.Find(context.TODO(), bson.D{{}}, options.Find())
	if err != nil {
		logging.SugaredLogger.Errorf("checking cache for collection '%s' of database '%s': %s", collection.Name(), collection.Database().Name(), err)
		return nil
	}

	var result []T
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

func emptyCollectionExists(ctx context.Context, coll *Collection) bool {
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

func updateCacheSingle[T any](ctx context.Context, content T, collection *Collection) {
	updateCache[T](ctx, []T{content}, collection)
}

func updateCache[T any](ctx context.Context, content []T, collection *Collection) {

	if collection.IsBroken() {
		return
	}

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

func persistCollection(ctx context.Context, collection *Collection, length int) {
	if length == 0 {
		_, _ = collection.InsertOne(ctx, bson.D{})
		_, _ = collection.DeleteOne(ctx, bson.D{})
	}
}
