package dto

type WebResponse struct {
	Status  int                    `json:"status" `
	Data    map[string]interface{} `json:"data" `
	Message string                 `json:"message" `
}
