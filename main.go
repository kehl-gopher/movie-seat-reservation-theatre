package main

import (
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/pkg/router"
)

func main() {
	config := env.SetEnv()
	router.Router(config)
}
