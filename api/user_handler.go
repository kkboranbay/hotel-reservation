package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kkboranbay/hotel-reservation/types"
)

// The fiber.Ctx type represents the context of an incoming HTTP request and the corresponding response. 
// It contains various methods and fields that allow you to interact with and handle the request and response.
// Here's an explanation of what fiber.Ctx represents and some of its key features:
// HTTP Request Information:
	// Ctx.Method(): Returns the HTTP method used in the request (e.g., "GET", "POST", "PUT", etc.).
	// Ctx.Path(): Returns the URL path of the request.
	// Ctx.Params(): Provides access to route parameters extracted from the URL path.
	// Ctx.Query(): Allows you to retrieve query parameters from the request's URL.
	// Ctx.Body(): Retrieves the request body data.
// Response Handling:
	// Ctx.Status(code int): Sets the HTTP status code for the response.
	// Ctx.JSON(v interface{}): Serializes the provided data as JSON and sends it as the response.
	// Ctx.Send(data []byte): Sends raw data as the response.
	// Ctx.Render(template string, data interface{}): Renders an HTML template with the provided data.
// HTTP Headers:
	// Ctx.Get(header string): Retrieves a specific header value from the request.
	// Ctx.Set(header string, value string): Sets a custom header in the response.
	// Ctx.Append(header string, value string): Appends a value to an existing response header.
// an so on...
func HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "Leo",
		LastName: "Ken",
	}

	return c.JSON(u)
}

func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON("Leo")
}