package response

import (
	"net/http"
	"field-service/constants"
	errConstants "field-service/constants/error"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Token   *string     `json:"token,omitempty"`
}

type ParamHTTPRes struct {
	Code    int
	Err     error
	Message *string
	Gin     *gin.Context
	Data    interface{}
	Token   *string
}

func HttpResponse(param ParamHTTPRes) {
	if param.Err == nil {
		param.Gin.JSON(param.Code, Response{
			Status:  constants.Success,
			Message: http.StatusText(http.StatusOK),
			Data:    param.Data,
			Token:   param.Token,
		})
		return
	}

	message := errConstants.ErrInternalServerError.Error()
	if param.Message != nil {
		message = *param.Message
	} else if errConstants.ErrMapping(param.Err) {
		message = param.Err.Error()
	}

	param.Gin.JSON(param.Code, Response{
		Status:  constants.Error,
		Message: message,
		Data:    param.Data,
	})
}
