package dao

type CvEvaluatorJob struct {
	Id              int    `gorm:"column:id;primaryKey;autoIncrement"`
	FileId          string `gorm:"column:file_id;type:varchar(50)"`
	JobId           string `gorm:"column:job_id;type:varchar(50)"`
	JobTitle        string `gorm:"column:job_title;type:text"`
	Status          string `gorm:"column:status;type:enum('queued', 'processing', 'completed', 'failed')"`
	CvMatchRate     string `gorm:"column:cv_match_rate;type:varchar(10)"`
	CvFeedback      string `gorm:"column:cv_feedback;type:text"`
	ProjectScore    string `gorm:"column:project_score;type:varchar(10)"`
	ProjectFeedback string `gorm:"column:project_feedback;type:text"`
	OverallSummary  string `gorm:"column:overall_summary;type:text"`
}

func (CvEvaluatorJob) TableName() string { return "cv_evaluator_job" }
