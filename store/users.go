package store

import (
	"context"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/rchamarthy/chata"
	"github.com/rchamarthy/chata/auth"
)

type Users map[string]*auth.User

func (users Users) Add(user *auth.User) error {
	if e := user.Validate(); e != nil {
		return e
	}

	users[user.ID] = user
	return nil
}

func (users Users) Copy() Users {
	copyUsers := Users{}
	for id, user := range users {
		copyUsers[id] = user
	}

	return copyUsers
}

type userError struct {
	u *auth.User
	e error
}

type UserDB struct {
	usersDir string
	users    Users
	lock     *sync.RWMutex
}

func NewUserDB(db string) *UserDB {
	return &UserDB{
		usersDir: db,
		users:    Users{},
		lock:     &sync.RWMutex{},
	}
}

func (db *UserDB) Init() error {
	db.lock.Lock()
	defer db.lock.Unlock()

	return os.MkdirAll(db.usersDir, 0755)
}

func (db *UserDB) Destroy() error {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.users = Users{}
	return os.RemoveAll(db.usersDir)
}

func (db UserDB) Load(ctx context.Context) error {
	db.lock.Lock()
	defer db.lock.Unlock()

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
						u, e := auth.LoadUser(p)
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
	}(userChan, db.usersDir)

	log := chata.Log(ctx)

	var err error
	for user := range userChan {
		if user.e != nil {
			log.Error("error loading user", "error", user.e)
			err = user.e
			continue
		}

		e := db.users.Add(user.u)
		if e != nil {
			log.Error("invalid user", "user", user.u.ID, "error", e)
			err = e
			continue
		}
	}

	return err
}

func (db *UserDB) Add(user *auth.User) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if e := user.SaveUser(db.usersDir); e != nil {
		return e
	}

	return db.users.Add(user)
}

func (db *UserDB) GetUser(id string) *auth.User {
	db.lock.RLock()
	defer db.lock.RUnlock()

	return db.users[id]
}

func (db *UserDB) GetAllUsers() Users {
	db.lock.RLock()
	defer db.lock.RUnlock()

	return db.users.Copy()
}

func (db *UserDB) DeleteUser(id string) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	path := path.Join(db.usersDir, id)
	if e := os.Remove(path); e != nil {
		return e
	}

	delete(db.users, id)
	return nil
}

func (db *UserDB) HasUser(id string) bool {
	return db.GetUser(id) != nil
}

func (db *UserDB) IsEmpty() bool {
	db.lock.RLock()
	defer db.lock.RUnlock()

	return len(db.users) == 0
}
