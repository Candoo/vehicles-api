package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Candoo/vehicle-api/internal/models"
	"gorm.io/gorm"
)

// SeedDatabase seeds the database with initial vehicle data
func SeedDatabase(db *gorm.DB) error {
	// Check if database already has data
	var count int64
	if err := db.Model(&models.Vehicle{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count vehicles: %w", err)
	}

	if count > 0 {
		log.Printf("Database already contains %d vehicles, skipping seed", count)
		return nil
	}

	log.Println("Seeding database with vehicle data...")

	// Load vehicle data from JSON file
	vehicles, err := loadVehicleData()
	if err != nil {
		return fmt.Errorf("failed to load vehicle data: %w", err)
	}

	// Insert vehicles in batches for better performance
	batchSize := 100
	for i := 0; i < len(vehicles); i += batchSize {
		end := i + batchSize
		if end > len(vehicles) {
			end = len(vehicles)
		}

		batch := vehicles[i:end]
		if err := db.Create(&batch).Error; err != nil {
			return fmt.Errorf("failed to insert batch: %w", err)
		}

		log.Printf("Inserted batch %d-%d of %d vehicles", i+1, end, len(vehicles))
	}

	log.Printf("Successfully seeded database with %d vehicles", len(vehicles))
	return nil
}

// loadVehicleData loads vehicle data from the JSON file
func loadVehicleData() ([]models.Vehicle, error) {
	// Try multiple possible paths for the JSON file
	possiblePaths := []string{
		"scripts/nexuspoint_vehicles.json",
		"../scripts/nexuspoint_vehicles.json",
		"../../scripts/nexuspoint_vehicles.json",
		filepath.Join("vehicle-api", "scripts", "nexuspoint_vehicles.json"),
		// Fallback to old data
		"scripts/vehicle_data.json",
		"../scripts/vehicle_data.json",
	}

	var jsonFile *os.File
	var err error
	var usedPath string

	for _, path := range possiblePaths {
		jsonFile, err = os.Open(path)
		if err == nil {
			usedPath = path
			break
		}
	}

	if jsonFile == nil {
		return nil, fmt.Errorf("could not find vehicle_data.json in any of the expected locations: %v", possiblePaths)
	}
	defer jsonFile.Close()

	log.Printf("Loading vehicle data from: %s", usedPath)

	var vehicles []models.Vehicle
	decoder := json.NewDecoder(jsonFile)
	if err := decoder.Decode(&vehicles); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return vehicles, nil
}
