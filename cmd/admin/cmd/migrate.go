package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "View all db nodes",
	Long:  `View all db nodes`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("status - get current controller db status")
		fmt.Println("to [number] - migrate to migration")
		fmt.Println("add - add boilerplate for migration")
	},
}
