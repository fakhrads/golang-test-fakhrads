package routes

import (
	"github.com/fakhrads/golang-test/controllers"
	"github.com/fakhrads/golang-test/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	userGroup := app.Group("/user")

	// Middleware to validate API key
	userGroup.Use(middlewares.KeyValidator())

	userGroup.Post("/register", controllers.RegisterUser)
	userGroup.Get("/list", controllers.GetUserList)
	userGroup.Get("/:user_id", controllers.GetUserDetail)
	userGroup.Patch("/", controllers.UpdateUser)
}
