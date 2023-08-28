package main

import (
	"log"

	"github.com/AlexCorn999/notes/internal/app/apiserver"
)

func main() {
	s := apiserver.NewAPI()
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
