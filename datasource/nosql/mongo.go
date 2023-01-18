package nosql

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
	Database    string
	MaxPoolSize int
}

func NewMongoDB(config *MongoConfig) (*mongo.Database, error) {
	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%d/?maxPoolSize=%d",
		config.Username, config.Password, config.Host, config.Port, config.MaxPoolSize,
	)

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("error connect: %w", err)
	}

	return client.Database(config.Database), nil
}
