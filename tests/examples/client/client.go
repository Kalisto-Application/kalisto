package main

import (
	"context"
	"log"
	"time"

	pb "kalisto/tests/examples/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial(":9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewBookStoreClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetBook(ctx, &pb.GetBookRequest{Id: "1"})
	if err != nil {
		log.Fatalf("could not get a book: %v", err)
	}
	log.Printf("Book: %s\n", r.Name)

	fd := r.ProtoReflect().Descriptor().Fields().ByName("name")
	log.Printf("proto reflect \n")
	log.Printf("Book: %s\n", r.ProtoReflect().Get(fd).String())
}
