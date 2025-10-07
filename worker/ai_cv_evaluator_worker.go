package worker

import (
	"context"
	"fmt"
	"log"
	"path"
	"sync"
	"time"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/application/services"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/models"
	chromaclient "github.com/afrizalsebastian/ai-cv-evaluator-with-go/modules/chroma-client"
	geminiclient "github.com/afrizalsebastian/ai-cv-evaluator-with-go/modules/gemini-client"
)

type ICvEvaluatorWorker interface {
	Start(ctx context.Context)
	Enqueue(job *models.JobItem) error
	Stop(ctx context.Context)
}

type cvEvaluatorWorker struct {
	fileStore services.ILocalFileStore
	gemini    geminiclient.IGeminiClient
	chroma    chromaclient.IChromaClient
	ingest    services.IIngestFile

	queue chan *models.JobItem
	mu    sync.RWMutex
	wg    sync.WaitGroup

	maxWorkers int
}

func NewCvEvaluatorWorker(
	fileStore services.ILocalFileStore,
	gemini geminiclient.IGeminiClient,
	chroma chromaclient.IChromaClient,
	ingest services.IIngestFile,
	maxWorkers int,
) ICvEvaluatorWorker {
	if maxWorkers <= 0 {
		maxWorkers = 5
	}

	return &cvEvaluatorWorker{
		fileStore:  fileStore,
		gemini:     gemini,
		chroma:     chroma,
		ingest:     ingest,
		maxWorkers: maxWorkers,
		queue:      make(chan *models.JobItem, 100),
	}
}

func (w *cvEvaluatorWorker) Start(ctx context.Context) {
	log.Printf("Worker start in background")
	for i := 0; i < w.maxWorkers; i++ {
		w.wg.Add(1)
		go w.worker(ctx)
	}
}

func (w *cvEvaluatorWorker) Enqueue(job *models.JobItem) error {
	select {
	case w.queue <- job:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("queue full, timeout enqueueing job")
	}
}

func (w *cvEvaluatorWorker) Stop(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	close(w.queue)
	w.wg.Wait()
}

func (w *cvEvaluatorWorker) worker(ctx context.Context) {
	defer w.wg.Done()
	for {
		select {
		case job, ok := <-w.queue:
			if !ok {
				return
			}
			w.processJob(ctx, job)
		case <-ctx.Done():
			return
		}
	}
}

func (w *cvEvaluatorWorker) processJob(ctx context.Context, job *models.JobItem) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	job.Status = models.StatusProcessing

	// extract text from file
	extractCv, err := w.ingest.ExtractTextFromPdf(path.Join("uploaded-file", job.FileId, "cv_file.pdf"))
	if err != nil {
		w.jobFailToProcess(job, err)
		return
	}

	extractReport, err := w.ingest.ExtractTextFromPdf(path.Join("uploaded-file", job.FileId, "report_file.pdf"))
	if err != nil {
		w.jobFailToProcess(job, err)
		return
	}

	//TODO: process job
	job.Status = models.StatusCompleted
}

func (w *cvEvaluatorWorker) jobFailToProcess(job *models.JobItem, err error) {
	fmt.Printf("job with id %s failed to prcess: %s\n", job.Id, err.Error())
	job.Status = models.StatusFailed
}
