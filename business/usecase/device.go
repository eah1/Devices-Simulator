// Package usecase contains the logic of use cases.
package usecase

import (
	"device-simulator/business/core/models"
	"device-simulator/business/web/webmodels"
	"fmt"
)

// CreateDevice create device case-use.
func (u *UseCase) CreateDevice(
	createDevice webmodels.CreateDevice, userID string,
) (webmodels.InformationDevice, error) {
	device := models.CreateDeviceWebToDevice(createDevice)
	device.UserID = userID

	if err := u.core.Device.Create(&device); err != nil {
		return webmodels.InformationDevice{}, fmt.Errorf("usecase.device.Create: %w", err)
	}

	return models.DeviceModelToWeb(device), nil
}
