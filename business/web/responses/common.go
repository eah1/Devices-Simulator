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

type Success struct {
	Status string `json:"status"  example:"OK"`
}

// SuccessLogin response token information.
type SuccessLogin struct {
	Status string `json:"status"  example:"OK"`
	Token  string `json:"token"`
}
