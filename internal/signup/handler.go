package signup

import (
	"log"
	"net/http"
	"strings"

	"alertly/internal/database"
	"alertly/internal/emails"
	"alertly/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func RegisterUserHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("Error al decodificar JSON: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Invalid input. Please check the data and try again.", nil)
		return
	}

	// Log de debug para ver qué datos se están recibiendo
	log.Printf("DEBUG: Datos recibidos en signup - Email: %s, FirstName: %s, LastName: %s, Password: %s (length: %d)",
		user.Email, user.FirstName, user.LastName,
		strings.Repeat("*", len(user.Password)), len(user.Password))

	if err := validate.Struct(user); err != nil {
		log.Printf("Error de validación: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Some form fields are invalid. Please review and try again.", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	registeredUser, code, err := service.RegisterUser(user)
	if err != nil {
		log.Printf("Error creating account: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "We couldn’t register your account. Please try again later.", nil)
		return
	}
	emails.SendTemplate(user.Email, "Activate your Alertly account", "new_account_activation_code", map[string]string{
		"Name": user.FirstName,
		"Code": code,
	})
	response.Send(c, http.StatusOK, false, "Your account has been created successfully. Please check your email to activate it.", registeredUser)
}
