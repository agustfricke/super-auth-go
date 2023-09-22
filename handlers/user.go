package handlers

import (
	"github.com/agustfricke/super-auth-go/models"
	"github.com/gofiber/fiber/v2"
)


func Home(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": user}})
}
