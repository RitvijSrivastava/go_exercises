package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "task is a CLI application to manage your TODOs.",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
}
