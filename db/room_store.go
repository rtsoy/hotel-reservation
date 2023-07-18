package db

import (
	"context"
	"github.com/rtsoy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RoomStore interface {
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	GetRooms(context.Context, *RoomQueryParams, *Pagination) ([]*types.Room, error)
}

type MongoRoomStore struct {
	client     *mongo.Client
	collection *mongo.Collection
	hotelStore HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		collection: client.Database(DBNAME).Collection(roomCollection),
		hotelStore: hotelStore,
	}
}

func NewMongoTestRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		collection: client.Database(TestDBNAME).Collection(roomCollection),
		hotelStore: hotelStore,
	}
}

type RoomQueryParams struct {
	Pagination

	Size      string
	Seaside   *bool
	FromPrice float64
	ToPrice   float64
	HotelID   primitive.ObjectID
}

func (s *MongoRoomStore) GetRooms(ctx context.Context, queryParams *RoomQueryParams, pagination *Pagination) ([]*types.Room, error) {
	// Default Pagination Values
	if pagination.Page == 0 {
		pagination.Page = int64(defaultPaginationPage)
	}
	if pagination.Limit == 0 {
		pagination.Limit = int64(defaultPaginationLimit)
	}

	// Check for empty values in filter
	filter := bson.M{}

	if queryParams.Size == "small" || queryParams.Size == "medium" || queryParams.Size == "large" {
		filter["size"] = queryParams.Size
	}
	if queryParams.Seaside != nil {
		filter["seaside"] = queryParams.Seaside
	}
	if queryParams.FromPrice != 0 {
		filter["price"] = bson.M{
			"$gte": queryParams.FromPrice,
		}
	}
	if queryParams.ToPrice != 0 {
		filter["price"] = bson.M{
			"$lte": queryParams.ToPrice,
		}
	}
	if queryParams.FromPrice != 0 && queryParams.ToPrice != 0 {
		filter["price"] = bson.M{
			"$lte": queryParams.ToPrice,
			"$gte": queryParams.FromPrice,
		}
	}
	if queryParams.HotelID.Hex() != "000000000000000000000000" {
		filter["hotelID"] = queryParams.HotelID
	}

	opts := &options.FindOptions{}

	opts.SetSkip((pagination.Page - 1) * pagination.Limit)
	opts.SetLimit(pagination.Limit)

	cur, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var rooms []*types.Room
	if err := cur.All(ctx, &rooms); err != nil {
		return nil, err
	}

	if len(rooms) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return rooms, nil
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := s.collection.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}

	room.ID = res.InsertedID.(primitive.ObjectID)

	filter := bson.M{"_id": room.HotelID}
	update := bson.M{"$push": bson.M{"rooms": room.ID}}

	if err := s.hotelStore.UpdateHotel(ctx, filter, update); err != nil {
		return nil, err
	}

	return room, nil
}
