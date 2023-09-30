package handlers

import (
	"fmt"
	"time"

	"github.com/agustfricke/super-auth-go/config"
	"github.com/agustfricke/super-auth-go/database"
	"github.com/agustfricke/super-auth-go/models"
	"github.com/agustfricke/super-auth-go/socials"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
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

  if err := db.First(&user, googleResponse.ID).Error; err != nil {
    user = models.User{
      SocialID:       googleResponse.ID,
      Email:          googleResponse.Email,
    }
    db.Create(&user)
    c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "created", "user created": user, "google res": googleResponse})
  } else {
    c.Status(fiber.StatusFound).JSON(fiber.Map{"status": "found", "user in db": user, "google res": googleResponse})
  }

  tokenByte := jwt.New(jwt.SigningMethodHS256)

  now := time.Now().UTC()
  claims := tokenByte.Claims.(jwt.MapClaims)
  expDuration := time.Hour * 24

  claims["sub"] = user.ID
  claims["exp"] = now.Add(expDuration).Unix()
  claims["iat"] = now.Unix()
  claims["nbf"] = now.Unix()

  tokenString, err := tokenByte.SignedString([]byte(config.Config("SECRET_KEY")))

  if err != nil {
    return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("generating JWT Token failed: %v", err)})
  }

  c.Cookie(&fiber.Cookie{
    Name:     "token",
    Value:    tokenString,
    Path:     "/",
    MaxAge:   24 * 60 * 60,
    Secure:   false,
    HTTPOnly: true,
    Domain:   "localhost",
  })

  return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "token": tokenString})
}

func AuthGitHub(c *fiber.Ctx) error {
  return c.Status(200).JSON(fiber.Map{"state": "success"})
}

func CallbackGitHub(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{"state": "success"})
}

