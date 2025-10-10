package cli

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/bootstrap"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/handlers"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/infrastructure/middleware"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/internal/generated"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

type appKeyType struct{}

var appKey = appKeyType{}

const apiPrefix = "/api/v1"

func init() {
	rootCmd.AddCommand(serveCommand)
}

var serveCommand = &cobra.Command{
	Use:   "serve",
	Short: "Start The HTTP server Go Evaluator",
	PreRun: func(cmd *cobra.Command, args []string) {
		app := bootstrap.NewApp()
		ctx := context.WithValue(cmd.Context(), appKey, app)
		cmd.SetContext(ctx)
	},
	Run: func(cmd *cobra.Command, args []string) {
		app := cmd.Context().Value(appKey).(*bootstrap.Application)
		startHTTPServer(app)

	},
}

func startHTTPServer(app *bootstrap.Application) {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	mainRouter := mux.NewRouter()
	mainRouter.Use(middleware.RecoveryMiddleware())
	mainRouter.Use(middleware.CORSMiddleware)

	// Not found handler
	mainRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("404 --- route not found")
		http.Error(w, "404 - route not found", http.StatusNotFound)
	})

	apiV1Router := mainRouter.PathPrefix(apiPrefix).Subrouter()

	apiServer, err := handlers.NewServer(app)
	if err != nil {
		log.Fatal("failed to init server")
		os.Exit(1)
	}

	apiV1Router.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Service is ready"))
	}).Methods("GET")

	oApiHandlers := generated.HandlerWithOptions(apiServer, generated.GorillaServerOptions{
		BaseRouter: apiV1Router,
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "invalid request"}`))
		},
	})

	mainRouter.PathPrefix(apiPrefix).Handler(registerCommonMiddleware(app, oApiHandlers))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", app.ENV.AppPort),
		Handler: mainRouter,
	}

	// Gracefull shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// start server
	go func() {
		log.Printf("🚀 Server is running on port %v\n", app.ENV.AppPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln("server failed to start")
			os.Exit(1)
		}
	}()

	// wait for signal
	_ = <-signalChan
	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Println("Server forced to shutdown")
	}

	log.Println("✅ Server exited gracefully")
}

func registerCommonMiddleware(_ *bootstrap.Application, handler http.Handler) http.Handler {
	middleware := []middleware.MiddlewareFunc{
		middleware.RequestTracing(),
	}

	for _, m := range middleware {
		handler = m(handler)
	}
	return handler
}
