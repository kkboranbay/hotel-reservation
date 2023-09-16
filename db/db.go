package db

const (
	DBNAME     = "hotel-reservation"
	TestDBNAME = "hotel-reservation-test"
	DBURI      = "mongodb+srv://kkboranbay:Mongodb12@cluster0.elxfrcu.mongodb.net/?retryWrites=true&w=majority"
)

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
	Book  BookingStore
}
