package main

import "github.com/gin-gonic/gin"

type ChatHandler struct {
}

func NewChatHandler(e *gin.Engine, config *Config) *ChatHandler {
	c := &ChatHandler{}

	return c
}
