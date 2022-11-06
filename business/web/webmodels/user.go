// Package webmodels contains core business API.
package webmodels

// RegisterUser contains information needed to create a new User.
type RegisterUser struct {
	FirstName string `json:"firstName" conform:"trim" validate:"required"`
	LastName  string `json:"lastName"  conform:"trim" validate:"required"`
	Email     string `json:"email"     conform:"trim" validate:"required,email"`
	Password  string `json:"password"  conform:"trim" validate:"required,min=8,max=64,printascii"`
	Company   string `json:"company"   conform:"trim" validate:"required"`
	Language  string `json:"language"  conform:"trim" validate:"required,oneof=es en fr pt,max=2"`
}
