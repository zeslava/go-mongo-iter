package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	miter "github.com/zeslava/go-mongo-iter"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestIterSingle(t *testing.T) {
	ctx := context.Background()

	container, err := startContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	endpoint, err := container.Endpoint(ctx, "mongodb")
	if err != nil {
		t.Error(fmt.Errorf("failed to get endpoint: %w", err))
	}

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
	if err != nil {
		t.Fatal(fmt.Errorf("error creating mongo client: %w", err))
	}

	t.Cleanup(func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			t.Logf("disconnect mongo: %v", err)
		}

		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	const docsCount = 10
	col, docs, err := mongoInitData(ctx, mongoClient, docsCount)
	if err != nil {
		t.Error(fmt.Errorf("mongo init data: %w", err))
		return
	}
	assert.Len(t, docs, docsCount, "docs count")

	var failed bool
	cur, err := col.Find(ctx, bson.D{})
	if err != nil {
		t.Error(fmt.Errorf("mongo find data: %w", err))
		return
	}
	iterCount := 0
	iter := miter.NewMongoIterSingle[Document](cur)
	defer iter.Close(ctx)
	for iter.Next(ctx) {
		doc := iter.Item()
		d, ok := docs[doc.ID]
		if ok && doc == d {
			iterCount++
			t.Logf("OK: %#v - %#v", d, doc)
			continue
		}
		t.Logf("BAD: %#v - %#v", d, doc)
		failed = true
		break
	}
	if err = iter.Err(); err != nil {
		t.Error(fmt.Errorf("mongo iter data: %w", err))
		return
	}
	assert.Equal(t, docsCount, iterCount, "iter count")

	if failed {
		t.Fatal("failed")
	}
}
