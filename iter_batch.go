package miter

import (
	"context"
)

// MongoIterBatch interface for batch items iteration
type MongoIterBatch[T any] interface {
	Next(ctx context.Context) bool
	Items() []T
	Close(ctx context.Context) error
	Err() error
}

func NewMongoIterBatch[T any](cur Cursor, size int) MongoIterBatch[T] {
	return &MongoIterBatchImpl[T]{
		cur:   cur,
		size:  size,
		items: make([]T, 0, size),
	}
}

type MongoIterBatchImpl[T any] struct {
	cur   Cursor
	size  int
	items []T
	err   error
}

func (i *MongoIterBatchImpl[T]) Next(ctx context.Context) bool {
	if i.err != nil {
		return false
	}
	i.items = i.items[:0]
	for i.cur.Next(ctx) {
		var item T
		if err := i.cur.Decode(&item); err != nil {
			i.err = err
			return len(i.items) > 0
		}
		i.items = append(i.items, item)
		if len(i.items) >= i.size {
			break
		}
	}

	return len(i.items) > 0
}

func (i *MongoIterBatchImpl[T]) Items() []T {
	return i.items
}

func (i *MongoIterBatchImpl[T]) Err() error {
	return i.cur.Err()
}

func (i *MongoIterBatchImpl[T]) Close(ctx context.Context) error {
	return i.cur.Close(ctx)
}
