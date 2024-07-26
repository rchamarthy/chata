package chat_test

import (
	"os"
	"testing"

	"github.com/rchamarthy/chata/chat"
	"github.com/stretchr/testify/require"
)

func TestSession(t *testing.T) {
	t.Parallel()

	require := require.New(t)
	s := chat.NewSession("user1", "user2")
	require.NotNil(s)
	require.NotEmpty(s.User1)
	require.NotEmpty(s.User2)
	require.NotEmpty(s.StartTime)
	require.NotEmpty(s.LastMsg)
	require.Empty(s.Messages)
	require.Equal("user1-user2", s.ID)

	s.AddMessage("user1", "hello")
	require.Len(s.Messages, 1)
	require.Equal("hello", s.Messages[0].Body)
	require.Equal("user1", s.Messages[0].Sender)
	require.Equal("hello", s.LastNMessages(1)[0].Body)
	require.Len(s.GetMessages(0), 1)
	require.Len(s.LastNMessages(2), 1)

	require.Nil(chat.NewSession("user1", "user1"))

	s = chat.NewSession("user2", "user1")
	require.NotNil(s)
	require.Equal("user1-user2", s.ID)
}

func TestSessionSave(t *testing.T) {
	t.Parallel()

	require := require.New(t)
	s := chat.NewSession("user1", "user2")
	require.NotNil(s)
	require.NoError(s.Save("./"))

	bs, err := chat.LoadSession("./user1-user3")
	require.Error(err)
	require.Nil(bs)

	s1, err := chat.LoadSession("./user1-user2")
	require.NoError(err)
	require.NotNil(s1)
	require.Equal(s.User1, s1.User1)
	require.Equal(s.User2, s1.User2)
	require.Equal(len(s.Messages), len(s1.Messages))

	require.NoError(os.WriteFile("crappy-session", []byte("blah"), 0600))
	defer os.Remove("crappy-session")
	s2, err := chat.LoadSession("crappy-session")
	require.Error(err)
	require.Nil(s2)

	require.NoError(s.Delete("./"))
}
