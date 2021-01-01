package main

import (
	"testing"
)

// go test -v url*
func TestBuildURLWithoutFormat(t *testing.T) {
	arg := "https://example.com/foo/bar"
	format := ""

	url, err := buildURL(arg, format)
	want := "https://example.com/foo/bar"

	if url != want || err != nil {
		t.Fatalf(`buildURL("https://example.com/foo/bar", "") = %q, %v, want %#q, nil`, url, err, want)
	}
}

func TestBuildURLWithFormat(t *testing.T) {
	arg := "foo"
	format := "https://example.com/%s/bar"

	url, err := buildURL(arg, format)
	want := "https://example.com/foo/bar"

	if url != want || err != nil {
		t.Fatalf(`buildURL("foo", "https://example.com/%%s/bar") = %q, %v, want %#q, nil`, url, err, want)
	}
}

func TestBuildURLError(t *testing.T) {
	arg := "foo"
	format := ""

	url, err := buildURL(arg, format)
	want := ""

	if url != want || err == nil {
		t.Fatalf(`buildURL("foo", "") = %q, %v, want %#q, error`, url, err, want)
	}
}
