package main

import (
	"testing"
)

func TestGenerateKey(t *testing.T) {
	in := []byte("Hello, World!")
	key := GenerateKey(in)
	t.Logf("%v", key)
}
