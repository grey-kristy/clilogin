package cmd

import (
	"clilogin/login"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Login status",
		Run: func(cmd *cobra.Command, args []string) {
			user, err := login.ReadUser()
			if err != nil {
				fmt.Println("Application is not logged")
			} else {
				fmt.Printf("Application is logged (user %s, id %s)\n", user.Email, user.UserID)
			}
		},
	}
	rootCmd.AddCommand(statusCmd)
}
