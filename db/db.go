package db

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	DBNAME     string
	DBURI      string
	TestDBNAME string
	TestDBURI  string
)

const (
	bookingCollection = "bookings"
	hotelCollection   = "hotels"
	roomCollection    = "rooms"
	userCollection    = "users"

	defaultPaginationPage  = 1
	defaultPaginationLimit = 10
)

type Pagination struct {
	Limit int64
	Page  int64
}

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

func init() {
	if err := godotenv.Load(); err != nil { // ../.env for tests
		log.Fatal("No .env file found")
	}

	DBNAME = os.Getenv("MONGO_DB_NAME")
	DBURI = os.Getenv("MONGO_DB_URI")
	TestDBNAME = os.Getenv("MONGO_TEST_DB_NAME")
	TestDBURI = os.Getenv("MONGO_TEST_DB_URI")
}
