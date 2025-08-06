package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/services"
	"github.com/turahe/master-data-rest-api/pkg/response"
)

// CityHTTPHandler handles HTTP requests for city operations
type CityHTTPHandler struct {
	cityService *services.CityService
}

// NewCityHTTPHandler creates a new CityHTTPHandler instance
func NewCityHTTPHandler(cityService *services.CityService) *CityHTTPHandler {
	return &CityHTTPHandler{
		cityService: cityService,
	}
}

// CreateCity handles POST /api/v1/cities
// @Summary Create a new city
// @Description Create a new city with the provided information
// @Tags cities
// @Accept json
// @Produce json
// @Param request body object true "City information"
// @Success 201 {object} response.Response{data=entities.City} "City created successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/cities [post]
func (h *CityHTTPHandler) CreateCity(c *gin.Context) {
	var req struct {
		ProvinceID uuid.UUID `json:"province_id" binding:"required"`
		Name       string    `json:"name" binding:"required"`
		Type       *string   `json:"type"`
		Code       *string   `json:"code"`
		Latitude   *string   `json:"latitude"`
		Longitude  *string   `json:"longitude"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	city, err := h.cityService.CreateCity(req.ProvinceID, req.Name)
	if err != nil {
		response.InternalServerError(c, "Failed to create city: "+err.Error())
		return
	}

	// Set optional fields
	if req.Type != nil {
		city.SetType(*req.Type)
	}
	if req.Code != nil {
		city.SetCode(*req.Code)
	}
	if req.Latitude != nil && req.Longitude != nil {
		city.SetCoordinates(*req.Latitude, *req.Longitude)
	}

	// Update the city with the additional information
	if err := h.cityService.UpdateCity(city); err != nil {
		response.InternalServerError(c, "Failed to update city: "+err.Error())
		return
	}

	response.Created(c, city, "City created successfully")
}

// GetCityByID handles GET /api/v1/cities/:id
// @Summary Get city by ID
// @Description Get a city by its UUID
// @Tags cities
// @Accept json
// @Produce json
// @Param id path string true "City ID" format(uuid)
// @Success 200 {object} response.Response{data=entities.City} "City found"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 404 {object} response.Response "City not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/cities/{id} [get]
func (h *CityHTTPHandler) GetCityByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid city ID: "+err.Error())
		return
	}

	city, err := h.cityService.GetCityByID(id)
	if err != nil {
		response.NotFound(c, "City not found: "+err.Error())
		return
	}

	response.Success(c, city, "City retrieved successfully")
}

// GetCityByCode handles GET /api/v1/cities/code/:code
// @Summary Get city by code
// @Description Get a city by its code
// @Tags cities
// @Accept json
// @Produce json
// @Param code path string true "City code"
// @Success 200 {object} response.Response{data=entities.City} "City found"
// @Failure 404 {object} response.Response "City not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/cities/code/{code} [get]
func (h *CityHTTPHandler) GetCityByCode(c *gin.Context) {
	code := c.Param("code")
	city, err := h.cityService.GetCityByCode(code)
	if err != nil {
		response.NotFound(c, "City not found: "+err.Error())
		return
	}

	response.Success(c, city, "City retrieved successfully")
}

// GetCityByName handles GET /api/v1/cities/name/:name
// @Summary Get city by name
// @Description Get a city by name
// @Tags cities
// @Accept json
// @Produce json
// @Param name path string true "City name"
// @Success 200 {object} response.Response{data=entities.City} "City found"
// @Failure 404 {object} response.Response "City not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/cities/name/{name} [get]
func (h *CityHTTPHandler) GetCityByName(c *gin.Context) {
	name := c.Param("name")
	city, err := h.cityService.GetCityByName(name)
	if err != nil {
		response.NotFound(c, "City not found: "+err.Error())
		return
	}

	response.Success(c, city, "City retrieved successfully")
}

// GetCitiesByProvince handles GET /api/v1/cities/province/:province_id
// @Summary Get cities by province
// @Description Get all cities for a specific province
// @Tags cities
// @Accept json
// @Produce json
// @Param province_id path string true "Province ID" format(uuid)
// @Success 200 {object} response.Response{data=[]entities.City} "Cities found"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 404 {object} response.Response "Cities not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/cities/province/{province_id} [get]
func (h *CityHTTPHandler) GetCitiesByProvince(c *gin.Context) {
	provinceIDStr := c.Param("province_id")
	provinceID, err := uuid.Parse(provinceIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid province ID: "+err.Error())
		return
	}

	cities, err := h.cityService.GetCitiesByProvinceID(provinceID)
	if err != nil {
		response.NotFound(c, "Cities not found: "+err.Error())
		return
	}

	response.Success(c, cities, "Cities retrieved successfully")
}

// GetAllCities handles GET /api/v1/cities
// @Summary Get all cities
// @Description Get all cities with optional pagination
// @Tags cities
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]entities.City} "Cities retrieved successfully"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/cities [get]
func (h *CityHTTPHandler) GetAllCities(c *gin.Context) {
	cities, err := h.cityService.GetAllCities()
	if err != nil {
		response.InternalServerError(c, "Failed to get cities: "+err.Error())
		return
	}

	response.Success(c, cities, "Cities retrieved successfully")
}

// UpdateCity handles PUT /api/v1/cities/:id
// @Summary Update a city
// @Description Update a city with the provided information
// @Tags cities
// @Accept json
// @Produce json
// @Param id path string true "City ID" format(uuid)
// @Param request body object true "City information"
// @Success 200 {object} response.Response{data=entities.City} "City updated successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 404 {object} response.Response "City not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/cities/{id} [put]
func (h *CityHTTPHandler) UpdateCity(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid city ID: "+err.Error())
		return
	}

	var req struct {
		ProvinceID uuid.UUID `json:"province_id" binding:"required"`
		Name       string    `json:"name" binding:"required"`
		Type       *string   `json:"type"`
		Code       *string   `json:"code"`
		Latitude   *string   `json:"latitude"`
		Longitude  *string   `json:"longitude"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	// Get existing city
	city, err := h.cityService.GetCityByID(id)
	if err != nil {
		response.NotFound(c, "City not found: "+err.Error())
		return
	}

	// Update fields
	city.ProvinceID = req.ProvinceID
	city.Update(req.Name)
	if req.Type != nil {
		city.SetType(*req.Type)
	}
	if req.Code != nil {
		city.SetCode(*req.Code)
	}
	if req.Latitude != nil && req.Longitude != nil {
		city.SetCoordinates(*req.Latitude, *req.Longitude)
	}

	if err := h.cityService.UpdateCity(city); err != nil {
		response.InternalServerError(c, "Failed to update city: "+err.Error())
		return
	}

	response.Success(c, city, "City updated successfully")
}

// DeleteCity handles DELETE /api/v1/cities/:id
// @Summary Delete a city
// @Description Delete a city by its UUID
// @Tags cities
// @Accept json
// @Produce json
// @Param id path string true "City ID" format(uuid)
// @Success 200 {object} response.Response "City deleted successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 404 {object} response.Response "City not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/cities/{id} [delete]
func (h *CityHTTPHandler) DeleteCity(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid city ID: "+err.Error())
		return
	}

	if err := h.cityService.DeleteCity(id); err != nil {
		response.NotFound(c, "City not found: "+err.Error())
		return
	}

	response.Success(c, nil, "City deleted successfully")
}

// GetCityCount handles GET /api/v1/cities/count
// @Summary Get city count
// @Description Get the total number of cities
// @Tags cities
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=object} "City count retrieved successfully"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/cities/count [get]
func (h *CityHTTPHandler) GetCityCount(c *gin.Context) {
	count, err := h.cityService.GetCityCount()
	if err != nil {
		response.InternalServerError(c, "Failed to get city count: "+err.Error())
		return
	}

	response.Success(c, gin.H{"count": count}, "City count retrieved successfully")
} 