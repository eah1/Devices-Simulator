package email_test

import (
	"testing"

	"device-simulator/foundation/email"
	"github.com/stretchr/testify/assert"
)

func TestSubject(t *testing.T) {
	t.Parallel()

	t.Log("Given the need to obtain the subject of an email template.")
	{
		t.Logf("\tWhen we retrieve the subject of a account-validation.")
		{
			assert.Equal(t, "Welcome", email.Subject("account-validation", "en"))
			assert.Equal(t, "Bienvenidos", email.Subject("account-validation", "es"))
			assert.Equal(t, "Bienvenu", email.Subject("account-validation", "fr"))
			assert.Equal(t, "Bem-vindo", email.Subject("account-validation", "pt"))
		}
	}
}
