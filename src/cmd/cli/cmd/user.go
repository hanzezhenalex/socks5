package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hanzezhenalex/socks5/src/agent/client/operations"
	"github.com/hanzezhenalex/socks5/src/agent/models"
)

var roles []string

var user = &cobra.Command{
	Use:   "user",
	Short: "user related operations",
	Long:  `create/delete/restrict etc`,
}

var create = &cobra.Command{
	Use:   "create",
	Short: "create new user",
	Long:  `only admin user could create new user`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := operations.NewPostV1AuthUserCreateParams()
		params.Authorization = &tokenC.token
		params.User = &models.User{
			Username: username,
			Password: password,
			Roles:    roles,
		}
		_, err := socksClient.Operations.PostV1AuthUserCreate(params)
		if err != nil {
			return fmt.Errorf("fail to create user %s, %w", username, err)
		}
		fmt.Printf("successfully create new user %s", username)
		fmt.Println()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(user)

	user.AddCommand(create)
	create.Flags().StringVarP(&username, "username", "u", "", "username for new user")
	create.Flags().StringVarP(&password, "password", "p", "", "password for new user")
	create.Flags().StringSliceP("roles", "r", roles, "roles for new user")
}
