package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rchamarthy/chata/auth"
	"github.com/rchamarthy/chata/store"
)

type UserHandler struct {
	usersDir string
	db       *store.UserDB
}

type UserError struct {
	Error string `json:"error" yaml:"error"`
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
	newUser := auth.NewUser("", "")

	if err := c.BindJSON(newUser); err != nil {
		c.JSON(http.StatusBadRequest, UserError{Error: err.Error()})
		return
	}

	if h.db.IsEmpty() {
		// This is the first user, first is always an admin
		newUser.Roles.Add(auth.ADMIN)
	}

	if h.db.HasUser(newUser.ID) {
		c.JSON(http.StatusConflict, UserError{
			Error: fmt.Sprintf("user id: %s already exists", newUser.ID),
		})
		return
	}

	// Add chatter and self roles to everyone else
	newUser.Roles.Add(auth.CHATTER)
	newUser.Roles.Add(auth.SELF)

	err := h.db.Add(newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, UserError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Params.ByName("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, UserError{
			Error: "id is not specified",
		})
		return
	}

	if !h.db.HasUser(id) {
		c.JSON(http.StatusNotFound, UserError{
			Error: fmt.Sprintf("id %s does not exist", id),
		})
		return
	}

	if err := h.db.DeleteUser(id); err != nil {
		c.JSON(http.StatusForbidden, UserError{Error: err.Error()})
	}

	c.JSON(http.StatusOK, "")
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	allUsers := h.db.GetAllUsers()
	c.JSON(http.StatusOK, allUsers)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Params.ByName("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, UserError{
			Error: "id is not specified",
		})
		return
	}

	user := h.db.GetUser(id)
	if user == nil {
		c.JSON(http.StatusNotFound, UserError{
			Error: fmt.Sprintf("id %s does not exist", id),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Params.ByName("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, UserError{
			Error: "id is not specified",
		})
		return
	}

	user := h.db.GetUser(id)
	if user == nil {
		c.JSON(http.StatusNotFound, UserError{
			Error: fmt.Sprintf("id %s does not exist", id),
		})
		return
	}

	if err := c.BindJSON(user); err != nil {
		c.JSON(http.StatusBadRequest, UserError{Error: err.Error()})
		return
	}

	err := h.db.Add(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, UserError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
