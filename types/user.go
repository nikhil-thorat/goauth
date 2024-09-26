package types

import (
	"time"
)

type User struct {
	ID               int    `json:"id"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	IsVerified       bool
	VerificationCode string
	CreatedAt        time.Time `json:"createdAt"`
}

type UserStore interface {
	CreateUser(user User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(userID int) (*User, error)
	UpdateUserDetails(userId int, firstName string, lastName string) error
	UpdateUserPassword(userID int, password string) error
	UpdateUserVerificationStatus(userID int, isVerified bool) error
	DeleteUser(userID int) error
}

type RegisterUserPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}
type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type VerifyUserPayload struct {
	Email            string `json:"email" validate:"required,email"`
	VerificationCode string `json:"verificationCode" validate:"required"`
}

type UpdateUserDetailsPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type UpdateUserPasswordPayload struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}
