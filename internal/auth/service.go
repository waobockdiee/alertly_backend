package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GenerateSessionToken(user User) (TokenResponse, error)
	AuthenticateUser(email, password string) (User, error)
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
		return TokenResponse{}, errors.New("error al generar token")
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

// func AccountID(tokenString string) (int64, error) {
// 	claims := &Claims{}
// 	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 		return jwtSecret, nil
// 	})
// 	if err != nil {
// 		return 0, err
// 	}
// 	if !token.Valid {
// 		return 0, errors.New("token inválido")
// 	}
// 	return claims.AccountID, nil
// }
