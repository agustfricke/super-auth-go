package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"

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

func SignInForm(c *fiber.Ctx) error {
  return c.Render("signin_form", fiber.Map{})
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
		return c.Status(fiber.StatusInternalServerError).SendString("Error creating user in database")
	}

	tokenByte := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)
  expDuration := time.Hour * 24

  claims["sub"] = user.ID
  claims["exp"] = now.Add(expDuration).Unix()
  claims["iat"] = now.Unix()
  claims["nbf"] = now.Unix()

	tokenString, err := tokenByte.SignedString([]byte(config.Config("JWT_SECRET")))

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("generating JWT Token failed: %v", err)})
	}

  SendEmail(tokenString, email)

  // Render mensaje de que se envio un correo para verificarlo
  return c.Render("verify", fiber.Map{
      "Email":    email,
  })
}

func VerifyEmail(c *fiber.Ctx) error {
	tokenString := c.Params("token")

	if tokenString == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Token faltante")
	}

	tokenByte, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
		}

		return []byte(config.Config("JWT_SECRET")), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("invalidate token: %v", err)})
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "invalid token claim"})

	}

	db := database.DB
	var user models.User

	if err := db.First(&user, "id = ?", fmt.Sprint(claims["sub"])).Error; err != nil {
		return c.Status(404).JSON(err)
	}

	if float64(user.ID) != claims["sub"] {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "the user belonging to this token no logger exists"})
	}

  user.Verified = new(bool)
  *user.Verified = true

	if err := db.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error al guardar los cambios en el usuario")
	}
  
  // Mensaje de email verificado
	return c.Render("success_verify", fiber.Map{})
}

func SignIn(c *fiber.Ctx) error {

	time.Sleep(2 * time.Second)

	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" || password == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing name, email or password")
	}

  var user models.User

  db := database.DB

	result := db.First(&user, "email = ?", strings.ToLower(email))
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid email or Password"})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid email or Password"})
	}

	if *user.Verified == false {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "Account not verificada"})
	}

	tokenByte := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)
  expDuration := time.Hour * 24

  claims["sub"] = user.ID
  claims["exp"] = now.Add(expDuration).Unix()
  claims["iat"] = now.Unix()
  claims["nbf"] = now.Unix()

	tokenString, err := tokenByte.SignedString([]byte(config.Config("JWT_SECRET")))

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("generating JWT Token failed: %v", err)})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   60 * 60,
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	return c.Render("home", fiber.Map{})
}

func LogoutUser(c *fiber.Ctx) error {
	expired := time.Now().Add(-time.Hour * 24)
  fmt.Println(expired)
	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   "",
		Expires: expired,
	})
	return c.Render("home", fiber.Map{})
}

func SignInGitHub(c *fiber.Ctx) error {
  // callback check si exite, si no, crearlo
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success"})
}

func SignInGoogle(c *fiber.Ctx) error {
  // callback check si exite, si no, crearlo
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success"})
}

