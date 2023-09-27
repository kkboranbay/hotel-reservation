package db

var MongoDB_NAME_ENV_NAME = "MONGO_DB_NAME"

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
	Book  BookingStore
}

type Pagination struct {
	Page  int64
	Limit int64
}
