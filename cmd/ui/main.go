package main

import (
	"context"
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

	log.Printf("creating client (%s)...", version)

	c, err := client.NewClient(target)
	if err != nil {
		log.Fatalf("error while creating client: %v", err)
	}

	log.Printf("sending scalar message: %s", message)
	m, err := c.Scalar(context.Background(), message)
	if err != nil {
		log.Fatalf("error while running client: %v", err)
	}

	log.Printf("received response: %s", m)
}
