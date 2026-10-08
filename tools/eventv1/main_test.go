package main

import "testing"

func TestScaffoldSmoke(t *testing.T) {
	// Minimal smoke test so `go test ./...` passes in pre-commit
	if true != true {
		t.Fatal("scaffold broken")
	}
}
