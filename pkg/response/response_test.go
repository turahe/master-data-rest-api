package response

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestSuccess(t *testing.T) {
	// Given
	app := fiber.New()
	data := map[string]string{"key": "value"}
	message := "Operation successful"

	app.Get("/test", func(c *fiber.Ctx) error {
		return Success(c, data, message)
	})

	// When
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, message, response.Message)
	assert.NotNil(t, response.Data)
	assert.Empty(t, response.Error)
}

func TestCreated(t *testing.T) {
	// Given
	app := fiber.New()
	data := map[string]string{"id": "123"}
	message := "Resource created"

	app.Post("/test", func(c *fiber.Ctx) error {
		return Created(c, data, message)
	})

	// When
	req := httptest.NewRequest("POST", "/test", nil)
	resp, err := app.Test(req)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, message, response.Message)
	assert.NotNil(t, response.Data)
	assert.Empty(t, response.Error)
}

func TestError(t *testing.T) {
	// Given
	app := fiber.New()
	statusCode := 500
	message := "Internal server error"

	app.Get("/test", func(c *fiber.Ctx) error {
		return Error(c, statusCode, message)
	})

	// When
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, statusCode, resp.StatusCode)

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Empty(t, response.Message)
	assert.Nil(t, response.Data)
	assert.Equal(t, message, response.Error)
}

func TestBadRequest(t *testing.T) {
	// Given
	app := fiber.New()
	message := "Invalid request"

	app.Get("/test", func(c *fiber.Ctx) error {
		return BadRequest(c, message)
	})

	// When
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, message, response.Error)
}

func TestNotFound(t *testing.T) {
	// Given
	app := fiber.New()
	message := "Resource not found"

	app.Get("/test", func(c *fiber.Ctx) error {
		return NotFound(c, message)
	})

	// When
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, message, response.Error)
}

func TestInternalServerError(t *testing.T) {
	// Given
	app := fiber.New()
	message := "Something went wrong"

	app.Get("/test", func(c *fiber.Ctx) error {
		return InternalServerError(c, message)
	})

	// When
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, message, response.Error)
}

func TestUnauthorized(t *testing.T) {
	// Given
	app := fiber.New()
	message := "Access denied"

	app.Get("/test", func(c *fiber.Ctx) error {
		return Unauthorized(c, message)
	})

	// When
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, message, response.Error)
}

func TestForbidden(t *testing.T) {
	// Given
	app := fiber.New()
	message := "Forbidden access"

	app.Get("/test", func(c *fiber.Ctx) error {
		return Forbidden(c, message)
	})

	// When
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 403, resp.StatusCode)

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, message, response.Error)
}
