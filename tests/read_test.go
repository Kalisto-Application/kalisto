package tests

import (
	"encoding/json"
	"kalisto/src/assembly"
	"log"
	"testing"
)

func BenchmarkRead(b *testing.B) {
	app, err := assembly.NewApp()
	if err != nil {
		panic(err)
	}

	for n := 0; n < b.N; n++ {
		ws, err := app.Api.FindWorkspaces()
		if err != nil {
			panic(err)
		}
		log.Default().Println("found: ", len(ws))
		data, err := json.Marshal(ws)
		if err != nil {
			panic(err)
		}
		_ = data
	}
}
