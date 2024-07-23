package auth

import (
	"errors"
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type User struct {
	ID    string    `json:"id"    yaml:"id"`
	Name  string    `json:"name"  yaml:"name"`
	Key   *Identity `json:"key"   yaml:"key"`
	Roles *Roles    `json:"roles" yaml:"roles"`
}

func NewUser(name string, id string, roles ...Role) *User {
	r := NewRoles(roles...)
	r.Add(SELF) // Everyone has their own role!

	return &User{
		Name:  name,
		ID:    id,
		Key:   GenerateIdentity(),
		Roles: r,
	}
}

func (user *User) Validate() error {
	if user.ID == "" {
		return errors.New("id cannot be empty")
	}

	if user.Name == "" {
		return errors.New("user cannot be empty")
	}

	if user.Key == nil {
		return errors.New("user public key cannot be empty")
	}

	return nil
}

func LoadUser(userFile string) (*User, error) {
	b, err := os.ReadFile(userFile)
	if err != nil {
		return nil, err
	}

	user := &User{ID: "", Name: "", Key: nil, Roles: NewRoles()}
	if e := yaml.Unmarshal(b, user); e != nil {
		return nil, e
	}

	return user, nil
}

func (user *User) SaveUser(usersDir string) error {
	if e := user.Validate(); e != nil {
		return e
	}

	userFile := path.Join(usersDir, user.ID)
	b, e := yaml.Marshal(user)
	if e != nil {
		return e
	}

	e = os.WriteFile(userFile, b, 0600)
	if e != nil {
		e = fmt.Errorf("error writing user file: %w", e)
	}

	return e
}
