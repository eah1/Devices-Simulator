// Package models contains structs of the model application.
package models

import (
	"device-simulator/business/web/webmodels"
	"time"
)

// Device represents the structure we need for moving data.
type Device struct {
	ID             string    `xorm:"pk not null 'id'" json:"id"`
	Name           string    `xorm:"name"             json:"name"`
	UserID         string    `xorm:"user_id"          json:"userId"`
	EnvironmentID  string    `xorm:"environment_id"   json:"environmentId"`
	DeviceConfigID string    `xorm:"device_config_id" json:"deviceConfigId"`
	CreatedAt      time.Time `xorm:"created"          json:"-"`
	UpdatedAt      time.Time `xorm:"updated"          json:"-"`
}

func (*Device) TableName() string {
	return "devices"
}

// CreateDeviceWebToDevice convert struct web model webmodels.CreateDevice to Device model.
func CreateDeviceWebToDevice(createDevice webmodels.CreateDevice) Device {
	device := new(Device)
	device.Name = createDevice.Name
	device.EnvironmentID = createDevice.EnvironmentID
	device.DeviceConfigID = createDevice.DeviceConfigID

	return *device
}

// DeviceModelToWeb convert struct model Device to webmodels.InformationDevice.
func DeviceModelToWeb(device Device) webmodels.InformationDevice {
	deviceInfo := new(webmodels.InformationDevice)
	deviceInfo.ID = device.ID
	deviceInfo.Name = device.Name

	return *deviceInfo
}
