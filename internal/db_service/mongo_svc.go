package db_service

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tiveqq/cv-ambulance-webapi/internal/ambulance_wl"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBService represents a service for interacting with MongoDB
type MongoDBService struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

// NewMongoDBService creates a new MongoDB service
func NewMongoDBService() (*MongoDBService, error) {
	// Get MongoDB connection details from environment variables
	mongoURI := os.Getenv("AMBULANCE_API_MONGODB_URI")
	if mongoURI == "" {
		// Default MongoDB URI if not provided
		username := os.Getenv("AMBULANCE_API_MONGODB_USERNAME")
		if username == "" {
			username = "root"
		}
		password := os.Getenv("AMBULANCE_API_MONGODB_PASSWORD")
		if password == "" {
			password = "example"
		}
		mongoURI = fmt.Sprintf("mongodb://%s:%s@mongo_db:27017", username, password)
	}

	dbName := os.Getenv("AMBULANCE_API_MONGODB_DATABASE")
	if dbName == "" {
		dbName = "ambulance"
	}

	collectionName := os.Getenv("AMBULANCE_API_MONGODB_COLLECTION")
	if collectionName == "" {
		collectionName = "patients"
	}

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Connected to MongoDB!")

	// Get database and collection
	database := client.Database(dbName)
	collection := database.Collection(collectionName)

	return &MongoDBService{
		client:     client,
		database:   database,
		collection: collection,
	}, nil
}

// Close closes the MongoDB connection
func (s *MongoDBService) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.client.Disconnect(ctx)
}

// GetAllPatients retrieves all patients from the database
func (s *MongoDBService) GetAllPatients(ctx context.Context) ([]ambulance_wl.Patient, error) {
	var patients []ambulance_wl.Patient

	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find patients: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &patients); err != nil {
		return nil, fmt.Errorf("failed to decode patients: %w", err)
	}

	return patients, nil
}

// GetPatientByID retrieves a patient by ID
func (s *MongoDBService) GetPatientByID(ctx context.Context, id string) (*ambulance_wl.Patient, error) {
	var patient ambulance_wl.Patient

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// If the ID is not a valid ObjectID, try to find by string ID
		err = s.collection.FindOne(ctx, bson.M{"id": id}).Decode(&patient)
	} else {
		// Try to find by ObjectID
		err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&patient)
		// If found by ObjectID but the ID field is empty, set it
		if err == nil && patient.Id == "" {
			patient.Id = id
		}
	}

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to find patient: %w", err)
	}

	return &patient, nil
}

// CreatePatient creates a new patient
func (s *MongoDBService) CreatePatient(ctx context.Context, patient ambulance_wl.PatientInput) (*ambulance_wl.Patient, error) {
	// Generate a new ObjectID
	objID := primitive.NewObjectID()

	// Create a new Patient from the input
	newPatient := ambulance_wl.Patient{
		Id:                     objID.Hex(),
		Name:                   patient.Name,
		Condition:              patient.Condition,
		DiagnosisDate:          patient.DiagnosisDate,
		TreatmentStartDate:     patient.TreatmentStartDate,
		ExpectedCompletionDate: patient.ExpectedCompletionDate,
		Status:                 patient.Status,
		// Set default values if needed
	}

	if newPatient.Status == "" {
		newPatient.Status = "new"
	}

	// Set a default doctor ID if not provided
	newPatient.DoctorId = "doctor1"

	// Insert the patient into the database
	_, err := s.collection.InsertOne(ctx, newPatient)
	if err != nil {
		return nil, fmt.Errorf("failed to create patient: %w", err)
	}

	return &newPatient, nil
}

// UpdatePatient updates an existing patient
func (s *MongoDBService) UpdatePatient(ctx context.Context, id string, patient ambulance_wl.PatientInput) (*ambulance_wl.Patient, error) {
	// First, check if the patient exists
	existingPatient, err := s.GetPatientByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existingPatient == nil {
		return nil, nil // Not found
	}

	// Update the patient fields
	updatedPatient := ambulance_wl.Patient{
		Id:                     id,
		Name:                   patient.Name,
		Condition:              patient.Condition,
		DiagnosisDate:          patient.DiagnosisDate,
		TreatmentStartDate:     patient.TreatmentStartDate,
		ExpectedCompletionDate: patient.ExpectedCompletionDate,
		Status:                 patient.Status,
		DoctorId:               existingPatient.DoctorId, // Preserve the doctor ID
	}

	// If status is empty, keep the existing status
	if updatedPatient.Status == "" {
		updatedPatient.Status = existingPatient.Status
	}

	// Update the patient in the database
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// If the ID is not a valid ObjectID, update by string ID
		_, err = s.collection.ReplaceOne(ctx, bson.M{"id": id}, updatedPatient)
	} else {
		// Try to update by ObjectID
		_, err = s.collection.ReplaceOne(ctx, bson.M{"_id": objID}, updatedPatient)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to update patient: %w", err)
	}

	return &updatedPatient, nil
}

// ArchivePatient archives a patient (soft delete)
func (s *MongoDBService) ArchivePatient(ctx context.Context, id string) error {
	// First, check if the patient exists
	existingPatient, err := s.GetPatientByID(ctx, id)
	if err != nil {
		return err
	}
	if existingPatient == nil {
		return nil // Not found
	}

	// Set the patient status to "archived"
	update := bson.M{
		"$set": bson.M{
			"status": "archived",
		},
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// If the ID is not a valid ObjectID, update by string ID
		_, err = s.collection.UpdateOne(ctx, bson.M{"id": id}, update)
	} else {
		// Try to update by ObjectID
		_, err = s.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	}

	if err != nil {
		return fmt.Errorf("failed to archive patient: %w", err)
	}

	return nil
}
