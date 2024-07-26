package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rchamarthy/chata/chat"
	"github.com/rchamarthy/chata/store"
)

type ChatHandler struct {
	chatDir string
	db      *store.ChatDB
	userDB  *store.UserDB
}

func NewChatHandler(e *gin.Engine, config *Config, userDB *store.UserDB) *ChatHandler {
	c := &ChatHandler{
		chatDir: config.ChatsDir,
		db:      store.NewChatDB(config.ChatsDir),
		userDB:  userDB,
	}

	if e := c.db.Init(); e != nil {
		panic(e)
	}

	if e := c.db.Load(context.Background()); e != nil {
		panic(e)
	}

	e.GET("/chats/:from/:to", c.GetChatForUserAndPeer)
	e.GET("/chats/:from", c.GetAllChatsForUser)
	e.POST("/chats/:from/:to", c.AddChat)
	e.DELETE("/chats/:from/:to", c.DeleteChat)
	e.POST("/message/:from/:to", c.SendMessage)

	return c
}

func (h *ChatHandler) GetChatForUserAndPeer(c *gin.Context) {
	from := c.Param("from")
	to := c.Param("to")
	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "from or to is not specified",
		})
		return
	}

	chat := h.db.Get(from, to)
	if chat == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "chat not found",
		})
		return
	}

	c.JSON(http.StatusOK, chat)
}

func (h *ChatHandler) GetAllChatsForUser(c *gin.Context) {
	from := c.Param("from")
	if from == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "from is not specified",
		})
		return
	}

	chats := h.db.GetSessionsByUser(from)

	if len(chats) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "chats not found",
		})
		return
	}
	c.JSON(http.StatusOK, chats)
}

func (h *ChatHandler) AddChat(c *gin.Context) {
	from := c.Param("from")
	to := c.Param("to")
	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "from or to is not specified",
		})
		return
	}

	if from == to {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot chat with  yourself",
		})
		return
	}

	if !h.ValidUsers(c, from, to) {
		return
	}

	chatSession := h.db.Get(from, to)
	if chatSession != nil {
		c.JSON(http.StatusOK, gin.H{
			"id": chatSession.ID,
		})
		return
	}

	chatSession = chat.NewSession(from, to)
	if chatSession == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot chat with  yourself",
		})
		return
	}

	if e := h.db.Add(chatSession); e != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": e.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": chatSession.ID,
	})
}

func (h *ChatHandler) DeleteChat(c *gin.Context) {
	from := c.Param("from")
	to := c.Param("to")
	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "from or to is not specified",
		})
		return
	}

	if e := h.db.Delete(from, to); e != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": e.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func (h *ChatHandler) SendMessage(c *gin.Context) {
	from := c.Param("from")
	to := c.Param("to")
	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "from or to is not specified",
		})
		return
	}

	session := h.db.Get(from, to)
	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "chat not found",
		})
		return
	}

	message := struct {
		Text string `json:"text"`
	}{}

	if e := c.BindJSON(&message); e != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": e.Error(),
		})
		return
	}

	if message.Text == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "message is not specified",
		})
		return
	}

	session.AddMessage(from, message.Text)

	if e := session.Save(h.chatDir); e != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": e.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func (h *ChatHandler) ValidUsers(c *gin.Context, from string, to string) bool {
	if u := h.userDB.GetUser(from); u == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("user %s does not exist", from),
		})
		return false
	}

	if u := h.userDB.GetUser(to); u == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("user %s does not exist", to),
		})
		return false
	}

	return true
}
