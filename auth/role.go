package auth

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Role int

const MaxRoles = 3

const (
	SELF Role = iota
	ADMIN
	CHATTER
)

func (r Role) MarshalText() ([]byte, error) {
	switch r {
	case ADMIN:
		return []byte("admin"), nil
	case CHATTER:
		return []byte("chatter"), nil
	case SELF:
		return []byte("self"), nil
	}

	return nil, errors.New("unknown role")
}

func (r *Role) UnmarshalText(text []byte) error {
	role := string(text)
	switch role {
	case "admin":
		*r = ADMIN
		return nil
	case "chatter":
		*r = CHATTER
		return nil
	case "self":
		*r = SELF
		return nil
	}

	return fmt.Errorf("unknown role: %s", role)
}

type Roles struct {
	roles map[Role]any
}

func NewRoles(roles ...Role) *Roles {
	r := &Roles{roles: make(map[Role]any)}
	for _, role := range roles {
		r.Add(role)
	}

	return r
}

func (r *Roles) Equal(anotherRoles *Roles) bool {
	if len(r.roles) != len(anotherRoles.roles) {
		return false
	}

	for anotherRole := range anotherRoles.roles {
		if !r.HasRole(anotherRole) {
			return false
		}
	}

	return true
}

func (r *Roles) HasRole(role Role) bool {
	_, ok := r.roles[role]

	return ok
}

func (r *Roles) Add(role Role) {
	r.roles[role] = nil
}

func (r *Roles) Remove(role Role) {
	delete(r.roles, role)
}

func (r *Roles) MarshalText() ([]byte, error) {
	b := bytes.Buffer{}
	i := 0
	delim := []byte(",")
	for role := range r.roles {
		rb, err := role.MarshalText()
		if err != nil {
			return nil, err
		}

		if i != 0 {
			b.Write(delim)
		}

		b.Write(rb)
		i++
	}

	return b.Bytes(), nil
}

func (r *Roles) UnmarshalText(text []byte) error {
	rolesText := string(text)
	r.roles = make(map[Role]any)
	tokens := strings.Split(rolesText, ",")
	for _, token := range tokens {
		role := SELF
		t := strings.TrimSpace(token)
		t = strings.ToLower(t)
		if e := role.UnmarshalText([]byte(t)); e != nil {
			return e
		}

		r.Add(role)
	}

	return nil
}
