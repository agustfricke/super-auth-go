package handlers

import (
	  "github.com/gofiber/fiber/v2"
    "log"
    "github.com/gofiber/fiber/v2"
    "golang.org/x/oauth2"
    "github.com/google/go-github/v38/github"
)

// Auth fiber handler
func AuthGoogle(c *fiber.Ctx) error {
	path := auth.ConfigGoogle()
	url := path.AuthCodeURL("state")
	return c.Redirect(url)

}

// Callback to receive google's response
func CallbackGoogle(c *fiber.Ctx) error {
	token, error := auth.ConfigGoogle().Exchange(c.Context(), c.FormValue("code"))
	if error != nil {
		panic(error)
	}
	email := auth.GetEmail(token.AccessToken)
	return c.Status(200).JSON(fiber.Map{"email": email, "login": true})
}

