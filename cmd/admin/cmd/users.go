package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(usersCmd)
}

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Control Users",
	Long:  `Control Users`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add\t\t\tadd new user")
		fmt.Print("\n\n")
	},
}
