package signup

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(user User) (User, error, string)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) RegisterUser(user User) (User, error, string) {
	code, err := generateActivationCode()

	if err != nil {
		return User{}, err, ""
	}
	user.ActivationCode = code

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("error al hashear la contraseña: %v", err), ""
	}
	user.Password = string(hashedPassword)
	user.Nickname = generateNickName(user.FirstName, user.BirthYear)

	id, err := s.repo.InsertUser(user)
	if err != nil {
		return User{}, err, ""
	}
	res, err := s.repo.GetUserByID(id)
	return res, err, user.ActivationCode
}

func generateActivationCode() (string, error) {
	code := ""
	for i := 0; i < 5; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code += n.String()
	}
	return code, nil
}

func generateNickName(firstName, yearBirth string) string {
	var nickname string = firstName + yearBirth
	return nickname
}
