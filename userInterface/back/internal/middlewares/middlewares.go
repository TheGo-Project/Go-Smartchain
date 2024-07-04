package middlewares

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rosty-git/Smartchain-backend/internal/models"
)

// AuthMiddleware Middleware JWT function
func AuthMiddleware(secret string) fiber.Handler {
	slog.Info("AuthMiddleware")

	//jwtTtl := time.Duration(604800) * 1_000_000_000
	field := "Authorization"

	return func(c *fiber.Ctx) error {
		// Get the cookie off the request
		tokenHeader, ok := c.GetReqHeaders()[field]
		if !ok {
			return fiber.ErrUnauthorized
		}

		var tokenString string
		if len(tokenHeader) > 0 {
			tokenString = tokenHeader[0]
		} else {
			return fiber.ErrUnauthorized
		}

		// Decode/validate it
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(secret), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check the expiry date
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				return fiber.ErrUnauthorized
			}

			userJson, ok := claims["user"]
			if !ok {
				slog.Error("AuthMiddleware")

				return fiber.ErrUnauthorized
			}

			var user models.User
			err := json.Unmarshal([]byte(userJson.(string)), &user)
			if err != nil {
				slog.Error("AuthMiddleware", "error", err)

				return fiber.ErrUnauthorized
			}

			slog.Info("AuthMiddleware", "user", user)

			c.Locals("User", user)

			return c.Next()
		} else {
			c.SendStatus(fiber.StatusUnauthorized)

			return c.Send(nil)
		}
	}
}
