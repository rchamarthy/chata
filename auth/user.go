package auth

import (
	"context"
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
	Name  string    `json:"name" yaml:"name"`
	ID    string    `json:"id" yaml:"name"`
	Key   *Identity `json:"key" yaml:"key"`
	Roles Roles     `json:"roles" yaml:"roles"`
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
		return fmt.Errorf("id cannot be empty")
	} else if user.Name == "" {
		return fmt.Errorf("user cannot be empty")
	} else if user.Key == nil {
		return fmt.Errorf("user public key cannot be empty")
	}

	return nil
}

func LoadUser(userFile string) (*User, error) {
	b, err := os.ReadFile(userFile)
	if err != nil {
		return nil, err
	}

	user := &User{}
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
		e = fmt.Errorf("error writing user file: %x", e)
	}

	return e
}

type Users map[string]*User

type user_error struct {
	u *User
	e error
}

func (users Users) Add(user *User) error {
	if e := user.Validate(); e != nil {
		return e
	}

	users[user.Name] = user
	return nil
}

func LoadUsers(ctx context.Context, userDir string) (Users, error) {
	userChan := make(chan user_error, 16)

	go func(channel chan user_error, d string) {
		waitGroup := &sync.WaitGroup{}
		err := filepath.WalkDir(d,
			func(p string, d fs.DirEntry, err error) error {
				waitGroup.Add(1)
				go func() {
					defer waitGroup.Done()
					if d.Type().IsRegular() {
						u, e := LoadUser(p)
						channel <- user_error{u, e}
					}
				}()

				return nil
			})

		if err != nil {
			fmt.Println(err)
		}

		waitGroup.Wait()

		close(channel)
	}(userChan, userDir)

	users := Users{}
	log := chata.Log(ctx)
	for user := range userChan {
		if user.e != nil {
			log.Error("error loading user", "error", user.e)
			continue
		}

		e := users.Add(user.u)
		if e != nil {
			log.Error("invalid user", "user", user.u.ID, "error", e)
			continue
		}
	}

	return users, nil
}
