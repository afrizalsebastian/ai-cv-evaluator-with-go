package repository

import (
	"context"
	"log"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/bootstrap"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/domain/models/dao"
	"gorm.io/gorm"
)

type ICvEvaluatorJobRepository interface {
	CreateJobItem(ctx context.Context, job *dao.CvEvaluatorJob) error
	GetByJobId(ctx context.Context, jobId string) (*dao.CvEvaluatorJob, error)
}

type cvEvaluatorJobRepository struct {
	db *gorm.DB
}

func NewCvEvaluatorJobRepository(app *bootstrap.Application) ICvEvaluatorJobRepository {
	return &cvEvaluatorJobRepository{
		db: app.DB,
	}
}

func (c *cvEvaluatorJobRepository) CreateJobItem(ctx context.Context, job *dao.CvEvaluatorJob) error {
	if err := c.db.WithContext(ctx).Create(&job).Error; err != nil {
		log.Println("failed to create job")
		return err
	}
	return nil
}

func (c *cvEvaluatorJobRepository) GetByJobId(ctx context.Context, jobId string) (*dao.CvEvaluatorJob, error) {
	var jobItem dao.CvEvaluatorJob
	if err := c.db.WithContext(ctx).Model(&dao.CvEvaluatorJob{}).Where("job_id = ?", jobId).First(&jobItem).Error; err != nil {
		log.Println("failed to get job by job id")
		return nil, err
	}

	return &jobItem, nil
}
