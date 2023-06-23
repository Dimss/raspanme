package cmd

import (
	"github.com/Dimss/raspanme/pkg/parser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func init() {
	parseCmd.PersistentFlags().StringP("file", "f", "/Users/dkartsev/Downloads/raspan.xlsx", "path to raspan xlsx file")
	viper.BindPFlag("file", parseCmd.PersistentFlags().Lookup("file"))

	rootCmd.AddCommand(parseCmd)
}

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "parse and load raspan xlsx file into the db",
	Run: func(cmd *cobra.Command, args []string) {

		zap.S().Info("Parsing xlsx file into the DB")
		p := parser.NewParser(viper.GetString("file"))
		p.LoadCategory()
		p.LoadQuestionsAndAnswers()

	},
}
