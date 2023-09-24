package handlers

import (
	"github.com/agustfricke/super-auth-go/models"
	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	return c.Render("home", fiber.Map{})
}

func User(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	return c.Render("profile", fiber.Map{
    "User": user,
  })
}
