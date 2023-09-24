package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kkboranbay/hotel-reservation/db/fixtures"
	"github.com/kkboranbay/hotel-reservation/types"
)

func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		adminUser = fixtures.AddUser(tdb.Store, "admin", "admin", true)
		user      = fixtures.AddUser(tdb.Store, "leo", "ken", false)
		hotel     = fixtures.AddHotel(tdb.Store, "Rixos", "Almaty", 5, nil)
		room      = fixtures.AddRoom(tdb.Store, "large", true, 100.12, hotel.ID)
		from      = time.Now()
		till      = time.Now().AddDate(0, 0, 4)
		booking   = fixtures.AddBooking(tdb.Store, room.ID, user.ID, from, till)
		app       = fiber.New(fiber.Config{
			ErrorHandler: ErrorHandler,
		})
		admin          = app.Group("/", JWTAuthentication(tdb.User), AdminAccess)
		bookingHandler = NewBookingHandler(tdb.Store)
	)

	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Auth-Token", createTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("not 200 response %d", resp.StatusCode)
	}

	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}

	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking but got %d", len(bookings))
	}

	have := bookings[0]
	if booking.ID != have.ID {
		t.Fatal("expected bookings to be equal")
	}

	if booking.UserId != have.UserId {
		t.Fatalf("expected %s but got %s", booking.UserId, have.UserId)
	}

	// test non-admin cannot access the bookings
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Auth-Token", createTokenFromUser(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expect unauthorized code but got %d", resp.StatusCode)
	}
}

func TestUserGetBooking(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		user    = fixtures.AddUser(tdb.Store, "leo", "ken", false)
		john    = fixtures.AddUser(tdb.Store, "john", "john", false)
		hotel   = fixtures.AddHotel(tdb.Store, "Rixos", "Almaty", 5, nil)
		room    = fixtures.AddRoom(tdb.Store, "small", true, 100, hotel.ID)
		booking = fixtures.AddBooking(tdb.Store, room.ID, user.ID, time.Now(), time.Now().AddDate(0, 0, 2))
		app     = fiber.New(fiber.Config{
			ErrorHandler: ErrorHandler,
		})
		bookingHandler = NewBookingHandler(tdb.Store)
		api            = app.Group("/", JWTAuthentication(tdb.User))
	)

	api.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Auth-Token", createTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("not 200 response %d", resp.StatusCode)
	}

	var bookingResp *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}

	if booking.ID != bookingResp.ID {
		t.Fatal("expected bookings to be equal")
	}

	if booking.UserId != bookingResp.UserId {
		t.Fatalf("expected %s but got %s", booking.UserId, bookingResp.UserId)
	}

	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Auth-Token", createTokenFromUser(john))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expect unauthorized code but got %d", resp.StatusCode)
	}
}
