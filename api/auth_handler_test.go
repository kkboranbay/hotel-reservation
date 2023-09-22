package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kkboranbay/hotel-reservation/db/fixtures"
)

func TestAuthSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	insertedUser := fixtures.AddUser(tdb.Store, "leo", "ken", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "leo@ken.com",
		Password: "leo_ken",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 but got %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}

	if authResp.Token == "" {
		t.Fatalf("expected the JWT token to be presented in the auth response")
	}

	// Set the encrypted password to an empty string, because we don't return
	// that in any JSON response.
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, authResp.User) {
		t.Fatalf("expected the user to be the inserted user")
	}
}

func TestAuthWithWrongPassword(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	fixtures.AddUser(tdb.Store, "leo", "ken", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "leoken@gmail.com",
		Password: "passwordNotCorrect",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 but got %d", resp.StatusCode)
	}

	var genericResp genericResponse
	if err := json.NewDecoder(resp.Body).Decode(&genericResp); err != nil {
		t.Fatal(err)
	}

	if genericResp.Type != "error" {
		t.Fatalf("expected type to be error but got %s", genericResp.Type)
	}

	if genericResp.Msg != "invalid credentials" {
		t.Fatalf("expected msg to be <invalid credentials> but got %s", genericResp.Msg)
	}
}
