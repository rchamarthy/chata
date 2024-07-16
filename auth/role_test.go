package auth_test

import (
	"testing"

	"github.com/rchamarthy/chata/auth"
	"github.com/stretchr/testify/assert"
)

func TestRole(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	b, e := auth.ADMIN.MarshalText()
	assert.NoError(e)
	assert.Equal("admin", string(b))

	b, e = auth.CHATTER.MarshalText()
	assert.NoError(e)
	assert.Equal("chatter", string(b))

	b, e = auth.SELF.MarshalText()
	assert.Error(e)
	assert.Empty(b)

	r := auth.SELF
	e = r.UnmarshalText([]byte("admin"))
	assert.NoError(e)
	assert.Equal(auth.ADMIN, r)

	e = r.UnmarshalText([]byte("chatter"))
	assert.NoError(e)
	assert.Equal(auth.CHATTER, r)

	assert.Error(r.UnmarshalText([]byte("unknown bad role")))
}

func TestRoles(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	r := auth.NewRoles(auth.ADMIN)
	assert.NotEmpty(r)
	assert.True(r.HasRole(auth.ADMIN))
	assert.False(r.HasRole(auth.CHATTER))

	r.Add(auth.ADMIN)
	assert.NotEmpty(r)
	assert.True(r.HasRole(auth.ADMIN))
	assert.False(r.HasRole(auth.CHATTER))

	r.Add(auth.CHATTER)
	assert.NotEmpty(r)
	assert.True(r.HasRole(auth.ADMIN))
	assert.True(r.HasRole(auth.CHATTER))

	r.Remove(auth.ADMIN)
	assert.NotEmpty(r)
	assert.False(r.HasRole(auth.ADMIN))
	assert.True(r.HasRole(auth.CHATTER))

	r.Remove(auth.CHATTER)
	assert.Empty(r)
	assert.False(r.HasRole(auth.ADMIN))
	assert.False(r.HasRole(auth.CHATTER))
}

func TestRolesMarshal(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	r := auth.NewRoles(auth.ADMIN, auth.CHATTER)
	b, e := r.MarshalText()
	assert.NoError(e)

	newRoles := auth.NewRoles()
	e = newRoles.UnmarshalText(b)
	assert.NoError(e)
	assert.True(r.Equal(newRoles))

	r.Add(auth.SELF)
	b, e = r.MarshalText()
	assert.Error(e)
	assert.Nil(b)

	e = r.UnmarshalText([]byte("bad roles"))
	assert.Error(e)
}

func TestRolesEqual(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	roles := auth.NewRoles(auth.ADMIN)
	anotherRoles := auth.NewRoles(auth.CHATTER)
	assert.False(roles.Equal(anotherRoles))

	anotherRoles = auth.NewRoles(auth.ADMIN)
	assert.True(roles.Equal(anotherRoles))
}
