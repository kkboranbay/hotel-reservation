package main

import (
	"flag"

	"github.com/kkboranbay/hotel-reservation/api"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// This line defines a command-line flag named "listenAddr". 
	// The flag allows you to specify the port on which the API server will listen. 
	// The default value is ":5000".
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	app := fiber.New()
	apiv1 := app.Group("/api/v1")

	apiv1.Get("/user", api.HandleGetUsers)
	apiv1.Get("/user/:id", api.HandleGetUser)
	app.Listen(*listenAddr)

	// GET /user
	// {
	// 	"firstName": "Leo",
	// 	"lastName": "Ken"
	// }
}
