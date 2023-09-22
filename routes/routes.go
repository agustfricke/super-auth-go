package routes

import (
	"github.com/agustfricke/super-auth-go/handlers"
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	app.Post("/signin", handlers.SignIn)
	app.Post("/signup", handlers.SignUp)
	app.Post("/verify", handlers.VerifyEmail)
	app.Post("/signin/google", handlers.SignInGitHub)
	app.Post("/signin/github", handlers.SignInGoogle)
}


