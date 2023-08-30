package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kkboranbay/hotel-reservation/db"
	"github.com/kkboranbay/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	return &testdb{
		UserStore: db.NewMongoUserStore(client, db.TestDBNAME),
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: "Leo",
		LastName:  "Ken",
		Email:     "leoken@google.com",
		Password:  "dasasdasda",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req, 5000)
	if err != nil {
		t.Error(err)
	}
	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	if len(user.ID) == 0 {
		t.Errorf("expecting a user id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expecting the EncryptedPassword not to be included in the json response")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected FirstName %s but got %s", user.FirstName, params.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected LastName %s but got %s", user.LastName, params.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected Email %s but got %s", user.Email, params.Email)
	}
}

func TestGetUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	userHandler := NewUserHandler(tdb.UserStore)

	app := fiber.New()
	app.Get("/:id", userHandler.HandleGetUser)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: "Leo",
		LastName:  "Ken",
		Email:     "leoken@google.com",
		Password:  "dasasdasda",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req, 5000)
	if err != nil {
		t.Error(err)
	}
	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)

	id := user.ID.Hex()
	req = httptest.NewRequest("GET", "/"+string(id), nil)
	req.Header.Add("Content-Type", "application/json")
	resp, err = app.Test(req, 5000)
	if err != nil {
		t.Error(err)
	}
	var gotUser types.User
	json.NewDecoder(resp.Body).Decode(&gotUser)

	if gotUser.ID != user.ID {
		t.Errorf("expected %s but got %s", user.ID, gotUser.ID)
	}
}
