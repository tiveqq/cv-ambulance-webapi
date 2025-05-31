// PatientService defines the interface for patient data operations
package ambulance_wl

import (
	"context"
)

// PatientService defines the interface for patient data operations
type PatientService interface {
	// GetAllPatients retrieves all patients from the database
	GetAllPatients(ctx context.Context) ([]Patient, error)

	// GetPatientByID retrieves a patient by ID
	GetPatientByID(ctx context.Context, id string) (*Patient, error)

	// CreatePatient creates a new patient
	CreatePatient(ctx context.Context, patient PatientInput) (*Patient, error)

	// UpdatePatient updates an existing patient
	UpdatePatient(ctx context.Context, id string, patient PatientInput) (*Patient, error)

	// ArchivePatient archives a patient (soft delete)
	ArchivePatient(ctx context.Context, id string) error
}