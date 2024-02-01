package main

import (
	"flag"
	"log"

	"github.com/mchmarny/grpcme/pkg/client"
)

var (
	target  string
	message string

	// set at build time
	version = "v0.0.1-default"
)

func main() {
	flag.StringVar(&target, "target", "localhost:50051", "Server address (host:port)")
	flag.StringVar(&message, "message", "hello from client", "Message to send to server")
	flag.Parse()

	if err := client.Run(version, target, message); err != nil {
		log.Fatalf("error while running client: %v", err)
	}
}
