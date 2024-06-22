package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	migrateCmd.AddCommand(migrateStepCmd)
	migrateStepCmd.Flags().StringVarP(&migrationsPath, "migration-path", "p", "./migrations/", "Source directory for migrations")
}

var migrateStepCmd = &cobra.Command{
	Use:   "to [step]",
	Short: "Migrate Controller DB to the specified level",
	Long:  `Migrate Controller DB to the specified level`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires one arg")
		}
		arg, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return errors.New("arg must be an integer")
		}
		if arg < 0 {
			return errors.New("arg must be a positive intiger")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		// Get target step
		target, _ := strconv.ParseInt(args[0], 10, 64)

		// Setup new context
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// Get user version
		userVersion, err := Repo.CtrlDB.GetVersion()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\nController Initial DB Version: %d\n", userVersion)

		// If target is the same as db version, exit
		if target == userVersion {
			fmt.Printf("Target matches current DB Version.\n\n")
			return
		}

		// Get from and to
		from := userVersion + 1
		to := target
		var add int64 = 1
		compare := func(a, b int64) bool {
			return a <= b
		}

		if from > to {
			from = userVersion
			to = target + 1
			add = -1
			compare = func(a, b int64) bool {
				return a >= b
			}
		}

		// Store migrations
		migrations := make([]string, 0)

		// Get all migrations relative to current db version
		for i := from; compare(i, to); i += add {

			// Get file path
			path := migrationsPath + fmt.Sprintf("ctrl-%d-up.sql", i)
			if userVersion > target {
				path = migrationsPath + fmt.Sprintf("ctrl-%d-down.sql", i)
			}

			fmt.Println(path)

			// Try to get up file
			file, err := os.ReadFile(path)

			// Exit if migration not found
			if os.IsNotExist(err) {
				break
			}

			// Panic if other error
			if err != nil {
				log.Fatal(err)
			}

			// Add file to migration
			migrations = append(migrations, string(file))
		}

		// Start transaction
		tx, err := Repo.CtrlDBConn.SQL.Begin()
		if err != nil {
			log.Fatal(err)
		}
		defer tx.Rollback()

		// Run migrations from current version up
		for _, migration := range migrations {

			// Run query
			_, err = tx.ExecContext(ctx, migration)
			if err != nil {
				log.Fatal(err)
			}
		}

		// Print migration message
		fmt.Printf("Migrations performed: %d\n", target)

		// Get user version
		userVersion, err = Repo.CtrlDB.GetVersion()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Controller DB Version: %d\n\n", userVersion)

		fmt.Printf("\nController Initial DB Version: %d\n", userVersion)

		// Commit migrations
		tx.Commit()
	},
}
