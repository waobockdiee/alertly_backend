package auth

import (
	"alertly/internal/database"
	"alertly/internal/response"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SignIn(c *gin.Context) {
	var req LoginRequest

	// Vincula el JSON recibido al struct LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Send(c, http.StatusBadRequest, true, "Invalid input", err.Error())
		return
	}

	// (Opcional) Valida el request con el validator si lo requieres:
	// if err := validate.Struct(req); err != nil {
	//     response.Send(c, http.StatusBadRequest, true, "Validation error", err.Error())
	//     return
	// }

	// Crea el repositorio y el servicio utilizando la DB inicializada
	repo := NewRepository(database.DB)
	service := NewService(repo)

	// Autentica al usuario
	user, err := service.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		response.Send(c, http.StatusUnauthorized, true, "Invalid credentials", err.Error())
		return
	}

	// Genera el token de sesión
	tokenResp, err := service.GenerateSessionToken(user)
	if err != nil {
		response.Send(c, http.StatusInternalServerError, true, "Error generating token", err.Error())
		return
	}

	fmt.Println(tokenResp.Token)

	response.Send(c, http.StatusOK, false, "Success", tokenResp)
}

func ValidateSession(c *gin.Context) {
	response.Send(c, http.StatusOK, false, "Token válido", nil)
}
func CheckPasswordMatch(c *gin.Context) {
	var pM PasswordMatch
	var err error

	if err = c.ShouldBindJSON(&pM); err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Invalid input", nil)
		return
	}

	pM.AccountID, err = GetUserFromContext(c)

	// fmt.Println("ID:", pM.AccountID)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "We couldn’t verify your session. Please log in again.", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	pM, err = service.CheckPasswordMatch(pM.Password, pM.Email, pM.AccountID)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error validating password. Please try later", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "success", pM.Email)
}
