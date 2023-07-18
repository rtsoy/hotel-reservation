package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/rtsoy/hotel-reservation/api"
	"github.com/rtsoy/hotel-reservation/api/errors"
	"github.com/rtsoy/hotel-reservation/api/middleware"
	"github.com/rtsoy/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func main() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	var (
		app = fiber.New(fiber.Config{
			ErrorHandler: errors.ErrorHandler,
		})

		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		userStore    = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)

		store = &db.Store{
			User:    userStore,
			Hotel:   hotelStore,
			Room:    roomStore,
			Booking: bookingStore,
		}

		apiv1 = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
		auth  = app.Group("/api")
		admin = apiv1.Group("/admin", middleware.AdminAuth)

		userHandler    = api.NewUserHandler(userStore)
		authHandler    = api.NewAuthHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
	)

	// Auth Handlers

	auth.Post("/auth", authHandler.HandleAuthenticate)

	// User Handlers

	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("user/:id", userHandler.HandleDeleteUser)

	// Hotel Handlers

	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// Room Handlers

	apiv1.Get("/room", roomHandler.HandleGetRooms)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	// Bookings Handlers

	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	// Admin Routes

	admin.Get("/booking", bookingHandler.HandleGetBookings)

	listenAddr := os.Getenv("LISTEN_ADDR")
	log.Fatal(app.Listen(listenAddr))
}
