package main

import (
	"kalisto/tests/examples/server"
	"log"
)

func main() {
	_, closed, err := server.Run(":9000")
	if err != nil {
		log.Fatalln(err)
	}
	<-closed
}
