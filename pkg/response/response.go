package response

import (
	"github.com/gofiber/fiber/v2"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success sends a successful response
func Success(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a created response
func Created(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends an error response
func Error(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Error:   message,
	})
}

// BadRequest sends a bad request response
func BadRequest(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusBadRequest, message)
}

// NotFound sends a not found response
func NotFound(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, message)
}

// InternalServerError sends an internal server error response
func InternalServerError(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusInternalServerError, message)
}

// Unauthorized sends an unauthorized response
func Unauthorized(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnauthorized, message)
}

// Forbidden sends a forbidden response
func Forbidden(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusForbidden, message)
}

// TooManyRequests sends a too many requests response
func TooManyRequests(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusTooManyRequests).JSON(Response{
		Success: false,
		Error:   message,
		Data:    data,
	})
}
