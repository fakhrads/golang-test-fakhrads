package controllers

import (
	"database/sql"
	"fmt"

	"github.com/fakhrads/golang-test/models"
	"github.com/fakhrads/golang-test/services"
	"github.com/gofiber/fiber/v2"
)

func RegisterUser(c *fiber.Ctx) error {
	var newUser models.User
	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	if err := services.ValidateCreditCard(newUser.CreditCard); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Credit card data invalid"})
	}
	if err := services.ValidateRequest(&newUser); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	userID, err := services.SaveUserToDB(newUser)
	if err != nil {
		fmt.Println("Error saving user:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Something went wrong. Please try again later."})
	}
	return c.Status(200).JSON(fiber.Map{"user_id": userID})

}

func GetUserList(c *fiber.Ctx) error {
	q := c.Query("q")
	ob := c.Query("ob")
	sb := c.Query("sb")
	of := c.Query("of")
	lt := c.Query("lt")
	users, err := services.GetUsersFromDatabase(q, ob, sb, of, lt)

	if err != nil {
		fmt.Println("Error retrieving user list:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Something went wrong. Please try again later."})
	}
	response := fiber.Map{
		"count": len(users),
		"rows":  users,
	}

	return c.Status(200).JSON(response)
}

func GetUserDetail(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	user, err := services.GetUserDetailFromDatabase(userID)
	if err != nil {
		fmt.Println("Error retrieving user detail:", err)
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "User not found."})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Something went wrong. Please try again later."})
	}
	response := fiber.Map{
		"user_id":    user.UserID,
		"name":       user.Name,
		"email":      user.Email,
		"address":    user.Address,
		"photos":     user.Photos,
		"creditcard": user.CreditCard,
	}
	return c.Status(200).JSON(response)
}

func UpdateUser(c *fiber.Ctx) error {
	var updateUserRequest models.User
	if err := c.BodyParser(&updateUserRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Validate required parameters
	if updateUserRequest.UserID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Please provide user_id."})
	}

	if err := services.ValidateRequest(&updateUserRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if err := services.ValidateCreditCard(updateUserRequest.CreditCard); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if err := services.UpdateUserInDatabase(&updateUserRequest); err != nil {
		fmt.Println("Error updating user:", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	response := fiber.Map{"success": true}
	return c.Status(200).JSON(response)
}
