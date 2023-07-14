package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rtsoy/hotel-reservation/db"
	"log"
	"net/http"
	"os"
	"time"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			return c.Status(http.StatusUnauthorized).JSON("no token")
		}

		claims, err := validateToken(token)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(err.Error())
		}

		expires := claims["expires"]
		parsedTime, err := time.Parse(time.RFC3339, expires.(string))
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON("token expired")
		}
		if time.Now().After(parsedTime) {
			return c.Status(http.StatusUnauthorized).JSON("token expired")
		}

		userID := claims["userID"].(string)
		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON("unauthorized")
		}

		// Set the current authenticated user to the context value
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("invalid signing method:", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}

		secret := os.Getenv("JWT_SECRET")

		return []byte(secret), nil
	})

	if err != nil {
		log.Println("failed to parse jwt token:", err)
		return nil, fmt.Errorf("unauthorized")
	}

	if !token.Valid {
		log.Println("invalid token:", err)
		return nil, fmt.Errorf("unauthorized")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("unauthorized")
}
