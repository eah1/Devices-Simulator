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

func NewEnvironment(testName, userID string) models.Environment {
	environment := new(models.Environment)
	environment.ID = uuid.NewString()
	environment.UserID = userID
	environment.Name = testName + "_" + faker.Name().Name()
	environment.Vars = map[string]interface{}{}

	return *environment
}
