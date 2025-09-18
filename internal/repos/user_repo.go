package repos

import (
	"context"
	"fmt"

	"github.com/soa-team-11/auth-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/google/uuid"
	mg "github.com/soa-team-11/auth-service/internal/providers/mongo"
)

type UserRepo interface {
	GetAll() ([]models.User, error)
	GetByUsername(username string) (*models.User, error)
	Create(user models.User) (*models.User, error)
	Update(user models.User) (*models.User, error)
	Delete(user models.User) bool
	DeleteByID(userID string) error
}

type UserRepoImpl struct {
	users *mongo.Collection
}

func NewUserRepo() *UserRepoImpl {
	return &UserRepoImpl{
		users: mg.GetDatabase().Collection("users"),
	}
}

func (r *UserRepoImpl) GetAll() ([]models.User, error) {
	cur, err := r.users.Find(context.Background(), bson.M{})

	if err != nil {
		return nil, err
	}

	var users []models.User
	err = cur.All(context.Background(), &users)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepoImpl) GetByUsername(username string) (*models.User, error) {
	var user models.User

	err := r.users.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepoImpl) Create(user models.User) (*models.User, error) {
	_, err := r.users.InsertOne(context.Background(), user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepoImpl) Update(user models.User) (*models.User, error) {
	// Update the user document in MongoDB by user_id
	filter := bson.M{"user_id": user.UserID}
	update := bson.M{"$set": user}
	_, err := r.users.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ToggleBlockUser toggles the Blocked field for the specified user and returns the new status
func (r *UserRepoImpl) ToggleBlockUser(userID string) (bool, error) {
	uuidVal, err := uuid.Parse(userID)
	if err != nil {
		return false, err // invalid UUID
	}
	var user models.User
	err = r.users.FindOne(context.Background(), bson.M{"user_id": uuidVal}).Decode(&user)
	if err != nil {
		return false, err
	}
	newStatus := !user.Blocked
	update := bson.M{"$set": bson.M{"blocked": newStatus}}
	_, err = r.users.UpdateOne(context.Background(), bson.M{"user_id": uuidVal}, update)
	if err != nil {
		return user.Blocked, err
	}
	return newStatus, nil
}

// IsUserBlocked checks if the user is blocked using user_id
func (r *UserRepoImpl) IsUserBlocked(userID string) (bool, error) {
	uuidVal, err := uuid.Parse(userID)
	if err != nil {
		return false, err // invalid UUID
	}
	var user models.User
	err = r.users.FindOne(context.Background(), bson.M{"user_id": uuidVal}).Decode(&user)
	if err != nil {
		return false, err
	}
	return user.Blocked, nil
}

func (r *UserRepoImpl) Delete(user models.User) bool {
	return false
}

func (r *UserRepoImpl) DeleteByID(userID string) error {
	ctx := context.Background()

	uuidVal, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	result, err := r.users.DeleteOne(ctx, bson.M{"user_id": uuidVal})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("user with ID %s not found", userID)
	}

	return nil
}
