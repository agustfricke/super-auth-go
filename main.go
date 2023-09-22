package main

import (
	"github.com/agustfricke/super-auth-go/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
  database.ConnectDB()

  engine := html.New("./templates", ".html")

	app := fiber.New(fiber.Config{
	  Views: engine, 
    ViewsLayout: "layouts/main", 
	})

  app.Static("/", "./public")

	app.Listen(":3000")
}
