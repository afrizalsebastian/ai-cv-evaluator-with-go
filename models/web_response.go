package models

type WebResponse struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data,omitempty"`
}

func CreateWebResponse(message string, status int, data interface{}) *WebResponse {
	return &WebResponse{
		Message: message,
		Status:  status,
		Data:    data,
	}
}
