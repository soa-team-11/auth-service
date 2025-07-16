package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type User struct {
	UserID   uuid.UUID `json:"user_id" bson:"user_id"`
	Username string    `json:"username" bson:"username"`
	Email    string    `json:"email" bson:"email"`
	Password string    `json:"password,omitempty" bson:"password"`
	Role     UserRole  `json:"role" bson:"role"`
}

type UserRole string

const (
	Tourist   UserRole = "tourist"
	Tourguide UserRole = "tourguide"
	Admin     UserRole = "admin"
)

func (u *User) IsValid() bool {
	return u.Username != "" && u.Email != "" && u.Password != "" && isRole(u.Role)
}

func (u *User) ToJSON() []byte {
	r, _ := json.Marshal(u)
	return r
}

func UserFromJSON(jsonString []byte) (*User, error) {
	var user User
	err := json.Unmarshal(jsonString, &user)
	return &user, err
}

func isRole(role UserRole) bool {
	return role == Tourist || role == Tourguide || role == Admin
}
