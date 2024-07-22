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
		panic(e)
	}

	if e := server.Run(); e != nil {
		panic(e)
	}
}
