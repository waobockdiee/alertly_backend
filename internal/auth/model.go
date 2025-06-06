package auth

import (
	"database/sql"

	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	AccountID           int64          `json:"account_id" validate:"required"`
	Email               string         `json:"email" validate:"required,email"`
	FirstName           string         `json:"first_name" validate:"required"`
	LastName            string         `json:"last_name" validate:"required"`
	Password            string         `json:"password"`
	PhoneNumber         sql.NullString `json:"phone_number"`
	BirthYear           int            `json:"birth_year"`
	BirthMonth          int            `json:"birth_month"`
	BirthDay            int            `json:"birth_day"`
	Status              string         `json:"status" validate:"required"`
	IsPremium           bool           `json:"is_premium" validate:"required"`
	HasFinishedTutorial bool           `json:"has_finished_tutorial"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type Claims struct {
	AccountID int64  `json:"account_id"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}

type PasswordMatch struct {
	AccountID int64  `json:"account_id"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Match     bool
}
