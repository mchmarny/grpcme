package client

import (
	"context"
	"time"

	"github.com/google/uuid"
	pb "github.com/mchmarny/grpcme/pkg/api/v1"
	"github.com/pkg/errors"
	anypb "google.golang.org/protobuf/types/known/anypb"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
	wrbp "google.golang.org/protobuf/types/known/wrapperspb"
)

// Scalar sends a message to the server and returns the response.
func (c *SimpleClient) Scalar(ctx context.Context, message string) (string, error) {
	if message == "" {
		return "", errors.New("message is required")
	}

	if !c.Connected() {
		return "", errors.New("client is not connected")
	}

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	// create message with wrapper
	a, err := anypb.New(wrbp.String(message))
	if err != nil {
		return "", errors.Wrap(err, "unable to create message")
	}

	// Scalar example
	r, err := c.service.Scalar(ctx, &pb.Request{
		Content: &pb.Content{
			Id:   uuid.New().String(),
			Data: a,
		},
		Sent: tspb.Now(),
	})
	if err != nil {
		return "", errors.Wrap(err, "could not send scalar message")
	}

	return r.GetProcessingDetails(), nil
}
