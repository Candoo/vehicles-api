package handler

import (
	"net/http"
	"os"

	"github.com/Candoo/vehicles-api/internal/database"
	"github.com/Candoo/vehicles-api/internal/handlers"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	// Initialize DB once for the serverless instance
	var err error
	dsn := os.Getenv("POSTGRES_URL")
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
}

// Handler is the entry point for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)
	
	router := gin.New()
	router.Use(gin.Recovery())

	// Replicating your routes from main.go
	vehicleHandler := handlers.NewVehicleHandler(db)
	
	api := router.Group("/api")
	{
		api.GET("/vehicles", vehicleHandler.GetVehicles)
		api.GET("/vehicles/:id", vehicleHandler.GetVehicleByID)
	}

	// This sends the request through the Gin router
	router.ServeHTTP(w, r)
}