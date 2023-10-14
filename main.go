package main

import (
	"fiber-api/database"
	"fiber-api/routes"
	"github.com/gofiber/fiber/v2"
	"log"
)

func welcome(ctx *fiber.Ctx) error {
	return ctx.SendString("Welcome to the home page")
}

func setUpRoutes(app *fiber.App) {
	// welcome endpoint
	app.Get("/api", welcome)

	// User endpoint
	app.Post("/api/users", routes.CreateUser)
	app.Get("/api/users", routes.GetUsers)
	app.Get("/api/users/:id", routes.GetUser)
	app.Put("/api/users/:id", routes.UpdateUser)
	app.Delete("/api/users/:id", routes.DeleteUser)

	// Product route
	app.Post("/api/products", routes.CreateProduct)
	app.Get("/api/products", routes.GetProducts)
	app.Get("/api/products/:id", routes.GetProduct)
	app.Put("/api/products/:id", routes.UpdateProduct)
	app.Delete("/api/products/:id", routes.DeleteProduct)

	// Order route
	app.Get("/api/orders", routes.GetOrders)
	app.Get("/api/orders/:id", routes.GetOrder)
	app.Post("/api/orders", routes.CreateOrder)
}

func main() {
	database.ConnectDb()
	app := fiber.New()

	setUpRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
