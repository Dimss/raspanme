package cmd

import (
	"github.com/Dimss/raspanme/pkg/translate"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func init() {
	translateCmd.PersistentFlags().StringP("apikey", "k", "", "API key for Azure Translate service")
	translateCmd.PersistentFlags().StringP("endpoint", "e", "", "Endpoint for Azure Translate service")
	viper.BindPFlag("apikey", translateCmd.PersistentFlags().Lookup("apikey"))
	viper.BindPFlag("endpoint", translateCmd.PersistentFlags().Lookup("endpoint"))

	rootCmd.AddCommand(translateCmd)
}

var translateCmd = &cobra.Command{
	Use:   "translate",
	Short: "translate raspan db",
	Run: func(cmd *cobra.Command, args []string) {

		zap.S().Info("Translating Raspan DB")
		t := translate.NewTranslator(viper.GetString("apikey"), viper.GetString("endpoint"))
		t.Run()
	},
}
