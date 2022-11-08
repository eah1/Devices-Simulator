package email

var subject = map[string]interface{}{
	"en": map[string]string{
		"account-validation": "Welcome",
	},
	"es": map[string]string{
		"account-validation": "Bienvenidos",
	},
	"pt": map[string]string{
		"account-validation": "Bem-vindo",
	},
	"fr": map[string]string{
		"account-validation": "Bienvenu",
	},
}

// Subject return type notification subject title.
// nolint: forcetypeassert
func Subject(template, language string) string {
	if _, ok := subject[language]; ok {
		if _, ok := subject[language].(map[string]string)[template]; ok {
			return subject[language].(map[string]string)[template]
		}
	}

	return ""
}
