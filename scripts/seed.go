package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kkboranbay/hotel-reservation/db"
	"github.com/kkboranbay/hotel-reservation/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(context.TODO()); err != nil {
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
}
