package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
	"github.com/turahe/master-data-rest-api/internal/domain/repositories"
	"github.com/turahe/master-data-rest-api/internal/domain/services"
	"github.com/turahe/master-data-rest-api/pkg/response"
)

// GeodirectoryHTTPHandler handles HTTP requests for geodirectory operations
type GeodirectoryHTTPHandler struct {
	geodirectoryService *services.GeodirectoryService
	searchService       repositories.SearchRepository
}

// NewGeodirectoryHTTPHandler creates a new GeodirectoryHTTPHandler instance
func NewGeodirectoryHTTPHandler(geodirectoryService *services.GeodirectoryService, searchService repositories.SearchRepository) *GeodirectoryHTTPHandler {
	return &GeodirectoryHTTPHandler{
		geodirectoryService: geodirectoryService,
		searchService:       searchService,
	}
}

// CreateGeodirectory handles POST /api/v1/geodirectories
// @Summary Create a new geodirectory
// @Description Create a new geodirectory with hierarchical support
// @Tags geodirectories
// @Accept json
// @Produce json
// @Param geodirectory body CreateGeodirectoryRequest true "Geodirectory data"
// @Success 201 {object} response.Response "Geodirectory created successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/geodirectories [post]
func (h *GeodirectoryHTTPHandler) CreateGeodirectory(c *fiber.Ctx) error {
	var req CreateGeodirectoryRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body: "+err.Error())
	}

	geodirectory, err := h.geodirectoryService.CreateGeodirectory(
		c.Context(),
		req.Name,
		entities.GeoType(req.Type),
		req.ParentID,
	)
	if err != nil {
		return response.InternalServerError(c, "Failed to create geodirectory: "+err.Error())
	}

	// Set optional fields
	if req.Code != nil {
		geodirectory.SetCode(*req.Code)
	}
	if req.PostalCode != nil {
		geodirectory.SetPostalCode(*req.PostalCode)
	}
	if req.Latitude != nil && req.Longitude != nil {
		geodirectory.SetCoordinates(*req.Latitude, *req.Longitude)
	}

	// Update with additional information
	if err := h.geodirectoryService.UpdateGeodirectory(c.Context(), geodirectory); err != nil {
		return response.InternalServerError(c, "Failed to update geodirectory: "+err.Error())
	}

	return response.Created(c, geodirectory, "Geodirectory created successfully")
}

// GetGeodirectoryByID handles GET /api/v1/geodirectories/:id
// @Summary Get geodirectory by ID
// @Description Get a geodirectory by its UUID
// @Tags geodirectories
// @Produce json
// @Param id path string true "Geodirectory ID (UUID)"
// @Success 200 {object} response.Response "Geodirectory retrieved successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Geodirectory not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/geodirectories/{id} [get]
func (h *GeodirectoryHTTPHandler) GetGeodirectoryByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid geodirectory ID: "+err.Error())
	}

	geodirectory, err := h.geodirectoryService.GetGeodirectoryByID(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Geodirectory not found: "+err.Error())
	}

	return response.Success(c, geodirectory, "Geodirectory retrieved successfully")
}

// GetGeodirectoryWithHierarchy handles GET /api/v1/geodirectories/:id/hierarchy
// @Summary Get geodirectory with hierarchy
// @Description Get a geodirectory with its parent and children
// @Tags geodirectories
// @Produce json
// @Param id path string true "Geodirectory ID (UUID)"
// @Success 200 {object} response.Response "Geodirectory with hierarchy retrieved successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Geodirectory not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/geodirectories/{id}/hierarchy [get]
func (h *GeodirectoryHTTPHandler) GetGeodirectoryWithHierarchy(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid geodirectory ID: "+err.Error())
	}

	geodirectory, err := h.geodirectoryService.GetGeodirectoryWithHierarchy(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Geodirectory not found: "+err.Error())
	}

	return response.Success(c, geodirectory, "Geodirectory with hierarchy retrieved successfully")
}

// GetAllGeodirectories handles GET /api/v1/geodirectories
// @Summary Get all geodirectories
// @Description Get all geodirectories with pagination
// @Tags geodirectories
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response "Geodirectories retrieved successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/geodirectories [get]
func (h *GeodirectoryHTTPHandler) GetAllGeodirectories(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	geodirectories, err := h.geodirectoryService.GetAllGeodirectories(c.Context(), limit, offset)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve geodirectories: "+err.Error())
	}

	return response.Success(c, geodirectories, "Geodirectories retrieved successfully")
}

