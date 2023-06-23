package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hanzezhenalex/socks5/src/agent/client/operations"
	"github.com/hanzezhenalex/socks5/src/agent/models"
)

var (
	username string
	password string
)

var login = &cobra.Command{
	Use:   "login",
	Short: "login to the system",
	Long:  `login to the system`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if strings.ReplaceAll(username, " ", "") == "" || strings.ReplaceAll(password, " ", "") == "" {
			return fmt.Errorf("username/password can't be empty")
		}

		response, err := socksClient.Operations.PostLogin(&operations.PostLoginParams{
			User: &models.User{
				Username: username,
				Password: password,
			},
		})
		if err != nil {
			return fmt.Errorf("fail to login, %w", err)
		}

		file, err := os.OpenFile(tokenFilePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return fmt.Errorf("fail to persist token locally")
		}
		defer func() { _ = file.Close() }()

		_, err = file.WriteString(response.Payload)
		if err != nil {
			return fmt.Errorf("fail to persist token locally")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(login)

	login.Flags().StringVarP(&username, "username", "u", "", "username for user")
	login.Flags().StringVarP(&password, "password", "p", "", "password for user")
}
