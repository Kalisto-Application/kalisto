package tests

import (
	"kalisto/src/assembly"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestYo(t *testing.T) {
	app := assembly.NewApp()
	api := app.Api
	newWs, err := api.NewWorkspace()
	require.NoError(t, err)
	_ = newWs
	println(newWs.Spec.Services[0].Methods[0].RequestExample)
}
