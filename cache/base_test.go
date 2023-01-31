package cache

import (
	"context"
	"github.com/a-grasso/deprec/configuration"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

var config, _ = configuration.Load("./../test.ut.config.json")

var cache, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI(config.MongoDB.URI).SetAuth(options.Credential{
	Username: config.MongoDB.Username,
	Password: config.MongoDB.Password,
}))

func TestMain(m *testing.M) {
	CleanDatabase()
	defer CleanDatabase()

	m.Run()
}

func CleanDatabase() {
	databases, err := cache.ListDatabases(context.TODO(), bson.D{})
	if err != nil {
		return
	}

	for _, database := range databases.Databases {
		err := cache.Database(database.Name).Drop(context.TODO())
		if err != nil {
			continue
		}
	}
}
