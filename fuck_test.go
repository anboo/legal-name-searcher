package main

import (
	"testing"
)

func TestFuck(t *testing.T) {
	if Fuck() != "fuck" {
		t.Fatalf("not fuck")
	}
}
