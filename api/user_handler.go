package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	errors2 "github.com/rtsoy/hotel-reservation/api/errors"
	"github.com/rtsoy/hotel-reservation/db"
	"github.com/rtsoy/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		values types.UpdateUserParams
		id     = c.Params("id")
	)

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors2.ErrInvalidID()
	}

	if err := c.BodyParser(&values); err != nil {
		return errors2.ErrBadRequest()
	}

	filter := bson.M{"_id": oid}
	update := bson.M{
		"$set": values.ToBSON(),
	}

	if err := h.userStore.UpdateUser(c.Context(), filter, update); err != nil {
		// TODO: 404
		return err
	}

	return c.JSON(map[string]string{
		"updated": id,
	})
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	var id = c.Params("id")

	if err := h.userStore.DeleteUser(c.Context(), id); err != nil {
		// TODO: 404
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("not found")
		}

		return err
	}

	return c.JSON(map[string]string{
		"deleted": id,
	})
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return errors2.ErrBadRequest()
	}

	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}

	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}

	return c.JSON(insertedUser)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var id = c.Params("id")

	user, err := h.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		// TODO: 404
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("not found")
		}

		return err
	}

	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	// TODO : 404
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors2.ErrResourceNotFound()
		}

		return err
	}

	return c.JSON(users)
}
