package datatypes

import (
	"database/sql/driver"
	"encoding/json"
)

type JSONB map[string]interface{}

func (j *JSONB) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), j)
}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}
