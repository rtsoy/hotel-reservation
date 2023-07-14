package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	myErrors "github.com/rtsoy/hotel-reservation/api/errors"
	"github.com/rtsoy/hotel-reservation/db"
	"github.com/rtsoy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return myErrors.ErrResourceNotFound()
		}

		return err
	}

	return c.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params types.BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return myErrors.ErrBadRequest()
	}

	if err := params.Validate(); err != nil {
		return myErrors.NewError(http.StatusBadRequest, err.Error())
	}

	roomID := c.Params("id")
	roomOID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return myErrors.ErrInvalidID()
	}

	user, ok := getAuthUser(c)
	if !ok {
		return myErrors.ErrUnauthorized()
	}

	available, err := isRoomAvailableForBooking(c.Context(), h.store.Booking, roomOID, params)
	if err != nil {
		return err
	}
	if !available {
		return myErrors.NewError(http.StatusBadRequest, fmt.Sprintf("Room %s is already booked", roomID))
	}

	booking := &types.Booking{
		UserID:     user.ID,
		RoomID:     roomOID,
		NumPersons: params.NumPersons,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
	}

	insertedBooking, err := h.store.Booking.InsertBooking(c.Context(), booking)
	if err != nil {
		return err
	}

	return c.JSON(insertedBooking)
}

func isRoomAvailableForBooking(ctx context.Context, bookingStore db.BookingStore, roomID primitive.ObjectID, params types.BookRoomParams) (bool, error) {
	filter := bson.M{
		"roomID": roomID,
		"fromDate": bson.M{
			"$lte": params.TillDate,
		},
		"tillDate": bson.M{
			"$gte": params.FromDate,
		},
	}

	bookings, err := bookingStore.GetBookings(ctx, filter)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return false, err
	}

	ok := len(bookings) == 0

	return ok, nil
}
