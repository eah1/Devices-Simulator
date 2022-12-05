// Package responses contains core return API.
package responses

import "device-simulator/business/web/webmodels"

// SuccessLogin response token information.
type SuccessLogin struct {
	Status string `json:"status"  example:"OK"`
	Token  string `json:"token"`
}

// SuccessUser response user information.
type SuccessUser struct {
	Status string                    `json:"status"  example:"OK"`
	User   webmodels.InformationUser `json:"user"`
}
