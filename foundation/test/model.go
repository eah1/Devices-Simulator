package test_test

import (
	"device-simulator/business/core/models"
	"device-simulator/business/sys/auth"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"syreclabs.com/go/faker"
)

const (
	sizeRandom  = 18
	tokenRandom = 200
)

// NewUser created new user model.
func NewUser(testName string) models.User {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	user := new(models.User)
	user.ID = uuid.NewString()
	user.FirstName = faker.Name().FirstName() + "_" + testName
	user.LastName = faker.Name().LastName() + "_" + testName
	user.Email = faker.Internet().Email()
	user.Password = string(hash)
	user.Validated = false
	user.ValidationToken = faker.RandomString(sizeRandom)
	user.Language = faker.RandomChoice([]string{"en", "es", "fr", "pt"})
	user.Company = faker.Company().Name()

	return *user
}

// NewCustomClaims created new custom claims.
func NewCustomClaims(testName string) auth.CustomClaims {
	claims := new(auth.CustomClaims)
	claims.StandardClaims = auth.NewStandardClaims()
	claims.ID = uuid.NewString()
	claims.Email = testName + faker.RandomString(sizeRandom) + "@" + faker.Internet().DomainName()

	return *claims
}

// NewAuthentication created new authentication model.
func NewAuthentication(testName, userID string) models.Authentication {
	authentication := new(models.Authentication)
	authentication.ID = uuid.NewString()
	authentication.Token = testName + "_" + faker.RandomString(tokenRandom)
	authentication.UserID = userID
	authentication.Valid = true
	authentication.LoginAt = time.Now()

	return *authentication
}

// NewEnvironment created new environment model.
func NewEnvironment(testName, userID string) models.Environment {
	environment := new(models.Environment)
	environment.ID = uuid.NewString()
	environment.UserID = userID
	environment.Name = testName + "_" + faker.Name().Name()
	environment.Vars = map[string]interface{}{
		"testName": testName,
	}

	return *environment
}

// NewDeviceConfig created new device config model.
func NewDeviceConfig(testName, userID string) models.DeviceConfig {
	deviceConfig := new(models.DeviceConfig)
	deviceConfig.ID = uuid.NewString()
	deviceConfig.Name = testName + "_" + faker.Name().Name()

	deviceConfig.Vars = map[string]interface{}{
		"testName": testName,
	}

	deviceConfig.MetricsFixed = map[string]interface{}{
		"Voltage": map[string]interface{}{
			"minValue": "219",
			"maxValue": "220",
		},
	}

	deviceConfig.MetricsAccumulated = map[string]interface{}{
		"PH_CON_TOT": map[string]interface{}{
			"minValue": "5",
			"maxValue": "10",
		},
	}

	deviceConfig.TypeSend = "MQTT"
	deviceConfig.Payload = "{\"id\": \"$SERIAL_NUMBER\",\"telemetry\":" +
		"{\"PH_CON_TOT\": $PH_CON_TOT},\"deviceType\": \"DEVICE\":}"
	deviceConfig.UserID = userID

	return *deviceConfig
}

// NewDevice created new device model.
func NewDevice(testName, userID, environmentID, deviceConfigID string) models.Device {
	device := new(models.Device)
	device.ID = uuid.NewString()
	device.Name = testName + "_" + faker.Name().Name()
	device.UserID = userID
	device.EnvironmentID = environmentID
	device.DeviceConfigID = deviceConfigID

	return *device
}
