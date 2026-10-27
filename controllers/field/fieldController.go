package controllers

import (
	"field-service/common/response"
	"field-service/domain/dto"
	"field-service/services"
	"net/http"

	errValidation "field-service/common/error"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type FieldController struct {
	service services.IServiceRegistry
}

type IFieldController interface {
	GetAllWithPagination(*gin.Context)
	GetAllWithoutPagination(*gin.Context)
	GetByUUID(*gin.Context)
	Create(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
}

func NewFieldController(service services.IServiceRegistry) IFieldController {
	return &FieldController{service: service}
}

func (c *FieldController) GetAllWithPagination(g *gin.Context) {
	var params dto.FieldRequestParam
	err := g.ShouldBindQuery(&params)
	if err != nil {
		response.HttpResponse(response.ParamHTTPRes{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  g,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(params)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPRes{
			Code:    http.StatusBadRequest,
			Err:     err,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     g,
		})
		return
	}

	result, err := c.service.GetField().GetAllWithPagination(g, &params)
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

func (c *FieldController) GetAllWithoutPagination(g *gin.Context) {
	result, err := c.service.GetField().GetAllWithoutPagination(g)
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

func (c *FieldController) GetByUUID(g *gin.Context) {
	result, err := c.service.GetField().GetByUUID(g, g.Param("uuid"))
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

func (c *FieldController) Create(g *gin.Context) {
	var req dto.FieldRequest
	err := g.ShouldBindWith(&req, binding.FormMultipart)
	if err != nil {
		response.HttpResponse(response.ParamHTTPRes{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  g,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPRes{
			Code:    http.StatusBadRequest,
			Err:     err,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     g,
		})
		return
	}

	result, err := c.service.GetField().Create(g, &req)
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

func (c *FieldController) Update(g *gin.Context) {
	var req dto.UpdateFieldRequest
	err := g.ShouldBindWith(&req, binding.FormMultipart)
	if err != nil {
		response.HttpResponse(response.ParamHTTPRes{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  g,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPRes{
			Code:    http.StatusBadRequest,
			Err:     err,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     g,
		})
		return
	}

	result, err := c.service.GetField().Update(g, g.Param("uuid"), &req)
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

func (c *FieldController) Delete(g *gin.Context) {
	err := c.service.GetField().Delete(g, g.Param("uuid"))
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
		Data: nil,
		Gin:  g,
	})
}
