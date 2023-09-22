package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	RoomId     primitive.ObjectID `bson:"room_id,omitempty" json:"room_id,omitempty"`
	UserId     primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
	NumPersons int                `bson:"num_persons,omitempty" json:"num_persons,omitempty"`
	FromDate   time.Time          `bson:"from_date,omitempty" json:"from_date,omitempty"`
	TillDate   time.Time          `bson:"till_date,omitempty" json:"till_date,omitempty"`
	Canceled   bool               `bson:"canceled" json:"canceled"`
}
