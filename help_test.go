package main

import (
	"strings"
	"testing"
)

// go test help*
func TestHelp(t *testing.T) {
	msg := help()
	substrings := strings.Split(msg, "\n")

	result := substrings[0]
	const expected string = "Act is a tool for AtCoder."

	if result != expected {
		t.Errorf("help() = %s, want %s", result, expected)
	}
}
