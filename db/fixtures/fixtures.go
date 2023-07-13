package fixtures

import (
	"context"
	"github.com/rtsoy/hotel-reservation/db"
	"github.com/rtsoy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

func AddBooking(store *db.Store, userID, roomID primitive.ObjectID, numPersons int, fromDate, tillDate time.Time, cancelled bool) *types.Booking {
	booking := &types.Booking{
		UserID:     userID,
		RoomID:     roomID,
		NumPersons: numPersons,
		FromDate:   fromDate,
		TillDate:   tillDate,
		Canceled:   cancelled,
	}

	insertBooking, err := store.Booking.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}

	return insertBooking
}

func AddRoom(store *db.Store, size string, seaside bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: seaside,
		Price:   price,
		HotelID: hotelID,
	}

	insertRoom, err := store.Room.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}

	return insertRoom
}

func AddHotel(store *db.Store, name, location string, rooms []primitive.ObjectID, rating int) *types.Hotel {
	roomIDs := rooms
	if rooms == nil {
		roomIDs = []primitive.ObjectID{}
	}

	hotel := &types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    roomIDs,
		Rating:   rating,
	}

	insertHotel, err := store.Hotel.InsertHotel(context.Background(), hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertHotel
}

func AddUser(store *db.Store, firstName, lastName, email, password string, isAdmin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	})
	if err != nil {
		log.Fatal(err)
	}

	user.IsAdmin = isAdmin

	insertUser, err := store.User.InsertUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	return insertUser
}
