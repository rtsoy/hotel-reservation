package api

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/rtsoy/hotel-reservation/api/middleware"
	"github.com/rtsoy/hotel-reservation/db/fixtures"
	"github.com/rtsoy/hotel-reservation/types"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestUserCancelNotOwnBooking(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t, tdb.client)

	var (
		creator = fixtures.AddUser(tdb.store, "user", "user",
			"user@example.org", "user", false)
		user = fixtures.AddUser(tdb.store, "user", "user",
			"user@example.org", "user", false)

		hotel   = fixtures.AddHotel(tdb.store, "testHotel", "Testestan", nil, 4)
		room    = fixtures.AddRoom(tdb.store, "medium", true, 199.9, hotel.ID)
		booking = fixtures.AddBooking(tdb.store, creator.ID, room.ID, 3,
			time.Now().AddDate(0, 0, 1).UTC(), time.Now().AddDate(0, 0, 8).UTC(), false)

		app   = fiber.New()
		route = app.Group("/", middleware.JWTAuthentication(tdb.store.User))

		bookingHandler = NewBookingHandler(tdb.store)
	)

	route.Get("/:id", bookingHandler.HandleCancelBooking)

	targetURL := "/" + booking.ID.Hex()
	req := httptest.NewRequest(http.MethodGet, targetURL, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Token", createTokenFromUser(user))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected status code 403 but got %d", resp.StatusCode)
	}

	var response genericResp
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	expectedResponse := genericResp{
		Type: "error",
		Msg:  "not allowed",
	}

	if !reflect.DeepEqual(response, expectedResponse) {
		t.Fatal("the response does not match an expected response")
	}
}

func TestAdminCancelNotOwnBooking(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t, tdb.client)

	var (
		user = fixtures.AddUser(tdb.store, "user", "user",
			"user@example.org", "user", false)
		admin = fixtures.AddUser(tdb.store, "admin", "admin",
			"admin@example.org", "admin", true)

		hotel   = fixtures.AddHotel(tdb.store, "testHotel", "Testestan", nil, 4)
		room    = fixtures.AddRoom(tdb.store, "medium", true, 199.9, hotel.ID)
		booking = fixtures.AddBooking(tdb.store, user.ID, room.ID, 3,
			time.Now().AddDate(0, 0, 1).UTC(), time.Now().AddDate(0, 0, 8).UTC(), false)

		app   = fiber.New()
		route = app.Group("/", middleware.JWTAuthentication(tdb.store.User))

		bookingHandler = NewBookingHandler(tdb.store)
	)

	route.Get("/:id", bookingHandler.HandleCancelBooking)

	targetURL := "/" + booking.ID.Hex()
	req := httptest.NewRequest(http.MethodGet, targetURL, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Token", createTokenFromUser(admin))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200 but got %d", resp.StatusCode)
	}

	var response genericResp
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	expectedResponse := genericResp{
		Type: "msg",
		Msg:  "updated",
	}

	if !reflect.DeepEqual(response, expectedResponse) {
		t.Fatal("the response does not match an expected response")
	}
}

func TestUserCancelOwnBooking(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t, tdb.client)

	var (
		user = fixtures.AddUser(tdb.store, "user", "user",
			"user@example.org", "user", false)

		hotel   = fixtures.AddHotel(tdb.store, "testHotel", "Testestan", nil, 4)
		room    = fixtures.AddRoom(tdb.store, "medium", true, 199.9, hotel.ID)
		booking = fixtures.AddBooking(tdb.store, user.ID, room.ID, 3,
			time.Now().AddDate(0, 0, 1).UTC(), time.Now().AddDate(0, 0, 8).UTC(), false)

		app   = fiber.New()
		route = app.Group("/", middleware.JWTAuthentication(tdb.store.User))

		bookingHandler = NewBookingHandler(tdb.store)
	)

	route.Get("/:id", bookingHandler.HandleCancelBooking)

	targetURL := "/" + booking.ID.Hex()
	req := httptest.NewRequest(http.MethodGet, targetURL, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Token", createTokenFromUser(user))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200 but got %d", resp.StatusCode)
	}

	var response genericResp
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	expectedResponse := genericResp{
		Type: "msg",
		Msg:  "updated",
	}

	if !reflect.DeepEqual(response, expectedResponse) {
		t.Fatal("the response does not match an expected response")
	}
}

