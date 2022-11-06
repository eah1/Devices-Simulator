package test_test

import (
	"device-simulator/business/web/webmodels"

	"syreclabs.com/go/faker"
)

// NewRegistrationUser created new user registration model.
func NewRegistrationUser(testName string) webmodels.RegisterUser {
	userRegistration := new(webmodels.RegisterUser)
	userRegistration.FirstName = faker.Name().FirstName() + "_" + testName
	userRegistration.LastName = faker.Name().LastName() + "_" + testName
	userRegistration.Email = faker.Internet().Email()
	userRegistration.Password = faker.Internet().Password(8, 64)
	userRegistration.Language = faker.RandomChoice([]string{"en", "es", "fr", "pt"})
	userRegistration.Company = faker.Company().Name()

	return *userRegistration
}
