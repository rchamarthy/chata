package store_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/rchamarthy/chata/chat"
	"github.com/rchamarthy/chata/store"
	"github.com/stretchr/testify/require"
)

func TestSessions(t *testing.T) {
	t.Parallel()

	require := require.New(t)
	s := store.NewSessions()
	require.NotNil(s)

	session := chat.NewSession("user1", "user2")
	require.NotNil(session)
	require.NotNil(s)
	s.Add(session)

	require.NotNil(s.Get("user1", "user2"))
	require.Nil(s.Get("user1", "user3"))
	require.Nil(s.Get("user3", "user2"))
	require.Nil(s.Get("user3", "user3"))
	require.Len(s.GetSessionsByUser("user1"), 1)
	require.Len(s.GetSessionsByUser("user2"), 1)
	require.Empty(s.GetSessionsByUser("user3"))

	require.NoError(s.Delete("user1", "user2"))
	require.Nil(s.Get("user1", "user2"))
	require.Empty(s.GetSessionsByUser("user1"))
	require.Empty(s.GetSessionsByUser("user2"))

	require.Error(s.Delete("user1", "user2"))

	s.Add(session)
	require.NoError(s.Delete("user2", "user1"))
}

func TestChatDB(t *testing.T) {
	t.Parallel()

	require := require.New(t)
	db := store.NewChatDB("./test-session-db")
	require.NotNil(db)
	require.Error(db.Add(chat.NewSession("a", "b")))
	require.NoError(db.Init())

	defer os.RemoveAll("./test-session-db")

	// Add
	session := chat.NewSession("user1", "user2")
	require.NoError(db.Add(session))
	require.NoError(db.Add(session)) // Idempotent

	// Get
	require.NotNil(db.Get("user1", "user2"))
	require.Nil(db.Get("user1", "user3"))
	require.NotNil(db.GetSessionsByUser("user1"))

	// Delete
	require.NoError(db.Delete("user1", "user2"))
	require.Nil(db.Get("user1", "user2"))
	require.Error(db.Delete("user1", "user2"))

	// Destroy
	require.NoError(db.Destroy())
}

func TestSessionLoad(t *testing.T) {
	t.Parallel()

	require := require.New(t)
	db := store.NewChatDB("./test-session-load")
	require.NotNil(db)
	require.Error(db.Load(context.Background())) // Will fail before init
	require.NoError(db.Init())
	defer os.RemoveAll("./test-session-load")
	require.NoError(db.Load(context.Background()))

	// Add 10 sessions
	for i := range 10 {
		for j := range 10 {
			user1 := fmt.Sprintf("user%d", i)
			user2 := fmt.Sprintf("user%d", j)
			session := chat.NewSession(user1, user2)
			if session != nil {
				require.NoError(db.Add(session))
			}
		}
	}

	require.NoError(db.Load(context.Background()))
	for i := range 10 {
		user := fmt.Sprintf("user%d", i)
		require.Len(db.GetSessionsByUser(user), 9)
	}
}
