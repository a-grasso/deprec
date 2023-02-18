package extraction

import (
	"context"
	"github.com/a-grasso/deprec/cache"
	"github.com/a-grasso/deprec/configuration"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

var config, _ = configuration.Load("./../config/config.json", "./../config/ut.env")

var mongoCache, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI(config.MongoDB.URI).SetAuth(options.Credential{
	Username: config.MongoDB.Username,
	Password: config.MongoDB.Password,
}))

var cacheClient = &cache.Cache{
	Client: mongoCache,
}

func TestMain(m *testing.M) {
	CleanDatabase()
	defer CleanDatabase()

	m.Run()
}

func CleanDatabase() {
	databases, err := mongoCache.ListDatabases(context.TODO(), bson.D{})
	if err != nil {
		return
	}

	for _, database := range databases.Databases {
		err := mongoCache.Database(database.Name).Drop(context.TODO())
		if err != nil {
			continue
		}
	}
}

func CheckNoDatabase(t *testing.T) {
	databases, err := mongoCache.ListDatabases(context.TODO(), bson.D{})
	if err != nil {
		assert.FailNow(t, "Could not fetch databases")
	}
	assert.True(t, len(databases.Databases) == 3) // 3 databases can not be dropped (admin, local, config)
}
