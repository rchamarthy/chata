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

}

func (h *UserHandler) DeleteUser(c *gin.Context) {

}

func (h *UserHandler) GetAllUsers(c *gin.Context) {

}

func (h *UserHandler) GetUser(c *gin.Context) {

}

func (h *UserHandler) UpdateUser(c *gin.Context) {

}
