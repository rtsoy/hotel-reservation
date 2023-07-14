package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rtsoy/hotel-reservation/api/errors"
	"github.com/rtsoy/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
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
		return errors.ErrInvalidID()
	}

	// TODO : 404
	booking, err := h.store.Booking.GetBookingByID(c.Context(), oid)
	if err != nil {
		return errors.ErrResourceNotFound()
	}

	user, ok := getAuthUser(c)
	if !ok {
		return errors.ErrUnauthorized()
	}

	if !user.IsAdmin && booking.UserID != user.ID {
		return errors.ErrForbidden()
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

	return c.JSON(genericResp{
		Type: "msg",
		Msg:  "updated",
	})
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	// TODO: 404
	bookings, err := h.store.Booking.GetBookings(c.Context(), nil)
	if err != nil {
		return errors.ErrResourceNotFound()
	}

	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.ErrInvalidID()
	}

	// TODO : 404
	booking, err := h.store.Booking.GetBookingByID(c.Context(), oid)
	if err != nil {
		return errors.ErrResourceNotFound()
	}

	user, ok := getAuthUser(c)
	if !ok {
		return c.Status(http.StatusUnauthorized).JSON("unauthorized")
	}
	if !user.IsAdmin && booking.UserID != user.ID {
		return c.Status(http.StatusForbidden).JSON(genericResp{
			Type: "error",
			Msg:  "not allowed",
		})
	}

	return c.JSON(booking)
}
