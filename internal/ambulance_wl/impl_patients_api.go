package ambulance_wl

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// implPatientsAPI implements the PatientsAPI interface
type implPatientsAPI struct {
	dbService PatientService
}

// NewPatientsAPI creates a new PatientsAPI implementation
func NewPatientsAPI(dbService PatientService) PatientsAPI {
	return &implPatientsAPI{
		dbService: dbService,
	}
}

// ArchivePatient handles DELETE /api/patients/:patientId
func (api *implPatientsAPI) ArchivePatient(c *gin.Context) {
	patientID := c.Param("patientId")
	if patientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Patient ID is required"})
		return
	}

	err := api.dbService.ArchivePatient(c.Request.Context(), patientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// CreatePatient handles POST /api/patients
func (api *implPatientsAPI) CreatePatient(c *gin.Context) {
	var patientInput PatientInput
	if err := c.ShouldBindJSON(&patientInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if patientInput.Name == "" || patientInput.Condition == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and condition are required"})
		return
	}

	patient, err := api.dbService.CreatePatient(c.Request.Context(), patientInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, patient)
}

// GetPatient handles GET /api/patients/:patientId
func (api *implPatientsAPI) GetPatient(c *gin.Context) {
	patientID := c.Param("patientId")
	if patientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Patient ID is required"})
		return
	}

	patient, err := api.dbService.GetPatientByID(c.Request.Context(), patientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if patient == nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, patient)
}

// GetPatients handles GET /api/patients
func (api *implPatientsAPI) GetPatients(c *gin.Context) {
	patients, err := api.dbService.GetAllPatients(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, patients)
}

// UpdatePatient handles PUT /api/patients/:patientId
func (api *implPatientsAPI) UpdatePatient(c *gin.Context) {
	patientID := c.Param("patientId")
	if patientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Patient ID is required"})
		return
	}

	var patientInput PatientInput
	if err := c.ShouldBindJSON(&patientInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if patientInput.Name == "" || patientInput.Condition == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and condition are required"})
		return
	}

	patient, err := api.dbService.UpdatePatient(c.Request.Context(), patientID, patientInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if patient == nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, patient)
}
