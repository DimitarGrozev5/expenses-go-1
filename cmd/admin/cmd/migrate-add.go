package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	migrateCmd.AddCommand(migrateAddCmd)
	migrateAddCmd.Flags().StringVarP(&migrationsPath, "migration-path", "p", "./migrations/", "Source directory for migrations")
}

var migrateAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add boilerplate for migration",
	Long:  `Add boilerplate for migration`,
	Run: func(cmd *cobra.Command, args []string) {

		// Get last migration
		i := 0
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

		// Create files
		up := `
/*
 * Migration Code Goes Here
 */
 
 
/*
 * Set user version
 */
PRAGMA user_version = ` + fmt.Sprint(i) + `;`
		down := `
/*
 * Disable foreign key constraints just in case
 */
PRAGMA foreign_keys = OFF;

/*
 * Migration Code Goes Here
 */

/*
 * Enable foreign key constraints
 */
PRAGMA foreign_keys = ON;
 
/*
 * Set user version
 */
PRAGMA user_version = ` + fmt.Sprint(i-1) + `;`

		// Write files
		err := os.WriteFile(migrationsPath+fmt.Sprintf("ctrl-%d-up.sql", i), []byte(up), 0644)
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(migrationsPath+fmt.Sprintf("ctrl-%d-down.sql", i), []byte(down), 0644)
		if err != nil {
			log.Fatal(err)
		}
	},
}
