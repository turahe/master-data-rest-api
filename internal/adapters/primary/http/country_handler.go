package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/services"
	"github.com/turahe/master-data-rest-api/pkg/response"
)

// CountryHTTPHandler handles HTTP requests for country operations
type CountryHTTPHandler struct {
	countryService *services.CountryService
}

// NewCountryHTTPHandler creates a new CountryHTTPHandler instance
func NewCountryHTTPHandler(countryService *services.CountryService) *CountryHTTPHandler {
	return &CountryHTTPHandler{
		countryService: countryService,
	}
}

// CreateCountry handles POST /api/v1/countries
// @Summary Create a new country
// @Description Create a new country with the provided information
// @Tags countries
// @Accept json
// @Produce json
// @Param request body object true "Country information"
// @Success 201 {object} response.Response{data=entities.Country} "Country created successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/countries [post]
func (h *CountryHTTPHandler) CreateCountry(c *gin.Context) {
	var req struct {
		CountryCode     string  `json:"country_code" binding:"required"`
		ISO31662        string  `json:"iso_3166_2" binding:"required"`
		ISO31663        string  `json:"iso_3166_3" binding:"required"`
		Name            string  `json:"name" binding:"required"`
		EEA             bool    `json:"eea"`
		CallingCode     string  `json:"calling_code" binding:"required"`
		Capital         *string `json:"capital"`
		Citizenship     *string `json:"citizenship"`
		CurrencyName    *string `json:"currency_name"`
		CurrencyCode    *string `json:"currency_code"`
		CurrencySubUnit *string `json:"currency_sub_unit"`
		CurrencySymbol  *string `json:"currency_symbol"`
		FullName        *string `json:"full_name"`
		RegionCode      *string `json:"region_code"`
		SubRegionCode   *string `json:"sub_region_code"`
		Flag            *string `json:"flag"`
		Latitude        *string `json:"latitude"`
		Longitude       *string `json:"longitude"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	country, err := h.countryService.CreateCountry(
		req.CountryCode, req.ISO31662, req.ISO31663, req.Name, req.EEA, req.CallingCode,
	)
	if err != nil {
		response.InternalServerError(c, "Failed to create country: "+err.Error())
		return
	}

	// Set optional fields
	if req.Capital != nil {
		country.SetCapital(*req.Capital)
	}
	if req.Citizenship != nil {
		country.SetCitizenship(*req.Citizenship)
	}
	if req.CurrencyName != nil && req.CurrencyCode != nil && req.CurrencySubUnit != nil && req.CurrencySymbol != nil {
		country.SetCurrencyInfo(*req.CurrencyName, *req.CurrencyCode, *req.CurrencySubUnit, *req.CurrencySymbol)
	}
	if req.FullName != nil {
		country.SetFullName(*req.FullName)
	}
	if req.RegionCode != nil && req.SubRegionCode != nil {
		country.SetRegionInfo(*req.RegionCode, *req.SubRegionCode)
	}
	if req.Flag != nil {
		country.SetFlag(*req.Flag)
	}
	if req.Latitude != nil && req.Longitude != nil {
		country.SetCoordinates(*req.Latitude, *req.Longitude)
	}

	// Update the country with the additional information
	if err := h.countryService.UpdateCountry(country); err != nil {
		response.InternalServerError(c, "Failed to update country: "+err.Error())
		return
	}

	response.Created(c, country, "Country created successfully")
}

// GetCountryByID handles GET /api/v1/countries/:id
// @Summary Get country by ID
// @Description Get a country by its UUID
// @Tags countries
// @Accept json
// @Produce json
// @Param id path string true "Country ID" format(uuid)
// @Success 200 {object} response.Response{data=entities.Country} "Country found"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 404 {object} response.Response "Country not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/countries/{id} [get]
func (h *CountryHTTPHandler) GetCountryByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid country ID: "+err.Error())
		return
	}

	country, err := h.countryService.GetCountryByID(id)
	if err != nil {
		response.NotFound(c, "Country not found: "+err.Error())
		return
	}

	response.Success(c, country, "Country retrieved successfully")
}

// GetCountryByCode handles GET /api/v1/countries/code/:code
func (h *CountryHTTPHandler) GetCountryByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.BadRequest(c, "Country code is required")
		return
	}

	country, err := h.countryService.GetCountryByCode(code)
	if err != nil {
		response.InternalServerError(c, "Failed to get country: "+err.Error())
		return
	}

	if country == nil {
		response.NotFound(c, "Country not found")
		return
	}

	response.Success(c, country, "Country retrieved successfully")
}

// GetCountryByISO31662 handles GET /api/v1/countries/iso31662/:code
func (h *CountryHTTPHandler) GetCountryByISO31662(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.BadRequest(c, "ISO 3166-2 code is required")
		return
	}

	country, err := h.countryService.GetCountryByISO31662(code)
	if err != nil {
		response.InternalServerError(c, "Failed to get country: "+err.Error())
		return
	}

	if country == nil {
		response.NotFound(c, "Country not found")
		return
	}

	response.Success(c, country, "Country retrieved successfully")
}

// GetCountryByISO31663 handles GET /api/v1/countries/iso31663/:code
func (h *CountryHTTPHandler) GetCountryByISO31663(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.BadRequest(c, "ISO 3166-3 code is required")
		return
	}

	country, err := h.countryService.GetCountryByISO31663(code)
	if err != nil {
		response.InternalServerError(c, "Failed to get country: "+err.Error())
		return
	}

	if country == nil {
		response.NotFound(c, "Country not found")
		return
	}

	response.Success(c, country, "Country retrieved successfully")
}

// GetCountryByName handles GET /api/v1/countries/name/:name
func (h *CountryHTTPHandler) GetCountryByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.BadRequest(c, "Country name is required")
		return
	}

	country, err := h.countryService.GetCountryByName(name)
	if err != nil {
		response.InternalServerError(c, "Failed to get country: "+err.Error())
		return
	}

	if country == nil {
		response.NotFound(c, "Country not found")
		return
	}

	response.Success(c, country, "Country retrieved successfully")
}

// GetAllCountries handles GET /api/v1/countries
func (h *CountryHTTPHandler) GetAllCountries(c *gin.Context) {
	countries, err := h.countryService.GetAllCountries()
	if err != nil {
		response.InternalServerError(c, "Failed to get countries: "+err.Error())
		return
	}

	response.Success(c, countries, "Countries retrieved successfully")
}

// GetCountriesByRegion handles GET /api/v1/countries/region/:regionCode
func (h *CountryHTTPHandler) GetCountriesByRegion(c *gin.Context) {
	regionCode := c.Param("regionCode")
	if regionCode == "" {
		response.BadRequest(c, "Region code is required")
		return
	}

	countries, err := h.countryService.GetCountriesByRegion(regionCode)
	if err != nil {
		response.InternalServerError(c, "Failed to get countries by region: "+err.Error())
		return
	}

	response.Success(c, countries, "Countries retrieved successfully")
}

// GetCountriesBySubRegion handles GET /api/v1/countries/subregion/:subRegionCode
func (h *CountryHTTPHandler) GetCountriesBySubRegion(c *gin.Context) {
	subRegionCode := c.Param("subRegionCode")
	if subRegionCode == "" {
		response.BadRequest(c, "Sub-region code is required")
		return
	}

	countries, err := h.countryService.GetCountriesBySubRegion(subRegionCode)
	if err != nil {
		response.InternalServerError(c, "Failed to get countries by sub-region: "+err.Error())
		return
	}

	response.Success(c, countries, "Countries retrieved successfully")
}

// GetEEACountries handles GET /api/v1/countries/eea
func (h *CountryHTTPHandler) GetEEACountries(c *gin.Context) {
	countries, err := h.countryService.GetEEACountries()
	if err != nil {
		response.InternalServerError(c, "Failed to get EEA countries: "+err.Error())
		return
	}

	response.Success(c, countries, "EEA countries retrieved successfully")
}

// UpdateCountry handles PUT /api/v1/countries/:id
func (h *CountryHTTPHandler) UpdateCountry(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid country ID: "+err.Error())
		return
	}

	var req struct {
		Capital         *string `json:"capital"`
		Citizenship     *string `json:"citizenship"`
		CountryCode     string  `json:"country_code" binding:"required"`
		CurrencyName    *string `json:"currency_name"`
		CurrencyCode    *string `json:"currency_code"`
		CurrencySubUnit *string `json:"currency_sub_unit"`
		CurrencySymbol  *string `json:"currency_symbol"`
		FullName        *string `json:"full_name"`
		ISO31662        string  `json:"iso_3166_2" binding:"required"`
		ISO31663        string  `json:"iso_3166_3" binding:"required"`
		Name            string  `json:"name" binding:"required"`
		RegionCode      *string `json:"region_code"`
		SubRegionCode   *string `json:"sub_region_code"`
		EEA             bool    `json:"eea"`
		CallingCode     string  `json:"calling_code" binding:"required"`
		Flag            *string `json:"flag"`
		Latitude        *string `json:"latitude"`
		Longitude       *string `json:"longitude"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	country, err := h.countryService.GetCountryByID(id)
	if err != nil {
		response.InternalServerError(c, "Failed to get country: "+err.Error())
		return
	}

	if country == nil {
		response.NotFound(c, "Country not found")
		return
	}

	// Update fields
	country.CountryCode = req.CountryCode
	country.ISO31662 = req.ISO31662
	country.ISO31663 = req.ISO31663
	country.Name = req.Name
	country.EEA = req.EEA
	country.CallingCode = req.CallingCode

	if req.Capital != nil {
		country.SetCapital(*req.Capital)
	}
	if req.Citizenship != nil {
		country.SetCitizenship(*req.Citizenship)
	}
	if req.CurrencyName != nil && req.CurrencyCode != nil && req.CurrencySubUnit != nil && req.CurrencySymbol != nil {
		country.SetCurrencyInfo(*req.CurrencyName, *req.CurrencyCode, *req.CurrencySubUnit, *req.CurrencySymbol)
	}
	if req.FullName != nil {
		country.SetFullName(*req.FullName)
	}
	if req.RegionCode != nil && req.SubRegionCode != nil {
		country.SetRegionInfo(*req.RegionCode, *req.SubRegionCode)
	}
	if req.Flag != nil {
		country.SetFlag(*req.Flag)
	}
	if req.Latitude != nil && req.Longitude != nil {
		country.SetCoordinates(*req.Latitude, *req.Longitude)
	}

	if err := h.countryService.UpdateCountry(country); err != nil {
		response.InternalServerError(c, "Failed to update country: "+err.Error())
		return
	}

	response.Success(c, country, "Country updated successfully")
}

// DeleteCountry handles DELETE /api/v1/countries/:id
func (h *CountryHTTPHandler) DeleteCountry(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid country ID: "+err.Error())
		return
	}

	exists, err := h.countryService.CountryExists(id)
	if err != nil {
		response.InternalServerError(c, "Failed to check country existence: "+err.Error())
		return
	}

	if !exists {
		response.NotFound(c, "Country not found")
		return
	}

	if err := h.countryService.DeleteCountry(id); err != nil {
		response.InternalServerError(c, "Failed to delete country: "+err.Error())
		return
	}

	response.Success(c, nil, "Country deleted successfully")
}

// GetCountryCount handles GET /api/v1/countries/count
func (h *CountryHTTPHandler) GetCountryCount(c *gin.Context) {
	count, err := h.countryService.GetCountryCount()
	if err != nil {
		response.InternalServerError(c, "Failed to get country count: "+err.Error())
		return
	}

	response.Success(c, gin.H{"count": count}, "Country count retrieved successfully")
}
