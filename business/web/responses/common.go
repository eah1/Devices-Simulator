// Package responses contains core return API.
package responses

// Validator ResponseFailed contains information needed to create a new Validator.
type Validator struct {
	Status  string   `json:"status"  example:"ERROR"`
	Details []string `json:"details"`
}

// Failed ResponseFailed contains information needed to create a new Failed.
type Failed struct {
	Status string `json:"status"  example:"ERROR"`
	Error  string `json:"error"`
}

// Success response success request.
type Success struct {
	Status string `json:"status"  example:"OK"`
}
