package cmd

import (
	"fmt"
	"github.com/hanzezhenalex/socks5/src/agent/client"
	"net"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	ip, port    string
	socksClient *client.SocksAgentAPIV1
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "socks-ctl",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if net.ParseIP(ip) == nil {
			return fmt.Errorf("illeagal ip: %s", ip)
		}
		if err := checkPort(port); err != nil {
			return err
		}
		cfg := &client.TransportConfig{}
		socksClient = client.NewHTTPClientWithConfig(nil, cfg)
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(func() {})

	rootCmd.PersistentFlags().StringVar(&ip, "addr", "127.0.0.1", "socks agent control server ip")
	rootCmd.PersistentFlags().StringVar(&port, "port", "8090", "socks agent control server port")
}

func checkPort(port string) error {
	p, err := strconv.Atoi(port)
	if err != nil {
		return err
	}
	if p < 1000 || p > 65535 {
		return fmt.Errorf("port out of range")
	}
	return nil
}
