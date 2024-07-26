package chat

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Message struct {
	Sender string    `json:"sender" yaml:"sender"`
	Body   string    `json:"body"   yaml:"body"`
	Time   time.Time `json:"time"   yaml:"time"`
}

type Session struct {
	ID        string    `json:"id"        yaml:"id"`
	User1     string    `json:"user1"     yaml:"user1"`
	User2     string    `json:"user2"     yaml:"user2"`
	StartTime time.Time `json:"startTime" yaml:"startTime"`
	LastMsg   time.Time `json:"lastMsg"   yaml:"lastMsg"`
	Messages  []Message `json:"messages"  yaml:"messages"`
}

func NewSession(user1 string, user2 string) *Session {
	if strings.Compare(user1, user2) == 0 {
		return nil
	}

	id := fmt.Sprintf("%s-%s", user1, user2)
	if strings.Compare(user1, user2) > 0 {
		id = fmt.Sprintf("%s-%s", user2, user1)
	}

	return &Session{
		ID:        id,
		User1:     user1,
		User2:     user2,
		StartTime: time.Now(),
		LastMsg:   time.Now(),
		Messages:  []Message{},
	}
}

func (s *Session) AddMessage(from, msg string) {
	s.LastMsg = time.Now()
	s.Messages = append(s.Messages, Message{
		Sender: from,
		Body:   msg,
		Time:   s.LastMsg,
	})
}

func (s *Session) GetMessages(index int) []Message {
	return s.Messages[index:]
}

func (s *Session) LastNMessages(n int) []Message {
	if n > len(s.Messages) {
		n = len(s.Messages)
	}

	return s.Messages[len(s.Messages)-n:]
}

func (s *Session) Save(sessionDir string) error {
	sessionFile := path.Join(sessionDir, s.ID)

	b, _ := yaml.Marshal(s)

	return os.WriteFile(sessionFile, b, 0600)
}

func LoadSession(sessionFile string) (*Session, error) {
	b, err := os.ReadFile(sessionFile)
	if err != nil {
		return nil, err
	}

	session := &Session{}
	if e := yaml.Unmarshal(b, session); e != nil {
		return nil, e
	}

	return session, nil
}

func (s *Session) Delete(sessionDir string) error {
	return os.Remove(path.Join(sessionDir, s.ID))
}
