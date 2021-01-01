package main

import (
	"fmt"
	"net/url"
)

func buildURL(arg string, format string) (string, error) {
	rawurl := ""
	if format == "" {
		rawurl = arg
	} else {
		rawurl = fmt.Sprintf(format, arg)
	}

	_, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return "", err
	} else {
		return rawurl, err
	}
}
