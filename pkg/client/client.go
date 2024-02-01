package client

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	pb "github.com/mchmarny/grpcme/pkg/api/v1"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	anypb "google.golang.org/protobuf/types/known/anypb"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
	wrbp "google.golang.org/protobuf/types/known/wrapperspb"
)

func Run(version, target, message string) error {
	if target == "" {
		return errors.New("target is required")
	}
	if message == "" {
		return errors.New("message is required")
	}

	log.Printf("starting client (%s)...", version)

	// Set up a connection to the server.
	creds := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial(target, creds, grpc.WithBlock())
	if err != nil {
		return errors.Wrap(err, "unable to connect")
	}
	defer conn.Close()
	c := pb.NewServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// create message with wrapper
	a, err := anypb.New(wrbp.String(message))
	if err != nil {
		return errors.Wrap(err, "unable to create message")
	}

	// Ping example
	r, err := c.Ping(ctx, &pb.PingRequest{
		Content: &pb.Content{
			Id:   uuid.New().String(),
			Data: a,
		},
		Sent: tspb.Now(),
	})
	if err != nil {
		return errors.Wrap(err, "could not ping")
	}
	log.Printf("ping Response: %s", r.GetProcessingDetails())

	// Stream example
	stream, err := c.Stream(ctx)
	if err != nil {
		return errors.Wrap(err, "could not create stream")
	}

	// Send messages to the stream
	for i := 0; i < 5; i++ {
		if err := stream.Send(&pb.PingRequest{
			Content: &pb.Content{
				Id:   uuid.New().String(),
				Data: a,
			},
			Sent: tspb.Now(),
		}); err != nil {
			return errors.Wrap(err, "failed to send a message")
		}

		reply, err := stream.Recv()
		if err != nil {
			return errors.Wrap(err, "failed to receive a reply")
		}

		log.Printf("stream response: %s", reply.GetProcessingDetails())
	}

	if err := stream.CloseSend(); err != nil {
		return errors.Wrap(err, "failed to close stream")
	}
	return nil
}
