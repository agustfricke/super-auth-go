package handlers

import (
	"github.com/agustfricke/super-auth-go/database"
	"github.com/agustfricke/super-auth-go/models"
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
    googleResponse := socials.GetGoogleResponse(token.AccessToken)

    db := database.DB 
    var user models.User
    db.Find(&user, googleResponse.ID)

    if err := db.First(&user, "SocialID = ?", googleResponse.ID).Error; err != nil {
        user = models.User{
            SocialID:       googleResponse.ID,
            Email:          googleResponse.Email,
        }
        db.Create(&user)
      return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "created", "user created": user, "google res": googleResponse})
    } else {
      return c.Status(fiber.StatusFound).JSON(fiber.Map{"status": "found", "user in db": user, "google res": googleResponse})
    }
}

func AuthGitHub(c *fiber.Ctx) error {
  return c.Status(200).JSON(fiber.Map{"state": "success"})
}

func CallbackGitHub(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{"state": "success"})
}

