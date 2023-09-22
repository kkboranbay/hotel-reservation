package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kkboranbay/hotel-reservation/db"
	"github.com/kkboranbay/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(store *db.Store, fn, ln string, isAdmin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fn,
		LastName:  ln,
		Email:     fmt.Sprintf("%s@%s.com", fn, ln),
		Password:  fmt.Sprintf("%s_%s", fn, ln),
	})
	if err != nil {
		log.Fatal(err)
	}

	user.IsAdmin = isAdmin

	insertedUser, err := store.User.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	return insertedUser
}

func AddHotel(store *db.Store, name, loc string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	roomsID := rooms
	if rooms == nil {
		roomsID = []primitive.ObjectID{}
	}
	hotel := types.Hotel{
		Name:     name,
		Location: loc,
		Rooms:    roomsID,
		Rating:   rating,
	}

	insertedHotel, err := store.Hotel.InsertHotel(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func AddRoom(store *db.Store, size string, ss bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := types.Room{
		Size:    size,
		Seaside: ss,
		Price:   price,
		HotelID: hotelID,
	}

	insertedRoom, err := store.Room.InsertRoom(context.TODO(), &room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func AddBooking(store *db.Store, rid, uid primitive.ObjectID, from, till time.Time) *types.Booking {
	booking := types.Booking{
		RoomId:   rid,
		UserId:   uid,
		FromDate: from,
		TillDate: till,
	}

	insertedBooking, err := store.Book.InsertBooking(context.TODO(), &booking)
	if err != nil {
		log.Fatal(err)
	}

	return insertedBooking
}
