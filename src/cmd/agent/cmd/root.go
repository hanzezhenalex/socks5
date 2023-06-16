package cmd

import (
	"fmt"
	"net"
	"strconv"

	"github.com/hanzezhenalex/socks5/src/agent"

	"github.com/spf13/cobra"
)

var (
	config       agent.Config
	mode         string
	ip           string
	socksSrvPort string
	controlPort  string
	commands     []string
	auths        []string
)

var rootCmd = &cobra.Command{
	Use: "agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		switch mode {
		case agent.LocalMode:
			config.Mode = agent.LocalMode
		case agent.ClusterMode:
			config.Mode = agent.ClusterMode
		default:
			return fmt.Errorf("unknown mode: %s", mode)
		}

		if res := net.ParseIP(ip); res == nil {
			return fmt.Errorf("illegal socks server addr, %s:%s", ip, socksSrvPort)
		}
		config.Socks5Config.IP = ip

		if err := checkPort(socksSrvPort); err != nil {
			return fmt.Errorf("illegal socksSrvPort, %s", err.Error())
		}
		config.Socks5Config.Port = socksSrvPort

		if err := checkPort(controlPort); err != nil {
			return fmt.Errorf("illegal control server Port, %s", err.Error())
		}
		config.ControlServerPort = controlPort

		config.Socks5Config.Command = commands
		config.Socks5Config.Auth = auths

		agent := agent.NewAgent(config)
		return agent.Run()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(func() {})

	rootCmd.Flags().StringVarP(&mode, "mode", "m", "local", "agent mode, local or cluster")
	rootCmd.Flags().StringVar(&ip, "ip", "0.0.0.0", "socks server ip")
	rootCmd.Flags().StringVar(&socksSrvPort, "socks-port", "1080", "socks server socksSrvPort")
	rootCmd.Flags().StringVar(&controlPort, "control-port", "8090", "agent control server socksSrvPort")
	rootCmd.Flags().StringSliceVarP(&commands, "commands", "c", []string{"connect"},
		"commands for socks server, supported=[connect,]")
	rootCmd.Flags().StringSliceVarP(&auths, "auths", "a", []string{"noAuth"},
		"auth methods for socks server, supported=[noAuth,]")
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
