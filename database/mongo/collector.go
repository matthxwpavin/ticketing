package mongo

// import (
// 	"context"
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// type Doctype interface {
// 	SetUpdatedAt(t time.Time)
// }

// type Collector[T Doctype] struct {
// 	c *mongo.Collection
// }

// func (c *Collector[T]) FindOne(
// 	ctx context.Context,
// 	filter interface{},
// 	opts ...*options.FindOneOptions,
// ) *SingleResult[T] {
// 	return &SingleResult[T]{SingleResult: c.c.FindOne(ctx, filter, opts...)}
// }

// func (c *Collector[T]) Find(ctx context.Context, filter any) ([]*T, error) {
// 	cur, err := c.c.Find(ctx, filter)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var res []*T
// 	return res, cur.All(ctx, &res)
// }

// func (c *Collector[T]) InsertOne(
// 	ctx context.Context,
// 	document *T,
// 	opts ...*options.InsertOneOptions,
// ) error {
// 	_, err := c.c.InsertOne(ctx, document, opts...)
// 	return err
// }

// func (c *Collector[T]) UpdateByID(
// 	ctx context.Context,
// 	hexID string,
// 	update *T,
// 	opts ...*options.UpdateOptions,
// ) (*mongo.UpdateResult, error) {
// 	oid, err := primitive.ObjectIDFromHex(hexID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	(*update).SetUpdatedAt(time.Now())
// 	return c.c.UpdateByID(ctx, oid, bson.D{{"$set", update}})
// }

// func (c *Collection[T]) DeleteByID(ctx context.Context, hexID string) error {
// 	oid, err := primitive.ObjectIDFromHex(hexID)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = c.c.DeleteOne(ctx, bson.D{{"_id", oid}})
// 	return err
// }

// func (c *Collection[T]) Drop(ctx context.Context) error {
// 	return c.c.Drop(ctx)
// }

// type SingleResult[T Doctype] struct {
// 	*mongo.SingleResult
// }

// func (r *SingleResult[_]) Err() error {
// 	return r.handleError(r.SingleResult.Err())
// }

// func (r *SingleResult[T]) Decode() (*T, error) {
// 	doc := new(T)
// 	err := r.SingleResult.Decode(doc)
// 	if err == mongo.ErrNoDocuments {
// 		return nil, nil
// 	}
// 	return doc, nil
// }

// func (r *SingleResult[_]) Raw() (bson.Raw, error) {
// 	raw, err := r.SingleResult.Raw()
// 	return raw, r.handleError(err)
// }

// func (r *SingleResult[_]) handleError(err error) error {
// 	if err == mongo.ErrNoDocuments {
// 		return nil
// 	}
// 	return err
// }
