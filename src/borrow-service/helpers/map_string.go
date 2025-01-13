package helpers

import (
	"borrow-service/dto"
	"encoding/json"
)

// UnmarshalJSONToWebResponse parses a JSON string into a WebResponse.
func UnmarshalJSONToWebResponse(stringJSON string) dto.WebResponse {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(stringJSON), &data); err != nil {
		return dto.WebResponse{
			Status:  500, // Indicate an internal server error
			Message: "Error unmarshalling string JSON to map[string]",
		}
	}

	// Extract relevant fields
	status := getInt(data["status"], 200)
	message := getString(data["message"], "Operation successful")
	dataContent := filterOutKeys(data, []string{"status", "message"})

	return dto.WebResponse{
		Status:  status,
		Message: message,
		Data:    dataContent,
	}
}

// Helper function to get an int from an interface
func getInt(value interface{}, defaultValue int) int {
	if v, ok := value.(float64); ok {
		return int(v)
	}
	return defaultValue
}

// Helper function to get a string from an interface
func getString(value interface{}, defaultValue string) string {
	if v, ok := value.(string); ok {
		return v
	}
	return defaultValue
}

// Helper function to filter out specific keys from a map
func filterOutKeys(original map[string]interface{}, keysToExclude []string) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range original {
		// Exclude keys in the keysToExclude list
		exclude := false
		for _, excludeKey := range keysToExclude {
			if key == excludeKey {
				exclude = true
				break
			}
		}
		if !exclude {
			result[key] = value
		}
	}
	return result
}
