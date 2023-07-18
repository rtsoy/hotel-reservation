package db

import (
	"context"
	"github.com/rtsoy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type BookingStore interface {
	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetBookings(context.Context, *BookingQueryParams, *Pagination) ([]*types.Booking, error)
	GetBookingByID(context.Context, primitive.ObjectID) (*types.Booking, error)
	UpdateBooking(context.Context, bson.M, bson.M) error
}

type MongoBookingStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	return &MongoBookingStore{
		client:     client,
		collection: client.Database(DBNAME).Collection(bookingCollection),
	}
}

func NewMongoTestBookingStore(client *mongo.Client) *MongoBookingStore {
	return &MongoBookingStore{
		client:     client,
		collection: client.Database(TestDBNAME).Collection(bookingCollection),
	}
}

func (s *MongoBookingStore) UpdateBooking(ctx context.Context, filter bson.M, update bson.M) error {
	res, err := s.collection.UpdateOne(ctx, filter, update)

	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return err
}

func (s *MongoBookingStore) GetBookingByID(ctx context.Context, oid primitive.ObjectID) (*types.Booking, error) {
	var booking types.Booking
	if err := s.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking); err != nil {
		return nil, err
	}

	return &booking, nil
}

type BookingQueryParams struct {
	Pagination

	UserID     primitive.ObjectID
	RoomID     primitive.ObjectID
	NumPersons int
	FromDate   time.Time
	TillDate   time.Time
	Canceled   *bool
}

func (s *MongoBookingStore) GetBookings(ctx context.Context, queryParams *BookingQueryParams, pagination *Pagination) ([]*types.Booking, error) {
	// Default Pagination Values
	if pagination.Page == 0 {
		pagination.Page = int64(defaultPaginationPage)
	}
	if pagination.Limit == 0 {
		pagination.Limit = int64(defaultPaginationLimit)
	}

	// Check for empty values in filter
	filter := bson.M{}

	if queryParams.UserID.Hex() != "000000000000000000000000" {
		filter["userID"] = queryParams.UserID
	}
	if queryParams.RoomID.Hex() != "000000000000000000000000" {
		filter["roomID"] = queryParams.RoomID
	}
	if queryParams.NumPersons != 0 {
		filter["numPersons"] = queryParams.NumPersons
	}
	if queryParams.TillDate.String() != "0001-01-01 00:00:00 +0000 UTC" {
		filter["tillDate"] = bson.M{
			"$lte": queryParams.TillDate,
		}
	}
	if queryParams.FromDate.String() != "0001-01-01 00:00:00 +0000 UTC" {
		filter["fromDate"] = bson.M{
			"$gte": queryParams.FromDate,
		}
	}
	if queryParams.TillDate.String() != "0001-01-01 00:00:00 +0000 UTC" && queryParams.FromDate.String() != "0001-01-01 00:00:00 +0000 UTC" {
		filter["fromDate"] = bson.M{
			"$lte": queryParams.TillDate,
		}
		filter["tillDate"] = bson.M{
			"$gte": queryParams.FromDate,
		}
	}
	if queryParams.Canceled != nil {
		filter["canceled"] = queryParams.Canceled
	}

	opts := &options.FindOptions{}

	opts.SetSkip((pagination.Page - 1) * pagination.Limit)
	opts.SetLimit(pagination.Limit)

	cur, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var bookings []*types.Booking
	if err := cur.All(ctx, &bookings); err != nil {
		return nil, err
	}

	if len(bookings) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return bookings, nil
}

func (s *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	res, err := s.collection.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}

	booking.ID = res.InsertedID.(primitive.ObjectID)

	return booking, nil
}
