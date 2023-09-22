package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/dgrijalva/jwt-go"

	"github.com/agustfricke/super-auth-go/config"
	"github.com/agustfricke/super-auth-go/database"
	"github.com/agustfricke/super-auth-go/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
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

func SignUpForm(c *fiber.Ctx) error {
  return c.Render("signup_form", fiber.Map{})
}


func SignUp(c *fiber.Ctx) error {
  db := database.DB

	time.Sleep(2 * time.Second)

	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	if name == "" || email == "" || password == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing name, email or password")
	}

  user := models.User{Name: name, Email: email, Password: string(hashedPassword)}
	if err := db.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error creating task in database")
	}

	tokenByte := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)

	claims["sub"] = user.ID
	claims["exp"] = now.Add(30).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := tokenByte.SignedString([]byte(config.Config("JWT_SECRET")))

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("generating JWT Token failed: %v", err)})
	}

  SendEmail(tokenString, email)

  return c.Render("verify", fiber.Map{
      "user":    user,
  })
}

func VerifyEmail(c *fiber.Ctx) error {
	tokenString := c.Params("token")

	if tokenString == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Token o ID faltante")
	}

	tokenKey := []byte(config.Config("JWT_SECRET"))
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return tokenKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return c.Status(fiber.StatusUnauthorized).SendString("Token JWT no válido")
		}
		return c.Status(fiber.StatusBadRequest).SendString("Token JWT no válido")
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		return c.Status(fiber.StatusBadRequest).SendString("ID de usuario no encontrado en el token")
	}

  fmt.Println(userID)

	db := database.DB
	var user models.User

	if err := db.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(err)
	}

  user.Verified = new(bool)
  *user.Verified = true

	if err := db.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error al guardar los cambios en el usuario")
	}

	return c.Render("success_verify", fiber.Map{
    "user": user,
  })
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
