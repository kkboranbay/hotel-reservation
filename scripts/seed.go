package main

import (
	"context"
	"log"

	"github.com/kkboranbay/hotel-reservation/db"
	"github.com/kkboranbay/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	ctx        = context.Background()
)

func seedHotel(name string, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []types.Room{
		{
			Size:  "small",
			Price: 99.9,
		},
		{
			Size:  "normal",
			Price: 159.9,
		},
		{
			Size:  "kingsize",
			Price: 222.9,
		},
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	seedHotel("Rixos", "Almaty", 5)
	seedHotel("Royal Tulip", "Almaty", 4)
	seedHotel("The Rits", "Astana", 5)
}

// the init function is a special function that is used to perform initialization tasks before
// the program's execution begins. The init function is automatically called
// by the Go runtime before the main function is executed.
func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	client.Database(db.DBNAME).Collection("hotels").Drop(ctx)
	client.Database(db.DBNAME).Collection("rooms").Drop(ctx)

	// if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
	// log.Fatal(err) // (AtlasError) user is not allowed to do action [dropDatabase] on [hotel-reservation.]
	// }

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
}
