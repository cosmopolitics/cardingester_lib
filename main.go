package main

import (
	"net/http"
)

type config struct {
	client *http.Client
}

func main() {
	cfg := &config{
		client: &http.Client{},
	}
	startRepl(cfg)
}
