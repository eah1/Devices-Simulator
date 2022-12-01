// Package core contains core business API.
package core

import (
	"device-simulator/app/config"
	"device-simulator/business/db/store"
	"device-simulator/business/sys/emailsender"
	emailConfig "device-simulator/foundation/email"
	"fmt"

	"github.com/jhillyerd/enmime"
	"go.uber.org/zap"
)

// EmailSenderCore manages the set of API for email sender access.
type EmailSenderCore struct {
	log         *zap.SugaredLogger
	config      config.Config
	store       store.Store
	emailSender *enmime.SMTPSender
	core        *Core
}

// NewEmailSenderCore constructs a core for user API access.
func NewEmailSenderCore(
	log *zap.SugaredLogger, config config.Config, store store.Store, emailSender *enmime.SMTPSender, core *Core,
) EmailSenderCore {
	return EmailSenderCore{
		log:         log,
		config:      config,
		store:       store,
		emailSender: emailSender,
		core:        core,
	}
}

// SendValidationEmail send validation email.
func (c *EmailSenderCore) SendValidationEmail(email, validationToken, language string) error {
	receiver := []string{email}

	urlValidation := "/api/v1/users/activate/"

	emailType := "account-validation"
	template := language + "/account-validation.html"
	title := emailConfig.Subject(emailType, language)

	data := struct {
		Email, ValidationURI string
	}{Email: email, ValidationURI: c.config.BaseURL + urlValidation + validationToken}

	emailStructure := emailsender.NewEmailStructure(
		c.emailSender, emailType, title, template, c.config.SMTPFrom, receiver, data, c.log)

	if err := emailsender.SendEmail(c.config.TemplateFolder, emailStructure); err != nil {
		return fmt.Errorf("core.emailsender.SendValidationEmail: %w", err)
	}

	return nil
}
