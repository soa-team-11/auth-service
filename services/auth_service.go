package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/soa-team-11/auth-service/api/external"
	"github.com/soa-team-11/auth-service/internal/repos"
	"github.com/soa-team-11/auth-service/models"
)

type AuthService struct {
	userRepo            repos.UserRepo
	stakeholdersService external.StakeholderService
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo:            repos.NewUserRepo(),
		stakeholdersService: external.StakeholderService{},
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
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	_, err = s.stakeholdersService.CreateProfile(created_user.UserID)

	if err != nil {
		return nil, fmt.Errorf("failed to create profile: %w", err)
	}

	return created_user, nil
}
