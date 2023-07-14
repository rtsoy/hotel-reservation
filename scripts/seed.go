package main

import (
	"context"
	"fmt"
	"github.com/rtsoy/hotel-reservation/db"
	"github.com/rtsoy/hotel-reservation/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var (
	client *mongo.Client
	store  *db.Store
	ctx    = context.Background()
)

func main() {
	user := fixtures.AddUser(store, "Kyrie", "Irving",
		"KyrieIrving@example.org", "uncle_drew11", false)

	admin := fixtures.AddUser(store, "Admin", "Admin",
		"admin@example.org", "admin", true)

	fmt.Printf("%s - user (%s : %s)\n", user.ID.Hex(), user.Email, "uncle_drew11")
	fmt.Printf("%s - admin (%s : %s)\n", admin.ID.Hex(), admin.Email, "admin")

	hotel := fixtures.AddHotel(store, "The Beverly Hills Hotel",
		"Los Angeles", nil, 5)

	fmt.Printf("%s - hotel\n", hotel.ID.Hex())

	mediumRoom := fixtures.AddRoom(store, "medium", true,
		799.99, hotel.ID)
	largeRoom := fixtures.AddRoom(store, "large", true,
		1299.9, hotel.ID)

	fmt.Printf("%s - mediumRoom\n", mediumRoom.ID.Hex())
	fmt.Printf("%s - largeRoom\n", largeRoom.ID.Hex())

	booking := fixtures.AddBooking(store, user.ID, largeRoom.ID, 2,
		time.Now(), time.Now().AddDate(0, 0, 7), false)

	fmt.Printf("%s - booking\n", booking.ID.Hex())
}

func init() {
	var err error

	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	store = &db.Store{}

	store.User = db.NewMongoUserStore(client)
	store.Hotel = db.NewMongoHotelStore(client)
	store.Room = db.NewMongoRoomStore(client, store.Hotel)
	store.Booking = db.NewMongoBookingStore(client)
}
