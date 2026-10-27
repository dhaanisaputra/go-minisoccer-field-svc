package controllers

import (
	"field-service/common/response"
	"field-service/domain/dto"
	"field-service/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TimeController struct {
	service services.IServiceRegistry
}

type ITimeController interface {
	GetAll(*gin.Context)
	GetByUUID(*gin.Context)
	Create(*gin.Context)
}

func NewTimeController(service services.IServiceRegistry) ITimeController {
	return &TimeController{service: service}
}

func (c *TimeController) GetAll(g *gin.Context) {
	result, err := c.service.GetTime().GetAll(g)
	if err != nil {
		response.HttpResponse(response.ParamHTTPRes{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  g,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPRes{
		Code: http.StatusOK,
		Data: result,
		Gin:  g,
	})
}

func (c *TimeController) GetByUUID(g *gin.Context) {
	result, err := c.service.GetTime().GetByUUID(g, g.Param("uuid"))
	if err != nil {
		response.HttpResponse(response.ParamHTTPRes{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  g,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPRes{
		Code: http.StatusOK,
		Data: result,
		Gin:  g,
	})
}

func (c *TimeController) Create(g *gin.Context) {
	var request dto.TimeRequest
	result, err := c.service.GetTime().Create(g, &request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPRes{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  g,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPRes{
		Code: http.StatusOK,
		Data: result,
		Gin:  g,
	})
}
