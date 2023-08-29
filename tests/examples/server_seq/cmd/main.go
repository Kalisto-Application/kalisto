package main

import (
	server "kalisto/tests/examples/server_seq"
	"log"
)

func main() {
	_, closed, err := server.Run(":9000")
	if err != nil {
		log.Fatalln(err)
	}
	<-closed
}
