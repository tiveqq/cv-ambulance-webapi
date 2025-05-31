# MongoDB Database Service

This package provides a MongoDB service implementation for the ambulance-webapi application.

## Overview

The MongoDB service provides CRUD operations for patient data, including:

- Creating new patients
- Retrieving patients by ID
- Retrieving all patients
- Updating existing patients
- Archiving patients (soft delete)

## Configuration

The service can be configured using the following environment variables:

- `AMBULANCE_API_MONGODB_URI`: The MongoDB connection URI (optional)
- `AMBULANCE_API_MONGODB_USERNAME`: The MongoDB username (default: "root")
- `AMBULANCE_API_MONGODB_PASSWORD`: The MongoDB password (default: "example")
- `AMBULANCE_API_MONGODB_DATABASE`: The MongoDB database name (default: "ambulance")
- `AMBULANCE_API_MONGODB_COLLECTION`: The MongoDB collection name (default: "patients")

If `AMBULANCE_API_MONGODB_URI` is not provided, the service will construct a URI using the username and password.

## Usage

```go
// Initialize the MongoDB service
mongoService, err := db_service.NewMongoDBService()
if err != nil {
    log.Fatalf("Failed to initialize MongoDB service: %v", err)
}
defer mongoService.Close()

// Use the service to perform operations
patients, err := mongoService.GetAllPatients(context.Background())
if err != nil {
    log.Fatalf("Failed to get patients: %v", err)
}

// Create a new patient
newPatient, err := mongoService.CreatePatient(context.Background(), patientInput)
if err != nil {
    log.Fatalf("Failed to create patient: %v", err)
}

// Get a patient by ID
patient, err := mongoService.GetPatientByID(context.Background(), "patient-id")
if err != nil {
    log.Fatalf("Failed to get patient: %v", err)
}

// Update a patient
updatedPatient, err := mongoService.UpdatePatient(context.Background(), "patient-id", patientInput)
if err != nil {
    log.Fatalf("Failed to update patient: %v", err)
}

// Archive a patient
err = mongoService.ArchivePatient(context.Background(), "patient-id")
if err != nil {
    log.Fatalf("Failed to archive patient: %v", err)
}
```

## Error Handling

The service returns appropriate errors for various failure scenarios:

- Connection failures
- Database operation failures
- Not found errors (returns nil for patient and nil error)

## Dependencies

This service depends on the official MongoDB Go driver:

```
go.mongodb.org/mongo-driver v1.14.0
```