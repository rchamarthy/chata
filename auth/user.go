package auth

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/rchamarthy/chata"
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

type Users map[string]*User

type userError struct {
	u *User
	e error
}

func (users Users) Add(user *User) error {
	if e := user.Validate(); e != nil {
		return e
	}

	users[user.ID] = user
	return nil
}

func LoadUsers(ctx context.Context, userDir string) (Users, error) {
	userChan := make(chan userError, 16)

	go func(channel chan userError, d string) {
		waitGroup := &sync.WaitGroup{}
		err := filepath.WalkDir(d,
			func(p string, d fs.DirEntry, err error) error {
				waitGroup.Add(1)
				go func(p string, f fs.DirEntry) {
					defer waitGroup.Done()
					if err != nil {
						channel <- userError{nil, err}
						return
					}

					if f.Type().IsRegular() {
						u, e := LoadUser(p)
						channel <- userError{u, e}
					}
				}(p, d)

				return err
			})

		if err != nil {
			chata.Log(ctx).Error("error loading users", "error", err)
		}

		waitGroup.Wait()

		close(channel)
	}(userChan, userDir)

	users := Users{}
	log := chata.Log(ctx)

	var err error
	for user := range userChan {
		if user.e != nil {
			log.Error("error loading user", "error", user.e)
			err = user.e
			continue
		}

		e := users.Add(user.u)
		if e != nil {
			log.Error("invalid user", "user", user.u.ID, "error", e)
			err = e
			continue
		}
	}

	return users, err
}
