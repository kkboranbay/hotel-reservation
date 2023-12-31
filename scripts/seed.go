package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kkboranbay/hotel-reservation/db"
	"github.com/kkboranbay/hotel-reservation/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGO_DB_URL")))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(os.Getenv(db.MongoDB_NAME_ENV_NAME)).Drop(context.TODO()); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)
	store := &db.Store{
		User:  db.NewMongoUserStore(client),
		Hotel: hotelStore,
		Room:  db.NewMongoRoomStore(client, hotelStore),
		Book:  db.NewMongoBookingStore(client),
	}
	user := fixtures.AddUser(store, "leo", "ken", false)
	fmt.Println("user ->", user)

	admin := fixtures.AddUser(store, "admin", "admin", true)
	fmt.Println("admin ->", admin)

	hotel := fixtures.AddHotel(store, "Rixos", "Almaty", 5, nil)
	fmt.Println("hotel ->", hotel)

	room := fixtures.AddRoom(store, "large", true, 200.5, hotel.ID)
	fmt.Println("room ->", room)

	booking := fixtures.AddBooking(store, room.ID, user.ID, time.Now(), time.Now().AddDate(0, 0, 5))
	fmt.Println("booking ->", booking)

	for i := 0; i < 30; i++ {
		name := fmt.Sprintf("some random name %d", i)
		location := fmt.Sprintf("location %d", i)

		fixtures.AddHotel(store, name, location, rand.Intn(5)+1, nil)
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
