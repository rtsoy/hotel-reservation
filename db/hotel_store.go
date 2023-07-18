package db

import (
	"context"
	"github.com/rtsoy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HotelStore interface {
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	UpdateHotel(context.Context, bson.M, bson.M) error
	GetHotels(context.Context, *HotelQueryParams, *Pagination) ([]*types.Hotel, error)
	GetHotelByID(context.Context, primitive.ObjectID) (*types.Hotel, error)
}

type MongoHotelStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client:     client,
		collection: client.Database(DBNAME).Collection(hotelCollection),
	}
}

func NewMongoTestHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client:     client,
		collection: client.Database(TestDBNAME).Collection(hotelCollection),
	}
}

func (s *MongoHotelStore) GetHotelByID(ctx context.Context, oid primitive.ObjectID) (*types.Hotel, error) {
	var hotel types.Hotel
	if err := s.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&hotel); err != nil {
		return nil, err
	}
	return &hotel, nil
}

type HotelQueryParams struct {
	Pagination

	Rating   int
	Name     string
	Location string
}

func (s *MongoHotelStore) GetHotels(ctx context.Context, queryParams *HotelQueryParams, pagination *Pagination) ([]*types.Hotel, error) {
	// Default Pagination Values
	if pagination.Page == 0 {
		pagination.Page = int64(defaultPaginationPage)
	}
	if pagination.Limit == 0 {
		pagination.Limit = int64(defaultPaginationLimit)
	}

	// Check for empty values in filter
	filter := bson.M{}

	if queryParams.Rating >= 1 && queryParams.Rating <= 5 {
		filter["rating"] = queryParams.Rating
	}
	if len(queryParams.Name) > 1 {
		filter["name"] = queryParams.Name
	}
	if len(queryParams.Location) > 1 {
		filter["location"] = queryParams.Location
	}

	opts := &options.FindOptions{}

	opts.SetSkip((pagination.Page - 1) * pagination.Limit)
	opts.SetLimit(pagination.Limit)

	cur, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var hotels []*types.Hotel
	if err := cur.All(ctx, &hotels); err != nil {
		return nil, err
	}

	if len(hotels) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return hotels, nil
}

func (s *MongoHotelStore) UpdateHotel(ctx context.Context, filter bson.M, update bson.M) error {
	res, err := s.collection.UpdateOne(ctx, filter, update)

	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return err
}

func (s *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := s.collection.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}

	hotel.ID = res.InsertedID.(primitive.ObjectID)

	return hotel, nil
}
