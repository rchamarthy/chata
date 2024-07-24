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

	c.rootCmd.PersistentFlags().StringP("server", "s", ":8888", "server address")

	c.rootCmd.AddCommand(registerCmd())
	return c
}

func (c *Client) Execute() {
	if err := c.rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
