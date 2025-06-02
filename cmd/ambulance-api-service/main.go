package main

import (
	"github.com/gin-contrib/cors"
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

	// Initialize router
	engine := gin.New()
	engine.Use(gin.Recovery())

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost"}, // frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// OpenAPI documentation
	//engine.GET("/openapi", api.HandleOpenApi)

	handleFunctions := &ambulance_wl.ApiHandleFunctions{
		PatientsAPI: ambulance_wl.NewPatientsAPI(mongoService),
		//AmbulanceConditionsAPI: ambulance_wl.NewAmbulanceConditionsApi(),
	}

	// Initialize API routes
	ambulance_wl.NewRouterWithGinEngine(engine, *handleFunctions)

	// Start the server
	log.Printf("Server listening on port %s", port)
	engine.GET("/openapi", api.HandleOpenApi)
	if err := engine.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
