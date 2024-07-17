package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collector[DOCTYPE any] struct {
	*mongo.Collection
}

func (c *Collector[DOCTYPE]) FindOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) *SingleResult[DOCTYPE] {
	return &SingleResult[DOCTYPE]{SingleResult: c.Collection.FindOne(ctx, filter, opts...)}
}

func (c *Collector[DOCTYPE]) InsertOne(
	ctx context.Context,
	document interface{},
	opts ...*options.InsertOneOptions,
) error {
	_, err := c.Collection.InsertOne(ctx, document, opts...)
	return err
}

type SingleResult[DOCTYPE any] struct {
	*mongo.SingleResult
}

func (r *SingleResult[_]) Err() error {
	return r.handleError(r.SingleResult.Err())
}

func (r *SingleResult[DOCTYPE]) Decode() (*DOCTYPE, error) {
	doc := new(DOCTYPE)
	err := r.SingleResult.Decode(doc)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return doc, nil
}

func (r *SingleResult[_]) Raw() (bson.Raw, error) {
	raw, err := r.SingleResult.Raw()
	return raw, r.handleError(err)
}

func (r *SingleResult[_]) handleError(err error) error {
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}
