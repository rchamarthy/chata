package auth_test

import (
	"context"
	"fmt"
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

func TestUsers(t *testing.T) {
	t.Parallel()

	require := require.New(t)
	u := auth.Users{}

	require.NoError(u.Add(auth.NewUser("user1", "user1")))
	require.NotEmpty(u)
	require.NotNil(u["user1"])

	require.Error(u.Add(auth.NewUser("", "")))
}

func TestLoadUsers(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	// Create 10 users and save them as files
	require.NoError(os.MkdirAll("./test-dir", 0777))
	defer os.RemoveAll("./test-dir")
	for i := range 10 {
		n := fmt.Sprintf("user-%d", i)
		u := auth.NewUser(n, n)
		require.NoError(u.SaveUser("./test-dir"))
	}

	// Load all users
	ctx := context.Background()
	u, e := auth.LoadUsers(ctx, "dir-doesnt-exist")
	require.Error(e)
	require.Empty(u)

	u, e = auth.LoadUsers(ctx, "./test-dir")
	require.NoError(e)
	require.Len(u, 10)

	e = os.WriteFile("./test-dir/bad-user", []byte("name: blah"), 0600)
	require.NoError(e)
	_, e = auth.LoadUsers(ctx, "./test-dir")
	require.Error(e)
}
