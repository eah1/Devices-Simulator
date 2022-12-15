// Package responses contains core return API.
package responses

import "device-simulator/business/web/webmodels"

// SuccessDeviceConfig response device config information.
type SuccessDeviceConfig struct {
	Status       string                             `json:"status"  example:"OK"`
	DeviceConfig webmodels.InformationDevicesConfig `json:"deviceConfig"`
}
