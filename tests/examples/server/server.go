package main

import (
	"context"
	"encoding/json"
	pb "kalisto/tests/examples/proto"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	pb.UnimplementedBookStoreServer
}

func (s *server) GetBook(ctx context.Context, in *pb.GetBookRequest) (*pb.GetBookRequest, error) {
	log.Printf("Received: %s\n", in.String())
	md, _ := metadata.FromIncomingContext(ctx)
	log.Printf("meta: %s", md.Get("k"))

	res := &pb.GetBookRequest{
		Id:       "1",
		Double:   999.999,
		Float:    999.999,
		Int32:    33,
		Int64:    33,
		Uint32:   33,
		Uint64:   33,
		Sint32:   33,
		Sint64:   33,
		Fixed32:  33,
		Fixed64:  33,
		Sfixed32: 33,
		Sfixed64: 33,
		Bool:     true,
		Bytes:    []byte(`{"just": "bytes"}`),
		Enum:     3,
		Book: &pb.Book{
			Id: "1",
		},
		StrToStr: map[string]string{"str": "str"},
		IntToBooks: map[int32]*pb.Book{1: &pb.Book{
			Id: "1",
		}},
		Etoe:    map[bool]pb.GetBookRequest_Enum{true: 4},
		Strings: []string{"str1", "str2"},
		Enums:   []pb.GetBookRequest_Enum{3, 4},
		Uints:   []uint32{1, 2},
		Books: []*pb.Book{&pb.Book{
			Id: "111",
		}, &pb.Book{
			Id: "222",
		}},
		DeepNestedBook: &pb.DeepNestedBook{
			HasNested: &pb.BookHasDeepNested{
				DeepNested: &pb.DeepNestedBook{
					HasNested: nil,
				},
			},
		},
		RepeatedNestedBook: []*pb.DeepNestedBook{
			&pb.DeepNestedBook{
				HasNested: &pb.BookHasDeepNested{
					DeepNested: &pb.DeepNestedBook{
						HasNested: nil,
					},
				},
			},
		},
		SomeBook: &pb.GetBookRequest_OneEnum{OneEnum: 4},
		AnotherBook: &pb.GetBookRequest_AnotherBookObject{AnotherBookObject: &pb.Book{
			Id: "333",
		}},
		Dur:  durationpb.New(time.Hour),
		Time: timestamppb.New(time.Now()),
	}

	grpc.SetTrailer(ctx, metadata.New(map[string]string{
		"content": "yes",
	}))
	return res, nil
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
