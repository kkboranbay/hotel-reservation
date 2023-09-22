package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kkboranbay/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Book.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}

	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Book.GetBookingById(c.Context(), id)
	if err != nil {
		return err
	}

	authUser, err := getAuthUser(c)
	if err != nil {
		return err
	}

	if booking.UserId != authUser.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResponse{
			Type: "error",
			Msg:  "not authorized",
		})
	}

	return c.JSON(booking)
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Book.GetBookingById(c.Context(), id)
	if err != nil {
		return err
	}

	authUser, err := getAuthUser(c)
	if err != nil {
		return err
	}

	if booking.UserId != authUser.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResponse{
			Type: "error",
			Msg:  "not authorized",
		})
	}

	if err = h.store.Book.UpdateBooking(c.Context(), id, bson.M{"canceled": true}); err != nil {
		return err
	}

	return c.JSON(genericResponse{
		Type: "msg",
		Msg:  "updated",
	})
}
