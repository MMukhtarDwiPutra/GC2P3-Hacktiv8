package dto

type WebResponse struct {
	Status  int                    `json:"status" `
	Data    map[string]interface{} `json:"data" `
	Message string                 `json:"message" `
}

// ErrorResponse represents a standardized error response structure.
type ErrResponse struct {
	Code    int    `json:"code"`              // HTTP status code
	Message string `json:"message"`           // Error message
	Details any    `json:"details,omitempty"` // Optional details for the error
}
