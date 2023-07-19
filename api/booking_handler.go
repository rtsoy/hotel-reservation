package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	myErrors "github.com/rtsoy/hotel-reservation/api/errors"
	"github.com/rtsoy/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return myErrors.ErrInvalidID()
	}

	booking, err := h.store.Booking.GetBookingByID(c.Context(), oid)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return myErrors.ErrResourceNotFound()
		}

		return err
	}

	user, ok := getAuthUser(c)
	if !ok {
		return myErrors.ErrUnauthorized()
	}

	if !user.IsAdmin && booking.UserID != user.ID {
		return myErrors.ErrForbidden()
	}

	filter := bson.M{"_id": booking.ID}
	update := bson.M{
		"$set": bson.M{
			"canceled": true,
		},
	}
	if err := h.store.Booking.UpdateBooking(c.Context(), filter, update); err != nil {
		return err
	}

	return c.JSON(map[string]string{
		"updated": id,
	})
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	var bookingQueryParams db.BookingQueryParams
	if err := c.QueryParser(&bookingQueryParams); err != nil {
		return myErrors.ErrBadRequest()
	}

	bookings, err := h.store.Booking.GetBookings(c.Context(), &bookingQueryParams, &bookingQueryParams.Pagination)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return myErrors.ErrResourceNotFound()
		}

		return err
	}

	response := &resourceResponse{
		Results: len(bookings),
		Page:    bookingQueryParams.Page,
		Data:    bookings,
	}

	return c.JSON(response)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return myErrors.ErrInvalidID()
	}

	booking, err := h.store.Booking.GetBookingByID(c.Context(), oid)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return myErrors.ErrResourceNotFound()
		}

		return myErrors.ErrResourceNotFound()
	}

	user, ok := getAuthUser(c)
	if !ok {
		return myErrors.ErrUnauthorized()
	}
	if !user.IsAdmin && booking.UserID != user.ID {
		return myErrors.ErrForbidden()
	}

	return c.JSON(booking)
}
