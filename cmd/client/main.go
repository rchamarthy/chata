package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	NewClient().Execute()
}

type Client struct {
	rootCmd *cobra.Command
}

func NewClient() *Client {
	c := &Client{
		rootCmd: &cobra.Command{
			Use:   "chata",
			Short: "command line chat client",
			Long:  "A simple chat client for chata server",
		},
	}

	c.rootCmd.PersistentFlags().StringP("server", "s", "http://127.0.0.1:8888", "server address")

	userCmd := &cobra.Command{
		Use:   "user",
		Short: "chata user management",
		Long:  "manage the lifecycle of the users",
	}

	c.rootCmd.AddCommand(userCmd)

	userCmd.AddCommand(registerCmd())
	userCmd.AddCommand(unregisterCmd())
	userCmd.AddCommand(listCmd())
	userCmd.AddCommand(updateCmd())

	return c
}

func (c *Client) Execute() {
	if err := c.rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
