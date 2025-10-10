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
	Result   JobResult `json:"result"`
}

type JobResult struct {
	CvMatchRate     string `json:"cv_match_rate"`
	CvFeedback      string `json:"cv_feedback"`
	ProjectScore    string `json:"project_score"`
	ProjectFeedback string `json:"project_feedback"`
	OverallSummary  string `json:"overall_summary"`
}
