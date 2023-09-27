package db

import (
	"context"
	"os"

	"github.com/kkboranbay/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	GetRooms(context.Context, Map) ([]*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection
	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		coll:       client.Database(os.Getenv(MongoDB_NAME_ENV_NAME)).Collection("rooms"),
		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	insertedRoom, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = insertedRoom.InsertedID.(primitive.ObjectID)

	filter := Map{"_id": room.HotelID}

	update := Map{"$push": bson.M{"rooms": room.ID}}
	if err := s.HotelStore.Update(ctx, filter, update); err != nil {
		return nil, err
	}

	return room, nil
}

func (s *MongoRoomStore) GetRooms(ctx context.Context, filter Map) ([]*types.Room, error) {
	if hotelIDStr, ok := filter["hotel_id"].(string); ok {
		oid, err := primitive.ObjectIDFromHex(hotelIDStr)
		if err != nil {
			return nil, err
		}
		filter["hotel_id"] = oid
	}

	resp, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var rooms []*types.Room
	if err := resp.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}
