package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func init() {
	dbnodeCmd.AddCommand(dbnodeAddCmd)
}

var dbnodeAddCmd = &cobra.Command{
	Use:   "new",
	Short: "View all db nodes",
	Long:  `View all db nodes`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get nodes
		id, err := Repo.CtrlDB.NewNode()
		if err != nil {
			log.Fatalf("Can't add new DB Node: %s", err)
		}

		fmt.Printf("Inserted new DB Node with ID: %d", id)
	},
}
