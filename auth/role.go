package auth

import (
	"bytes"
	"fmt"
	"strings"
)

type Role int

const MAX_ROLES = 3

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
	}

	return nil, fmt.Errorf("unknown role")
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
	}

	return fmt.Errorf("unknown role: %s", role)
}

type Roles map[Role]any

func NewRoles(roles ...Role) Roles {
	r := make(Roles, MAX_ROLES)
	for _, role := range roles {
		r.Add(role)
	}

	return r
}

func (roles Roles) Equal(anotherRoles Roles) bool {
	if len(roles) != len(anotherRoles) {
		return false
	}

	for anotherRole := range anotherRoles {
		if !roles.HasRole(anotherRole) {
			return false
		}
	}

	return true
}

func (roles Roles) HasRole(role Role) bool {
	_, ok := roles[role]
	return ok
}

func (roles Roles) Add(role Role) {
	roles[role] = nil
}

func (roles Roles) Remove(role Role) {
	delete(roles, role)
}

func (roles Roles) MarshalText() ([]byte, error) {
	b := bytes.Buffer{}
	i := 0
	for role := range roles {
		rb, err := role.MarshalText()
		if err != nil {
			return nil, err
		}

		if i != 0 {
			b.Write([]byte(","))
		}

		b.Write(rb)
		i++
	}

	return b.Bytes(), nil
}

func (roles Roles) UnmarshalText(text []byte) error {
	rolesText := string(text)

	tokens := strings.Split(rolesText, ",")
	for _, token := range tokens {
		r := SELF
		t := strings.TrimSpace(token)
		t = strings.ToLower(t)
		if e := r.UnmarshalText([]byte(t)); e != nil {
			return e
		}

		roles.Add(r)
	}

	return nil
}
