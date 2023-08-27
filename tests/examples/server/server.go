package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	pb "kalisto/tests/examples/proto"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedBookStoreServer
}

func (s *server) GetBook(ctx context.Context, in *pb.GetBookRequest) (*pb.GetBookRequest, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	fmt.Println("time: ", in.Time.AsTime().String())
	fmt.Println("dur: ", in.Dur.AsDuration().String())
	fmt.Println("oneEnum: ", in.GetOneEnum())

	grpc.SetHeader(ctx, md)
	return in, nil
}

func (s *server) Empty(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	data, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	log.Printf("Received: %s\n", string(data))
	return &emptypb.Empty{}, nil
}

func (s *server) Any(ctx context.Context, in *anypb.Any) (*anypb.Any, error) {
	data, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	log.Printf("Received: %s\n", string(data))
	return &anypb.Any{TypeUrl: "google.protobuf.Empty", Value: nil}, nil
}

func (s *server) Error(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Error(codes.InvalidArgument, "message")
}

func Run(port string) (func() error, <-chan struct{}, error) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return nil, nil, err
	}

	s := grpc.NewServer()
	pb.RegisterBookStoreServer(s, &server{})

	closed := make(chan struct{})
	go func() {
		defer close(closed)
		if err := s.Serve(lis); err != nil && !errors.Is(err, net.ErrClosed) {
			log.Printf("failed to serve: %v\n", err)
		}
		log.Println("server closed")
	}()
	return lis.Close, closed, nil
}
