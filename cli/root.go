package cli

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "Go CV Evaluator with Gemini",
	Short: "A Go Application for CV Evalutor with Gemini LLM",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("error execute the command: %v", err)
	}
}
