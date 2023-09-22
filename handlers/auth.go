package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func SignUp(c *fiber.Ctx) error {
  // instancia db
  // Agarar input desde formulrio
  // encriptar password con bcrypt
  // crear usuario en db
  // mandar email con token
  // retornar html de aviso para verificar
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success"})
}

func VerifyEmail(c *fiber.Ctx) error {
  // obtener token desde url
  // si esta en la url activar la cuenta
  // retornar html de exito
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success"})
}

func SignIn(c *fiber.Ctx) error {
  // instancia db
  // check si el usuario esta activado si no mostrar error 
  // Agarar input desde formulrio
  // desencriptar password
  // si ok crear token
  // retornar html a profile route
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success"})
}

func SignInGitHub(c *fiber.Ctx) error {
  // callback check si exite, si no, crearlo
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success"})
}

func SignInGoogle(c *fiber.Ctx) error {
  // callback check si exite, si no, crearlo
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success"})
}
