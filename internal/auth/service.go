package auth

import (
	"errors"
	"fmt"
	"os"
	"strconv"
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

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// RegisterUser registra un nuevo usuario: inserta el registro y lo retorna con el ID asignado.
func (s *service) GenerateSessionToken(user User) (TokenResponse, error) {
	// Configuración de duración del token (90 días por defecto)
	expHours := 90 * 24 // 90 días = 2160 horas
	if envHours := os.Getenv("JWT_EXPIRATION_HOURS"); envHours != "" {
		if hours, err := strconv.Atoi(envHours); err == nil && hours > 0 {
			expHours = hours
		}
	}
	expirationTime := time.Now().Add(time.Duration(expHours) * time.Hour)
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
	// log.Printf("DEBUG: Token generated for user %s: %s", user.Email, tokenString)
	if err != nil {
		return TokenResponse{}, errors.New("we couldn’t start your session. Please try again shortly")
	}

	return TokenResponse{Token: tokenString}, nil
}

func (s *service) AuthenticateUser(email, password string) (User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		fmt.Printf("❌ [AUTH-SVC] GetUserByEmail failed: %v\n", err)
		return User{}, err
	}
	fmt.Printf("🔑 [AUTH-SVC] Comparing password for %s (hash starts with: %s...)\n", email, user.Password[:10])
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		fmt.Printf("❌ [AUTH-SVC] Password mismatch for %s: %v\n", email, err)
		return User{}, errors.New("invalid credentials")
	}
	fmt.Printf("✅ [AUTH-SVC] Password matches for %s\n", email)

	// Verificar estado de la cuenta
	switch user.Status {
	case "active":
		// OK, continuar con login
	case "pending_activation":
		return User{}, errors.New("your account is not activated yet. Please check your email for the activation code")
	case "inactive":
		return User{}, errors.New("your account has been deactivated. Please contact support at support@alertly.ca")
	case "blocked":
		return User{}, errors.New("your account has been suspended. Please contact support at support@alertly.ca")
	default:
		return User{}, errors.New("invalid account status. Please contact support at support@alertly.ca")
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
