package models

type JobStatus string

const (
	StatusQueued     JobStatus = "queued"
	StatusProcessing JobStatus = "processing"
	StatusCompleted  JobStatus = "completed"
	StatusFailed     JobStatus = "failed"
)

type JobItem struct {
	Id       string    `json:"id"`
	JobTitle string    `json:"job_title"`
	FileId   string    `json:"file_id"`
	Status   JobStatus `json:"status"`
}
