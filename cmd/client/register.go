package main

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/go-resty/resty/v2"
	"github.com/rchamarthy/chata/auth"
	"github.com/spf13/cobra"
)

func registerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "register <id> <name>",
		Short: "register a new user",
		Long:  "register a new user with chata server",
		RunE:  RegisterUser,
		Args:  cobra.ExactArgs(2),
	}
}

func RegisterUser(cmd *cobra.Command, args []string) error {
	serverAddress, err := cmd.Flags().GetString("server")
	if err != nil {
		return err
	}

	user := auth.NewUser(args[1], args[0])

	if e := saveUser(user); e != nil {
		return e
	}

	user.Key = user.Key.Public()

	url := fmt.Sprintf("%s/users/%s", serverAddress, user.ID)
	client := resty.New()
	r, err := client.R().SetBody(user).Put(url)
	if err != nil {
		return err
	}

	if r.StatusCode() != http.StatusCreated {
		return fmt.Errorf("error registering user:\n %s", string(r.Body()))
	}

	fmt.Printf("user: %s id: %s is registered\n", user.Name, user.ID)
	return nil
}

func saveUser(user *auth.User) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	chatDir := path.Join(home, ".chata")
	if e := os.MkdirAll(chatDir, 0700); e != nil {
		return e
	}

	return user.SaveUser(chatDir)
}
