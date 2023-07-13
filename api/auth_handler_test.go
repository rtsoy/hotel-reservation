package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/rtsoy/hotel-reservation/db"
	"github.com/rtsoy/hotel-reservation/types"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func insertTestUser(userStore db.UserStore, params types.CreateUserParams) (*types.User, error) {
	userParams := types.CreateUserParams{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
		Password:  params.Password,
	}

	user, err := types.NewUserFromParams(userParams)
	if err != nil {
		return nil, err
	}

	insertedUser, err := userStore.InsertUser(context.TODO(), user)
	if err != nil {
		return nil, err
	}

	return insertedUser, nil
}

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t, tdb.client)

	insertedUser, err := insertTestUser(tdb.store.User, types.CreateUserParams{
		FirstName: "James",
		LastName:  "Harden",
		Email:     "jamesHarden13@example.com",
		Password:  "super_secret_password",
	})
	if err != nil {
		t.Fatal(err)
	}

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.store.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "jamesHarden13@example.com",
		Password: "super_secret_password",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if uint(resp.StatusCode) != http.StatusOK {
		t.Fatalf("expected http status code %d but got %d", http.StatusOK, resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}

	if authResp.Token == "" {
		t.Fatalf("expected the JWT to be present in the auth response")
	}

	insertedUser.EncryptedPassword = ""

	if !reflect.DeepEqual(insertedUser, authResp.User) {
		t.Fatal("expected the user to be inserted user")
	}
}

func TestAuthenticateWithWrongPasswordFailure(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t, tdb.client)

	_, err := insertTestUser(tdb.store.User, types.CreateUserParams{
		FirstName: "James",
		LastName:  "Harden",
		Email:     "jamesHarden13@example.com",
		Password:  "super_secret_password",
	})
	if err != nil {
		t.Fatal(err)
	}

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.store.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "jamesHarden13@example.com",
		Password: "wrong_super_secret_password",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if uint(resp.StatusCode) != fiber.StatusBadRequest {
		t.Fatalf("expected http status code %d but got %d", fiber.StatusInternalServerError, resp.StatusCode)
	}

	var genResp genericResp
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}

	if genResp.Type != "error" {
		t.Fatalf("expected response type to be `error` but got %s", genResp.Type)
	}

	if genResp.Msg != "invalid credentials" {
		t.Fatalf("expected response message to be `invalid credentials` but got %s", genResp.Msg)
	}
}
