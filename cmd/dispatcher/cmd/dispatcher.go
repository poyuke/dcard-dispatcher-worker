package cmd

import (
	"dispatcher-worker/pkg/dispatcher"

	"github.com/spf13/cobra"
)

var dispatcherCmd = &cobra.Command{
	Use:   "dispatcher",
	Short: "SHA-1 file content",
	Long:  `Job worker use SHA-a hash file content`,
	Run: func(cmd *cobra.Command, args []string) {
		dispatcher.Execute(cmd, args)
	},
}

func init() {
	InitConfig()
	rootCmd.AddCommand(dispatcherCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// helloCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// helloCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
