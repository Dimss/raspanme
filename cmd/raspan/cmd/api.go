package cmd

import (
	"github.com/Dimss/raspanme/pkg/rest"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {

	rootCmd.AddCommand(apiCmd)
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "start API server",
	Run: func(cmd *cobra.Command, args []string) {

		zap.S().Info("Starting Rest API Server")
		s := rest.NewServer()
		s.Start()

	},
}
