package comments

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

func SaveClusterComment(c *gin.Context) {
	var accountID int64
	var err error
	var incoID int64

	accountID, err = auth.GetUserFromContext(c)
	if err != nil {
		fmt.Println("error1", err)
		response.Send(c, http.StatusInternalServerError, true, "We couldn’t verify your session. Please log in again.", err.Error())
		return
	}

	var comment InComment
	if err = c.BindJSON(&comment); err != nil {
		fmt.Println("error1", err)
		response.Send(c, http.StatusBadRequest, true, "Invalid input format. Please check the data and try again.", err.Error())
		return
	}

	if err = validate.Struct(comment); err != nil {
		fmt.Println("error1", err)
		response.Send(c, http.StatusBadRequest, true, "Some fields are missing or incorrect. Please review the form and try again.", err.Error())
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	comment.AccountID = accountID

	incoID, err = service.Save(comment)

	if err != nil {
		fmt.Println("error1", err)
		response.Send(c, http.StatusInternalServerError, true, "We couldn’t save your comment. Please try again later.", 0)
		return
	}
	var commentOut Comment
	commentOut, err = service.GetCommentById(incoID)

	if err != nil {
		fmt.Println("error1", err)
		response.Send(c, http.StatusInternalServerError, true, "Comment was saved, but we couldn’t retrieve it. Please refresh or try again later.", 0)
		return
	}
	response.Send(c, http.StatusOK, false, "Comment sent", commentOut)

}

func GetClusterComments(c *gin.Context) {
	var result []Comment
	var err error
	var inclID int64
	tmpInclID := c.Param("incl_id")
	inclID, err = strconv.ParseInt(tmpInclID, 10, 64)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Invalid ID format. Please try again.", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)
	result, err = service.GetClusterCommentsByID(inclID)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "We couldn’t load the comments. Please try again later.", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "success", result)
}
