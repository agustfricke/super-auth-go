package routes

import (
	"github.com/agustfricke/super-auth-go/handlers"
	"github.com/agustfricke/super-auth-go/middleware"
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	app.Get("/", handlers.Home)
	app.Post("/signin", handlers.SignIn)
	app.Post("/signup", handlers.SignUp)
	app.Get("/signup/form", handlers.SignUpForm)
	app.Get("/signin/form", handlers.SignInForm)
  app.Get("/:token", handlers.VerifyEmail)
	app.Post("/signin/google", handlers.SignInGitHub)
	app.Post("/signin/github", handlers.SignInGoogle)

  app.Get("/:token", handlers.VerifyEmail)
	app.Get("/users/me", middleware.Auth, handlers.Home)
}


