package handler

import (
    "net/http"
    "os"
    "sync"

    "github.com/Candoo/vehicles-api/internal/handlers"
    "github.com/Candoo/vehicles-api/internal/repository"
    "github.com/gin-gonic/gin"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var (
    db   *gorm.DB
    once sync.Once
)

func getDB() *gorm.DB {
    once.Do(func() {
        dsn := os.Getenv("POSTGRES_URL")
        if dsn == "" {
            return
        }
        var err error
        db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err != nil {
            return
        }
    })
    return db
}

func Handler(w http.ResponseWriter, r *http.Request) {
    gin.SetMode(gin.ReleaseMode)
    router := gin.New()
    router.Use(gin.Recovery())

    databaseInstance := getDB()
    if databaseInstance == nil {
        http.Error(w, "Database connection failed", http.StatusInternalServerError)
        return
    }

    vehicleRepo := repository.NewVehicleRepository(databaseInstance)
    vehicleHandler := handlers.NewVehicleHandler(vehicleRepo)

    api := router.Group("/api")
    {
        api.GET("/vehicles", vehicleHandler.GetVehicles)
        api.GET("/vehicles/:id", vehicleHandler.GetVehicleByID)
    }

    router.ServeHTTP(w, r)
}