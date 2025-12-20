package main

import (
	"net/http"
	"time"

	"github.com/cosmopolitics/cardingester/internal"
)

func main() {
	cache := cardingester.NewCache(time.Minute * 20)
	cfg := &Config{
		client: &http.Client{},
		cache: &cache,
	}
	startRepl(cfg)
}
