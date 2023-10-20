package tests

import (
	"kalisto/src/definitions/proto/interpreter"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertJsObjectsAreEqual(t *testing.T, v1, v2 string) bool {
	t.Helper()

	ip := interpreter.NewInterpreter("")
	ex1, err := ip.ExportValue(v1, "")
	assert.NoError(t, err)
	ex2, err := ip.ExportValue(v2, "")
	assert.NoError(t, err)

	return assert.EqualValues(t, ex1, ex2)
}
