package repos

import (
	"context"

	"github.com/soa-team-11/auth-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	mg "github.com/soa-team-11/auth-service/internal/providers/mongo"
)

type UserRepo interface {
	GetAll() ([]models.User, error)
	GetByUsername(username string) (*models.User, error)
	Create(user models.User) (*models.User, error)
	Update(user models.User) (*models.User, error)
	Delete(user models.User) bool
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
	return nil, nil
}

func (r *UserRepoImpl) Delete(user models.User) bool {
	return false
}
