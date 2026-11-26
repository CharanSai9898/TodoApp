package controllers

import (
	"encoding/json"
	"os"
	"strings"
	"time"
	"todo/src/db"
	"todo/src/models"
	"todo/src/services"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func getJWTSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

func RegisterUser(c fiber.Ctx) error {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body request

	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	if body.Email == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email and password are required"})
	}

	if db.Db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database not connected"})
	}

	var existingUser models.User
	err := db.Db.Collection("users").FindOne(c.Context(), bson.M{"email": body.Email}).Decode(&existingUser)

	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User Already exists"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create user"})
	}

	userDoc := bson.M{
		"email":    body.Email,
		"password": string(hashedPassword),
	}
	_, err = db.Db.Collection("users").InsertOne(c.Context(), userDoc)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create user"})

	}
	return c.JSON(fiber.Map{"message": "user registered successfully"})
}

func LoginUser(c fiber.Ctx) error {
	secret := getJWTSecret()
	if len(secret) == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server misconfigured"})
	}
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body request

	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	if db.Db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database not connected"})
	}

	var user models.User

	err := db.Db.Collection("users").FindOne(c.Context(), bson.M{"email": body.Email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Credentials"})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Credentials"})
	}

	// Generate 6-digit OTP
	otp, err := services.GenerateOTP()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate verification code"})
	}

	cleanEmail := strings.ToLower(strings.TrimSpace(user.Email))

	// Store OTP in database (valid for 5 minutes)
	_, _ = db.Db.Collection("otps").DeleteMany(c.Context(), bson.M{"email": cleanEmail})
	_, err = db.Db.Collection("otps").InsertOne(c.Context(), models.OTP{
		ID:        primitive.NewObjectID(),
		Email:     cleanEmail,
		OTP:       otp,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save verification code"})
	}

	// Dispatch OTP Email in background so the HTTP request does not hang if SMTP is slow
	go func(targetEmail, targetOTP string) {
		_ = services.SendOTPEmail(targetEmail, targetOTP)
	}(cleanEmail, otp)

	// Issue short-lived temporary token for 2FA verification step
	tempTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID.Hex(),
		"email":  cleanEmail,
		"type":   "2fa_temp",
		"exp":    time.Now().Add(5 * time.Minute).Unix(),
	})

	signedTempToken, err := tempTokenObj.SignedString(secret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to initiate 2FA"})
	}

	return c.JSON(fiber.Map{
		"requires2FA": true,
		"tempToken":   signedTempToken,
		"email":       cleanEmail,
		"message":     "Verification code sent to your email",
	})
}

func VerifyOTP(c fiber.Ctx) error {
	type request struct {
		TempToken string `json:"tempToken"`
		OTP       string `json:"otp"`
	}
	var body request
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	trimmedOTP := strings.TrimSpace(body.OTP)
	if body.TempToken == "" || trimmedOTP == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Token and code are required"})
	}

	secret := getJWTSecret()
	token, err := jwt.Parse(body.TempToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return secret, nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Verification session expired, please log in again"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "2fa_temp" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid verification session"})
	}

	userId, _ := claims["userId"].(string)
	email, _ := claims["email"].(string)
	cleanEmail := strings.ToLower(strings.TrimSpace(email))

	if db.Db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database not connected"})
	}

	var otpRecord models.OTP
	err = db.Db.Collection("otps").FindOne(c.Context(), bson.M{
		"email": cleanEmail,
		"otp":   trimmedOTP,
	}).Decode(&otpRecord)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid verification code. Please check and try again"})
	}

	if time.Now().After(otpRecord.ExpiresAt) {
		_, _ = db.Db.Collection("otps").DeleteMany(c.Context(), bson.M{"email": cleanEmail})
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Verification code has expired. Please click Resend code"})
	}

	// Delete used OTP so it cannot be re-used
	_, _ = db.Db.Collection("otps").DeleteMany(c.Context(), bson.M{"email": cleanEmail})

	// Issue full 24-hour session JWT
	sessionToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	})

	t, err := sessionToken.SignedString(secret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not complete login"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    t,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		SameSite: "Lax",
	})

	return c.JSON(fiber.Map{"message": "Login successful"})
}

