package main

import (
	"github.com/andfenastari/gui/backend/wayland"
	"log"
)

func main() {
	b, err := wayland.NewBackend()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer b.Close()
}
