// Package responses contains core return API.
package responses

// Validator ResponseFailed contains information needed to create a new Validator.
type Validator struct {
	Status  string      `json:"status"`
	Details interface{} `json:"details"`
}

// Failed ResponseFailed contains information needed to create a new Failed.
type Failed struct {
	Status string      `json:"status"`
	Error  interface{} `json:"error"`
}
