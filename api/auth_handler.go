package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	myErrors "github.com/rtsoy/hotel-reservation/api/errors"
	"github.com/rtsoy/hotel-reservation/db"
	"github.com/rtsoy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"time"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return myErrors.ErrBadRequest()
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return myErrors.ErrWrongCredentials()
		}

		return err
	}

	if !types.IsPasswordValid(user.EncryptedPassword, params.Password) {
		return myErrors.ErrWrongCredentials()
	}

	response := AuthResponse{
		User:  user,
		Token: createTokenFromUser(user),
	}

	return c.JSON(response)
}

func createTokenFromUser(user *types.User) string {
	claims := jwt.MapClaims{
		"userID":  user.ID,
		"email":   user.Email,
		"expires": time.Now().Add(time.Hour * 4),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")

	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Println(err)
	}

	return tokenStr
}
