package api

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/kkboranbay/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.client.Database(os.Getenv(db.MongoDB_NAME_ENV_NAME)).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGO_DB_URL_TEST")))
	if err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)
	return &testdb{
		client: client,
		Store: &db.Store{
			User:  db.NewMongoUserStore(client),
			Hotel: hotelStore,
			Room:  db.NewMongoRoomStore(client, hotelStore),
			Book:  db.NewMongoBookingStore(client),
		},
	}
}
