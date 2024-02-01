package server

import (
	"context"
	"log"
	"net"
	"sync/atomic"

	pb "github.com/mchmarny/grpcme/pkg/api/v1"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const (
	protocol = "tcp" // network protocol
)

var (
	counter atomic.Uint64 // counter for messages
)

// server is used to implement your Service.
type server struct {
	pb.UnimplementedServiceServer
}

func processContent(in *pb.Content) *pb.Response {
	counter.Add(1)
	return &pb.Response{
		RequestId:         in.GetId(),
		MessageCount:      int64(counter.Load()),
		MessagesProcessed: int64(counter.Load()),
		ProcessingDetails: "processed successfully",
	}
}

// Scalar implements the single method of the Service.
func (s *server) Scalar(_ context.Context, in *pb.Request) (*pb.Response, error) {
	c := in.GetContent()
	log.Printf("received scalar: %v", c.GetData())
	return processContent(c), nil
}

// Stream implements the Stream method of the Service.
func (s *server) Stream(stream pb.Service_StreamServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			return errors.Wrap(err, "failed to receive")
		}

		c := in.GetContent()
		log.Printf("received stream: %v", c.GetData())

		if err := stream.Send(processContent(c)); err != nil {
			return errors.Wrap(err, "failed to send")
		}
	}
}

func Run(version, address string) error {
	if address == "" {
		return errors.New("address is required")
	}

	log.Printf("starting server (%s)...", version)

	// Create a listener on the specified address.
	lis, err := net.Listen(protocol, address)
	if err != nil {
		return errors.Wrap(err, "failed to listen")
	}
	s := grpc.NewServer()
	pb.RegisterServiceServer(s, &server{})
	log.Printf("server (%s) listening at %v", version, lis.Addr())
	if err := s.Serve(lis); err != nil {
		return errors.Wrap(err, "failed to serve")
	}
	return nil
}
