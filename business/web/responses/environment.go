// Package responses contains core return API.
package responses

import "device-simulator/business/web/webmodels"

// SuccessEnvironment response environment information.
type SuccessEnvironment struct {
	Status      string                           `json:"status"  example:"OK"`
	Environment webmodels.InformationEnvironment `json:"environment"`
}
