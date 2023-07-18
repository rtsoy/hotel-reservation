package main

import (
	"context"
	"fmt"
	"github.com/jaswdr/faker"
	"github.com/rtsoy/hotel-reservation/api"
	"github.com/rtsoy/hotel-reservation/db"
	"github.com/rtsoy/hotel-reservation/db/fixtures"
	"github.com/rtsoy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var (
	client *mongo.Client
	store  *db.Store
	ctx    = context.Background()
	fake   faker.Faker
)

func main() {
	// Users

	user1 := fixtures.AddUser(store, "Kyrie", "Irving",
		"KyrieIrving@example.org", "uncle_drew11", false)

	user2 := fixtures.AddUser(store, "Kevin", "Durant",
		"KevinDurant@example.org", "kd_trey35", false)

	admin := fixtures.AddUser(store, "Admin", "Admin",
		"admin@example.org", "admin", true)

	fmt.Printf("%s - user1 (%s : %s)\n", user1.ID.Hex(), user1.Email, "uncle_drew11")
	fmt.Printf("%s - user2 (%s : %s)\n", user2.ID.Hex(), user2.Email, "kd_trey35")
	fmt.Printf("%s - admin (%s : %s)\n", admin.ID.Hex(), admin.Email, "admin")

	users := []*types.User{
		user1,
		user2,
	}

	// Hotels

	var rooms []*types.Room

	for i := 0; i < 100; i++ {
		hotel := fixtures.AddHotel(store, fake.Company().Name(), fake.Address().City(),
			nil, fake.IntBetween(1, 5))

		fmt.Printf("%s - hotel\n", hotel.ID.Hex())

		// Rooms

		smallRoom := fixtures.AddRoom(store, "small", fake.Bool(),
			float64(fake.IntBetween(100, 500)), hotel.ID)
		rooms = append(rooms, smallRoom)

		mediumRoom := fixtures.AddRoom(store, "medium", fake.Bool(),
			float64(fake.IntBetween(500, 1000)), hotel.ID)
		rooms = append(rooms, mediumRoom)

		largeRoom := fixtures.AddRoom(store, "large", fake.Bool(),
			float64(fake.IntBetween(1000, 2000)), hotel.ID)
		rooms = append(rooms, largeRoom)
	}

	// Bookings

	for i := 0; i < 100; i++ {
		fromDate := time.Now().AddDate(0, fake.IntBetween(0, 3), fake.IntBetween(0, 30))
		tillDate := fromDate.AddDate(0, fake.IntBetween(0, 0), fake.IntBetween(3, 21))

		user := users[fake.IntBetween(0, len(users)-1)].ID
		room := rooms[fake.IntBetween(0, len(rooms)-1)].ID
		numPersons := fake.IntBetween(1, 4)

		params := types.BookRoomParams{
			FromDate:   fromDate,
			TillDate:   tillDate,
			NumPersons: numPersons,
		}

		isAvailable, err := api.IsRoomAvailableForBooking(context.Background(), store.Booking, room, params)
		if err != nil {
			log.Fatal(err)
		}
		if !isAvailable {
			i--
			continue
		}

		booking := fixtures.AddBooking(store, user, room, numPersons, fromDate, tillDate, fake.Bool())

		fmt.Printf("%s - booking\n", booking.ID.Hex())
	}
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

	fake = faker.New()
}
