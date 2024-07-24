package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <config-file>\n", os.Args[0])
		os.Exit(1)
	}

	server, e := NewServer(os.Args[1])
	if e != nil {
		fmt.Printf("Server Error: %v\n", e)
		os.Exit(2)
	}

	if e := server.Run(); e != nil {
		fmt.Printf("Server Run Error: %v\n", e)
		os.Exit(3)
	}
}