func TestUserGetNotOwnBooking(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t, tdb.client)

	var (
		creator = fixtures.AddUser(tdb.store, "user", "user",
			"user@example.org", "user", false)
		user = fixtures.AddUser(tdb.store, "user", "user",
			"user@example.org", "user", false)

		hotel   = fixtures.AddHotel(tdb.store, "testHotel", "Testestan", nil, 4)
		room    = fixtures.AddRoom(tdb.store, "medium", true, 199.9, hotel.ID)
		booking = fixtures.AddBooking(tdb.store, creator.ID, room.ID, 3,
			time.Now().AddDate(0, 0, 1).UTC(), time.Now().AddDate(0, 0, 8).UTC(), false)

		app   = fiber.New()
		route = app.Group("/", middleware.JWTAuthentication(tdb.store.User))

		bookingHandler = NewBookingHandler(tdb.store)
	)

	route.Get("/:id", bookingHandler.HandleGetBooking)

	targetURL := "/" + booking.ID.Hex()
	req := httptest.NewRequest(http.MethodGet, targetURL, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Token", createTokenFromUser(user))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected status code 403 but got %d", resp.StatusCode)
	}

	var response genericResp
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	expectedResponse := genericResp{
		Type: "error",
		Msg:  "not allowed",
	}

	if !reflect.DeepEqual(response, expectedResponse) {
		t.Fatal("the response does not match an expected response")
	}
}

func TestAdminGetNotOwnBooking(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t, tdb.client)

	var (
		user = fixtures.AddUser(tdb.store, "user", "user",
			"user@example.org", "user", false)
		admin = fixtures.AddUser(tdb.store, "admin", "admin",
			"admin@example.org", "admin", true)

		hotel   = fixtures.AddHotel(tdb.store, "testHotel", "Testestan", nil, 4)
		room    = fixtures.AddRoom(tdb.store, "medium", true, 199.9, hotel.ID)
		booking = fixtures.AddBooking(tdb.store, user.ID, room.ID, 3,
			time.Now().AddDate(0, 0, 1).UTC(), time.Now().AddDate(0, 0, 8).UTC(), false)

		app   = fiber.New()
		route = app.Group("/", middleware.JWTAuthentication(tdb.store.User))

		bookingHandler = NewBookingHandler(tdb.store)
	)

	route.Get("/:id", bookingHandler.HandleGetBooking)

	targetURL := "/" + booking.ID.Hex()
	req := httptest.NewRequest(http.MethodGet, targetURL, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Token", createTokenFromUser(admin))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200 but got %d", resp.StatusCode)
	}

	var bookingResponse *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResponse); err != nil {
		t.Fatal(err)
	}

	if bookingResponse.ID != booking.ID {
		t.Fatalf("expected booking id %s but got %s", booking.ID.Hex(), bookingResponse.ID.Hex())
	}
}

func TestUserGetOwnBooking(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t, tdb.client)

	var (
		user = fixtures.AddUser(tdb.store, "user", "user",
			"user@example.org", "user", false)

		hotel   = fixtures.AddHotel(tdb.store, "testHotel", "Testestan", nil, 4)
		room    = fixtures.AddRoom(tdb.store, "medium", true, 199.9, hotel.ID)
		booking = fixtures.AddBooking(tdb.store, user.ID, room.ID, 3,
			time.Now().AddDate(0, 0, 1).UTC(), time.Now().AddDate(0, 0, 8).UTC(), false)

		app   = fiber.New()
		route = app.Group("/", middleware.JWTAuthentication(tdb.store.User))

		bookingHandler = NewBookingHandler(tdb.store)
	)

	route.Get("/:id", bookingHandler.HandleGetBooking)

	targetURL := "/" + booking.ID.Hex()
	req := httptest.NewRequest(http.MethodGet, targetURL, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Token", createTokenFromUser(user))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200 but got %d", resp.StatusCode)
	}

	var bookingResponse *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResponse); err != nil {
		t.Fatal(err)
	}

	if bookingResponse.ID != booking.ID {
		t.Fatalf("expected booking id %s but got %s", booking.ID.Hex(), bookingResponse.ID.Hex())
	}
}

func TestNoTokenGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t, tdb.client)

	var (
		user = fixtures.AddUser(tdb.store, "user", "user",
			"user@example.org", "user", false)

		hotel = fixtures.AddHotel(tdb.store, "testHotel", "Testestan", nil, 4)
		room  = fixtures.AddRoom(tdb.store, "medium", true, 199.9, hotel.ID)

		_ = fixtures.AddBooking(tdb.store, user.ID, room.ID, 3,
			time.Now().AddDate(0, 0, 1).UTC(), time.Now().AddDate(0, 0, 8).UTC(), false)
		_ = fixtures.AddBooking(tdb.store, user.ID, room.ID, 3,
			time.Now().AddDate(0, 0, 10).UTC(), time.Now().AddDate(0, 0, 17).UTC(), false)

		app   = fiber.New()
		admin = app.Group("/", middleware.JWTAuthentication(tdb.store.User), middleware.AdminAuth)

		bookingHandler = NewBookingHandler(tdb.store)
	)

	admin.Get("/", bookingHandler.HandleGetBookings)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Content-Type", "application/json")
	// req.Header.Add("X-Api-Token", createTokenFromUser(user))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status code 403 but got %d", resp.StatusCode)
	}

	var response string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response != "no token" {
		t.Fatalf("expected `no token` response message but got %s", response)
	}
}

func TestNonAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t, tdb.client)

	var (
		user = fixtures.AddUser(tdb.store, "user", "user",
			"user@example.org", "user", false)

		hotel = fixtures.AddHotel(tdb.store, "testHotel", "Testestan", nil, 4)
		room  = fixtures.AddRoom(tdb.store, "medium", true, 199.9, hotel.ID)

		_ = fixtures.AddBooking(tdb.store, user.ID, room.ID, 3,
			time.Now().AddDate(0, 0, 1).UTC(), time.Now().AddDate(0, 0, 8).UTC(), false)
		_ = fixtures.AddBooking(tdb.store, user.ID, room.ID, 3,
			time.Now().AddDate(0, 0, 10).UTC(), time.Now().AddDate(0, 0, 17).UTC(), false)

		app   = fiber.New()
		admin = app.Group("/", middleware.JWTAuthentication(tdb.store.User), middleware.AdminAuth)

		bookingHandler = NewBookingHandler(tdb.store)
	)

	admin.Get("/", bookingHandler.HandleGetBookings)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Token", createTokenFromUser(user))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected status code 403 but got %d", resp.StatusCode)
	}

	var response string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response != "forbidden" {
		t.Fatalf("expected `forbidden` response message but got %s", response)
	}
}

func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t, tdb.client)

	var (
		adminUser = fixtures.AddUser(tdb.store, "admin", "admin",
			"admin@example.org", "admin", true)
		user = fixtures.AddUser(tdb.store, "user", "user",
			"user@example.org", "user", false)

		hotel = fixtures.AddHotel(tdb.store, "testHotel", "Testestan", nil, 4)
		room  = fixtures.AddRoom(tdb.store, "medium", true, 199.9, hotel.ID)

		booking1 = fixtures.AddBooking(tdb.store, user.ID, room.ID, 3,
			time.Now().AddDate(0, 0, 1).UTC(), time.Now().AddDate(0, 0, 8).UTC(), false)
		booking2 = fixtures.AddBooking(tdb.store, user.ID, room.ID, 3,
			time.Now().AddDate(0, 0, 10).UTC(), time.Now().AddDate(0, 0, 17).UTC(), false)

		app   = fiber.New()
		admin = app.Group("/", middleware.JWTAuthentication(tdb.store.User), middleware.AdminAuth)

		bookingHandler = NewBookingHandler(tdb.store)
	)

	admin.Get("/", bookingHandler.HandleGetBookings)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Token", createTokenFromUser(adminUser))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200 but got %d", resp.StatusCode)
	}

	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}

	if len(bookings) != 2 {
		t.Fatalf("expected response length of 2 but got %d", len(bookings))
	}

	if booking1.ID != bookings[0].ID {
		t.Fatalf("expected booking id %s but got %s", booking1.ID.Hex(), bookings[0].ID.Hex())
	}
	if booking2.ID != bookings[1].ID {
		t.Fatalf("expected booking id %s but got %s", booking2.ID, bookings[1].ID)
	}
}
