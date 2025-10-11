package bootstrap

import (
	"context"
	"log"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/config"
	chromaclient "github.com/afrizalsebastian/ai-cv-evaluator-with-go/modules/chroma-client"
	geminiclient "github.com/afrizalsebastian/ai-cv-evaluator-with-go/modules/gemini-client"
	ingestdocument "github.com/afrizalsebastian/ai-cv-evaluator-with-go/modules/ingest-document"
	jobstore "github.com/afrizalsebastian/ai-cv-evaluator-with-go/modules/job-store"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/worker"
)

type Application struct {
	ENV      *config.Config
	Worker   worker.ICvEvaluatorWorker
	JobStore jobstore.IJobStore
}

func NewApp() *Application {
	ctx := context.Background()
	app := &Application{}

	if err := config.Init(); err != nil {
		log.Fatal("failed to initialize configuration")
	}

	app.ENV = config.Get()

	// Init Gemini Client
	geminiCient, err := geminiclient.NewGeminiAiCLient(ctx, app.ENV.GeminiApiKey, app.ENV.GeminiModel)
	if err != nil {
		log.Fatal("failed to init gemini client")
	}

	// Init chroma
	chromaClient, err := chromaclient.NewChromaClient(ctx, app.ENV.ChromaUrl)
	if err != nil {
		log.Fatalf("failed to init chroma client, %s", err.Error())
	}

	// Init ingestDocument
	ingesDocument := ingestdocument.NewIngestFile(chromaClient)

	// Assing Worker
	aiWorker := worker.NewCvEvaluatorWorker(geminiCient, chromaClient, ingesDocument, 5)
	app.Worker = aiWorker

	// Init Job Store
	jobStore := jobstore.NewJobStore()
	app.JobStore = jobStore

	return app
}
