package usecase

import (
	"device-simulator/business/core/models"
	"device-simulator/business/web/webmodels"
	"fmt"
)

// CreateEnvironment create environment case-use.
func (u *UseCase) CreateEnvironment(
	createEnvironment webmodels.CreateEnvironment, userID string,
) (webmodels.InformationEnvironment, error) {
	environment := models.CreateEnvironmentWebToEnvironment(createEnvironment)
	environment.UserID = userID

	if err := u.core.Environment.Create(&environment); err != nil {
		return webmodels.InformationEnvironment{}, fmt.Errorf("usecase.environment.Create: %w", err)
	}

	return models.EnvironmentModelToWeb(environment), nil
}
