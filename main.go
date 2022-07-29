package main

import (
	"github.com/andfenastari/gui/backend/x"
	"log"
)

func main() {
	b, err := x.NewBackend()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer b.Close()
}
