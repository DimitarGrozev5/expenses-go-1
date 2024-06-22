package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dbnodeCmd)
}

var dbnodeCmd = &cobra.Command{
	Use:   "dbnodes",
	Short: "View all db nodes",
	Long:  `View all db nodes`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get nodes
		nodes, err := Repo.CtrlDB.GetNodes()
		if err != nil {
			log.Fatalf("Can't get DB Nodes from Controlle DB: %s", err)
		}

		if len(nodes) == 0 {
			fmt.Println("You don't have db nodes yet.")
			return
		}

		fmt.Println("-----------------------------------------")
		fmt.Println("|\tID\t|\tAddress\t\t|")
		fmt.Println("-----------------------------------------")
		for _, node := range nodes {
			fmt.Printf("|\t%d\t|\t%s\t\t|\n", node.ID, node.RemoteAddress)
		}
		fmt.Println("-----------------------------------------")
	},
}
