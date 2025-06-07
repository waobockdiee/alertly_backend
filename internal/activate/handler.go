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
		response.Send(c, http.StatusBadRequest, true, "Invalid inputs. Please check the information and try again.", err.Error())
		return
	}

	if err := validate.Struct(user); err != nil {
		log.Printf("Error de validación: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Some fields are invalid. Please review the form and try again.", err.Error())
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	err := service.ActivateAccount(user)

	if err != nil {
		log.Printf("Error al activar  el usuario: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Something went wrong while creating your account. Please try again later.", err.Error())
		return
	}

	// fmt.Println("email: ", user.Email)
	authRepo := auth.NewRepository(database.DB)
	authService := auth.NewService(authRepo)
	userAuth, err := authRepo.GetUserByEmail(user.Email)
	// fmt.Println("debug: ", userAuth)
	if err != nil {
		fmt.Println(err.Error())
		response.Send(c, http.StatusInternalServerError, true, "We couldn’t load your user data. Please try again in a moment.", err.Error())
		return
	}
	tokenResp, err := authService.GenerateSessionToken(userAuth)
	if err != nil {
		response.Send(c, http.StatusInternalServerError, true, "We couldn’t start your session. Please try again shortly.", err.Error())
		return
	}

	response.Send(c, http.StatusOK, false, "Your account has been created successfully. Welcome aboard!", tokenResp)
}
