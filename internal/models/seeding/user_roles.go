package seeding

import (
	"fmt"

	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/models"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/utility"
)

type SeedData struct {
	Db *repository.Database
}

var roleList = []string{"admin", "user"}
var roleCount int64

func (s *SeedData) CheckRoleExists() bool {
	role := &models.Role{}
	if err := s.Db.Pdb.DB.Model(&role).Where("name IN  (?) ", roleList).Count(&roleCount).Error; err != nil {
		panic(err)
	}
	if roleCount > 0 {
		fmt.Println("user role already seeded skipping...")
		return true
	}
	return false

}

func (s *SeedData) SeedRole() {
	var permission models.Permission

	roleDescription := []string{
		"Manages users, theaters, movies, and bookings. Oversees content, handles reports, and ensures smooth platform operations",
		"Can browse movies, book tickets, view showtimes, and manage personal booking",
	}
	if roleCount > 0 {
		fmt.Println("user role already seeded skipping...")
		return
	}
	for ind, roleName := range roleList {
		role := models.Role{
			ID:          utility.GenerateUUID(),
			Name:        roleName,
			Description: roleDescription[ind],
		}
		if roleName == "admin" {
			permission = models.Permission{
				ID:     utility.GenerateUUID(),
				RoleID: role.ID,
				Role:   role,
				PermissionList: models.PermissionList{
					CanCreateMovies:        true,
					CanDeleteMovie:         true,
					CanUpdateMovie:         true,
					CanGetMovies:           true,
					CanBanUsers:            true,
					CanDeleteUsers:         true,
					CanGetUsers:            true,
					CanCreateUsers:         true,
					CanBookSeats:           true,
					CanCreateSeats:         true,
					CanRemoveSeats:         true,
					CanUpdateSeats:         true,
					CanCancelBooking:       true,
					CanGetBookings:         true,
					CanGetRoles:            true,
					CanCreateRoles:         true,
					CanUpdateRoles:         true,
					CanDeleteRoles:         true,
					CanCreateTickets:       true,
					CanBuyTickets:          true,
					CanDeleteTickets:       true,
					CanViewTickets:         true,
					CanCancelTickets:       true,
					CanVerifyTickets:       true,
					CanBanTickets:          true,
					CanRefundTickets:       true,
					CanRequestTicketRefund: false,
				},
			}
		} else if roleName == "user" {
			permission = models.Permission{
				ID:     utility.GenerateUUID(),
				RoleID: role.ID,
				Role:   role,
				PermissionList: models.PermissionList{
					CanCreateMovies:        false,
					CanDeleteMovie:         false,
					CanUpdateMovie:         false,
					CanGetMovies:           true,
					CanBanUsers:            false,
					CanDeleteUsers:         false,
					CanGetUsers:            true,
					CanCreateUsers:         false,
					CanBookSeats:           true,
					CanCreateSeats:         false,
					CanRemoveSeats:         false,
					CanUpdateSeats:         false,
					CanCancelBooking:       true,
					CanGetBookings:         true,
					CanGetRoles:            false,
					CanCreateRoles:         false,
					CanUpdateRoles:         false,
					CanDeleteRoles:         false,
					CanCreateTickets:       true,
					CanBuyTickets:          true,
					CanDeleteTickets:       false,
					CanViewTickets:         true,
					CanCancelTickets:       true,
					CanVerifyTickets:       false,
					CanBanTickets:          false,
					CanRefundTickets:       false,
					CanRequestTicketRefund: true,
				},
			}
		}
		if err := s.Db.Pdb.DB.Create(&role).Error; err != nil {
			panic(err)
		}
		if err := s.Db.Pdb.DB.Create(&permission).Error; err != nil {
			panic(err)
		}
	}

}
