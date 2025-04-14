package activate

import (
	"fmt"
	"log"
	"net/http"

	"alertly/internal/auth"
	"alertly/internal/database"
	"alertly/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ActivateAccount(c *gin.Context) {
	var user ActivateAccountRequest
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

	err := service.ActivateAccount(user)

	if err != nil {
		log.Printf("Error al activar  el usuario: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error al registrar usuario", err.Error())
		return
	}

	// fmt.Println("email: ", user.Email)
	authRepo := auth.NewRepository(database.DB)
	authService := auth.NewService(authRepo)
	userAuth, err := authRepo.GetUserByEmail(user.Email)
	// fmt.Println("debug: ", userAuth)
	if err != nil {
		fmt.Println(err.Error())
		response.Send(c, http.StatusInternalServerError, true, "Error al obtener usuario", err.Error())
		return
	}
	tokenResp, err := authService.GenerateSessionToken(userAuth)
	if err != nil {
		response.Send(c, http.StatusInternalServerError, true, "Error al generar token", err.Error())
		return
	}

	response.Send(c, http.StatusOK, false, "Usuario registrado exitosamente", tokenResp)
}
