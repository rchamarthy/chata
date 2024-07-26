package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/rchamarthy/chata/chat"
	"github.com/spf13/cobra"
)

func chatCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "chat",
		Short: "chat with another user",
		Long:  "manage the chat between two users",
	}

	c.AddCommand(connectCmd())
	c.AddCommand(deleteChatCmd())
	c.AddCommand(sendMessageCmd())
	c.AddCommand(showChatsCmd())

	return c
}

func connectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "connect <from> <to>",
		Short: "connect to another user",
		Long:  "connect to another user",
		RunE:  Connect,
		Args:  cobra.ExactArgs(2),
	}
}

func deleteChatCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <from> <to>",
		Short: "delete a chat between two users",
		Long:  "delete a chat between two users",
		RunE:  DeleteChat,
		Args:  cobra.ExactArgs(2),
	}
}

func sendMessageCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "send <from> <to> <message>",
		Short: "send message to another user",
		Long:  "send message to another user",
		RunE:  SendMessage,
		Args:  cobra.ExactArgs(3),
	}
}

func showChatsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "show all chats between two users",
		Long:  "show all chats between two users",
		RunE:  ShowChats,
		Args:  cobra.MaximumNArgs(2),
	}
}

func Connect(cmd *cobra.Command, args []string) error {
	serverAddress, err := cmd.Flags().GetString("server")
	if err != nil {
		return err
	}

	from := args[0]
	to := args[1]

	url := fmt.Sprintf("%s/chats/%s/%s", serverAddress, from, to)
	client := resty.New()
	r, err := client.R().Post(url)
	if err != nil {
		return err
	}

	if r.StatusCode() == http.StatusOK {
		fmt.Println("chat is already created")
		return nil
	} else if r.StatusCode() != http.StatusCreated {
		return fmt.Errorf("error creating chat:\n %s", string(r.Body()))
	}

	fmt.Printf("chat between %s and %s is created\n", from, to)
	return nil
}

func DeleteChat(cmd *cobra.Command, args []string) error {
	serverAddress, err := cmd.Flags().GetString("server")
	if err != nil {
		return err
	}

	from := args[0]
	to := args[1]

	url := fmt.Sprintf("%s/chats/%s/%s", serverAddress, from, to)
	client := resty.New()
	r, err := client.R().Delete(url)
	if err != nil {
		return err
	}

	if r.StatusCode() != http.StatusOK {
		return fmt.Errorf("error deleting chat:\n %s", string(r.Body()))
	}

	fmt.Printf("chat between %s and %s is deleted\n", from, to)
	return nil
}

func SendMessage(cmd *cobra.Command, args []string) error {
	serverAddress, err := cmd.Flags().GetString("server")
	if err != nil {
		return err
	}

	from := args[0]
	to := args[1]
	message := args[2]

	url := fmt.Sprintf("%s/message/%s/%s", serverAddress, from, to)
	client := resty.New()

	m := struct {
		Text string `json:"text"`
	}{Text: message}
	r, err := client.R().SetBody(&m).Post(url)
	if err != nil {
		return err
	}

	if r.StatusCode() != http.StatusOK {
		return fmt.Errorf("error sending message:\n %s", string(r.Body()))
	}

	fmt.Printf("message sent from %s to %s\n", from, to)
	return nil
}

func ShowChats(cmd *cobra.Command, args []string) error {
	serverAddress, err := cmd.Flags().GetString("server")
	if err != nil {
		return err
	}

	switch len(args) {
	case 0:
		return errors.New("need at least one user id")
	case 1:
		return showAllChats(serverAddress, args[0])
	case 2:
		return showOneChat(serverAddress, args[0], args[1])
	default:
		return errors.New("invalid number of arguments")
	}
}

func showAllChats(server string, from string) error {
	url := fmt.Sprintf("%s/chats/%s", server, from)
	client := resty.New()
	sessions := []*chat.Session{}
	r, err := client.R().SetResult(&sessions).Get(url)
	if err != nil {
		return err
	}

	if r.StatusCode() != http.StatusOK {
		return fmt.Errorf("error getting chats:\n %s", string(r.Body()))
	}

	for _, session := range sessions {
		otherUser := from
		if session.User1 == from {
			otherUser = session.User2
		}
		lastMessage := session.LastMsg.String()
		fmt.Printf("other user: %s messages: %d last message: %s\n", otherUser,
			len(session.Messages), lastMessage)
	}

	return nil
}

func showOneChat(server string, from string, to string) error {
	url := fmt.Sprintf("%s/chats/%s/%s", server, from, to)
	client := resty.New()
	session := &chat.Session{}
	r, err := client.R().SetResult(session).Get(url)
	if err != nil {
		return err
	}

	if r.StatusCode() != http.StatusOK {
		return fmt.Errorf("error getting chat:\n %s", string(r.Body()))
	}

	for _, msg := range session.Messages {
		fmt.Printf("%s: %s\n", msg.Sender, msg.Body)
	}

	return nil
}
