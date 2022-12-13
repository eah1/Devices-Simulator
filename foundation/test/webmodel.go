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
