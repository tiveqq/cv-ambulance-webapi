package ambulance_wl

import "github.com/gin-gonic/gin"

// AmbulanceConditionsAPI defines the interface for ambulance conditions operations
type AmbulanceConditionsAPI interface {
	// GetConditions retrieves the list of available medical conditions
	GetConditions(c *gin.Context)
}
