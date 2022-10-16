// Package responses contains core return API.
package responses

type SuccessHealth struct {
	Status string `json:"status" example:"OK"`
	Health Health `json:"health"`
}

type Health struct {
	BuildVersion string `json:"buildVersion" example:"localhost"`
}
