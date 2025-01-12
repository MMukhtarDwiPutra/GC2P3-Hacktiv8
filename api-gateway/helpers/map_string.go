package helpers

import (
	"encoding/json"
)

func UnmarshalJSONToMap(stringJSON string) map[string]interface{} {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(stringJSON), &data); err != nil {
		return map[string]interface{}{
			"message": "Error unmarshalling string JSON to map[string]",
		}
	}

	return data
}
