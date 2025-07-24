package services

import (
	"fmt"

	"github.com/soa-team-11/auth-service/internal/repos"
	"github.com/soa-team-11/auth-service/models"
)

type AccountService struct {
	userRepo repos.UserRepo
}

func NewAccountService() *AccountService {
	return &AccountService{
		userRepo: repos.NewUserRepo(),
	}
}

func (s *AccountService) ListAccounts() ([]models.User, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error retrieving users: %s", err.Error())
	}

	// Remove passwords
	for i := range users {
		users[i].Password = ""
	}

	if len(users) == 0 {
		return []models.User{}, nil
	}

	return users, nil
}

// ToggleBlockUser toggles the blocked status of a user by MongoDB _id and returns the new status (true if blocked, false if unblocked)
func (s *AccountService) ToggleBlockUser(userID string) (bool, error) {
	return s.userRepo.(*repos.UserRepoImpl).ToggleBlockUser(userID)
}
