package main

import (
	"log"

	"github.com/aparnasukesh/api-gateway/config"
	"github.com/aparnasukesh/api-gateway/internals/boot"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("error loading config file")
	}
	boot.Start(cfg)
}
