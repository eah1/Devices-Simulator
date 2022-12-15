// Package webmodels contains core business API.
//
//nolint:lll
package webmodels

// CreateDeviceConfig  contains information needed to create a new Device Config.
type CreateDeviceConfig struct {
	Name               string                             `json:"name"     conform:"trim" validate:"required"`
	Vars               []*DevicesConfigVars               `json:"vars"                    validate:"required,min=1,dive"`
	MetricsFixed       []*DevicesConfigMetricsFixed       `json:"metricsFixed"            validate:"required,min=1,dive"`
	MetricsAccumulated []*DevicesConfigMetricsAccumulated `json:"metricsAccumulated"      validate:"required,min=1,dive"`
	TypeSend           string                             `json:"typeSend" conform:"trim" validate:"required,oneof=MQTT"`
	Payload            string                             `json:"payload"  conform:"trim" validate:"required"`
}

// DevicesConfigVars fields to complete the CreateEnvironment struct.
type DevicesConfigVars struct {
	Key string `json:"key" conform:"trim" validate:"required"`
	Var string `json:"var" conform:"trim" validate:"required"`
}

// DevicesConfigMetricsFixed fields to generate random values in metrics fixed.
type DevicesConfigMetricsFixed struct {
	Metric       string                       `json:"metric" conform:"trim" validate:"required"`
	RandomValues []*DevicesConfigRandomValues `json:"randomValues"          validate:"required,min=1,dive"`
}

// DevicesConfigMetricsAccumulated fields to generate random values in metrics accumulated.
type DevicesConfigMetricsAccumulated struct {
	Metric       string                       `json:"metric" conform:"trim" validate:"required"`
	RandomValues []*DevicesConfigRandomValues `json:"randomValues"          validate:"required,min=1,dive"`
}

// DevicesConfigRandomValues fields to generate types values of metrics.
type DevicesConfigRandomValues struct {
	TypeValue string `json:"typeValue" conform:"trim" validate:"required"`
	Value     string `json:"value"     conform:"trim" validate:"required"`
}

// InformationDevicesConfig contains information.
type InformationDevicesConfig struct {
	ID                 string                 `json:"id"`
	Name               string                 `json:"name"`
	Vars               map[string]interface{} `json:"vars"               swaggertype:"object,string" example:"key:value,key2:value2"`
	MetricsFixed       map[string]interface{} `json:"metricsFixed"       swaggertype:"object,string" example:"key:value,key2:value2"`
	MetricsAccumulated map[string]interface{} `json:"metricsAccumulated" swaggertype:"object,string" example:"key:value,key2:value2"`
	TypeSend           string                 `json:"typeSend"`
	Payload            string                 `json:"payload"`
}
