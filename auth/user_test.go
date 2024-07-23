package auth_test

import (
	"os"
	"testing"

	"github.com/rchamarthy/chata/auth"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	u := auth.NewUser("user1", "user1")
	require.NotNil(u)
	require.NotEmpty(u.ID)
	require.NotEmpty(u.Name)
	require.NotNil(u.Key)
	require.NotEmpty(u.Roles)
	require.True(u.Roles.HasRole(auth.SELF))
	require.NoError(u.Validate())

	u = auth.NewUser("user1", "user1")
	u.Key = nil
	require.Error(u.Validate())

	u = auth.NewUser("user1", "")
	require.Error(u.Validate())

	u = auth.NewUser("", "user1")
	require.Error(u.Validate())
}

func TestLoadUser(t *testing.T) {
	t.Parallel()

	require := require.New(t)
	u, e := auth.LoadUser("file-doesnt-exist")
	require.Error(e)
	require.Nil(u)

	require.NoError(os.WriteFile("./blah", []byte("blah"), 0600))
	defer os.Remove("./blah")
	u, e = auth.LoadUser("./blah")
	require.Error(e)
	require.Nil(u)

	u = auth.NewUser("user1", "user1")
	require.NoError(u.SaveUser("."))
	defer os.Remove("./user1")

	u1, e := auth.LoadUser("./user1")
	require.NoError(e)
	require.NotNil(u1)
	require.Equal(u.Name, u1.Name)
	require.Equal(u.ID, u1.ID)
}

func TestSaveUser(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	u := auth.NewUser("", "")
	require.NotNil(u)
	require.Error(u.SaveUser("blah"))

	u = auth.NewUser("user1", "user1")
	u.Roles.Add(auth.Role(200))
	require.Error(u.SaveUser("blah"))

	u = auth.NewUser("user1", "user1")
	require.Error(u.SaveUser("blah"))

	require.NoError(u.SaveUser("."))
	defer os.Remove("user1")
}
