package test_test

import (
	"device-simulator/business/web/webmodels"

	"syreclabs.com/go/faker"
)

const (
	minPassword = 8
	maxPassword = 64
	randomise   = 10
)

// NewRegistrationUser created new user registration model.
func NewRegistrationUser(testName string) webmodels.RegisterUser {
	userRegistration := new(webmodels.RegisterUser)
	userRegistration.FirstName = faker.Name().FirstName() + "_" + testName
	userRegistration.LastName = faker.Name().LastName() + "_" + testName
	userRegistration.Email = faker.Internet().Email()
	userRegistration.Password = faker.Internet().Password(minPassword, maxPassword)
	userRegistration.Language = faker.RandomChoice([]string{"en", "es", "fr", "pt"})
	userRegistration.Company = faker.Company().Name()

	return *userRegistration
}

// NewUpdateUser created new user update model.
func NewUpdateUser(testName string) webmodels.UpdateUser {
	userUpdate := new(webmodels.UpdateUser)
	userUpdate.FirstName = faker.Name().FirstName() + "_" + testName
	userUpdate.LastName = faker.Name().LastName() + "_" + testName
	userUpdate.Language = faker.RandomChoice([]string{"en", "es", "fr", "pt"})
	userUpdate.Company = faker.Company().Name()

	return *userUpdate
}

// NewUpdatePasswordUser created new user update password model.
func NewUpdatePasswordUser(currentPassword string) webmodels.UpdatePasswordUser {
	userUpdatePassword := new(webmodels.UpdatePasswordUser)
	userUpdatePassword.CurrentPassword = currentPassword
	userUpdatePassword.NewPassword = faker.Internet().Password(minPassword, maxPassword)

	return *userUpdatePassword
}

// NewCreateEnvironment created new environment model.
func NewCreateEnvironment(testName string) webmodels.CreateEnvironment {
	createEnvironment := new(webmodels.CreateEnvironment)
	createEnvironment.Name = faker.Name().Name() + "_" + testName
	createEnvironment.Vars = append(createEnvironment.Vars,
		&webmodels.EnvironmentVars{Key: "URI", Var: faker.Internet().Url()})
	createEnvironment.Vars = append(createEnvironment.Vars,
		&webmodels.EnvironmentVars{Key: "SECRET_KEY", Var: faker.RandomString(randomise)})

	return *createEnvironment
}

// NewCreateDevicesConfig created new devices config model.
func NewCreateDevicesConfig(testName string) webmodels.CreateDeviceConfig {
	createdDevicesConfig := new(webmodels.CreateDeviceConfig)
	createdDevicesConfig.Name = faker.Name().Name() + "_" + testName
	createdDevicesConfig.Payload = "{\"id\": \"$SERIAL_NUMBER\",\"telemetry\":" +
		"{\"PH_CON_TOT\": $PH_CON_TOT},\"deviceType\": \"DEVICE\":}"
	createdDevicesConfig.TypeSend = "MQTT"

	createdDevicesConfig.Vars = append(createdDevicesConfig.Vars,
		&webmodels.DevicesConfigVars{Key: "testName", Var: testName})
	createdDevicesConfig.Vars = append(createdDevicesConfig.Vars,
		&webmodels.DevicesConfigVars{Key: "SERIAL_NUMBER", Var: faker.RandomString(randomise)})
	createdDevicesConfig.Vars = append(createdDevicesConfig.Vars,
		&webmodels.DevicesConfigVars{Key: "GATEWAY_TOKEN", Var: faker.RandomString(randomise)})

	randomValuesFixed := make([]*webmodels.DevicesConfigRandomValues, 0)
	randomValuesFixed = append(randomValuesFixed,
		&webmodels.DevicesConfigRandomValues{TypeValue: "minValue", Value: "210"})
	randomValuesFixed = append(randomValuesFixed,
		&webmodels.DevicesConfigRandomValues{TypeValue: "maxValue", Value: "220"})

	randomValuesAccumulated := make([]*webmodels.DevicesConfigRandomValues, 0)
	randomValuesAccumulated = append(randomValuesAccumulated,
		&webmodels.DevicesConfigRandomValues{TypeValue: "minValue", Value: "15"})
	randomValuesAccumulated = append(randomValuesAccumulated,
		&webmodels.DevicesConfigRandomValues{TypeValue: "maxValue", Value: "30"})

	createdDevicesConfig.MetricsFixed = append(createdDevicesConfig.MetricsFixed,
		&webmodels.DevicesConfigMetricsFixed{Metric: "Voltage", RandomValues: randomValuesFixed})

	createdDevicesConfig.MetricsAccumulated = append(createdDevicesConfig.MetricsAccumulated,
		&webmodels.DevicesConfigMetricsAccumulated{Metric: "PH_CON_TOT", RandomValues: randomValuesAccumulated})

	return *createdDevicesConfig
}
