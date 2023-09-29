package handlers

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"

	"github.com/agustfricke/super-auth-go/config"
	"github.com/agustfricke/super-auth-go/database"
	"github.com/agustfricke/super-auth-go/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)


func SendEmail(token string, email string) {
	secretPassword := config.Config("EMAIL_SECRET_KEY")
	host := config.Config("HOST")
	auth := smtp.PlainAuth(
		"",
		"agustfricke@gmail.com",
		secretPassword,
		"smtp.gmail.com",
	)

	tmpl, err := template.ParseFiles("templates/verify_email.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	data := struct {
		Token string
		Host  string
	}{
		Token: token,
		Host:  host,
	}

	var bodyContent bytes.Buffer
	err = tmpl.Execute(&bodyContent, data)
	if err != nil {
		fmt.Println(err)
		return
	}

	content := fmt.Sprintf("To: %s\r\n"+
		"Subject: Verify Your Email Address\r\n"+
		"Content-Type: text/html; charset=utf-8\r\n"+
		"\r\n"+
		"%s", email, bodyContent.String())

	err = smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"agustfricke@gmail.com",
		[]string{email},
		[]byte(content),
	)
	if err != nil {
		fmt.Println(err)
	}
}

func VerifyEmail(c *fiber.Ctx) error {
    tokenString := c.Params("token")

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "You are not logged in"})
	}

	tokenByte, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
		}

    return []byte(config.Config("SECRET_KEY")), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("invalidate token: %v", err)})
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "invalid token claim"})

	}

	var user models.User
  db := database.DB
	db.First(&user, "id = ?", fmt.Sprint(claims["sub"]))

	if float64(user.ID) != claims["sub"] {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "the user belonging to this token no logger exists"})
	}

    user.Verified = true

	  if err := db.Save(&user).Error; err != nil {
		  return c.Status(fiber.StatusInternalServerError).SendString("Error al guardar los cambios en el usuario")
	  }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": user}})
}

