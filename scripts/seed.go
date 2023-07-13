package main

import (
	"context"
	"fmt"
	"github.com/rtsoy/hotel-reservation/db"
	"github.com/rtsoy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	client     *mongo.Client
	userStore  db.UserStore
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	ctx        = context.Background()
)

func seedUser(firstName, lastName, email, password string) {
	userParams := types.CreateUserParams{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	}

	user, err := types.NewUserFromParams(userParams)
	if err != nil {
		log.Fatal(err)
	}

	insertUser, err := userStore.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(insertUser)
}

func seedHotel(name, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insertHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	rooms := []types.Room{
		{
			Size:    "small",
			Price:   99.9,
			HotelID: insertHotel.ID,
		},
		{
			Size:    "kingsize",
			Price:   199.9,
			HotelID: insertHotel.ID,
		},
		{
			Size:    "normal",
			Price:   129.9,
			HotelID: insertHotel.ID,
		},
	}

	for _, room := range rooms {
		insertRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(insertRoom)
	}

	fmt.Println(insertHotel)
}

func main() {
	client.Database(db.DBNAME).Collection("users").Drop(ctx)
	client.Database(db.DBNAME).Collection("hotels").Drop(ctx)
	client.Database(db.DBNAME).Collection("rooms").Drop(ctx)

	seedUser("Lebron", "James", "legoat6@gmail.com", "goat6")
	seedUser("Kevin", "Durant", "traykd35@gmail.com", "goat35")
	seedUser("Steph", "Curry", "chiefCurry@gmail.com", "goat30")

	seedHotel("Bellucia", "France", 4)
	seedHotel("The Cozy Hotel", "Amsterdam", 3)
	seedHotel("Don't die in your sleep", "London", 1)
	seedHotel("Eleon", "Moscow", 5)
}

func init() {
	var err error

	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	userStore = db.NewMongoUserStore(client)
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
}
