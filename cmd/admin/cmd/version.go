package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dbnodeCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Admin",
	Long:  `Print the version number of Admin`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Expenses Admin v0.1")
	},
}
