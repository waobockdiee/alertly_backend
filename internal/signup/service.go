package signup

import (
	cryptorand "crypto/rand"
	"fmt"
	"log"
	"math/big"
	mathrand "math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type Service interface {
	RegisterUser(user User) (User, string, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func randStringMath(n int) string {
	mathrand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[mathrand.Intn(len(letters))]
	}
	return string(b)
}

func (s *service) RegisterUser(user User) (User, string, error) {
	code, err := generateActivationCode()

	if err != nil {
		log.Printf("ERROR 1: %v", err)
		return User{}, "", err
	}
	user.ActivationCode = code

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("ERROR 2: %v", err)
		return User{}, "", fmt.Errorf("error al hashear la contraseña: %v", err)
	}

	nicknameUniqueString := randStringMath(3)
	user.Password = string(hashedPassword)
	user.Nickname = generateNickName(user.FirstName, nicknameUniqueString)

	id, err := s.repo.InsertUser(user)
	if err != nil {
		log.Printf("ERROR 3: %v", err)
		return User{}, "", err
	}
	res, err := s.repo.GetUserByID(id)

	if err != nil {
		log.Printf("ERROR 4: %v", err)
		return res, "", err
	}
	return res, user.ActivationCode, err
}

func generateActivationCode() (string, error) {
	code := ""
	for i := 0; i < 5; i++ {
		n, err := cryptorand.Int(cryptorand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code += n.String()
	}
	return code, nil
}

func generateNickName(firstName, randString string) string {
	var nickname string = firstName + "_" + randString
	return nickname
}
