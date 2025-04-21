package seeding

import (
	"fmt"

	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/models"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/utility"
)

var Genre = []string{
	"Action",
	"Adventure",
	"Animation",
	"Biography",
	"Comedy",
	"Crime",
	"Documentary",
	"Drama",
	"Family",
	"Fantasy",
	"History",
	"Horror",
	"Music",
	"Musical",
	"Mystery",
	"Romance",
	"Sci-Fi",
	"Short",
	"Sport",
	"Superhero",
	"Thriller",
	"War",
	"Western",
}

var GenreCount int64

func (s *SeedData) CheckGenreExists() bool {
	s.Db.Pdb.DB.Model(&models.Genre{}).Where("name IN (?)", Genre).Count(&GenreCount)
	genreCount := int64(len(Genre))
	return GenreCount == genreCount
}

// SeedGenre seeds the database with genre data
func (s *SeedData) SeedGenre() {
	for _, genreName := range Genre {
		genre := models.Genre{
			ID:   utility.GenerateUUID(),
			Name: genreName,
		}
		if err := s.Db.Pdb.DB.Create(&genre).Error; err != nil {
			continue // skip if error occurs
		}
	}
	fmt.Println("Genre seeded successfully...")
}
