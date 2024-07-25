package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rchamarthy/chata/store"
)

type UserHandler struct {
	usersDir string
	db       *store.UserDB
}

func NewUserHandler(e *gin.Engine, config *Config) *UserHandler {
	u := &UserHandler{
		usersDir: config.UsersDir,
		db:       store.NewUserDB(config.UsersDir),
	}

	if e := u.db.Init(); e != nil {
		panic(e)
	}

	if e := u.db.Load(context.Background()); e != nil {
		panic(e)
	}

	e.GET("/users", u.GetAllUsers)
	e.GET("/users/:id", u.GetUser)
	e.PUT("/users/:id", u.RegisterUser)
	e.DELETE("/users/:id", u.DeleteUser)
	e.POST("/users/:id", u.UpdateUser)

	return u
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "id is required"})
		return
	}
	if e := h.db.RegisterUser(id); e != nil {
		c.JSON(500, gin.H{"error": e.Error()})
		return
	}
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "id is required"})
		return
	}
	if e := h.db.DeleteUser(id); e != nil {
		c.JSON(500, gin.H{"error": e.Error()})
		return
	}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.db.GetAllUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, users)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "id is required"})
		return
	}
	user, err := h.db.GetUser(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "id is required"})
		return
	}
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if e := h.db.UpdateUser(id, user); e != nil {
		c.JSON(500, gin.H{"error": e.Error()})
		return
	}
}
