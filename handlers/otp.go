package handlers

import (
	"github.com/agustfricke/super-auth-go/database"
	"github.com/agustfricke/super-auth-go/models"
	"github.com/gofiber/fiber"
	"github.com/pquerna/otp/totp"
)

func GenerateOTP(c *fiber.Ctx) error {
	  tokenUser := c.Locals("user").(*models.User)

    key, err := totp.Generate(totp.GenerateOpts{
        Issuer:      "Tech con Agust",
        AccountName: tokenUser.Email,
        SecretSize:  15,
    })

    if err != nil {
        panic(err)
    }

    var user models.User
    db := database.DB
    result := db.First(&user, "id = ?", tokenUser.ID)
    if result.Error != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "fail",
            "message": "Correo electrónico o contraseña no válidos",
        })
    }

    dataToUpdate := models.User{
        Otp_secret:   key.Secret(),
        Otp_auth_url: key.URL(),
    }

    db.Model(&user).Updates(dataToUpdate)

    otpResponse := fiber.Map{
        "base32":      key.Secret(),
        "otpauth_url": key.URL(),
    }

    return c.JSON(otpResponse)
}

func VerifyOTP(c *fiber.Ctx) error {
    var payload *models.OTPInput
	  tokenUser := c.Locals("user").(*models.User)

    if err := c.BodyParser(&payload); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "fail",
            "message": err.Error(),
        })
    }

    var user models.User
    db := database.DB
    result := db.First(&user, "id = ?", tokenUser.ID)
    if result.Error != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "fail",
            "message": "El token no es válido o el usuario no existe",
        })
    }

    valid := totp.Validate(payload.Token, user.Otp_secret)
    if !valid {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "fail",
            "message": "El token no es válido o el usuario no existe",
        })
    }

    dataToUpdate := models.User{
        Otp_enabled:  true,
        Otp_verified: true,
    }

    db.Model(&user).Updates(dataToUpdate)

    userResponse := fiber.Map{
        "id":          user.ID,
        "name":        user.Name,
        "email":       user.Email,
        "otp_enabled": user.Otp_enabled,
    }

    return c.JSON(fiber.Map{
        "otp_verified": true,
        "user":         userResponse,
    })
}

func DisableOTP(c *fiber.Ctx) error {
	  tokenUser := c.Locals("user").(*models.User)

    var user models.User
    db := database.DB
    result := db.First(&user, "id = ?", tokenUser.ID)
    if result.Error != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "fail",
            "message": "El usuario no existe",
        })
    }

    user.Otp_enabled = false
    db.Save(&user)

    userResponse := fiber.Map{
        "id":          user.ID,
        "name":        user.Name,
        "email":       user.Email,
        "otp_enabled": user.Otp_enabled,
    }

    return c.JSON(fiber.Map{
        "otp_disabled": true,
        "user":         userResponse,
    })
}
