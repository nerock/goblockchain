package main

import (
	"log"

	"github.com/nerock/goblockchain/internal/cache"
	"github.com/nerock/goblockchain/internal/server"
)

func main() {
	c := cache.New()
	srv := server.New(8080, c)

	log.Fatal(srv.Run())
}
