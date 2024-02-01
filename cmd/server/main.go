package main

import (
	"flag"
	"log"

	"github.com/mchmarny/grpcme/pkg/server"
)

var (
	address string

	// set at build time
	version = "v0.0.1-default"
)

func main() {
	flag.StringVar(&address, "address", ":50051", "Server address (host:port)")
	flag.Parse()

	if err := server.Run(version, address); err != nil {
		log.Fatalf("error while running server: %v", err)
	}
}
