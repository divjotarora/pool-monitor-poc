package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbClient struct {
	ID       primitive.ObjectID // the Client ID
	client   *mongo.Client
	coll     *mongo.Collection
	numConns int
}

func newDbClient(ctx context.Context, uri string) (*dbClient, error) {
	newClient := &dbClient{
		ID: primitive.NewObjectID(),
	}

	monitor := &event.PoolMonitor{
		Event: newClient.HandlePoolEvent,
	}
	clientOpts := options.Client().ApplyURI(uri).SetPoolMonitor(monitor)

	var err error
	newClient.client, err = mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	newClient.coll = newClient.client.Database("foo").Collection("bar")
	return newClient, nil
}

func (d *dbClient) HandlePoolEvent(evt *event.PoolEvent) {
	switch evt.Type {
	case event.ConnectionCreated:
		d.numConns++
	case event.ConnectionClosed:
		d.numConns--
	}
}

func (d *dbClient) Close(ctx context.Context) {
	_ = d.coll.Drop(ctx)
	_ = d.client.Disconnect(ctx)
}

func (d *dbClient) Insert(ctx context.Context, document interface{}) error {
	_, err := d.coll.InsertOne(ctx, document)
	return err
}

func (d *dbClient) PrintStats() {
	fmt.Printf("Client Stats\nNum Connections Open: %d\n", d.numConns)
}
