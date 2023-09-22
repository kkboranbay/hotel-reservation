package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kkboranbay/hotel-reservation/db"
	"github.com/kkboranbay/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	roomStore    db.RoomStore
	hotelStore   db.HotelStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	ctx          = context.Background()
)

func seedUser(isAdmin bool, fname, lname, email string) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  "password",
	})
	if err != nil {
		log.Fatal(err)
	}

	user.IsAdmin = isAdmin

	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	return user
}

func seedHotel(name string, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func seedRoom(size string, ss bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := types.Room{
		Size:    size,
		Seaside: ss,
		Price:   price,
		HotelID: hotelID,
	}

	insertedRoom, err := roomStore.InsertRoom(ctx, &room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func seedBooking(roomID, userID primitive.ObjectID, from, till time.Time) {
	booking := types.Booking{
		RoomId:   roomID,
		UserId:   userID,
		FromDate: from,
		TillDate: till,
	}

	resp, err := bookingStore.InsertBooking(ctx, &booking)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("booking:", resp.ID)
}

func main() {
	hotel := seedHotel("Rixos", "Almaty", 5)
	seedHotel("Royal Tulip", "Almaty", 4)
	seedHotel("The Rits", "Astana", 5)

	room := seedRoom("small", true, 122.1, hotel.ID)
	seedRoom("big", false, 150.1, hotel.ID)
	seedRoom("middle", true, 200.1, hotel.ID)

	leo := seedUser(false, "Leo", "Ken", "leoken@gmail.com")
	seedUser(true, "Admin", "Admin", "admin@gmail.com")

	seedBooking(room.ID, leo.ID, time.Now(), time.Now().AddDate(0, 0, 2))
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
	client.Database(db.DBNAME).Collection("users").Drop(ctx)
	client.Database(db.DBNAME).Collection("bookings").Drop(ctx)

	// if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
	// log.Fatal(err) // (AtlasError) user is not allowed to do action [dropDatabase] on [hotel-reservation.]
	// }

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
	bookingStore = db.NewMongoBookingStore(client)
}
