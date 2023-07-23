package test

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/pkg/errors"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// mongodbContainer represents the mongodb container type used in the module
type mongodbContainer struct {
	testcontainers.Container
}

// startContainer creates an instance of the mongodb container type
func startContainer(ctx context.Context) (*mongodbContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:6",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("Waiting for connections"),
			wait.ForListeningPort("27017/tcp"),
		),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	return &mongodbContainer{Container: container}, nil
}

type Document struct {
	ID   string `bson:"_id" fake:"skip"`
	Name string `bson:"name" fake:"{name}"`
}

func randDoc() Document {
	return Document{
		ID:   uuid.NewString(),
		Name: gofakeit.Name(),
	}
}

func randDocs(n int) map[string]Document {
	docs := make(map[string]Document, n)
	for i := 0; i < n; i++ {
		doc := randDoc()
		docs[doc.ID] = doc
	}
	return docs
}

func mongoInitData(ctx context.Context, client *mongo.Client, n int) (*mongo.Collection, map[string]Document, error) {
	col := client.Database("test").Collection("data")
	_, err := col.DeleteMany(ctx, bson.D{})
	if err != nil {
		return nil, nil, errors.Wrap(err, "clear mongo collection")
	}
	docs := randDocs(n)
	idocs := make([]interface{}, 0, len(docs))
	for _, doc := range docs {
		idocs = append(idocs, doc)
	}
	_, err = col.InsertMany(ctx, idocs)
	if err != nil {
		return nil, nil, errors.Wrap(err, "insert docs")
	}

	return col, docs, nil
}
