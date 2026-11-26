package src

import (
	"log"
	"os"
	"todo/src/db"
	"todo/src/routes"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"
)

func SetupApp() *fiber.App {

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowCredentials: true,
	}))

	err := godotenv.Load()
	if err != nil {
		_ = godotenv.Load("Backend/.env")
	}
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET is missing; set it in the environment or .env before starting the server")
	}

	db.ConnectDB()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Welcome")
	})
	routes.AuthRoutes(app)

	return app
}
