package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/services"
	"github.com/turahe/master-data-rest-api/pkg/response"
)

// ProvinceHTTPHandler handles HTTP requests for province operations
type ProvinceHTTPHandler struct {
	provinceService *services.ProvinceService
}

// NewProvinceHTTPHandler creates a new ProvinceHTTPHandler instance
func NewProvinceHTTPHandler(provinceService *services.ProvinceService) *ProvinceHTTPHandler {
	return &ProvinceHTTPHandler{
		provinceService: provinceService,
	}
}

// CreateProvince handles POST /api/v1/provinces
// @Summary Create a new province
// @Description Create a new province with the provided information
// @Tags provinces
// @Accept json
// @Produce json
// @Param request body object true "Province information"
// @Success 201 {object} response.Response{data=entities.Province} "Province created successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/provinces [post]
func (h *ProvinceHTTPHandler) CreateProvince(c *gin.Context) {
	var req struct {
		CountryID uuid.UUID `json:"country_id" binding:"required"`
		Name      string    `json:"name" binding:"required"`
		Region    *string   `json:"region"`
		ISO31662  *string   `json:"iso_3166_2"`
		Code      *string   `json:"code"`
		Latitude  *string   `json:"latitude"`
		Longitude *string   `json:"longitude"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	province, err := h.provinceService.CreateProvince(req.CountryID, req.Name)
	if err != nil {
		response.InternalServerError(c, "Failed to create province: "+err.Error())
		return
	}

	// Set optional fields
	if req.Region != nil {
		province.SetRegion(*req.Region)
	}
	if req.ISO31662 != nil {
		province.SetISO31662(*req.ISO31662)
	}
	if req.Code != nil {
		province.SetCode(*req.Code)
	}
	if req.Latitude != nil && req.Longitude != nil {
		province.SetCoordinates(*req.Latitude, *req.Longitude)
	}

	// Update the province with the additional information
	if err := h.provinceService.UpdateProvince(province); err != nil {
		response.InternalServerError(c, "Failed to update province: "+err.Error())
		return
	}

	response.Created(c, province, "Province created successfully")
}

// GetProvinceByID handles GET /api/v1/provinces/:id
// @Summary Get province by ID
// @Description Get a province by its UUID
// @Tags provinces
// @Accept json
// @Produce json
// @Param id path string true "Province ID" format(uuid)
// @Success 200 {object} response.Response{data=entities.Province} "Province found"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 404 {object} response.Response "Province not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/provinces/{id} [get]
func (h *ProvinceHTTPHandler) GetProvinceByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid province ID: "+err.Error())
		return
	}

	province, err := h.provinceService.GetProvinceByID(id)
	if err != nil {
		response.NotFound(c, "Province not found: "+err.Error())
		return
	}

	response.Success(c, province, "Province retrieved successfully")
}

// GetProvinceByCode handles GET /api/v1/provinces/code/:code
// @Summary Get province by code
// @Description Get a province by its code
// @Tags provinces
// @Accept json
// @Produce json
// @Param code path string true "Province code"
// @Success 200 {object} response.Response{data=entities.Province} "Province found"
// @Failure 404 {object} response.Response "Province not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/provinces/code/{code} [get]
func (h *ProvinceHTTPHandler) GetProvinceByCode(c *gin.Context) {
	code := c.Param("code")
	province, err := h.provinceService.GetProvinceByCode(code)
	if err != nil {
		response.NotFound(c, "Province not found: "+err.Error())
		return
	}

	response.Success(c, province, "Province retrieved successfully")
}

// GetProvinceByName handles GET /api/v1/provinces/name/:name
// @Summary Get province by name
// @Description Get a province by name
// @Tags provinces
// @Accept json
// @Produce json
// @Param name path string true "Province name"
// @Success 200 {object} response.Response{data=entities.Province} "Province found"
// @Failure 404 {object} response.Response "Province not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/provinces/name/{name} [get]
func (h *ProvinceHTTPHandler) GetProvinceByName(c *gin.Context) {
	name := c.Param("name")
	province, err := h.provinceService.GetProvinceByName(name)
	if err != nil {
		response.NotFound(c, "Province not found: "+err.Error())
		return
	}

	response.Success(c, province, "Province retrieved successfully")
}

// GetProvincesByCountry handles GET /api/v1/provinces/country/:country_id
// @Summary Get provinces by country
// @Description Get all provinces for a specific country
// @Tags provinces
// @Accept json
// @Produce json
// @Param country_id path string true "Country ID" format(uuid)
// @Success 200 {object} response.Response{data=[]entities.Province} "Provinces found"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 404 {object} response.Response "Provinces not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/provinces/country/{country_id} [get]
func (h *ProvinceHTTPHandler) GetProvincesByCountry(c *gin.Context) {
	countryIDStr := c.Param("country_id")
	countryID, err := uuid.Parse(countryIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid country ID: "+err.Error())
		return
	}

	provinces, err := h.provinceService.GetProvincesByCountryID(countryID)
	if err != nil {
		response.NotFound(c, "Provinces not found: "+err.Error())
		return
	}

	response.Success(c, provinces, "Provinces retrieved successfully")
}

// GetAllProvinces handles GET /api/v1/provinces
// @Summary Get all provinces
// @Description Get all provinces with optional pagination
// @Tags provinces
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]entities.Province} "Provinces retrieved successfully"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/provinces [get]
func (h *ProvinceHTTPHandler) GetAllProvinces(c *gin.Context) {
	provinces, err := h.provinceService.GetAllProvinces()
	if err != nil {
		response.InternalServerError(c, "Failed to get provinces: "+err.Error())
		return
	}

	response.Success(c, provinces, "Provinces retrieved successfully")
}

// UpdateProvince handles PUT /api/v1/provinces/:id
// @Summary Update a province
// @Description Update a province with the provided information
// @Tags provinces
// @Accept json
// @Produce json
// @Param id path string true "Province ID" format(uuid)
// @Param request body object true "Province information"
// @Success 200 {object} response.Response{data=entities.Province} "Province updated successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 404 {object} response.Response "Province not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/provinces/{id} [put]
func (h *ProvinceHTTPHandler) UpdateProvince(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid province ID: "+err.Error())
		return
	}

	var req struct {
		CountryID uuid.UUID `json:"country_id" binding:"required"`
		Name      string    `json:"name" binding:"required"`
		Region    *string   `json:"region"`
		ISO31662  *string   `json:"iso_3166_2"`
		Code      *string   `json:"code"`
		Latitude  *string   `json:"latitude"`
		Longitude *string   `json:"longitude"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	// Get existing province
	province, err := h.provinceService.GetProvinceByID(id)
	if err != nil {
		response.NotFound(c, "Province not found: "+err.Error())
		return
	}

	// Update fields
	province.CountryID = req.CountryID
	province.Update(req.Name)
	if req.Region != nil {
		province.SetRegion(*req.Region)
	}
	if req.ISO31662 != nil {
		province.SetISO31662(*req.ISO31662)
	}
	if req.Code != nil {
		province.SetCode(*req.Code)
	}
	if req.Latitude != nil && req.Longitude != nil {
		province.SetCoordinates(*req.Latitude, *req.Longitude)
	}

	if err := h.provinceService.UpdateProvince(province); err != nil {
		response.InternalServerError(c, "Failed to update province: "+err.Error())
		return
	}

	response.Success(c, province, "Province updated successfully")
}

// DeleteProvince handles DELETE /api/v1/provinces/:id
// @Summary Delete a province
// @Description Delete a province by its UUID
// @Tags provinces
// @Accept json
// @Produce json
// @Param id path string true "Province ID" format(uuid)
// @Success 200 {object} response.Response "Province deleted successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 404 {object} response.Response "Province not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/provinces/{id} [delete]
func (h *ProvinceHTTPHandler) DeleteProvince(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid province ID: "+err.Error())
		return
	}

	if err := h.provinceService.DeleteProvince(id); err != nil {
		response.NotFound(c, "Province not found: "+err.Error())
		return
	}

	response.Success(c, nil, "Province deleted successfully")
}

// GetProvinceCount handles GET /api/v1/provinces/count
// @Summary Get province count
// @Description Get the total number of provinces
// @Tags provinces
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=object} "Province count retrieved successfully"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/provinces/count [get]
func (h *ProvinceHTTPHandler) GetProvinceCount(c *gin.Context) {
	count, err := h.provinceService.GetProvinceCount()
	if err != nil {
		response.InternalServerError(c, "Failed to get province count: "+err.Error())
		return
	}

	response.Success(c, gin.H{"count": count}, "Province count retrieved successfully")
} 