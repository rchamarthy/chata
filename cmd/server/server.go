package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"gojini.dev/config"
)

type Config struct {
	Address  string `json:"address"  yaml:"address"`
	UsersDir string `json:"usersDir" yaml:"usersDir"`
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

	cfg := &Config{Address: ":8080", UsersDir: "."}
	if e := store.Get("server", cfg); e != nil {
		return nil, e
	}

	engine := gin.Default()
	return &ChatServer{
		engine: engine,
		config: cfg,
		users:  NewUserHandler(engine, cfg),
		chats:  NewChatHandler(engine, cfg),
	}, nil
}

func (s *ChatServer) Run() error {
	return s.engine.Run(s.config.Address)
}
