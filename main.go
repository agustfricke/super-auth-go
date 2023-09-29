package main

import (
	"github.com/agustfricke/super-auth-go/database"
	"github.com/agustfricke/super-auth-go/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
  database.ConnectDB()

  app := fiber.New()

  app.Use(cors.New(cors.Config{
      AllowOrigins: "http://localhost:5173",
      AllowMethods: "GET, POST",
      AllowCredentials: true,
      AllowHeaders: "Origin, Content-Type, Accept",
  }))

  router.SetupRoutes(app)

	app.Listen(":8080")
}