func ResendOTP(c fiber.Ctx) error {
	type request struct {
		TempToken string `json:"tempToken"`
	}
	var body request
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	secret := getJWTSecret()
	token, err := jwt.Parse(body.TempToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return secret, nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Verification session expired, please log in again"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "2fa_temp" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid verification session"})
	}

	email, _ := claims["email"].(string)
	cleanEmail := strings.ToLower(strings.TrimSpace(email))

	otp, err := services.GenerateOTP()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate code"})
	}

	if db.Db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database not connected"})
	}

	_, _ = db.Db.Collection("otps").DeleteMany(c.Context(), bson.M{"email": cleanEmail})
	_, err = db.Db.Collection("otps").InsertOne(c.Context(), models.OTP{
		ID:        primitive.NewObjectID(),
		Email:     cleanEmail,
		OTP:       otp,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not save code"})
	}

	go func(targetEmail, targetOTP string) {
		_ = services.SendOTPEmail(targetEmail, targetOTP)
	}(cleanEmail, otp)

	return c.JSON(fiber.Map{"message": "New verification code sent"})
}

func LogoutUser(c fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		MaxAge:   -1,
		HTTPOnly: true,
		SameSite: "Lax",
	})

	return c.JSON(fiber.Map{"message": "Logout successful"})

}
func GetCurrentUser(c fiber.Ctx) error {

	// Read JWT cookie
	tokenString := c.Cookies("jwt")

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"error": "Not authenticated",
			},
		)
	}

	// Validate JWT
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}

			return getJWTSecret(), nil
		},
	)

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"error": "Invalid or expired token",
			},
		)
	}

	// Extract user ID
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	userID, ok := claims["userId"].(string)

	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	objectID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	if db.Db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database not connected"})
	}

	// Find user in MongoDB
	var user models.User

	err = db.Db.Collection("users").
		FindOne(
			c.Context(),
			bson.M{"_id": objectID},
		).
		Decode(&user)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"error": "User not found",
			},
		)
	}

	// Return safe user details
	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"id":    user.ID.Hex(),
			"email": user.Email,
		},
	})
}

func getGoogleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func GoogleLogin(c fiber.Ctx) error {
	config := getGoogleOAuthConfig()
	if config.ClientID == "" || config.ClientSecret == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Google OAuth is not configured. Please add GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET to .env",
		})
	}
	url := config.AuthCodeURL("random-state", oauth2.AccessTypeOffline)
	return c.Redirect().To(url)
}

func GoogleCallback(c fiber.Ctx) error {
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	code := c.Query("code")
	if code == "" {
		return c.Redirect().To(frontendURL + "/Login?error=missing_code")
	}

	config := getGoogleOAuthConfig()
	token, err := config.Exchange(c.Context(), code)
	if err != nil {
		return c.Redirect().To(frontendURL + "/Login?error=exchange_failed")
	}

	client := config.Client(c.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return c.Redirect().To(frontendURL + "/Login?error=fetch_user_failed")
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return c.Redirect().To(frontendURL + "/Login?error=decode_failed")
	}

	if googleUser.Email == "" {
		return c.Redirect().To(frontendURL + "/Login?error=missing_email")
	}

	if db.Db == nil {
		return c.Redirect().To(frontendURL + "/Login?error=database_unavailable")
	}

	var user models.User
	err = db.Db.Collection("users").FindOne(c.Context(), bson.M{"email": googleUser.Email}).Decode(&user)
	if err != nil {
		// New user from Google
		user = models.User{
			ID:       primitive.NewObjectID(),
			Email:    googleUser.Email,
			GoogleID: googleUser.ID,
			Todos:    []primitive.ObjectID{},
		}
		_, err = db.Db.Collection("users").InsertOne(c.Context(), user)
		if err != nil {
			return c.Redirect().To(frontendURL + "/Login?error=user_creation_failed")
		}
	} else if user.GoogleID == "" {
		// Existing email, link Google ID
		_, _ = db.Db.Collection("users").UpdateOne(
			c.Context(),
			bson.M{"_id": user.ID},
			bson.M{"$set": bson.M{"googleId": googleUser.ID}},
		)
	}

	secret := getJWTSecret()
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID.Hex(),
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	})
	t, err := jwtToken.SignedString(secret)
	if err != nil {
		return c.Redirect().To(frontendURL + "/Login?error=token_failed")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    t,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		SameSite: "Lax",
	})

	return c.Redirect().To(frontendURL + "/Dashboard")
}
