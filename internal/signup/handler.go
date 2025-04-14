package signup

import (
	"log"
	"net/http"

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
		response.Send(c, http.StatusBadRequest, true, "Datos de entrada inválidos", err.Error())
		return
	}

	if err := validate.Struct(user); err != nil {
		log.Printf("Error de validación: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Datos de formulario no válidos", err.Error())
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	registeredUser, err, code := service.RegisterUser(user)
	if err != nil {
		log.Printf("Error al registrar el usuario: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error al registrar usuario", err.Error())
		return
	}
	emails.Send(user.Email, "Activation Code", code)
	response.Send(c, http.StatusOK, false, "Usuario registrado exitosamente", registeredUser)
}
