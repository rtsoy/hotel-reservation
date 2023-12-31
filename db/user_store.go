package db

import (
	"context"
	"github.com/rtsoy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserStore interface {
	GetUserByID(context.Context, primitive.ObjectID) (*types.User, error)
	GetUserByEmail(context.Context, string) (*types.User, error)
	GetUsers(context.Context, *UserQueryParams, *Pagination) ([]*types.User, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
	DeleteUser(context.Context, primitive.ObjectID) error
	UpdateUser(context.Context, bson.M, bson.M) error
}

type MongoUserStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client:     client,
		collection: client.Database(DBNAME).Collection(userCollection),
	}
}

func NewMongoTestUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client:     client,
		collection: client.Database(TestDBNAME).Collection(userCollection),
	}
}

func (s *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	var user types.User
	if err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, filter bson.M, update bson.M) error {
	res, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, oid primitive.ObjectID) error {
	res, err := s.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (s *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = res.InsertedID.(primitive.ObjectID)

	return user, nil
}

type UserQueryParams struct {
	Pagination

	FirstName string
	LastName  string
	Email     string
	IsAdmin   *bool
}

func (s *MongoUserStore) GetUsers(ctx context.Context, queryParams *UserQueryParams, pagination *Pagination) ([]*types.User, error) {
	// Default Pagination Values
	if pagination.Page == 0 {
		pagination.Page = int64(defaultPaginationPage)
	}
	if pagination.Limit == 0 {
		pagination.Limit = int64(defaultPaginationLimit)
	}

	// Check for empty values in filter
	filter := bson.M{}

	if len(queryParams.FirstName) > 1 {
		filter["firstName"] = queryParams.FirstName
	}
	if len(queryParams.LastName) > 1 {
		filter["lastName"] = queryParams.LastName
	}
	if len(queryParams.Email) > 1 {
		filter["email"] = queryParams.Email
	}
	if queryParams.IsAdmin != nil {
		filter["isAdmin"] = queryParams.IsAdmin
	}

	opts := &options.FindOptions{}

	opts.SetSkip((pagination.Page - 1) * pagination.Limit)
	opts.SetLimit(pagination.Limit)

	cur, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var users []*types.User
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return users, nil

}

func (s *MongoUserStore) GetUserByID(ctx context.Context, oid primitive.ObjectID) (*types.User, error) {
	var user types.User
	if err := s.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
