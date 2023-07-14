package errors

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	if apiError, ok := err.(Error); ok {
		return c.Status(apiError.Code).JSON(apiError)
	}
	someError := NewError(http.StatusInternalServerError, err.Error())
	return c.Status(someError.Code).JSON(someError)
}

func (e Error) Error() string {
	return e.Message
}

func NewError(code int, msg string) Error {
	return Error{
		Code:    code,
		Message: msg,
	}
}

func ErrResourceNotFound() Error {
	return Error{
		Code:    http.StatusNotFound, // 404
		Message: "Resource Not Found",
	}
}

func ErrBadRequest() Error {
	return Error{
		Code:    http.StatusBadRequest, // 400
		Message: "Invalid JSON",
	}
}

func ErrForbidden() Error {
	return Error{
		Code:    http.StatusForbidden, // 403
		Message: "Access Forbidden",
	}
}

func ErrTokenExpired() Error {
	return Error{
		Code:    http.StatusUnauthorized, // 401
		Message: "Token is expired",
	}
}

func ErrInvalidToken() Error {
	return Error{
		Code:    http.StatusUnauthorized, // 401
		Message: "Invalid Token",
	}
}

func ErrNoToken() Error {
	return Error{
		Code:    http.StatusUnauthorized, // 401
		Message: "No token provided",
	}
}

func ErrUnauthorized() Error {
	return Error{
		Code:    http.StatusUnauthorized, // 401
		Message: "Authentication required",
	}
}

func ErrInvalidID() Error {
	return Error{
		Code:    http.StatusBadRequest, // 400
		Message: "Invalid ID",
	}
}
