package extraction

import (
	"context"
	"deprec/configuration"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

var config, _ = configuration.Load("./../test.ut.config.json")

func TestMain(m *testing.M) {
	cleanDatabase()
	defer cleanDatabase()

	m.Run()
}

func cleanDatabase() {
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

func checkNoDatabase(t *testing.T) {
	databases, err := cache.ListDatabases(context.TODO(), bson.D{})
	if err != nil {
		assert.FailNow(t, "Could not fetch databases")
	}
	assert.True(t, len(databases.Databases) == 3) // 3 databases can not be dropped (admin, local, config)
}
