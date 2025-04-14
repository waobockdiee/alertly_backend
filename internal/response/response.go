package response

import "github.com/gin-gonic/gin"

type APIResponse struct {
	IsError    bool        `json:"error"`
	StatusCode int         `json:"error_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func Send(c *gin.Context, statusCode int, isError bool, message string, data interface{}) {
	c.SecureJSON(statusCode, APIResponse{
		IsError:    isError,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}
