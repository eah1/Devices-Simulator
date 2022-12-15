// Package usecase contains the logic of use cases.
package usecase

import (
	"device-simulator/business/core/models"
	"device-simulator/business/web/webmodels"
	"fmt"
)

// CreateDeviceConfig create device config case-use.
func (u *UseCase) CreateDeviceConfig(
	createDevicesConfig webmodels.CreateDeviceConfig, userID string,
) (webmodels.InformationDevicesConfig, error) {
	devicesConfig := models.CreateDevicesConfWebToDevicesConfig(createDevicesConfig)
	devicesConfig.UserID = userID

	if err := u.core.DeviceConfig.Create(&devicesConfig); err != nil {
		return webmodels.InformationDevicesConfig{}, fmt.Errorf("usecase.device_config.Create: %w", err)
	}

	return models.DevicesConfigModelToWeb(devicesConfig), nil
}
