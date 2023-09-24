package main

import (
	"github.com/agustfricke/super-auth-go/database"
	"github.com/agustfricke/super-auth-go/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
)

func main() {
  database.ConnectDB()


  engine := html.New("./templates", ".html")

	app := fiber.New(fiber.Config{
	  Views: engine, 
    ViewsLayout: "layouts/main", 
	})

  app.Use(cors.New(cors.Config{
      AllowOrigins: "http://localhost:5173",
      AllowMethods: "GET, POST, PUT, DELETE",
      AllowCredentials: true,
      AllowHeaders: "Origin, Content-Type, Accept",
  }))

  app.Static("/", "./public")

  routes.Routes(app)

	app.Listen(":3000")
}
