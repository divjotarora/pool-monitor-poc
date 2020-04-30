package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbClient struct {
	client      *mongo.Client
	coll        *mongo.Collection
	poolMonitor *poolMonitor
}

func newDbClient(ctx context.Context, uri string) (*dbClient, error) {
	newClient := &dbClient{
		poolMonitor: newPoolMonitor(),
	}

	// Having PoolMonitor be an interface would have made this a little esaier because we could have instead done
	// SetPoolMonitor(newClient.poolMonitor).
	monitor := &event.PoolMonitor{
		Event: newClient.poolMonitor.HandleEvent,
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

func (d *dbClient) Close(ctx context.Context) {
	_ = d.coll.Drop(ctx)
	_ = d.client.Disconnect(ctx)
}

func (d *dbClient) Insert(ctx context.Context, document interface{}) error {
	_, err := d.coll.InsertOne(ctx, document)
	return err
}

func (d *dbClient) PrintStats() {
	fmt.Printf("Client Stats\nNum Connections Open: %d\n", d.poolMonitor.conns)
}
