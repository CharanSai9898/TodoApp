package routes

import (
	"todo/src/controllers"

	"github.com/gofiber/fiber/v3"
)

func AuthRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/register", controllers.RegisterUser)
	auth.Post("/login", controllers.LoginUser)
	auth.Post("/verify-otp", controllers.VerifyOTP)
	auth.Post("/resend-otp", controllers.ResendOTP)
	auth.Post("/logout", controllers.LogoutUser)
	auth.Get("/me", controllers.GetCurrentUser)
	auth.Get("/google", controllers.GoogleLogin)
	auth.Get("/google/callback", controllers.GoogleCallback)
}
