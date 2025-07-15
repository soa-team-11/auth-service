package models

type User struct {
	UserID   string
	Username string
	Email    string
	Password string
	Role     UserRole
}

type UserRole string

const (
	Tourist   UserRole = "tourist"
	Tourguide UserRole = "tourguide"
	Admin     UserRole = "admin"
)
