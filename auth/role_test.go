package auth_test

import (
	"testing"

	"github.com/rchamarthy/chata/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRole(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	require := require.New(t)

	b, e := auth.ADMIN.MarshalText()
	require.NoError(e)
	assert.Equal("admin", string(b))

	b, e = auth.CHATTER.MarshalText()
	require.NoError(e)
	assert.Equal("chatter", string(b))

	b, e = auth.SELF.MarshalText()
	require.Error(e)
	assert.Empty(b)

	r := auth.SELF
	e = r.UnmarshalText([]byte("admin"))
	require.NoError(e)
	assert.Equal(auth.ADMIN, r)

	e = r.UnmarshalText([]byte("chatter"))
	require.NoError(e)
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
	require := require.New(t)

	r := auth.NewRoles(auth.ADMIN, auth.CHATTER)
	b, e := r.MarshalText()
	require.NoError(e)

	newRoles := auth.NewRoles()
	e = newRoles.UnmarshalText(b)
	require.NoError(e)
	assert.True(r.Equal(newRoles))

	r.Add(auth.SELF)
	b, e = r.MarshalText()
	require.Error(e)
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
