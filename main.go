package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	ctx := context.TODO()
	uri := "mongodb://localhost:27017"

	client, err := newDbClient(ctx, uri)
	if err != nil {
		panic(err)
	}
	defer func() {
		client.Close(ctx)
		client.PrintStats()
	}()

	client.PrintStats()

	if err = client.Insert(ctx, bson.D{{"x", 1}}); err != nil {
		panic(err)
	}

	client.PrintStats()
}
