package router

import (
	"github.com/agustfricke/super-auth-go/handlers"
	"github.com/agustfricke/super-auth-go/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
  // Basic auth 
  app.Post("/signup", handlers.SignUp) 
  app.Post("/signin", handlers.SignIn) 
  app.Get("/logout", handlers.Logout)
  // OTP
  app.Post("/generate", middleware.DeserializeUser, handlers.GenerateOTP)
  app.Post("/verify", middleware.DeserializeUser, handlers.VerifyOTP)
  app.Post("/disable", middleware.DeserializeUser, handlers.DisableOTP) 
  // Email
	app.Get("/verify/:token", handlers.VerifyEmail)
  // Google auth
	app.Get("/auth/google", handlers.AuthGoogle)
	app.Get("/auth/google/callback", handlers.CallbackGoogle) 
  // GitHub auth
  app.Get("/auth/github", handlers.AuthGitHub)
  app.Get("/auth/github/callback", handlers.CallbackGitHub)

  app.Get("/users", handlers.GetUsers) 
}
