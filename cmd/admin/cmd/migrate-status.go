package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	migrateCmd.AddCommand(migrateStatusCmd)
	migrateStatusCmd.Flags().StringVarP(&migrationsPath, "migration-path", "p", "./migrations/", "Source directory for migrations")
}

var migrationsPath string

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get Controller DB migration status",
	Long:  `Get Controller DB migration status`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get user version
		userVersion, err := Repo.CtrlDB.GetVersion()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\nController DB Version: %d\n", userVersion)

		// Get last migration
		i := userVersion
		for {

			// Update count
			i++

			// Try to get file
			_, err := os.ReadFile(migrationsPath + fmt.Sprintf("ctrl-%d-up.sql", i))

			// Exit if migration not found
			if os.IsNotExist(err) {
				break
			}

			// Panic if other error
			if err != nil {
				log.Fatal(err)
			}
		}

		fmt.Printf("Latest migration: %d\n\n", i-1)
	},
}
