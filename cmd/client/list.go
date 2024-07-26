package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/rchamarthy/chata/auth"
	"github.com/rchamarthy/chata/store"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list users or a user",
		Long:  "list a user's information or all users if no user is specified",
		RunE:  ListUsers,
		Args:  cobra.MaximumNArgs(1),
	}
}

func ListUsers(cmd *cobra.Command, args []string) error {
	server, err := cmd.Flags().GetString("server")
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return listAllUsers(server)
	}

	user := &auth.User{}
	url := fmt.Sprintf("%s/users/%s", server, args[0])
	client := resty.New()
	r, err := client.R().SetResult(user).Get(url)
	if err != nil {
		return err
	}

	if r.StatusCode() != http.StatusOK {
		return fmt.Errorf("unable to list user: %s", args[0])
	}

	return PrintYaml(user)
}

func listAllUsers(server string) error {
	users := store.Users{}

	url := server + "/users"
	client := resty.New()
	r, err := client.R().SetResult(&users).Get(url)
	if err != nil {
		return err
	}

	if r.StatusCode() != http.StatusOK {
		return errors.New("error listing users")
	}

	for _, user := range users {
		if e := PrintYaml(user); e != nil {
			return e
		}
		fmt.Println()
	}

	return nil
}

func PrintYaml(obj any) error {
	out, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}

	fmt.Println(string(out))

	return nil
}
