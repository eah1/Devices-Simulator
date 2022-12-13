package models

import "time"

// DeviceConfig represents the structure we need for moving data.
type DeviceConfig struct {
	ID                 string                 `xorm:"pk not null 'id'"          json:"id"`
	Name               string                 `xorm:"name"                      json:"name"`
	Vars               map[string]interface{} `xorm:"jsonb vars"                json:"vars"`
	MetricsFixed       map[string]interface{} `xorm:"jsonb metrics_fixed"       json:"metricsFixed"`
	MetricsAccumulated map[string]interface{} `xorm:"jsonb metrics_accumulated" json:"metricsAccumulated"`
	TypeSend           string                 `xorm:"type_send"                 json:"typeSend"`
	Payload            string                 `xorm:"payload"                   json:"payload"`
	UserID             string                 `xorm:"user_id"                   json:"userId"`
	CreatedAt          time.Time              `xorm:"created"                   json:"-"`
	UpdatedAt          time.Time              `xorm:"updated"                   json:"-"`
}
