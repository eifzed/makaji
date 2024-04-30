package main

import (
	"context"
	"time"

	"github.com/eifzed/joona/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func getDBConnection(cfg *config.Config) (*mongo.Client, error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(cfg.Secrets.Data.MongoDBConfig.URI).SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, nil
}
