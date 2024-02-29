# go-mongo-iter
Go MongoDB generic iterator

## Usage

### Single item iteration
```go
package main

import (
	"context"
	"fmt"
	miter "github.com/zeslava/go-mongo-iter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func single() error {
	ctx := context.Background()

	// Creating mongo client
	mongoClient, _ := mongo.Connect(ctx, options.Client().ApplyURI("localhost:27017"))
	defer mongoClient.Disconnect(ctx)

	// Define simple document type
	type Document struct {
		ID string `bson:"_id"`
	}

	// Find all documents in collection
	cursor, _ := mongoClient.Database("db").Collection("collection").Find(ctx, bson.D{})

	// Now creating iterator over cursor
	iter := miter.NewMongoIterSingle[Document](cursor)

	// Do not forget to close iterator after work
	defer iter.Close(ctx)

	// Iterating by documents
	for iter.Next(ctx) {
		// Take document
		doc := iter.Item()
		// Print doc id for example
		fmt.Printf("doc id: %s\n", doc.ID)
	}

	// Handling iteration error
	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := single(); err != nil {
		log.Fatal(err)
	}
}
```

### Batch items iteration
```go
package main

import (
	"context"
	"fmt"
	miter "github.com/zeslava/go-mongo-iter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func batch() error {
	ctx := context.Background()

	// Creating mongo client
	mongoClient, _ := mongo.Connect(ctx, options.Client().ApplyURI("localhost:27017"))
	defer mongoClient.Disconnect(ctx)

	// Define simple document type
	type Document struct {
		ID string `bson:"_id"`
	}

	// Find all documents in collection
	cursor, _ := mongoClient.Database("db").Collection("collection").Find(ctx, bson.D{})

	// Now creating iterator over cursor with batch = 10
	iterator := miter.NewMongoIterBatch[Document](cursor, 10)

	// Do not forget to close iterator after work
	defer iterator.Close(ctx)

	// Iterating by documents
	for iterator.Next(ctx) {
		// Take items
		docs := iterator.Items()
		// Iterating again ny batch just for example
		for _, doc := range docs {
			fmt.Printf("doc id: %s\n", doc.ID)
		}
	}

	// Handling iteration error
	if err := iterator.Err(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := batch(); err != nil {
		log.Fatal(err)
	}
}
```