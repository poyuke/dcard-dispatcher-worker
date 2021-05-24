package cmd

import (
	"dispatcher-worker/pkg/resetter"

	"github.com/spf13/cobra"
)

var resetterCmd = &cobra.Command{
	Use:   "resetter",
	Short: "reset long-awaited jobs",
	Long:  `reset long-awaited jobs`,
	Run: func(cmd *cobra.Command, args []string) {
		resetter.Execute(cmd, args)
	},
}

func init() {
	InitConfig()
	rootCmd.AddCommand(resetterCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// helloCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// helloCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
