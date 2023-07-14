package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rtsoy/hotel-reservation/api/errors"
	"github.com/rtsoy/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.ErrInvalidID()
	}

	// TODO : 404
	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), oid)
	if err != nil {
		return err
	}

	return c.JSON(hotel)
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.ErrInvalidID()
	}
	filter := bson.M{"hotelID": oid}

	// TODO : 404
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return err
	}

	return c.JSON(rooms)
}

func QueriesToBSON(queries map[string]string) bson.M {
	res := bson.M{}

	for key, value := range queries {
		res[key] = value
	}

	return res
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	filter := QueriesToBSON(c.Queries())

	// TODO : 404
	hotels, err := h.store.Hotel.GetHotels(c.Context(), filter)
	if err != nil {
		return err
	}

	return c.JSON(hotels)
}
