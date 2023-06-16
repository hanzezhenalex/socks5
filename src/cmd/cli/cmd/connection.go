package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var connection = &cobra.Command{
	Use:   "connection",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

var list = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := socksClient.Operations.GetV1ConnectionList(nil)
		if err != nil {
			return fmt.Errorf("fail to make call to agent control server, %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(connection)
	connection.AddCommand(list)
}
