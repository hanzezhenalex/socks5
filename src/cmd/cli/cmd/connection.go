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
	Short: "Cli tool for agent connection operations",
	Long: `List/Shutdown(not implement yet) connections agent is proxying now,
auth is not working, in the future, users can only op on his/her own connections`,
}

var list = &cobra.Command{
	Use:   "list",
	Short: "list the connections agent is proxying now",
	Long: `List connections agent is proxying now,
auth is not working, in the future, users can only op on his/her own connections.`,
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
		t.AppendFooter(table.Row{"", "Total", len(response.Payload)})

		t.Render()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(connection)

	connection.AddCommand(list)
}
