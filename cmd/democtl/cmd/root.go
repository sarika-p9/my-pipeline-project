package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "democtl",
	Short: "CLI for managing the pipeline system",
	Long:  `democtl is a command-line tool for interacting with the distributed manufacturing pipeline system.`,
}

func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	return err
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(pipelineCmd)

}