// SearchGeodirectories handles GET /api/v1/geodirectories/search
// @Summary Search geodirectories
// @Description Search geodirectories by name, code, or postal code
// @Tags geodirectories
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response "Geodirectories found"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/geodirectories/search [get]
func (h *GeodirectoryHTTPHandler) SearchGeodirectories(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return response.BadRequest(c, "Search query is required")
	}

	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	geodirectories, err := h.geodirectoryService.SearchGeodirectories(c.Context(), query, limit, offset)
	if err != nil {
		return response.InternalServerError(c, "Failed to search geodirectories: "+err.Error())
	}

	return response.Success(c, geodirectories, "Geodirectories found")
}

// GetGeodirectoriesByType handles GET /api/v1/geodirectories/type/:type
// @Summary Get geodirectories by type
// @Description Get geodirectories filtered by type
// @Tags geodirectories
// @Produce json
// @Param type path string true "Geodirectory Type" Enums(CONTINENT,SUBCONTINENT,COUNTRY,STATE,PROVINCE,REGENCY,CITY,DISTRICT,SUBDISTRICT,VILLAGE)
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response "Geodirectories retrieved successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/geodirectories/type/{type} [get]
func (h *GeodirectoryHTTPHandler) GetGeodirectoriesByType(c *fiber.Ctx) error {
	geoType := entities.GeoType(c.Params("type"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	geodirectories, err := h.geodirectoryService.GetGeodirectoriesByType(c.Context(), geoType, limit, offset)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve geodirectories: "+err.Error())
	}

	return response.Success(c, geodirectories, "Geodirectories retrieved successfully")
}

// GetChildren handles GET /api/v1/geodirectories/:id/children
// @Summary Get children of a geodirectory
// @Description Get direct children of a geodirectory
// @Tags geodirectories
// @Produce json
// @Param id path string true "Parent Geodirectory ID (UUID)"
// @Param type query string false "Filter by child type"
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response "Children retrieved successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Geodirectory not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/geodirectories/{id}/children [get]
func (h *GeodirectoryHTTPHandler) GetChildren(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid geodirectory ID: "+err.Error())
	}

	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	childType := c.Query("type")

	var children []*entities.Geodirectory

	if childType != "" {
		children, err = h.geodirectoryService.GetChildrenByType(c.Context(), id, entities.GeoType(childType), limit, offset)
	} else {
		children, err = h.geodirectoryService.GetChildren(c.Context(), id, limit, offset)
	}

	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve children: "+err.Error())
	}

	return response.Success(c, children, "Children retrieved successfully")
}

// GetAncestors handles GET /api/v1/geodirectories/:id/ancestors
// @Summary Get ancestors of a geodirectory
// @Description Get all ancestors of a geodirectory using nested set model
// @Tags geodirectories
// @Produce json
// @Param id path string true "Geodirectory ID (UUID)"
// @Success 200 {object} response.Response "Ancestors retrieved successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Geodirectory not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/geodirectories/{id}/ancestors [get]
func (h *GeodirectoryHTTPHandler) GetAncestors(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid geodirectory ID: "+err.Error())
	}

	ancestors, err := h.geodirectoryService.GetAncestors(c.Context(), id)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve ancestors: "+err.Error())
	}

	return response.Success(c, ancestors, "Ancestors retrieved successfully")
}

// GetDescendants handles GET /api/v1/geodirectories/:id/descendants
// @Summary Get descendants of a geodirectory
// @Description Get all descendants of a geodirectory using nested set model
// @Tags geodirectories
// @Produce json
// @Param id path string true "Geodirectory ID (UUID)"
// @Param limit query int false "Limit" default(100)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response "Descendants retrieved successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Geodirectory not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/geodirectories/{id}/descendants [get]
func (h *GeodirectoryHTTPHandler) GetDescendants(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid geodirectory ID: "+err.Error())
	}

	limit, _ := strconv.Atoi(c.Query("limit", "100"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	descendants, err := h.geodirectoryService.GetDescendants(c.Context(), id, limit, offset)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve descendants: "+err.Error())
	}

	return response.Success(c, descendants, "Descendants retrieved successfully")
}

