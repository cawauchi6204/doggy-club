package utils

import (
	"encoding/json"
)

// ToJSON converts a struct to JSON string
func ToJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

// FromJSON converts JSON string to struct
func FromJSON(jsonStr string, v interface{}) error {
	return json.Unmarshal([]byte(jsonStr), v)
}