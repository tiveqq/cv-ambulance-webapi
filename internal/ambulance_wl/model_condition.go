package ambulance_wl

// Condition represents a medical condition
type Condition struct {
	// Unique identifier of the condition
	Id string `json:"id"`

	// Name of the condition
	Name string `json:"name"`

	// Description of the condition
	Description string `json:"description,omitempty"`

	// Severity level of the condition (1-10)
	Severity int `json:"severity,omitempty"`
}
