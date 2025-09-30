package main

import (
	"log"

	"github.com/flaskbush/seerx/c2/cli"
	"github.com/flaskbush/seerx/c2/server"
)

func main() {
	// Create and start the C2 server
	srv, err := server.NewServer("seerx.db")
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Start the server in a goroutine
	go func() {
		if err := srv.Start(":8080"); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Start the interactive CLI
	cli := cli.NewCLI(srv)
	cli.Start()
}
