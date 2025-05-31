package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tiveqq/cv-ambulance-webapi/api"
	"github.com/tiveqq/cv-ambulance-webapi/internal/ambulance_wl"
	"github.com/tiveqq/cv-ambulance-webapi/internal/db_service"
	"log"
	"os"
	"strings"
)

func main() {
	log.Printf("Server started")
	port := os.Getenv("AMBULANCE_API_PORT")
	if port == "" {
		port = "8080"
	}
	environment := os.Getenv("AMBULANCE_API_ENVIRONMENT")
	if !strings.EqualFold(environment, "production") { // case insensitive comparison
		gin.SetMode(gin.DebugMode)
	}

	// Initialize MongoDB service
	mongoService, err := db_service.NewMongoDBService()
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB service: %v", err)
	}
	defer mongoService.Close()

	// Initialize API handlers
	patientsAPI := ambulance_wl.NewPatientsAPI(mongoService)
	ambulanceConditionsAPI := ambulance_wl.NewAmbulanceConditionsApi()

	// Initialize router
	engine := gin.New()
	engine.Use(gin.Recovery())

	// OpenAPI documentation
	engine.GET("/openapi", api.HandleOpenApi)

	// Initialize API routes
	ambulance_wl.NewRouterWithGinEngine(engine, ambulance_wl.ApiHandleFunctions{
		PatientsAPI:            patientsAPI,
		AmbulanceConditionsAPI: ambulanceConditionsAPI,
	})

	// Start the server
	log.Printf("Server listening on port %s", port)
	if err := engine.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
