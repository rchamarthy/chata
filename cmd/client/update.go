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

func updateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "updates a user",
		Long:  "updates a user's information",
		RunE:  UpdateUser,
		Args:  cobra.ExactArgs(1),
	}

	cmd.Flags().StringP("name", "n", "", "name update")
	cmd.Flags().BoolP("update-key", "k", false, "update the key")
	cmd.Flags().BoolP("add-admin", "r", false, "update admin role")

	return cmd
}

func UpdateUser(cmd *cobra.Command, args []string) error {
	serverAddress, err := cmd.Flags().GetString("server")
	if err != nil {
		return err
	}

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}

	uKey, err := cmd.Flags().GetBool("update-key")
	if err != nil {
		return err
	}

	admin, err := cmd.Flags().GetBool("add-admin")
	if err != nil {
		return err
	}

	user, key, err := makeUser(args[0], name, uKey, admin)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/users/%s", serverAddress, user.ID)
	client := resty.New()
	r, err := client.R().SetBody(user).SetResult(user).Post(url)
	if err != nil {
		return err
	}

	if r.StatusCode() != http.StatusOK {
		return fmt.Errorf("error registering user:\n %s", string(r.Body()))
	}

	// Save user profile with private key
	if uKey {
		user.Key = key
	}

	if e := saveUser(user); e != nil {
		return e
	}

	fmt.Printf("user with id: %s is updated\n", user.ID)
	user.Key = user.Key.Public() // avoid printing private key
	return PrintYaml(user)
}

func makeUser(id string, name string, uKey bool, admin bool) (*auth.User, *auth.Identity, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, nil, err
	}

	user, err := auth.LoadUser(path.Join(home, ".chata", id))
	if err != nil {
		return nil, nil, err
	}

	var key *auth.Identity
	if uKey {
		key = auth.GenerateIdentity()
		user.Key = key.Public()
	} else {
		user.Key = user.Key.Public()
	}

	if admin {
		user.Roles.Add(auth.ADMIN)
	}

	if name != "" {
		user.Name = name
	}

	return user, key, nil
}
