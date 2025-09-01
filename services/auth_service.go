package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/soa-team-11/auth-service/api/external"
	"github.com/soa-team-11/auth-service/internal/repos"
	"github.com/soa-team-11/auth-service/models"
	"github.com/soa-team-11/auth-service/utils/jwt"
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

type LoginDTO struct {
	UserID   uuid.UUID       `json:"user_id" bson:"user_id"`
	Username string          `json:"username" bson:"username"`
	Role     models.UserRole `json:"role" bson:"role"`
	Token    string          `json:"token" bson:"token"`
}

func (s *AuthService) Login(username string, password string) (*LoginDTO, error) {
	retrieved_user, _ := s.userRepo.GetByUsername(username)

	if retrieved_user == nil {
		return nil, fmt.Errorf("user '%s' not found", username)
	}

	if retrieved_user.Password != password {
		return nil, fmt.Errorf("incorrect password")
	}

	if retrieved_user.Blocked {
		return nil, fmt.Errorf("user is blocked")
	}

	claims := map[string]interface{}{
		"user_id":  retrieved_user.UserID,
		"username": retrieved_user.Username,
		"role":     retrieved_user.Role,
		"exp":      time.Now().Add(time.Hour * 48).Unix(), // expires in 48h
	}

	_, tokenString, _ := jwt.GetTokenAuth().Encode(claims)

	return &LoginDTO{
		UserID:   retrieved_user.UserID,
		Username: retrieved_user.Username,
		Role:     retrieved_user.Role,
		Token:    tokenString}, nil
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
