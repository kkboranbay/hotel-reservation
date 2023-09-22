package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/kkboranbay/hotel-reservation/db/fixtures"
)

func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	user := fixtures.AddUser(tdb.Store, "leo", "ken", false)
	hotel := fixtures.AddHotel(tdb.Store, "Rixos", "Almaty", 5, nil)
	room := fixtures.AddRoom(tdb.Store, "large", true, 100.12, hotel.ID)

	from := time.Now()
	till := time.Now().AddDate(0, 0, 4)
	booking := fixtures.AddBooking(tdb.Store, room.ID, user.ID, from, till)
	fmt.Println(booking)
}
