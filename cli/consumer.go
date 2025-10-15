package cli

import (
	"context"
	"log"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/bootstrap"
	"github.com/spf13/cobra"
)

var topic string

func init() {
	consumerCommand.Flags().StringVar(&topic, "topic", "", "kafka consumer topic")
	rootCmd.AddCommand(consumerCommand)
}

var consumerCommand = &cobra.Command{
	Use:   "consumer",
	Short: "Start consumer for Go CV Evaluator",
	PreRun: func(cmd *cobra.Command, args []string) {
		if topic == "" {
			log.Fatal("--topic is required")
		}
		app := bootstrap.NewApp()
		ctx := context.WithValue(cmd.Context(), appKey, app)
		cmd.SetContext(ctx)
	},
	Run: func(cmd *cobra.Command, args []string) {
		app := cmd.Context().Value(appKey).(*bootstrap.Application)
		startConsumer(app)
	},
}

func startConsumer(_ *bootstrap.Application) {
	log.Println("Consumer Running with topic", topic)
}
