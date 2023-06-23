package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
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
		params := operations.NewPostLoginParams()
		params.User = &models.User{
			Username: username,
			Password: password,
		}
		response, err := socksClient.Operations.PostLogin(params)
		if err != nil {
			return fmt.Errorf("fail to login, %w", err)
		}

		c, err := NewTokenCollector()
		if err != nil {
			return fmt.Errorf("fail to persist token, %w", err)
		}
		if err := c.set(response.Payload); err != nil {
			return err
		}
		fmt.Printf("successfully login as %s", username)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(login)

	login.Flags().StringVarP(&username, "username", "u", "", "username for user")
	login.Flags().StringVarP(&password, "password", "p", "", "password for user")
}

type tokenCollector struct {
	token string
	path  string
}

func (c *tokenCollector) read() error {
	_, err := os.Stat(c.path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("please login first")
		}
		return fmt.Errorf("unable to check login status, please re-login")
	}
	data, err := os.ReadFile(c.path)
	if err != nil {
		return fmt.Errorf("unable to check login status, please re-login")
	}
	c.token = string(data)
	return nil
}

func (c *tokenCollector) set(token string) error {
	file, err := os.OpenFile(c.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		logrus.Errorf("fail to open token file: err=%s", err.Error())
		return fmt.Errorf("fail to persist token locally")
	}
	defer func() { _ = file.Close() }()

	_, err = file.WriteString(token)
	if err != nil {
		logrus.Errorf("fail to write to token file: err=%s", err.Error())
		return fmt.Errorf("fail to persist token locally")
	}
	c.token = token
	return nil
}

func NewTokenCollector() (*tokenCollector, error) {
	c := &tokenCollector{}
	switch runtime.GOOS {
	case "windows":
		c.path = tokenFilePathWindows
	case "linux":
		c.path = tokenFilePathLinux
	default:
		return nil, fmt.Errorf("unsupported os: %s", runtime.GOOS)
	}

	return c, nil
}
