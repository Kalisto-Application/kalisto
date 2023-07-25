package interpreter

import "testing"

func TestIT(t *testing.T) {
	ip := NewInterpreter("{token: 'yoba'}")

	md, err := ip.CreateMetadata("{auth: g.token}")
	if err != nil {
		t.Fatal(err)
	}
	if md.Get("auth")[0] != "yoba" {
		t.Fatal("auth is not yoba")
	}
}
