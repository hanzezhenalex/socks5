package main

import (
	"fmt"
	"net"
	"os"

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
	Use: "start",
	RunE: func(cmd *cobra.Command, args []string) error {
		switch mode {
		case agent.LocalMode:
			config.Mode = agent.LocalMode
		case agent.ClusterMode:
			config.Mode = agent.ClusterMode
		default:
			return fmt.Errorf("unknown mode: %s", mode)
		}

		if ip := net.ParseIP(fmt.Sprintf("%s:%s", ip, socksSrvPort)); ip == nil {
			return fmt.Errorf("illegal socks server addr, %s:%s", ip, socksSrvPort)
		}
		config.Socks5Config.IP = ip
		config.Socks5Config.Port = socksSrvPort

		if ip := net.ParseIP(fmt.Sprintf("%s:%s", ip, controlPort)); ip == nil {
			return fmt.Errorf("illegal agent control server addr, %s:%s", ip, controlPort)
		}
		config.ControlServerPort = controlPort

		config.Socks5Config.Command = commands
		config.Socks5Config.Auth = auths

		agent := agent.NewAgent(config)
		return agent.Run()
	},
}

func main() {
	rootCmd.Flags().StringVarP(&mode, "mode", "m", "local", "agent mode, local or cluster")
	rootCmd.Flags().StringVar(&ip, "ip", "0.0.0.0", "socks server ip")
	rootCmd.Flags().StringVar(&socksSrvPort, "socks-port", "1080", "socks server socksSrvPort")
	rootCmd.Flags().StringVar(&controlPort, "control-port", "8090", "agent control server socksSrvPort")
	rootCmd.Flags().StringSliceVarP(&commands, "commands", "c", []string{"connect"},
		"commands for socks server, supported=[connect,]")
	rootCmd.Flags().StringSliceVarP(&auths, "auths", "a", []string{"noAuth"},
		"auth methods for socks server, supported=[noAuth,]")

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("error exit, %s", err.Error())
		os.Exit(1)
	}
}
