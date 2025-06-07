package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GenerateSessionToken(user User) (TokenResponse, error)
	AuthenticateUser(email, password string) (User, error)
	CheckPasswordMatch(password, email string, accountID int64) (PasswordMatch, error)
}

type service struct {
	repo Repository
}

var jwtSecret = []byte("mi_clave_secreta")

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// RegisterUser registra un nuevo usuario: inserta el registro y lo retorna con el ID asignado.
func (s *service) GenerateSessionToken(user User) (TokenResponse, error) {
	expirationTime := time.Now().Add(72 * time.Hour)
	claims := &Claims{
		AccountID: user.AccountID,
		Email:     user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return TokenResponse{}, errors.New("we couldn’t start your session. Please try again shortly")
	}

	return TokenResponse{Token: tokenString}, nil
}

func (s *service) AuthenticateUser(email, password string) (User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return User{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return User{}, errors.New("invalid credentials")
	}
	return user, nil
}

func (s *service) CheckPasswordMatch(password, email string, accountID int64) (PasswordMatch, error) {
	var pm PasswordMatch
	var err error
	pm, err = s.repo.GetUserById(accountID)
	fmt.Println("DEBUG:", pm)

	if pm.Email == "" {
		return pm, errors.New("wrong user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(pm.Password), []byte(password)); err != nil {
		return pm, errors.New("wrong password")
	}

	return pm, err
}
