package core_test

import (
	"testing"

	"device-simulator/business/core"
	"device-simulator/business/db/store"
	tt "device-simulator/foundation/test"
	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func TestEmailSenderSendValidationEmail(t *testing.T) {
	t.Parallel()

	testName := "core-email-sender-send-validation-email"

	// Create store.
	newLog := tt.InitLogger(t, "t-"+testName)
	newConfig := tt.InitConfig()

	newConfig.TemplateFolder = "../../business/template/"

	newStore := store.NewStore(newLog, tt.InitDatabase(t, newConfig, newLog))
	newCore := core.NewCore(newLog, newConfig, newStore, tt.InitEmailConfig(t, newConfig))

	t.Log("Given the need to work with a send validation email.")
	{
		t.Logf("\tWhen a correct send validation email.")
		{
			assert.Nil(t, newCore.EmailSender.SendValidationEmail(
				faker.Internet().Email(), faker.RandomString(10),
				faker.RandomChoice([]string{"es", "en", "fr", "pt"})))
		}

		t.Logf("\tWhen a fail send validation email")
		{
			assert.Error(t, newCore.EmailSender.SendValidationEmail(
				faker.Internet().Email(), faker.RandomString(10), faker.RandomString(2)))

			assert.Error(t, newCore.EmailSender.SendValidationEmail(
				faker.RandomString(20), faker.RandomString(10), faker.RandomString(2)))
		}
	}
}
