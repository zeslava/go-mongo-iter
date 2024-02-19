package miter

import "context"

// Cursor interface of mongo.Cursor
type Cursor interface {
	Next(ctx context.Context) bool
	Decode(val any) error
	Err() error
	Close(ctx context.Context) error
}
