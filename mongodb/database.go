package mongodb

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/matthxwpavin/ticketing/fmts"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Connect(ctx context.Context, uri string, databaseName string) (*mongo.Database, DisconnectFunc) {
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(uri),
		&options.ClientOptions{
			BSONOptions: &options.BSONOptions{
				NilMapAsEmpty:   true,
				NilSliceAsEmpty: true,
			},
		},
	)
	if err != nil {
		fmts.Panicf("failed to connect mongo: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		fmts.Panicf("failed to ping mongo: %v", err)
	}

	return client.Database(databaseName), client.Disconnect
}

func ConnTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

type DisconnectFunc func(context.Context) error

type MigrationOptions struct {
	CollectionName string
	Opts           *options.CreateCollectionOptions
}

func Migrate(ctx context.Context, db *mongo.Database, opts []*MigrationOptions) {
	cnames, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		fmts.Panicf("failed to list collection names: %v", err)
	}
	for _, opt := range opts {
		if !slices.Contains(cnames, opt.CollectionName) {
			if err := db.CreateCollection(ctx, opt.CollectionName, opt.Opts); err != nil {
				fmts.Panicf(fmt.Sprintf("failed to create schema: %v", err))
			}
		}
	}
}

type Collection[T any] struct {
	*Collector[T]
}

func NewCollection[T any](db *mongo.Database, collectionName string) *Collection[T] {
	return &Collection[T]{&Collector[T]{Collection: db.Collection(collectionName)}}
}
