package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/rtsoy/hotel-reservation/db"
	"github.com/rtsoy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testdb struct {
	client *mongo.Client
	store  *db.Store
}

func (tdb *testdb) teardown(t *testing.T, client *mongo.Client) {
	if err := client.Database(db.TestDBNAME).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		t.Fatal(err)
	}

	return &testdb{
		client: client,
		store: &db.Store{
			User:  db.NewMongoTestUserStore(client),
			Hotel: db.NewMongoTestHotelStore(client),
			Room:  db.NewMongoTestRoomStore(client, db.NewMongoTestHotelStore(client)),
		},
	}
}

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

	_, err := insertTestUser(tdb.store.User, types.CreateUserParams{
		FirstName: "user1",
		LastName:  "user1",
		Email:     "user1@example.com",
		Password:  "user1password",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = insertTestUser(tdb.store.User, types.CreateUserParams{
		FirstName: "user2",
		LastName:  "user2",
		Email:     "user2@example.com",
		Password:  "user2password",
	})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	var users *[]types.User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		t.Fatal(err)
	}

	if len(*users) != 2 {
		t.Fatalf("expected %d length of response but got %d", 2, len(*users))
	}
}

func TestGetUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t, tdb.client)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.store.User)
	app.Get("/:id", userHandler.HandleGetUser)

	expectedUser, err := insertTestUser(tdb.store.User, types.CreateUserParams{
		FirstName: "James",
		LastName:  "Harden",
		Email:     "jamesHarden13@example.com",
		Password:  "qwerty123",
	})
	if err != nil {
		t.Fatal(err)
	}

	targetURL := "/" + expectedUser.ID.Hex()

	req := httptest.NewRequest(http.MethodGet, targetURL, nil)
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
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

	user, err := insertTestUser(tdb.store.User, types.CreateUserParams{
		FirstName: "James",
		LastName:  "Harden",
		Email:     "jamesHarden13@example.com",
		Password:  "qwerty123",
	})
	if err != nil {
		t.Fatal(err)
	}
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

	user, err := insertTestUser(tdb.store.User, types.CreateUserParams{
		FirstName: "James",
		LastName:  "Harden",
		Email:     "jamesHarden13@example.com",
		Password:  "qwerty123",
	})
	if err != nil {
		t.Fatal(err)
	}
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
