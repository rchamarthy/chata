package main

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"gojini.dev/config"
)

type Config struct {
	Address  string `json:"address"  yaml:"address"`
	UsersDir string `json:"usersDir" yaml:"usersDir"`
	ChatsDir string `json:"chatsDir" yaml:"chatsDir"`
}

func (c *Config) Validate() error {
	if c.Address == "" {
		return errors.New("address is not specified")
	}

	if c.UsersDir == "" {
		return errors.New("usersDir is not specified")
	}

	if c.ChatsDir == "" {
		return errors.New("chatsDir is not specified")
	}

	return nil
}

type ChatServer struct {
	engine *gin.Engine
	config *Config
	users  *UserHandler
	chats  *ChatHandler
}

func NewServer(configFile string) (*ChatServer, error) {
	ctx := context.Background()

	store := config.New()
	if e := store.LoadFromFile(ctx, configFile); e != nil {
		return nil, e
	}

	cfg := &Config{}
	if e := store.Get("server", cfg); e != nil {
		return nil, e
	}

	if e := cfg.Validate(); e != nil {
		return nil, e
	}

	engine := gin.Default()
	users := NewUserHandler(engine, cfg)
	return &ChatServer{
		engine: engine,
		config: cfg,
		users:  users,
		chats:  NewChatHandler(engine, cfg, users.db),
	}, nil
}

func (s *ChatServer) Run() error {
	return s.engine.Run(s.config.Address)
}
