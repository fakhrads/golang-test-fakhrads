// main.go

package main

import (
	"github.com/fakhrads/golang-test/routes"
	"github.com/fakhrads/golang-test/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	app.Use(logger.New())

	// Load environment variables
	utils.LoadEnv()

	// Initialize and migrate the database
	utils.InitDB()

	// Setup routes
	routes.SetupUserRoutes(app)

	// Start the server
	app.Listen(":3000")
}
