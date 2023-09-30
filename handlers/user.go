package handlers

import (
	"github.com/agustfricke/super-auth-go/database"
	"github.com/agustfricke/super-auth-go/models"
	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
    db := database.DB 
    var users []models.User
    db.Find(&users)
    return c.JSON(users)
}

func GetMe(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
    "status": "success", 
    "data": fiber.Map{"user": user},
  })
}
