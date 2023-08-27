package server

import (
	"context"
	"fmt"
	pb "kalisto/tests/examples/proto_sequence"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type server struct {
	pb.UnimplementedSequenceServiceServer
}

func (s *server) First(ctx context.Context, in *pb.Seq) (*pb.Seq, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if in.Rpc != "" || in.Value != 0 {
		return nil, fmt.Errorf("expected empty request in First")
	}
	in.Rpc = "First"
	in.Value = 1

	grpc.SetHeader(ctx, md)
	return in, nil
}

func (s *server) Second(ctx context.Context, in *pb.Seq) (*pb.Seq, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if in.Rpc != "First" || in.Value != 1 {
		return nil, fmt.Errorf("expected first request in Second")
	}
	in.Rpc = "Second"
	in.Value = 2

	grpc.SetHeader(ctx, md)
	return in, nil
}

func (s *server) Third(ctx context.Context, in *pb.Seq) (*pb.Seq, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if in.Rpc != "Second" || in.Value != 2 {
		return nil, fmt.Errorf("expected second request in Third")
	}
	in.Rpc = "Third"
	in.Value = 3

	grpc.SetHeader(ctx, md)
	return in, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterSequenceServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	log.Println("server closed")
}
