package middleware

import (
	"errors"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func getSecret() ([]byte, error) {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		return nil, errors.New("JWT_SECRET not set")
	}
	return []byte(s), nil
}

func AuthMiddleware(c fiber.Ctx) error {
	unauthorized := func(msg string) error {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": msg})
	}

	tokenString := c.Cookies("jwt")
	if tokenString == "" {
		return unauthorized("Unauthorized")
	}

	secret, err := getSecret()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server misconfigured"})
	}

	token, err := jwt.Parse(
		tokenString,
		func(t *jwt.Token) (interface{}, error) { return secret, nil },
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return unauthorized("Token Expired")
		}
		return unauthorized("Unauthorized")
	}
	if !token.Valid {
		return unauthorized("Unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return unauthorized("Unauthorized")
	}

	userID, ok := claims["userId"]
	if !ok {
		return unauthorized("Unauthorized")
	}

	c.Locals("userId", userID)
	return c.Next()
}