package router

import (
	"github.com/agustfricke/super-auth-go/handlers"
	"github.com/agustfricke/super-auth-go/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
  // Basic auth 
  app.Post("/signup", handlers.SignUp) // tested OK
  app.Post("/signin", handlers.SignIn) // tested OK
  app.Get("/logout", handlers.Logout) // tested OK
  // OTP
  app.Post("/generate", middleware.DeserializeUser, handlers.GenerateOTP) // tested OK
  app.Post("/verify", middleware.DeserializeUser, handlers.VerifyOTP) // tested OK
  app.Post("/disable", middleware.DeserializeUser, handlers.DisableOTP) // tested OK
  // Email
	app.Get("/verify/:token", handlers.VerifyEmail) // tested OK
  // Google auth
	app.Get("/auth/google", handlers.AuthGoogle) // tested OK
	app.Get("/auth/google/callback", handlers.CallbackGoogle) 
  // GitHub auth
  app.Get("/auth/github", handlers.AuthGitHub)
  app.Get("/auth/github/callback", handlers.CallbackGitHub)

  app.Get("/users", handlers.GetUsers) 
}
