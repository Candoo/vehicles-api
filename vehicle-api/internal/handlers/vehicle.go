package handlers

import (
	"net/http"
	"strconv"

	"github.com/Candoo/vehicle-api/internal/models"
	"github.com/Candoo/vehicle-api/internal/repository"
	"github.com/gin-gonic/gin"
)

// VehicleHandler handles HTTP requests for vehicles
type VehicleHandler struct {
	repo *repository.VehicleRepository
}

// NewVehicleHandler creates a new vehicle handler
func NewVehicleHandler(repo *repository.VehicleRepository) *VehicleHandler {
	return &VehicleHandler{repo: repo}
}

// GetVehicles godoc
// @Summary Get list of vehicles
// @Description Get a paginated list of vehicles with optional filtering
// @Tags vehicles
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param results_per_page query int false "Results per page" default(10)
// @Param advert_classification query string false "Advertisement classification (New, Used, All)" Enums(New, Used, All)
// @Param make query string false "Vehicle make"
// @Param model query string false "Vehicle model"
// @Param fuel_type query string false "Fuel type"
// @Param transmission query string false "Transmission type"
// @Param body_type query string false "Body type"
// @Param min_price query string false "Minimum price"
// @Param max_price query string false "Maximum price"
// @Param min_year query string false "Minimum year"
// @Param max_year query string false "Maximum year"
// @Success 200 {object} models.VehicleResponse
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /vehicles [get]
func (h *VehicleHandler) GetVehicles(c *gin.Context) {
	// Parse query parameters
	filters := models.VehicleFilters{
		Page:                   parseIntQuery(c, "page", 1),
		ResultsPerPage:         parseIntQuery(c, "results_per_page", 10),
		AdvertClassification:   c.Query("advert_classification"),
		Make:                   c.Query("make"),
		Model:                  c.Query("model"),
		FuelType:               c.Query("fuel_type"),
		Transmission:           c.Query("transmission"),
		BodyType:               c.Query("body_type"),
		MinPrice:               c.Query("min_price"),
		MaxPrice:               c.Query("max_price"),
		MinYear:                c.Query("min_year"),
		MaxYear:                c.Query("max_year"),
	}

	// Validate page and results_per_page
	if filters.Page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "page must be greater than 0",
		})
		return
	}

	if filters.ResultsPerPage < 1 || filters.ResultsPerPage > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "results_per_page must be between 1 and 100",
		})
		return
	}

	// Fetch vehicles from repository
	vehicles, metadata, err := h.repo.GetVehicles(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch vehicles",
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, models.VehicleResponse{
		Data: vehicles,
		Meta: *metadata,
	})
}

// GetVehicleByID godoc
// @Summary Get vehicle by ID
// @Description Get a single vehicle by its ID
// @Tags vehicles
// @Accept json
// @Produce json
// @Param id path int true "Vehicle ID"
// @Success 200 {object} models.Vehicle
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Vehicle not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /vehicles/{id} [get]
func (h *VehicleHandler) GetVehicleByID(c *gin.Context) {
	// Parse ID from path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid vehicle ID",
		})
		return
	}

	// Fetch vehicle from repository
	vehicle, err := h.repo.GetVehicleByID(uint(id))
	if err != nil {
		if err.Error() == "vehicle not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "vehicle not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch vehicle",
		})
		return
	}

	// Return vehicle
	c.JSON(http.StatusOK, vehicle)
}

// GetVehicleByVRM godoc
// @Summary Get vehicle by VRM
// @Description Get a single vehicle by its Vehicle Registration Mark (VRM)
// @Tags vehicles
// @Accept json
// @Produce json
// @Param vrm path string true "Vehicle Registration Mark"
// @Success 200 {object} models.Vehicle
// @Failure 404 {object} map[string]interface{} "Vehicle not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /vehicles/vrm/{vrm} [get]
func (h *VehicleHandler) GetVehicleByVRM(c *gin.Context) {
	vrm := c.Param("vrm")

	// Fetch vehicle from repository
	vehicle, err := h.repo.GetVehicleByVRM(vrm)
	if err != nil {
		if err.Error() == "vehicle not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "vehicle not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch vehicle",
		})
		return
	}

	// Return vehicle
	c.JSON(http.StatusOK, vehicle)
}

// GetAvailableMakes godoc
// @Summary Get available makes
// @Description Get a list of all available vehicle makes
// @Tags vehicles
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of makes"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /vehicles/makes [get]
func (h *VehicleHandler) GetAvailableMakes(c *gin.Context) {
	makes, err := h.repo.GetAvailableMakes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch makes",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"makes": makes,
	})
}

// GetAvailableModels godoc
// @Summary Get available models
// @Description Get a list of all available vehicle models, optionally filtered by make
// @Tags vehicles
// @Accept json
// @Produce json
// @Param make query string false "Filter by vehicle make"
// @Success 200 {object} map[string]interface{} "List of models"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /vehicles/models [get]
func (h *VehicleHandler) GetAvailableModels(c *gin.Context) {
	make := c.Query("make")

	models, err := h.repo.GetAvailableModels(make)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch models",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"models": models,
	})
}

// parseIntQuery parses an integer query parameter with a default value
func parseIntQuery(c *gin.Context, key string, defaultValue int) int {
	valueStr := c.Query(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
