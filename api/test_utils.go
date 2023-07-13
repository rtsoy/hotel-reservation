package api

import (
	"context"
	"github.com/rtsoy/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

type testdb struct {
	client *mongo.Client
	store  *db.Store
}

func (tdb *testdb) teardown(t *testing.T, client *mongo.Client) {
	if err := client.Database(db.TestDBNAME).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.TestDBURI))
	if err != nil {
		t.Fatal(err)
	}

	return &testdb{
		client: client,
		store: &db.Store{
			User:    db.NewMongoTestUserStore(client),
			Hotel:   db.NewMongoTestHotelStore(client),
			Room:    db.NewMongoTestRoomStore(client, db.NewMongoTestHotelStore(client)),
			Booking: db.NewMongoTestBookingStore(client),
		},
	}
}
