package mongo

import (
	"context"
	"fmt"

	"github.com/matthxwpavin/ticketing/logging/sugar"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Collection[T any] struct {
	Name string
	c    *mongo.Collection
}

func NewCollection[T any](db *DB, name string) *Collection[T] {
	return &Collection[T]{c: db.db.Collection(name)}
}

func (c *Collection[T]) Insert(ctx context.Context, document *T) (string, error) {
	logger := sugar.FromContext(ctx)
	res, err := c.c.InsertOne(ctx, document)
	if err != nil {
		logger.Errorw("could not insert the document", "error", err)
		return "", err
	}
	id, ok := res.InsertedID.(string)
	if !ok {
		logger.Errorf("inserted id is not an object id, type: %T", res.InsertedID)
		return "", fmt.Errorf("inserted id is not an object id, type: %T", res.InsertedID)
	}
	return id, nil
}

func (c *Collection[T]) FindByID(ctx context.Context, id string) (*T, error) {
	return c.FindOne(ctx, c.idFilter(id))
}

func (c *Collection[T]) FindOne(ctx context.Context, filter any) (*T, error) {
	logger := sugar.FromContext(ctx)

	r := c.c.FindOne(ctx, filter)
	notfound := r.Err() == mongo.ErrNoDocuments
	if r.Err() != nil && !notfound {
		logger.Errorw("unable to find document", "error", r.Err())
		return nil, r.Err()
	}
	if notfound {
		return nil, nil
	}

	var dst T
	if err := r.Decode(&dst); err != nil {
		logger.Errorw("could not decode result", "error", err)
		return nil, err
	}
	return &dst, nil
}

func (c *Collection[T]) FindAll(ctx context.Context) ([]*T, error) {
	return c.Find(ctx, bson.D{})
}

func (c *Collection[T]) Find(ctx context.Context, filter any) ([]*T, error) {
	logger := sugar.FromContext(ctx)

	var dst []*T
	cursor, err := c.c.Find(ctx, filter)
	if err != nil {
		logger.Errorw("could not find documents", "error", err)
		return nil, err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			logger.Errorw("a cursor failed to close", "error", err)
		}
	}()
	if err := cursor.All(ctx, &dst); err != nil {
		logger.Errorw("could not decode all results", "error", err)
		return nil, err
	}
	return dst, nil
}

func (c *Collection[T]) DeleteByID(ctx context.Context, id string) error {
	logger := sugar.FromContext(ctx)

	if _, err := c.c.DeleteOne(ctx, c.idFilter(id)); err != nil {
		logger.Errorw("could not delete a document", "error", err)
		return err
	}
	return nil
}

func (c *Collection[T]) UpdateByID(ctx context.Context, id string, update *T) error {
	logger := sugar.FromContext(ctx)

	if _, err := c.c.UpdateByID(ctx, id, bson.D{{"$set", update}}); err != nil {
		logger.Errorw("could not update the document", "error", err)
		return err
	}
	return nil
}

func (c *Collection[_]) DeleteAll(ctx context.Context) error {
	logger := sugar.FromContext(ctx)

	_, err := c.c.DeleteMany(ctx, bson.D{})
	if err != nil {
		logger.Errorw("could not delete all", "error", err)
	}
	return err
}

func (c *Collection[_]) Aggregate(ctx context.Context, stage bson.D, stages ...bson.D) (*mongo.Cursor, error) {
	allStages := append([]bson.D{stage}, stages...)
	return c.c.Aggregate(ctx, mongo.Pipeline(allStages))
}

func (c *Collection[_]) idFilter(id string) bson.D {
	return bson.D{{"_id", id}}
}
