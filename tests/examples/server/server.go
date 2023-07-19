package main

import (
	"context"
	"encoding/json"
	pb "kalisto/tests/examples/proto"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedBookStoreServer
}

func (s *server) GetBook(ctx context.Context, in *pb.GetBookRequest) (*pb.GetBookResponse, error) {
	dur := in.Dur.AsDuration().String()
	t := in.Time.AsTime().String()
	tt := time.UnixMilli(in.Dur.AsDuration().Milliseconds()).String()
	log.Println("dur: ", dur)
	log.Println("time: ", t)
	log.Println("since: ", tt)
	data, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	log.Printf("Received: %s\n", string(data))
	return &pb.GetBookResponse{Name: "Chuk und Gek"}, nil
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

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterBookStoreServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	log.Println("server closed")
}
