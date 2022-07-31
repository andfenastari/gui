package main

import (
	"github.com/andfenastari/gui/backend/x"
	"log"
)

func main() {
	var b x.Backend
	err := b.Init()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer b.Close()
}
