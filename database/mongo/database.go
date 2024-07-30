package mongo

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/matthxwpavin/ticketing/logging/sugar"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB struct {
	URI     string
	Name    string
	Options []*MigrationOptions
	client  *mongo.Client
	db      *mongo.Database
}

func (s *DB) Connect(ctx context.Context) error {
	logger := sugar.FromContext(ctx)

	var err error
	s.client, err = mongo.Connect(ctx, options.Client().ApplyURI(s.URI), &options.ClientOptions{
		BSONOptions: &options.BSONOptions{
			NilMapAsEmpty:   true,
			NilSliceAsEmpty: true,
		},
	})
	if err != nil {
		logger.Errorw("mongo db failed to connect", "error", err)
		return err
	}

	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := s.client.Ping(pingCtx, readpref.Primary()); err != nil {
		logger.Errorw("could not ping the database", "error", err)
		return err
	}
	s.db = s.client.Database(s.Name)
	return nil
}

func (s *DB) Migrate(ctx context.Context) error {
	cnames, err := s.db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return fmt.Errorf("could not list collection names: %v", err)
	}
	for _, opt := range s.Options {
		if !slices.Contains(cnames, opt.CollectionName) {
			createOpts := &options.CreateCollectionOptions{
				Validator: opt.Validator.MongoSchema(),
			}
			if err := s.db.CreateCollection(ctx, opt.CollectionName, createOpts); err != nil {
				return fmt.Errorf("could not create the schema: %v", err)
			}
		}
	}
	return nil
}

func (s *DB) Disconnect(ctx context.Context) error {
	logger := sugar.FromContext(ctx)
	disCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := s.client.Disconnect(disCtx)
	if err != nil {
		logger.Errorw("database failed to disconnect", "error", err)
	} else {
		logger.Infoln("database disconnected")
	}
	return err
}
