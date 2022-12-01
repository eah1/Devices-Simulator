package test_test

import (
	"device-simulator/business/web/webmodels"

	"syreclabs.com/go/faker"
)

const (
	minPassword = 8
	maxPassword = 64
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
