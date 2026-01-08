package repository

import (
	"fmt"
	"strings"

	"github.com/Candoo/vehicle-api/internal/models"
	"gorm.io/gorm"
)

// VehicleRepository handles database operations for vehicles
type VehicleRepository struct {
	db *gorm.DB
}

// NewVehicleRepository creates a new vehicle repository
func NewVehicleRepository(db *gorm.DB) *VehicleRepository {
	return &VehicleRepository{db: db}
}

// GetVehicles retrieves vehicles with pagination and filtering
func (r *VehicleRepository) GetVehicles(filters models.VehicleFilters) ([]models.Vehicle, *models.ResponseMetadata, error) {
	var vehicles []models.Vehicle
	var total int64

	// Build query with filters
	query := r.db.Model(&models.Vehicle{})

	// Apply filters
	if filters.AdvertClassification != "" && strings.ToLower(filters.AdvertClassification) != "all" {
		query = query.Where("LOWER(advert_classification) = ?", strings.ToLower(filters.AdvertClassification))
	}

	if filters.Make != "" {
		query = query.Where("LOWER(make) = ?", strings.ToLower(filters.Make))
	}

	if filters.Model != "" {
		query = query.Where("LOWER(model) LIKE ?", "%"+strings.ToLower(filters.Model)+"%")
	}

	if filters.FuelType != "" {
		query = query.Where("LOWER(fuel_type) = ?", strings.ToLower(filters.FuelType))
	}

	if filters.Transmission != "" {
		query = query.Where("LOWER(transmission) = ?", strings.ToLower(filters.Transmission))
	}

	if filters.BodyType != "" {
		query = query.Where("LOWER(body_type) = ?", strings.ToLower(filters.BodyType))
	}

	if filters.MinPrice != "" && filters.MinPrice != "0" {
		query = query.Where("CAST(price AS INTEGER) >= ?", filters.MinPrice)
	}

	if filters.MaxPrice != "" && filters.MaxPrice != "0" {
		query = query.Where("CAST(price AS INTEGER) <= ?", filters.MaxPrice)
	}

	if filters.MinYear != "" && filters.MinYear != "0" {
		query = query.Where("CAST(year AS INTEGER) >= ?", filters.MinYear)
	}

	if filters.MaxYear != "" && filters.MaxYear != "0" {
		query = query.Where("CAST(year AS INTEGER) <= ?", filters.MaxYear)
	}

	// Count total results
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to count vehicles: %w", err)
	}

	// Set defaults for pagination
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.ResultsPerPage < 1 {
		filters.ResultsPerPage = 10
	}

	// Calculate offset
	offset := (filters.Page - 1) * filters.ResultsPerPage

	// Fetch paginated results
	if err := query.
		Order("vehicle_id ASC").
		Limit(filters.ResultsPerPage).
		Offset(offset).
		Find(&vehicles).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to fetch vehicles: %w", err)
	}

	// Calculate metadata
	lastPage := int(total) / filters.ResultsPerPage
	if int(total)%filters.ResultsPerPage > 0 {
		lastPage++
	}

	// Get additional statistics
	var allTotal, totalNew, totalUsed int64
	r.db.Model(&models.Vehicle{}).Count(&allTotal)
	r.db.Model(&models.Vehicle{}).Where("LOWER(advert_classification) = ?", "new").Count(&totalNew)
	r.db.Model(&models.Vehicle{}).Where("LOWER(advert_classification) = ?", "used").Count(&totalUsed)
	// TODO: Add has_offer field to model if needed
	// r.db.Model(&models.Vehicle{}).Where("has_offer = ?", true).Count(&offerVehicles)

	metadata := &models.ResponseMetadata{
		CurrentPage:       filters.Page,
		LastPage:          lastPage,
		PerPage:           filters.ResultsPerPage,
		Total:             total,
		AllTotal:          allTotal,
		TotalNewVehicles:  totalNew,
		TotalUsedVehicles: totalUsed,
		OfferVehicles:     0, // TODO: Calculate when has_offer field is added
	}

	return vehicles, metadata, nil
}

// GetVehicleByID retrieves a single vehicle by ID
func (r *VehicleRepository) GetVehicleByID(id int) (*models.Vehicle, error) {
	var vehicle models.Vehicle

	if err := r.db.First(&vehicle, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("vehicle not found")
		}
		return nil, fmt.Errorf("failed to fetch vehicle: %w", err)
	}

	return &vehicle, nil
}

// GetVehicleByVRM retrieves a single vehicle by VRM (Vehicle Registration Mark)
func (r *VehicleRepository) GetVehicleByVRM(vrm string) (*models.Vehicle, error) {
	var vehicle models.Vehicle

	if err := r.db.Where("LOWER(vrm) = ?", strings.ToLower(vrm)).First(&vehicle).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("vehicle not found")
		}
		return nil, fmt.Errorf("failed to fetch vehicle: %w", err)
	}

	return &vehicle, nil
}

// GetAvailableMakes retrieves all unique makes
func (r *VehicleRepository) GetAvailableMakes() ([]string, error) {
	var makes []string

	if err := r.db.Model(&models.Vehicle{}).
		Distinct("make").
		Order("make ASC").
		Pluck("make", &makes).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch makes: %w", err)
	}

	return makes, nil
}

// GetAvailableModels retrieves all unique models for a given make
func (r *VehicleRepository) GetAvailableModels(make string) ([]string, error) {
	var modelList []string

	query := r.db.Model(&models.Vehicle{}).Distinct("model").Order("model ASC")

	if make != "" {
		query = query.Where("LOWER(make) = ?", strings.ToLower(make))
	}

	if err := query.Pluck("model", &modelList).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch models: %w", err)
	}

	return modelList, nil
}
