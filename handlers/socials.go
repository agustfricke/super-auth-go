package handlers

import (

	"github.com/agustfricke/super-auth-go/socials"
	"github.com/gofiber/fiber/v2"
)

func AuthGoogle(c *fiber.Ctx) error {
	path := socials.ConfigGoogle()
	url := path.AuthCodeURL("state")
	return c.Redirect(url)

}

func CallbackGoogle(c *fiber.Ctx) error {
	token, error := socials.ConfigGoogle().Exchange(c.Context(), c.FormValue("code"))
	if error != nil {
		panic(error)
	}
	email := socials.GetEmail(token.AccessToken)
  /*
  if userID match with user.SocialID {
    Create new user in the database with the info of google
    generate a new token
  } else {
    generate a new token with the info
  }
  */
	return c.Status(200).JSON(fiber.Map{"email": email, "login": true})
}


func AuthGitHub(c *fiber.Ctx) error {
  return c.Status(200).JSON(fiber.Map{"state": "success"})
}

func CallbackGitHub(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{"state": "success"})
}

