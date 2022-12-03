// Package emailsender contains system email config library.
package emailsender

import (
	"bytes"
	"crypto/tls"
	"device-simulator/app/config"
	mycErrors "device-simulator/business/sys/errors"
	"fmt"
	"html/template"
	"net"
	"net/smtp"
	"strconv"

	"github.com/jhillyerd/enmime"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const emailPort = 587

type loginAuth struct {
	username, password string
}

func authorization(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(*smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("unknown from server")
		}
	}

	return nil, nil
}

// InnitEmailConfig create a basic configuration email service.
func InnitEmailConfig(config config.Config) (*enmime.SMTPSender, error) {
	// connection tcp.
	conn, err := net.Dial(config.SMTPNetwork, config.SMTPHost+":"+config.SMTPPort)
	if err != nil {
		return nil, fmt.Errorf("emailsender.InnitEmailConfig.net.Dial(%s, %s): %w",
			config.SMTPNetwork, config.SMTPHost+":"+config.SMTPPort, err)
	}

	client, err := smtp.NewClient(conn, config.SMTPHost)
	if err != nil {
		return nil, fmt.Errorf("emailsender.InnitEmailConfig.smtp.NewClient(%+v, %s): %w",
			conn, config.SMTPHost, err)
	}

	port, err := strconv.Atoi(config.SMTPPort)
	if err != nil {
		return nil, fmt.Errorf("emailsender.InnitEmailConfig.strconv.Atoi(%s): %w", config.SMTPPort, err)
	}

	if port == emailPort {
		// TLS config.
		tlsConfig := new(tls.Config)
		tlsConfig.ServerName = config.SMTPHost

		if err = client.StartTLS(tlsConfig); err != nil {
			return nil, fmt.Errorf("emailsender.InnitEmailConfig.client.StartTLS(%+v): %w", tlsConfig, err)
		}
	}

	// authentication configuration.
	auth := authorization(config.PostmarkToken, config.PostmarkToken)

	if err = client.Auth(auth); err != nil {
		return nil, fmt.Errorf("emailsender.InnitEmailConfig.client.Auth(%+v): %w", auth, err)
	}

	sender := enmime.NewSMTP(config.SMTPHost+":"+config.SMTPPort, auth)

	return sender, nil
}

type EmailStructure struct {
	sender   *enmime.SMTPSender
	smtpFrom string
	tag      string
	title    string
	template string
	receiver []string
	data     interface{}
	log      *zap.SugaredLogger
}

func NewEmailStructure(emailSender *enmime.SMTPSender, emailType, title, template, from string, receiver []string,
	data interface{}, log *zap.SugaredLogger,
) EmailStructure {
	emailStructure := new(EmailStructure)
	emailStructure.sender = emailSender
	emailStructure.tag = emailType
	emailStructure.title = title
	emailStructure.template = template
	emailStructure.smtpFrom = from
	emailStructure.receiver = receiver
	emailStructure.data = data
	emailStructure.log = log

	return *emailStructure
}

// SendEmail send email to receivers.
func SendEmail(templateFolder string, emailStructure EmailStructure) error {
	temp, err := template.ParseFiles(templateFolder + emailStructure.template)
	if err != nil {
		return fmt.Errorf("emailsender.SendEmail.template.ParseFiles(%s): %w",
			templateFolder+emailStructure.template, mycErrors.ErrEmailSender)
	}

	var body bytes.Buffer

	if err := temp.Execute(&body, emailStructure.data); err != nil {
		return fmt.Errorf("emailsender.SendEmail.temp.Execute(%+v, %+v): %w",
			&body, emailStructure.data, mycErrors.ErrEmailSender)
	}

	master := enmime.Builder().
		From("CIRCUTOR", emailStructure.smtpFrom).
		Subject(emailStructure.title).
		HTML(body.Bytes()).Header("X-PM-Tag", emailStructure.tag)

	for _, emailAddress := range emailStructure.receiver {
		msg := master.To(emailAddress, emailAddress)

		if err := msg.Send(emailStructure.sender); err != nil {
			return fmt.Errorf("emailsender.SendEmail.msg.Send(%+v): %w",
				emailStructure.sender, mycErrors.ErrEmailSender)
		}

		emailStructure.log.Infow("email", "sendEmail", emailAddress)
	}

	return nil
}
