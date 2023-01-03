// Package responses contains core return API.
package responses

import "device-simulator/business/web/webmodels"

// SuccessDevice response device information.
type SuccessDevice struct {
	Status string                      `json:"status"  example:"OK"`
	Device webmodels.InformationDevice `json:"device"`
}
