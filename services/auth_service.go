package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/soa-team-11/auth-service/internal/repos"
	"github.com/soa-team-11/auth-service/models"
)

type AuthService struct {
	userRepo repos.UserRepo
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repos.NewUserRepo(),
	}
}

func (s *AuthService) Register(user models.User) (*models.User, error) {

	// Check if user is valid
	if !user.IsValid() {
		return nil, fmt.Errorf("user not valid")
	}

	// Check if roles are allowed
	if user.Role == models.Admin {
		return nil, fmt.Errorf("role '%s' is not allowed for registration", user.Role)
	}

	// Check if username already exists
	if retrieved_user, _ := s.userRepo.GetByUsername(user.Username); retrieved_user != nil {
		return nil, fmt.Errorf("username '%s' already taken", user.Username)
	}

	user.UserID = uuid.New()
	created_user, err := s.userRepo.Create(user)

	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	return created_user, nil
}
