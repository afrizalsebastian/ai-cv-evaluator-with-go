package models

type EvaluateRequest struct {
	JobTitle string `json:"job_title"`
	FileId   string `json:"file_id"`
}
