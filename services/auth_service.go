package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/soa-team-11/auth-service/api/external"
	"github.com/soa-team-11/auth-service/internal/repos"
	"github.com/soa-team-11/auth-service/models"
	"github.com/soa-team-11/auth-service/utils/jwt"
	"go.opentelemetry.io/otel"
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

func (s *AuthService) Login(ctx context.Context, username string, password string) (*LoginDTO, error) {
	tracer := otel.Tracer("auth-service")
	_, span := tracer.Start(ctx, "AuthService.Login")
	defer span.End()

	retrieved_user, err := s.userRepo.GetByUsername(username)

	if err != nil {
		span.RecordError(err)
	}

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
	span.End()

	return &LoginDTO{
		UserID:   retrieved_user.UserID,
		Username: retrieved_user.Username,
		Role:     retrieved_user.Role,
		Token:    tokenString}, nil
}

func (s *AuthService) Register(ctx context.Context, user models.User) (*models.User, error) {
	tracer := otel.Tracer("auth-service")
	_, span := tracer.Start(ctx, "AuthService.Register")
	defer span.End()

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
		span.RecordError(err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	_, err = s.stakeholdersService.CreateProfile(created_user.UserID)

	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to create profile: %w", err)
	}

	span.End()
	return created_user, nil
}

func (s *AuthService) DeleteUser(userID string) error {
	return s.userRepo.DeleteByID(userID)
}
