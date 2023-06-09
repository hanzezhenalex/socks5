package cmd

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/hanzezhenalex/socks5/src"
	"github.com/hanzezhenalex/socks5/src/agent/client"

	httpTransport "github.com/go-openapi/runtime/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	ip, port    string
	socksClient *client.SocksAgentAPIV1
	tokenC      *tokenCollector
	debug       bool
)

const (
	tokenFilePathWindows = ""
	tokenFilePathLinux   = "/tmp/socks_token"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "socks-ctl",
	Short: "Cli for socks5 agent",
	Long: `Socks-ctl is a cli tool for socks5 server. Users can check the status of agent/server
And also control the behaviors of socks server.
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		_tokenC, err := NewTokenCollector()
		if err != nil {
			logrus.Errorf("fail to create token collector, %s", err.Error())
			return fmt.Errorf("fail to read token")
		}

		if cmd.Use == "login" || cmd.Use == "logout" {

		} else {
			if err := _tokenC.read(); err != nil {
				logrus.Errorf("fail to read token, %s", err.Error())
				return fmt.Errorf("fail to read token")
			}
		}

		tokenC = _tokenC

		if net.ParseIP(ip) == nil {
			return fmt.Errorf("illeagal ip: %s", ip)
		}
		if err := checkPort(port); err != nil {
			return err
		}
		cfg := &client.TransportConfig{
			Host:    fmt.Sprintf("%s:%s", ip, port),
			Schemes: []string{"https"},
		}
		runtime := httpTransport.New(cfg.Host, cfg.BasePath, cfg.Schemes)
		dialer := &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}
		defaultTransport := &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialer.DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		runtime.Transport = defaultTransport
		socksClient = client.New(runtime, nil)

		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		} else {
			logrus.SetOutput(src.BlackHole{})
		}
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
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "run in debug mode")
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
