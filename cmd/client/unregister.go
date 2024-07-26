package main

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

func unregisterCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unregister register <id>",
		Short: "unregister a user",
		Long:  "unregister a user with chata server",
		RunE:  UnregisterUser,
		Args:  cobra.ExactArgs(1),
	}
}

func UnregisterUser(cmd *cobra.Command, args []string) error {
	serverAddress, err := cmd.Flags().GetString("server")
	if err != nil {
		return err
	}

	id := args[0]

	url := fmt.Sprintf("%s/users/%s", serverAddress, id)
	client := resty.New()
	r, err := client.R().Delete(url)
	if err != nil {
		return err
	}

	if r.StatusCode() != http.StatusOK {
		return fmt.Errorf("error deleting user:\n %s", string(r.Body()))
	}

	fmt.Printf("user id: %s is deleted\n", id)
	return nil
}
