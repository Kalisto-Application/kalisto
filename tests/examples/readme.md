### this is a service example repsitory 
it's created in order to tests the app against the target servie.

#### how to install
- download protoc: https://github.com/protocolbuffers/protobuf/releases
- copy the files, for example bin to /usr/local/bin/ and include/google to /usr/local/include/
- isntall go plugins:
  - go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
  - go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
