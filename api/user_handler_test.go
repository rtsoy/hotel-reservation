package api

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/rtsoy/hotel-reservation/db/fixtures"
	"github.com/rtsoy/hotel-reservation/types"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t, tdb.client)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.store.User)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: "James",
		LastName:  "Harden",
		Email:     "jamesHarden13@example.com",
		Password:  "qwerty123",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status code 201 but got %d", resp.StatusCode)
	}

	var user types.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Error(err)
	}

	if len(user.ID) == 0 {
		t.Errorf("expecting a user id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expecting the encryptedPassword not to be included in the json response")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected firstName %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected lastName %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email %s but got %s", params.Email, user.Email)
	}
}

func TestGetUsers(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t, tdb.client)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.store.User)
	app.Get("/", userHandler.HandleGetUsers)

	_ = fixtures.AddUser(tdb.store, "user1", "user1",
		"user1@example.com", "user1password", false)

	_ = fixtures.AddUser(tdb.store, "user2", "user2",
		"user2@example.com", "user2password", false)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200 but got %d", resp.StatusCode)
	}

	var response resourceResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.Results != 2 {
		t.Fatalf("expected results 2 but got %d", response.Results)
	}
}

func TestGetUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t, tdb.client)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.store.User)
	app.Get("/:id", userHandler.HandleGetUser)

	expectedUser := fixtures.AddUser(tdb.store, "James", "Harden",
		"jamesHarden13@example.com", "qwerty123", false)

	targetURL := "/" + expectedUser.ID.Hex()

	req := httptest.NewRequest(http.MethodGet, targetURL, nil)
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200 but got %d", resp.StatusCode)
	}

	var user *types.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Fatal(err)
	}

	if expectedUser.ID != user.ID {
		t.Fatalf("expected firstname %s but got %s", expectedUser.ID.Hex(), user.ID.Hex())
	}
	if expectedUser.Email != user.Email {
		t.Fatalf("expected email %s but got %s", expectedUser.Email, user.Email)
	}
	if expectedUser.LastName != user.LastName {
		t.Fatalf("expected lastname %s but got %s", expectedUser.LastName, user.LastName)
	}
	if expectedUser.FirstName != user.FirstName {
		t.Fatalf("expected firstname %s but got %s", expectedUser.FirstName, user.FirstName)
	}
}

func TestUpdateUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t, tdb.client)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.store.User)
	app.Put("/:id", userHandler.HandlePutUser)

	user := fixtures.AddUser(tdb.store, "James", "Harden",
		"jamesHarden13@example.com", "qwerty123", false)

	userID := user.ID.Hex()

	params := types.UpdateUserParams{
		FirstName: "Jarden",
		LastName:  "Hames",
	}
	b, _ := json.Marshal(params)

	targetURL := "/" + userID

	req := httptest.NewRequest(http.MethodPut, targetURL, bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200 but got %d", resp.StatusCode)
	}

	var updatedResponse map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&updatedResponse); err != nil {
		t.Fatal(err)
	}

	if _, ok := updatedResponse["updated"]; !ok {
		t.Fatal("expected response with map where updated is a key")
	}

	if updatedResponse["updated"] != userID {
		t.Fatalf("expected deleted user id %s but got %s", userID, updatedResponse["deleted"])
	}
}

func TestDeleteUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t, tdb.client)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.store.User)
	app.Delete("/:id", userHandler.HandleDeleteUser)

	user := fixtures.AddUser(tdb.store, "James", "Harden",
		"jamesHarden13@example.com", "qwerty123", false)

	userID := user.ID.Hex()

	targetURL := "/" + userID

	req := httptest.NewRequest(http.MethodDelete, targetURL, nil)
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	var deletedResponse map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&deletedResponse); err != nil {
		t.Fatal(err)
	}

	if _, ok := deletedResponse["deleted"]; !ok {
		t.Fatal("expected response with map where deleted is a key")
	}

	if deletedResponse["deleted"] != userID {
		t.Fatalf("expected deleted user id %s but got %s", userID, deletedResponse["deleted"])
	}
}
