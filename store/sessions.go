package store

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/rchamarthy/chata"
	"github.com/rchamarthy/chata/chat"
)

type Sessions struct {
	sessions        map[string]*chat.Session
	sessionsByUsers map[string]map[string]*chat.Session
}

func NewSessions() *Sessions {
	return &Sessions{
		sessions:        map[string]*chat.Session{},
		sessionsByUsers: map[string]map[string]*chat.Session{},
	}
}

func (s *Sessions) Add(session *chat.Session) {
	s.sessions[session.ID] = session

	user1Sessions := s.sessionsByUsers[session.User1]
	if user1Sessions == nil {
		user1Sessions = map[string]*chat.Session{}
		s.sessionsByUsers[session.User1] = user1Sessions
	}
	user1Sessions[session.User2] = session

	user2Sessions := s.sessionsByUsers[session.User2]
	if user2Sessions == nil {
		user2Sessions = map[string]*chat.Session{}
		s.sessionsByUsers[session.User2] = user2Sessions
	}
	user2Sessions[session.User1] = session
}

func (s *Sessions) Get(user1 string, user2 string) *chat.Session {
	user1Sessions := s.sessionsByUsers[user1]
	if user1Sessions == nil {
		return nil
	}
	return user1Sessions[user2]
}

func (s *Sessions) GetSessionsByUser(user string) []*chat.Session {
	user1Sessions := s.sessionsByUsers[user]
	if user1Sessions == nil {
		return nil
	}

	sessions := make([]*chat.Session, 0, len(user1Sessions))
	for _, session := range user1Sessions {
		sessions = append(sessions, session)
	}
	return sessions
}

func (s *Sessions) Delete(user1 string, user2 string) error {
	session := s.Get(user1, user2)
	if session == nil {
		return fmt.Errorf("session for users %s and %s not found", user1, user2)
	}

	delete(s.sessions, session.ID)
	delete(s.sessionsByUsers[user1], user2)
	delete(s.sessionsByUsers[user2], user1)

	return nil
}

type ChatDB struct {
	sessionsDir string
	sessions    *Sessions
	lock        *sync.RWMutex
}

func NewChatDB(db string) *ChatDB {
	return &ChatDB{
		sessionsDir: db,
		sessions:    NewSessions(),
		lock:        &sync.RWMutex{},
	}
}

func (db *ChatDB) Init() error {
	db.lock.Lock()
	defer db.lock.Unlock()

	return os.MkdirAll(db.sessionsDir, 0755)
}

func (db *ChatDB) Destroy() error {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.sessions = NewSessions()
	return os.RemoveAll(db.sessionsDir)
}

type sessionError struct {
	session *chat.Session
	err     error
}

func (db *ChatDB) Load(ctx context.Context) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	sessionChan := make(chan sessionError, 16)

	go func(channel chan sessionError, d string) {
		waitGroup := &sync.WaitGroup{}
		err := filepath.WalkDir(d,
			func(p string, d fs.DirEntry, err error) error {
				waitGroup.Add(1)
				go func(p string, f fs.DirEntry) {
					defer waitGroup.Done()
					if err != nil {
						channel <- sessionError{nil, err}
						return
					}

					if f.Type().IsRegular() {
						s, e := chat.LoadSession(p)
						channel <- sessionError{s, e}
					}
				}(p, d)

				return err
			})

		if err != nil {
			chata.Log(ctx).Error("error loading sessions", "error", err)
		}

		waitGroup.Wait()
		close(channel)
	}(sessionChan, db.sessionsDir)

	var err error
	log := chata.Log(ctx)
	for se := range sessionChan {
		if se.err == nil {
			db.sessions.Add(se.session)
		} else {
			err = se.err
			log.Error("error loading session", "error", se.err)
		}
	}

	return err
}

func (db *ChatDB) Add(session *chat.Session) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if e := session.Save(db.sessionsDir); e != nil {
		return e
	}

	db.sessions.Add(session)
	return nil
}

func (db *ChatDB) Get(id1 string, id2 string) *chat.Session {
	db.lock.RLock()
	defer db.lock.RUnlock()

	return db.sessions.Get(id1, id2)
}

func (db *ChatDB) Delete(user1 string, user2 string) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	session := db.sessions.Get(user1, user2)
	if session == nil {
		return fmt.Errorf("session for users %s and %s not found", user1, user2)
	}

	if e := session.Delete(db.sessionsDir); e != nil {
		return e
	}

	return db.sessions.Delete(user1, user2)
}

func (db *ChatDB) GetSessionsByUser(user string) []*chat.Session {
	db.lock.RLock()
	defer db.lock.RUnlock()

	sessions := db.sessions.GetSessionsByUser(user)
	return sessions
}
