package main

import (
	"context"
	"github.com/jhump/protoreflect/dynamic"
	"kalisto/src/filesystem"
	"kalisto/src/proto/client"
	"kalisto/src/proto/compiler"
	"log"
	"os"
	"path"
	"time"
)

func main() {
	c, err := client.NewClient(context.Background(), client.Config{Addr: ":9000"})
	defer func() { _ = c.Close() }()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not get working directory: %v", err)
	}

	protoFiles, err := filesystem.SearchProtoFiles(path.Join(wd, "tests/examples/proto/service.proto"))
	if err != nil {
		log.Fatalf("could not search proto files: %v", err)
	}

	registry, err := compiler.NewFileCompiler().Compile([]string{protoFiles.AbsoluteDirPath}, protoFiles.RelativeProtoPaths)
	if err != nil {
		log.Fatalf("could not compile proto files: %v", err)
	}

	serviceDesc := registry.Descriptors[0].GetServices()[0]
	getBookMethodDesc := serviceDesc.GetMethods()[0]
	request, response := dynamic.NewMessage(getBookMethodDesc.GetInputType()), dynamic.NewMessage(getBookMethodDesc.GetOutputType())

	request.SetFieldByName("id", "1")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Println("Full service name:", serviceDesc.GetFullyQualifiedName())
	log.Println("Method name:", getBookMethodDesc.GetName())

	err = c.Invoke(ctx, "/"+serviceDesc.GetFullyQualifiedName()+"/"+getBookMethodDesc.GetName(), request, response)
	if err != nil {
		log.Fatalf("could not get a book: %v", err)
	}

	b, err := response.MarshalJSONIndent()
	if err != nil {
		log.Fatalf("could not marshal response: %v", err)
	}

	log.Printf("Response: %s\n", string(b))
}
