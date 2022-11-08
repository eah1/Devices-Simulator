package models

// SendValidationEmail structure send validation email payload.
type SendValidationEmail struct {
	Email           string
	ValidationToken string
	Language        string
}
