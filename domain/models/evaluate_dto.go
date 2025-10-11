package models

type EvaluateRequest struct {
	JobTitle string `json:"job_title" validate:"required"`
	FileId   string `json:"file_id" validate:"required"`
}

type EvaluateResponse struct {
	JobId  string `json:"job_id"`
	Status string `json:"status"`
}
