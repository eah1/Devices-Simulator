package test_test

import (
	"device-simulator/business/core/models"
	"device-simulator/business/sys/auth"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"syreclabs.com/go/faker"
)

const sizeRandom = 18

// NewUser created new user model.
func NewUser(testName string) models.User {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	user := new(models.User)
	user.ID = uuid.New().String()
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
