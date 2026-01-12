package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Candoo/vehicles-api/internal/models"
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

	// Insert vehicles with explicit primary key values
	// Disable auto-increment behavior for primary key by using Omit with empty list
	for i, vehicle := range vehicles {
		// Use Omit("") to force GORM to include the primary key in INSERT
		if err := db.Omit("").Create(&vehicle).Error; err != nil {
			return fmt.Errorf("failed to insert vehicle %d: %w", vehicle.VehicleID, err)
		}

		if (i+1)%10 == 0 || i+1 == len(vehicles) {
			log.Printf("Inserted %d of %d vehicles", i+1, len(vehicles))
		}
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
