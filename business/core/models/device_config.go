package models

import (
	"device-simulator/business/web/webmodels"
	"time"
)

// DeviceConfig represents the structure we need for moving data.
type DeviceConfig struct {
	ID                 string                 `xorm:"pk not null 'id'"          json:"id"`
	Name               string                 `xorm:"name"                      json:"name"`
	Vars               map[string]interface{} `xorm:"jsonb vars"                json:"vars"`
	MetricsFixed       map[string]interface{} `xorm:"jsonb metrics_fixed"       json:"metricsFixed"`
	MetricsAccumulated map[string]interface{} `xorm:"jsonb metrics_accumulated" json:"metricsAccumulated"`
	TypeSend           string                 `xorm:"type_send"                 json:"typeSend"`
	Payload            string                 `xorm:"payload"                   json:"payload"`
	UserID             string                 `xorm:"user_id"                   json:"userId"`
	CreatedAt          time.Time              `xorm:"created"                   json:"-"`
	UpdatedAt          time.Time              `xorm:"updated"                   json:"-"`
}

func (*DeviceConfig) TableName() string {
	return "devices_config"
}

// CreateDevicesConfWebToDevicesConfig convert struct web model webmodels.CreateDeviceConfig to DeviceConfig model.
func CreateDevicesConfWebToDevicesConfig(createDeviceConfig webmodels.CreateDeviceConfig) DeviceConfig {
	devicesConfig := new(DeviceConfig)
	devicesConfig.Name = createDeviceConfig.Name
	devicesConfig.TypeSend = createDeviceConfig.TypeSend
	devicesConfig.Payload = createDeviceConfig.Payload

	devicesConfig.Vars = make(map[string]interface{})
	devicesConfig.MetricsFixed = make(map[string]interface{})
	devicesConfig.MetricsAccumulated = make(map[string]interface{})

	for _, vars := range createDeviceConfig.Vars {
		devicesConfig.Vars[vars.Key] = vars.Var
	}

	for _, metrics := range createDeviceConfig.MetricsFixed {
		randomValues := make(map[string]string)
		for _, random := range metrics.RandomValues {
			randomValues[random.TypeValue] = random.Value
		}

		devicesConfig.MetricsFixed[metrics.Metric] = randomValues
	}

	for _, metrics := range createDeviceConfig.MetricsAccumulated {
		randomValues := make(map[string]string)
		for _, random := range metrics.RandomValues {
			randomValues[random.TypeValue] = random.Value
		}

		devicesConfig.MetricsAccumulated[metrics.Metric] = randomValues
	}

	return *devicesConfig
}

// DevicesConfigModelToWeb convert struct model DevicesConfig to webmodels.InformationDevicesConfig.
func DevicesConfigModelToWeb(devicesConfig DeviceConfig) webmodels.InformationDevicesConfig {
	devicesConfigInfo := new(webmodels.InformationDevicesConfig)
	devicesConfigInfo.ID = devicesConfig.ID
	devicesConfigInfo.Name = devicesConfig.Name
	devicesConfigInfo.Vars = devicesConfig.Vars
	devicesConfigInfo.MetricsFixed = devicesConfig.MetricsFixed
	devicesConfigInfo.MetricsAccumulated = devicesConfig.MetricsAccumulated
	devicesConfigInfo.TypeSend = devicesConfig.TypeSend
	devicesConfigInfo.Payload = devicesConfig.Payload

	return *devicesConfigInfo
}
