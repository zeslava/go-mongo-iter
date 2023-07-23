package miter

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoIterSingle[T any] interface {
	Next(ctx context.Context) bool
	Item() T
	Close(ctx context.Context) error
	Err() error
}

func NewMongoIterSingle[T any](cur *mongo.Cursor) MongoIterSingle[T] {
	return &MongoIterSingleImpl[T]{cur: cur}
}

type MongoIterSingleImpl[T any] struct {
	cur  *mongo.Cursor
	item T
	err  error
}

func (i *MongoIterSingleImpl[T]) Next(ctx context.Context) bool {
	if i.err != nil {
		return false
	}
	if i.cur.Next(ctx) {
		var item T
		i.err = i.cur.Decode(&item)
		i.item = item
		return true
	}
	return false
}

func (i *MongoIterSingleImpl[T]) Item() T {
	return i.item
}

func (i *MongoIterSingleImpl[_]) Err() error {
	return i.cur.Err()
}

func (i *MongoIterSingleImpl[_]) Close(ctx context.Context) error {
	return i.cur.Close(ctx)
}
