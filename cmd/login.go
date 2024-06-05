package cmd

import (
	"clilogin/login"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	var loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Login into application",
		Run: func(cmd *cobra.Command, args []string) {
			stop := make(chan login.CallbackResponse)
			err := login.InitGoogleAuth(stop)
			if err != nil {
				fmt.Printf("Authorisation error: %s\n", err)
				return
			}

			resp := <-stop
			if resp.Error != nil {
				fmt.Printf("Authorisation error: %s\n", resp.Error)
				return
			}

			fmt.Printf("Login is successful for user %s, id %s\n", resp.User.Email, resp.User.UserID)
		},
	}

	rootCmd.AddCommand(loginCmd)
}
