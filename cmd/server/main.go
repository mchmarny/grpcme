package main

import (
	"io"
	"log"
	"net"

	"os"

	ptypes "github.com/golang/protobuf/ptypes"

	ev "github.com/mchmarny/gcputil/env"
	pb "github.com/mchmarny/grpc-sample/pkg/api/v1"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	logger   = log.New(os.Stdout, "", 0)
	grpcPort = ev.MustGetEnvVar("PORT", "8080")
)

type messageService struct{}

func (s *messageService) Send(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	return &pb.Response{
		Content: &pb.Content{
			Index:      1,
			Message:    req.GetMessage(),
			ReceivedOn: ptypes.TimestampNow(),
		},
	}, nil
}

func (s *messageService) SendStream(stream pb.MessageService_SendStreamServer) error {
	var i int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			logger.Println("Client disconnected")
			return nil
		}
		if err != nil {
			return errors.Wrap(err, "Failed to receive send")
		}

		m := req.GetMessage()
		i++
		logger.Printf("Replying to send[%d]: %+v", i, m)

		err = stream.Send(&pb.Response{
			Content: &pb.Content{
				Index:      i,
				Message:    m,
				ReceivedOn: ptypes.TimestampNow(),
			},
		})

		if err != nil {
			return errors.Wrap(err, "Failed to send pong")
		}
	}
}

func startGRPCServer(hostPort string) error {
	listener, err := net.Listen("tcp", hostPort)
	if err != nil {
		return errors.Wrapf(err, "Failed to listen on %s: %v", hostPort, err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterMessageServiceServer(grpcServer, &messageService{})
	return grpcServer.Serve(listener)
}

func main() {

	grpcHostPort := net.JoinHostPort("0.0.0.0", grpcPort)

	go func() {
		err := startGRPCServer(grpcHostPort)
		if err != nil {
			logger.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	logger.Println("Server started...")
	select {}
}
