package comments

import (
	"alertly/internal/auth"
	"alertly/internal/database"
	"alertly/internal/response"
	"fmt"
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
		response.Send(c, http.StatusInternalServerError, true, "error", err.Error())
		return
	}

	var comment InComment
	if err = c.BindJSON(&comment); err != nil {
		fmt.Println("error1", err)
		response.Send(c, http.StatusBadRequest, true, "Wrong data in", err.Error())
		return
	}

	if err = validate.Struct(comment); err != nil {
		fmt.Println("error1", err)
		response.Send(c, http.StatusBadRequest, true, "Bad request", err.Error())
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	comment.AccountID = accountID

	incoID, err = service.Save(comment)

	if err != nil {
		fmt.Println("error1", err)
		response.Send(c, http.StatusInternalServerError, true, "error", 0)
		return
	}
	var commentOut Comment
	commentOut, err = service.GetCommentById(incoID)

	if err != nil {
		fmt.Println("error1", err)
		response.Send(c, http.StatusInternalServerError, true, "error", 0)
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
		response.Send(c, http.StatusInternalServerError, true, "error parsing data", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)
	result, err = service.GetClusterCommentsByID(inclID)

	if err != nil {
		response.Send(c, http.StatusInternalServerError, true, "error fetching data", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "success", result)
}
