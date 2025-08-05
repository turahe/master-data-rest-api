package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success sends a successful response
func Success(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a created response
func Created(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error:   message,
	})
}

// BadRequest sends a bad request response
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// NotFound sends a not found response
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

// InternalServerError sends an internal server error response
func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}

// Unauthorized sends an unauthorized response
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

// Forbidden sends a forbidden response
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message)
} 