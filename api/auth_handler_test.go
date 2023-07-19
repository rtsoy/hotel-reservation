package api

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/rtsoy/hotel-reservation/api/errors"
	"github.com/rtsoy/hotel-reservation/db/fixtures"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t, tdb.client)

	insertedUser := fixtures.AddUser(tdb.store, "James", "Harden",
		"jamesHarden13@example.com", "super_secret_password", false)

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

	_ = fixtures.AddUser(tdb.store, "James", "Harden",
		"jamesHarden13@example.com", "super_secret_password", false)

	app := fiber.New(fiber.Config{ErrorHandler: errors.ErrorHandler})
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

	expectedError := errors.ErrWrongCredentials()

	if resp.StatusCode != expectedError.Code {
		t.Fatalf("expected http status code %d but got %d", expectedError.Code, resp.StatusCode)
	}

	var errorResponse errors.Error
	if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedError, errorResponse) {
		t.Fatal("the error does not match an expected error")
	}
}
