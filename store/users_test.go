package store_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/rchamarthy/chata/auth"
	"github.com/rchamarthy/chata/store"
	"github.com/stretchr/testify/require"
)

func TestUsers(t *testing.T) {
	t.Parallel()

	require := require.New(t)
	u := store.Users{}

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
	db := store.NewUserDB("dir-doesnt-exists")
	require.NotNil(db)

	e := db.Load(ctx)
	require.Error(e)

	db = store.NewUserDB("./test-dir")
	require.NotNil(db)
	e = db.Load(ctx)
	require.NoError(e)
	require.Len(db.GetAllUsers(), 10)

	e = os.WriteFile("./test-dir/bad-user", []byte("name: blah"), 0600)
	require.NoError(e)
	e = db.Load(ctx)
	require.Error(e)
}

func TestUserDB(t *testing.T) {
	t.Parallel()

	require := require.New(t)
	db := store.NewUserDB("./test-user-db")
	require.NotNil(db)
	require.NoError(db.Init())
	require.True(db.IsEmpty())

	defer os.RemoveAll("./test-user-db")

	// test Add
	require.NoError(db.Add(auth.NewUser("user1", "user1")))
	require.NoError(db.Add(auth.NewUser("user2", "user2")))
	require.Error(db.Add(auth.NewUser("user1", "")))
	require.False(db.IsEmpty())

	// test GetAllUsers
	users := db.GetAllUsers()
	require.Len(users, 2)

	// Test GetUser
	u := db.GetUser("user1")
	require.NotNil(u)
	require.Equal("user1", u.Name)
	require.Equal("user1", u.ID)
	require.True(db.HasUser("user1"))

	u = db.GetUser("user2")
	require.NotNil(u)
	require.Equal("user2", u.Name)
	require.Equal("user2", u.ID)
	require.True(db.HasUser("user1"))

	require.Nil(db.GetUser("user3"))
	require.False(db.HasUser("user3"))

	// Test delete
	require.NoError(db.DeleteUser("user1"))
	require.NoError(db.DeleteUser("user2"))
	require.Error(db.DeleteUser("user3"))

	// Test Destroy
	require.NoError(db.Destroy())
	require.True(db.IsEmpty())
}