// UpdateGeodirectory handles PUT /api/v1/geodirectories/:id
// @Summary Update a geodirectory
// @Description Update an existing geodirectory
// @Tags geodirectories
// @Accept json
// @Produce json
// @Param id path string true "Geodirectory ID (UUID)"
// @Param geodirectory body UpdateGeodirectoryRequest true "Updated geodirectory data"
// @Success 200 {object} response.Response "Geodirectory updated successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Geodirectory not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/geodirectories/{id} [put]
func (h *GeodirectoryHTTPHandler) UpdateGeodirectory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid geodirectory ID: "+err.Error())
	}

	var req UpdateGeodirectoryRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body: "+err.Error())
	}

	geodirectory, err := h.geodirectoryService.GetGeodirectoryByID(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Geodirectory not found: "+err.Error())
	}

	// Update fields
	if req.Name != nil {
		geodirectory.Name = *req.Name
	}
	if req.Code != nil {
		geodirectory.SetCode(*req.Code)
	}
	if req.PostalCode != nil {
		geodirectory.SetPostalCode(*req.PostalCode)
	}
	if req.Latitude != nil && req.Longitude != nil {
		geodirectory.SetCoordinates(*req.Latitude, *req.Longitude)
	}

	if err := h.geodirectoryService.UpdateGeodirectory(c.Context(), geodirectory); err != nil {
		return response.InternalServerError(c, "Failed to update geodirectory: "+err.Error())
	}

	return response.Success(c, geodirectory, "Geodirectory updated successfully")
}

// DeleteGeodirectory handles DELETE /api/v1/geodirectories/:id
// @Summary Delete a geodirectory
// @Description Delete a geodirectory and all its descendants
// @Tags geodirectories
// @Produce json
// @Param id path string true "Geodirectory ID (UUID)"
// @Success 200 {object} response.Response "Geodirectory deleted successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Geodirectory not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/geodirectories/{id} [delete]
func (h *GeodirectoryHTTPHandler) DeleteGeodirectory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid geodirectory ID: "+err.Error())
	}

	if err := h.geodirectoryService.DeleteGeodirectory(c.Context(), id); err != nil {
		return response.InternalServerError(c, "Failed to delete geodirectory: "+err.Error())
	}

	return response.Success(c, nil, "Geodirectory deleted successfully")
}

// MoveGeodirectory handles POST /api/v1/geodirectories/:id/move
// @Summary Move a geodirectory to a new parent
// @Description Move a geodirectory and its descendants to a new parent
// @Tags geodirectories
// @Accept json
// @Produce json
// @Param id path string true "Geodirectory ID (UUID)"
// @Param request body MoveGeodirectoryRequest true "Move request data"
// @Success 200 {object} response.Response "Geodirectory moved successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Geodirectory not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/geodirectories/{id}/move [post]
func (h *GeodirectoryHTTPHandler) MoveGeodirectory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid geodirectory ID: "+err.Error())
	}

	var req MoveGeodirectoryRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body: "+err.Error())
	}

	if err := h.geodirectoryService.MoveGeodirectory(c.Context(), id, req.NewParentID); err != nil {
		return response.InternalServerError(c, "Failed to move geodirectory: "+err.Error())
	}

	return response.Success(c, nil, "Geodirectory moved successfully")
}

// RebuildNestedSet handles POST /api/v1/geodirectories/rebuild
// @Summary Rebuild nested set structure
// @Description Rebuild the nested set structure for all geodirectories
// @Tags geodirectories
// @Produce json
// @Success 200 {object} response.Response "Nested set rebuilt successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/geodirectories/rebuild [post]
func (h *GeodirectoryHTTPHandler) RebuildNestedSet(c *fiber.Ctx) error {
	if err := h.geodirectoryService.RebuildNestedSet(c.Context()); err != nil {
		return response.InternalServerError(c, "Failed to rebuild nested set: "+err.Error())
	}

	return response.Success(c, nil, "Nested set rebuilt successfully")
}

// Request/Response DTOs

type CreateGeodirectoryRequest struct {
	Name       string     `json:"name" validate:"required"`
	Type       string     `json:"type" validate:"required"`
	Code       *string    `json:"code,omitempty"`
	PostalCode *string    `json:"postal_code,omitempty"`
	Latitude   *string    `json:"latitude,omitempty"`
	Longitude  *string    `json:"longitude,omitempty"`
	ParentID   *uuid.UUID `json:"parent_id,omitempty"`
}

type UpdateGeodirectoryRequest struct {
	Name       *string `json:"name,omitempty"`
	Code       *string `json:"code,omitempty"`
	PostalCode *string `json:"postal_code,omitempty"`
	Latitude   *string `json:"latitude,omitempty"`
	Longitude  *string `json:"longitude,omitempty"`
}

type MoveGeodirectoryRequest struct {
	NewParentID uuid.UUID `json:"new_parent_id" validate:"required"`
}
