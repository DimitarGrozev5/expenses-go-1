package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "admin",
	Short: "Admin is used to manage the Expense app",
	Long:  `Admin is used to manage the Expense app`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running and exiting")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
