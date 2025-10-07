package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/application/handlers"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/application/router"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/application/services"
	chromaclient "github.com/afrizalsebastian/ai-cv-evaluator-with-go/modules/chroma-client"
	geminiclient "github.com/afrizalsebastian/ai-cv-evaluator-with-go/modules/gemini-client"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/worker"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	geminiKey := os.Getenv("GEMINI_API_KEY")
	if geminiKey == "" {
		log.Fatal("GEMINI_API_KEY is required")
	}

	geminiModel := os.Getenv("GEMINI_MODEL")
	if geminiKey == "" {
		geminiModel = "gemini-2.5-flash"
	}

	chromaURL := os.Getenv("CHROMA_URL")
	if chromaURL == "" {
		chromaURL = "http://localhost:8000"
	}

	ctx := context.Background()

	// dependencies

	// modules
	geminiCient, err := geminiclient.NewGeminiAiCLient(ctx, geminiKey, geminiModel)
	if err != nil {
		log.Fatal("failed to init gemini client")
	}
	chromaClient, err := chromaclient.NewChromaClient(ctx, chromaURL)
	if err != nil {
		log.Fatalf("failed to init chroma client, %s", err.Error())
	}

	// services
	fileStore, err := services.NewLocalFileStore("./uploaded-file")
	if err != nil {
		log.Fatal("failed to init fileStore")
	}
	jobStore := services.NewJobStore()
	ingest := services.NewIngestFile(chromaClient)

	// worker
	aiWorker := worker.NewCvEvaluatorWorker(fileStore, geminiCient, chromaClient, 5)

	// handler
	uploadHandler := handlers.NewUploadHandler(fileStore)
	evaluateHandler := handlers.NewEvaluateHandler(aiWorker, jobStore)
	resultHandler := handlers.NewResultHandler(jobStore)
	chromaHandler := handlers.NewChromaHandler(ingest)

	// router
	r := router.NewRouter(uploadHandler, evaluateHandler, resultHandler, chromaHandler)

	// start worker
	go aiWorker.Start(ctx)

	// serve
	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      r.Router(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("🚀 Server is running on port %s", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	// gracefull shutdown
	<-stop
	log.Println("⚙️  Shutdown signal received, stopping server...")
	ctxShutDown, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	aiWorker.Stop(ctx)
	log.Println("✅ Server exited gracefully")
}
