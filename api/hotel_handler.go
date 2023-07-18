package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	myErrors "github.com/rtsoy/hotel-reservation/api/errors"
	"github.com/rtsoy/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		return myErrors.ErrInvalidID()
	}

	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), oid)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return myErrors.ErrResourceNotFound()
		}

		return err
	}

	return c.JSON(hotel)
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return myErrors.ErrInvalidID()
	}

	roomQueryParams := &db.RoomQueryParams{
		HotelID: oid,
	}

	rooms, err := h.store.Room.GetRooms(c.Context(), roomQueryParams, &roomQueryParams.Pagination)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return myErrors.ErrResourceNotFound()
		}

		return err
	}

	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var hotelQueryParams db.HotelQueryParams
	if err := c.QueryParser(&hotelQueryParams); err != nil {
		return myErrors.ErrBadRequest()
	}

	hotels, err := h.store.Hotel.GetHotels(c.Context(), &hotelQueryParams, &hotelQueryParams.Pagination)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return myErrors.ErrResourceNotFound()
		}

		return err
	}

	response := &resourceResponse{
		Results: len(hotels),
		Page:    hotelQueryParams.Page,
		Data:    hotels,
	}

	return c.JSON(response)
}
