package bootstrap

import (
	"log"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/config"
)

type Application struct {
	ENV *config.Config
}

func NewApp() *Application {
	app := &Application{}

	if err := config.Init(); err != nil {
		log.Fatal("failed to initialize configuration")
	}

	app.ENV = config.Get()

	return app
}
