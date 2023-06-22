package main

import (
	"context"
	"kalisto/src/api"
	"kalisto/src/models"
	"kalisto/src/proto/client"
	"kalisto/src/proto/compiler"
	"log"
	"os"
	"path"
)

func main() {
	newClient := func(ctx context.Context, addr string) (api.Client, error) {
		return client.NewClient(ctx, client.Config{
			Addr: addr,
		})
	}

	a := api.New(compiler.NewFileCompiler(), nil, nil, nil, newClient)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not get working directory: %v", err)
	}

	req := models.Request{
		Addr:            ":9000",
		ProtoPath:       path.Join(wd, "tests/examples/proto/service.proto"),
		FullServiceName: "kalisto.tests.examples.service.BookStore",
		MethodName:      "GetBook",
		Script: `
			a = "1"
 			request = {id: a}
			`,
	}

	resp, err := a.SendGrpc(req)
	if err != nil {
		log.Fatalf("could not send grpc request: %v", err)
	}

	log.Printf("Response: %s\n", resp.Body)
}
