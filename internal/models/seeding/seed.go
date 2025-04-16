package seeding

import "github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"

func StartSeeding(db *repository.Database) {
	s := &SeedData{
		Db: db,
	}

	if !s.CheckRoleExists() {
		s.SeedRole()
	}
}
