package db

const (
	DBNAME     = "hotel-reservation"
	TestDBNAME = "hotel-reservation-test"
	DBURI      = "mongodb://localhost:27017"
	TestDBURI  = "mongodb://localhost:27017"

	bookingCollection = "bookings"
	hotelCollection   = "hotels"
	roomCollection    = "rooms"
	userCollection    = "users"
)

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}
