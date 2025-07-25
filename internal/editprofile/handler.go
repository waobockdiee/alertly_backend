package editprofile

import (
	"alertly/internal/auth"
	"alertly/internal/database"
	"alertly/internal/media"
	"alertly/internal/response"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func GetMyProfile(c *gin.Context) {
	var account Account
	var err error

	account.AccountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "Unauthorized", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	data, err := service.GetAccountByID(account.AccountID)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error saving code", err)
		return
	}

	response.Send(c, http.StatusOK, false, "success", data)

}

func GenerateCodeUpdateEmail(c *gin.Context) {
	var account Account
	var err error

	if err = c.ShouldBindJSON(&account); err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Invalid entry data", nil)
		return
	}

	account.AccountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error editprofile/handler.go GenerateCodeUpdateEmail: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "Unauthorized", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	err = service.GenerateCodeUpdateEmail(account.AccountID)

	if err != nil {
		log.Printf("Error editprofile/handler.go GenerateCodeUpdateEmail: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error saving code", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "Success", nil)

}

func ValidateUpdateEmailCode(c *gin.Context) {
	var account Account
	var err error
	var match bool

	if err = c.ShouldBindJSON(&account); err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Invalid entry data", nil)
		return
	}

	account.AccountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "error", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	match, err = service.ValidateUpdateEmailCode(account.AccountID, account.Code)
	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "error", nil)
		return
	}

	if !match {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "error", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "success", nil)

}

func UpdateEmail(c *gin.Context) {

	var account Account
	var err error

	if err := c.ShouldBindJSON(&account); err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Invalid entry data", nil)
		return
	}

	account.AccountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "error", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	err = service.UpdateEmail(account.AccountID, account.Email)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "error", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "success", nil)

}

func UpdatePassword(c *gin.Context) {
	var account Account
	var err error

	if err = c.ShouldBindJSON(&account); err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Invalid entry data", nil)
		return
	}

	account.AccountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "error", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	err = service.CheckPasswordMatch(account.Password, account.NewPassword, account.AccountID)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Wrong password. Please try again", nil)
		return
	}

	err = service.UpdatePassword(account.AccountID, account.NewPassword)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Error updating password. Please try later", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "success", nil)
}

func UpdateBirthDate(c *gin.Context) {
	var account Account
	var err error

	if err = c.ShouldBindJSON(&account); err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Invalid format data", nil)
		return
	}

	account.AccountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "We couldn’t verify your session. Please log in again", nil)
		return
	}

	if account.BirthYear == "" || account.BirthMonth == "" || account.BirthDay == "" {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Invalid format data", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	err = service.UpdateBirthDate(account.AccountID, account.BirthYear, account.BirthMonth, account.BirthDay)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error while updating birthdate information. Please try later", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "success", nil)
}

func UpdateIsPremium(c *gin.Context) {
	var account Account
	var err error

	if err = c.ShouldBindJSON(&account); err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Invalid format data", nil)
		return
	}

	account.AccountID, err = auth.GetUserFromContext(c)
	account.IsPremium = true

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "error", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	err = service.UpdateIsPremium(account.AccountID, account.IsPremium)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error while updating birthdate information. Please try later", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "success", nil)
}

func UpdateIsPrivateProfile(c *gin.Context) {
	var account Account
	var err error

	if err = c.ShouldBindJSON(&account); err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Invalid format data", nil)
		return
	}

	account.AccountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "error", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	err = service.UpdateIsPrivateProfile(account.AccountID, account.IsPrivateProfile)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error while updating status information. Please try later", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "success", nil)
}

func UpdateFullName(c *gin.Context) {
	var account Account
	var err error

	if err = c.ShouldBindJSON(&account); err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Invalid format data", nil)
		return
	}

	account.AccountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "error", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	err = service.UpdateFullName(account.AccountID, account.FirstName, account.LastName)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error while updating status information. Please try later", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "success", nil)
}

func UpdatePhoneNumber(c *gin.Context) {
	var account Account
	var err error

	if err = c.ShouldBindJSON(&account); err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Invalid format data", nil)
		return
	}

	account.AccountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "error", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	err = service.UpdatePhoneNumber(account.AccountID, account.PhoneNumber)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error while updating phone number information. Please try later", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "success", nil)
}

func UpdateNickname(c *gin.Context) {
	var account Account
	var err error

	if err = c.ShouldBindJSON(&account); err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Invalid format data", nil)
		return
	}

	account.AccountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "error", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	err = service.UpdateNickname(account.AccountID, account.NickName)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error while updating nickname information. Please try later", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "success", nil)
}

func UpdateThumbnail(c *gin.Context) {
	// var media Media
	var accountID int64
	var err error
	// var tmpFilePath string
	// var uploadDir string

	accountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "Unauthorized. You need login with your account", nil)
		return
	}

	file, header, err := c.Request.FormFile("file")

	if err != nil {
		log.Printf("Error retrieving file: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Error fetching file", nil)
		return
	}

	defer file.Close()

	ext := filepath.Ext(header.Filename)
	tmpFile, err := os.CreateTemp("", "orig_*"+ext)

	if err != nil {
		log.Printf("Error creating temp file on /account/edit/picture: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error creating folder", err.Error())
		return
	}

	tmpFilePath := tmpFile.Name()

	if _, err := io.Copy(tmpFile, file); err != nil {
		log.Printf("Error saving temp file: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error saving tmpl file", err.Error())
		return
	}

	tmpFile.Close()
	defer os.Remove(tmpFilePath)

	uploadDir := "uploads/profile"

	processedFilePath, err := media.ProcessImage(tmpFilePath, uploadDir)

	if err != nil {
		log.Printf("Error processing image: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error processing file", err.Error())
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	err = service.UpdateThumbnail(accountID, processedFilePath)

	if err != nil {
		log.Printf("error saving profile picture in handler: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "error saving incident. please try again", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "Success", nil)

}

func UpdateReceiveNotifications(c *gin.Context) {
	var account Account
	var err error

	if err = c.ShouldBindJSON(&account); err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Invalid format data", nil)
		return
	}

	account.AccountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "error", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	err = service.UpdateReceiveNotifications(account.AccountID)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error while updating phone number information. Please try later", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "success", nil)
}

func DesactivateAccount(c *gin.Context) {
	var account Account
	var err error

	if err = c.ShouldBindJSON(&account); err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Invalid format data", nil)
		return
	}

	account.AccountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "error", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	err = service.DesactivateAccount(account)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error while updating phone number information. Please try later", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "success", nil)
}
