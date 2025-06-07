package myplaces

import (
	"alertly/internal/auth"
	"alertly/internal/database"
	"alertly/internal/response"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

var validate = validator.New()

func Get(c *gin.Context) {

	accountId := c.Param("account_id")

	if accountId == "" {
		response.Send(c, http.StatusBadRequest, true, "Bad request", nil)
		return
	}

	id, err := strconv.Atoi(accountId)

	if err != nil {
		response.Send(c, http.StatusBadRequest, true, "Error converting param", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	result, err := service.Get(id)

	if err != nil {
		log.Printf("We couldn’t load the categories. Please try again later: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "We couldn’t load the categories. Please try again later.", err.Error())
		return
	}

	response.Send(c, http.StatusOK, false, "Success", result)
}

// hide
func Add(c *gin.Context) {

	accountID, err := auth.GetUserFromContext(c)
	if err != nil {
		response.Send(c, http.StatusInternalServerError, true, "error", err.Error())
		return
	}

	var myPlace MyPlaces
	if err := c.BindJSON(&myPlace); err != nil {
		log.Printf("JSON error: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Wrong data in", err.Error())
		return
	}

	if err := validate.Struct(myPlace); err != nil {
		log.Printf("Error de validación: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Bad request", err.Error())
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	myPlace.AccountId = accountID
	result, err := service.Add(myPlace)

	if err != nil {
		response.Send(c, http.StatusInternalServerError, true, "Error saving place. Please try later", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "Saved! Incident alerts on.", result)
}

// update or remove(hide changing status = 'inactive')
func Update(c *gin.Context) {
	var myPlace MyPlaces
	if err := c.BindJSON(&myPlace); err != nil {
		log.Printf("JSON error: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Wrong data in", err.Error())
		return
	}

	if err := validate.Struct(myPlace); err != nil {
		log.Printf("Error de validación: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Bad request", err.Error())
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	err := service.Update(myPlace)
	if err != nil {
		response.Send(c, http.StatusBadRequest, true, "Bad request", "")
		return
	}

	response.Send(c, http.StatusOK, false, "success", myPlace)

}

func FullUpdate(c *gin.Context) {
	var err error
	accountID, err := auth.GetUserFromContext(c)
	if err != nil {
		response.Send(c, http.StatusInternalServerError, true, "error", err.Error())
		return
	}

	var myPlace MyPlaces
	if err := c.BindJSON(&myPlace); err != nil {
		log.Printf("JSON error: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Wrong data in", err.Error())
		return
	}

	if err := validate.Struct(myPlace); err != nil {
		log.Printf("Error de validación: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Bad request", err.Error())
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	myPlace.AccountId = accountID
	err = service.FullUpdate(myPlace)
	if err != nil {
		response.Send(c, http.StatusBadRequest, true, "Bad request", "")
		return
	}

	response.Send(c, http.StatusOK, false, "success", myPlace)

}

func GetByAccountId(c *gin.Context) {
	accountID, err := auth.GetUserFromContext(c)

	if err != nil {
		response.Send(c, http.StatusInternalServerError, true, "error", err.Error())
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	data, err := service.GetByAccountId(accountID)

	if err != nil {
		response.Send(c, http.StatusInternalServerError, true, "Error fetching data", err.Error())
		return
	}

	response.Send(c, http.StatusOK, false, "success", data)

}

func GetById(c *gin.Context) {
	accountID, err := auth.GetUserFromContext(c)

	if err != nil {
		response.Send(c, http.StatusInternalServerError, true, "error", err.Error())
		return
	}
	aflId := c.Param("afl_id")
	formattedAflId, err := strconv.ParseInt(aflId, 10, 64)
	if err != nil {
		response.Send(c, http.StatusInternalServerError, true, "Error converting data", err)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	data, err := service.GetById(accountID, formattedAflId)

	if err != nil {
		response.Send(c, http.StatusInternalServerError, true, "Error fetching data", err.Error())
		return
	}

	response.Send(c, http.StatusOK, false, "success", data)
}

func Delete(c *gin.Context) {
	accountID, err := auth.GetUserFromContext(c)

	fmt.Println("error getuserfromcontext", err)
	if err != nil {
		response.Send(c, http.StatusInternalServerError, true, "error", err.Error())
		return
	}

	aflId := c.Param("afl_id")
	formattedAflID, err := strconv.ParseInt(aflId, 10, 64)
	fmt.Println("error parseint", err)
	if err != nil {
		response.Send(c, http.StatusInternalServerError, true, "Error converting data", err)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	err = service.Delete(accountID, formattedAflID)
	fmt.Println("error deleting", err)
	if err != nil {
		response.Send(c, http.StatusBadRequest, true, "Bad request", "")
		return
	}

	response.Send(c, http.StatusOK, false, "success", "")
}
