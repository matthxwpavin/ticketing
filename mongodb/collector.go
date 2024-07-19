package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Doctype interface {
	SetUpdatedAt(t time.Time)
}

type Collector[T Doctype] struct {
	*mongo.Collection
}

func (c *Collector[Doctype]) FindOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) *SingleResult[Doctype] {
	return &SingleResult[Doctype]{SingleResult: c.Collection.FindOne(ctx, filter, opts...)}
}

func (c *Collector[Doctype]) InsertOne(
	ctx context.Context,
	document *Doctype,
	opts ...*options.InsertOneOptions,
) error {
	_, err := c.Collection.InsertOne(ctx, document, opts...)
	return err
}

func (c *Collector[Doctype]) UpdateByID(
	ctx context.Context,
	hexID string,
	update *Doctype,
	opts ...*options.UpdateOptions,
) (*mongo.UpdateResult, error) {
	oid, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return nil, err
	}
	(*update).SetUpdatedAt(time.Now())
	return c.Collection.UpdateByID(ctx, oid, bson.D{{"$set", update}})
}

func (c *Collection[Doctype]) DeleteByID(ctx context.Context, hexID string) error {
	oid, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return err
	}
	_, err = c.Collection.DeleteOne(ctx, bson.D{{"_id", oid}})
	return err
}

type SingleResult[T Doctype] struct {
	*mongo.SingleResult
}

func (r *SingleResult[_]) Err() error {
	return r.handleError(r.SingleResult.Err())
}

func (r *SingleResult[Doctype]) Decode() (*Doctype, error) {
	doc := new(Doctype)
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
