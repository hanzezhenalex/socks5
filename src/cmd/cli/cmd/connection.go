package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
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
		response, err := socksClient.Operations.GetV1ConnectionList(nil)
		if err != nil {
			return fmt.Errorf("fail to make call to agent control server, %w", err)
		}
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)

		t.AppendHeader(table.Row{"uuid", "source", "target"})

		var rows []table.Row
		for _, payload := range response.Payload {
			rows = append(rows, table.Row{
				payload.UUID, payload.Source, payload.Target,
			})
		}
		t.AppendRows(rows)

		t.AppendSeparator()
		t.AppendFooter(table.Row{"", "", "Total", len(response.Payload)})

		t.Render()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(connection)

	connection.AddCommand(list)
}
