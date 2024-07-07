package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/mail"
	"os"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func init() {
	usersCmd.AddCommand(usersAddCmd)
}

var usersAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new user",
	Long:  `Add new user`,
	Run: func(cmd *cobra.Command, args []string) {

		// Get DB Nodes
		nodes, err := Repo.CtrlDB.GetActiveNodes()
		if err != nil {
			log.Fatal(err)
		}

		// There must be nodes
		if len(nodes) == 0 {
			fmt.Printf("There are no active DB Nodes. Start a DB Node before adding users")
			return
		}

		// Get min user db version
		minVer, err := Repo.CtrlDB.GetMinUserVersion()
		if err != nil {
			log.Fatal(err)
		}

		// Get max user db version
		maxVer, err := Repo.CtrlDB.GetMaxUserVersion()
		if err != nil {
			log.Fatal(err)
		}

		// Get max migration
		i := 0
		for {

			// Update count
			i++

			// Try to get file
			_, err := os.ReadFile(migrationsPath + fmt.Sprintf("user-%d-up.sql", i))

			// Exit if migration not found
			if os.IsNotExist(err) {
				break
			}

			// Panic if other error
			if err != nil {
				log.Fatal(err)
			}
		}
		i = i - 1

		// Get user email
		emailPrompt := promptui.Prompt{
			Label:    "Enter user email",
			Validate: isEmail,
		}
		email, err := emailPrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		// Get user password (yeah I know it's stupid)
		passPrompt := promptui.Prompt{
			Label:    "Enter user password",
			Validate: validPass,
		}
		pass, err := passPrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		// Get user version
		dbVersionPrompot := promptui.Prompt{
			Label:    fmt.Sprintf("Enter db version (%d:%d; newest %d)", minVer, maxVer, i),
			Validate: validDBVersion,
		}
		dbVersionString, err := dbVersionPrompot.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		// Parse input
		parsedMail, _ := mail.ParseAddress(email)
		userEmail := parsedMail.Address
		dbVersion, _ := strconv.ParseInt(dbVersionString, 10, 64)

		// Add user
		err = Repo.CtrlDB.AddNewUser(userEmail, pass, dbVersion)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("User added")

		fmt.Printf("\n\n")
	},
}

// Is email
func isEmail(input string) error {
	_, err := mail.ParseAddress(input)
	if err != nil {
		return errors.New("invalid email")
	}
	return nil
}

// Valid password
func validPass(input string) error {
	if len(input) < 3 {
		return errors.New("invalid password")
	}
	return nil
}

// Valid password
func validDBVersion(input string) error {
	_, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return errors.New("invalid number")
	}
	return nil
}
