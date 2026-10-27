package controllers

import (
	errValidation "field-service/common/error"
	"field-service/common/response"
	"field-service/domain/dto"
	"field-service/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type FieldScheduleController struct {
	service services.IServiceRegistry
}

type IFieldScheduleController interface {
	GetAllWithPagination(*gin.Context)
	GetAllByFieldIDAndDate(*gin.Context)
	GetByUUID(*gin.Context)
	Create(*gin.Context)
	Update(*gin.Context)
	UpdateStatus(*gin.Context)
	Delete(*gin.Context)
	GenerateScheduleForOneMonth(*gin.Context)
}

func NewFieldScheduleController(service services.IServiceRegistry) IFieldScheduleController {
	return &FieldScheduleController{service: service}
}

func (c *FieldScheduleController) GetAllWithPagination(g *gin.Context) {
	var params dto.FieldScheduleRequestParam
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

	result, err := c.service.GetFieldSchedule().GetAllWithPagination(g, &params)
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

func (c *FieldScheduleController) GetAllByFieldIDAndDate(g *gin.Context) {
	var params dto.FieldScheduleByFieldIDAndDateRequestParam
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

	result, err := c.service.GetFieldSchedule().GetAllByFieldIDAndDate(g, g.Param("uuid"), params.Date)
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

func (c *FieldScheduleController) GetByUUID(g *gin.Context) {
	result, err := c.service.GetFieldSchedule().GetByUUID(g, g.Param("uuid"))
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

func (c *FieldScheduleController) Create(g *gin.Context) {
	var params dto.FieldScheduleRequest
	err := g.ShouldBindJSON(&params)
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

	err = c.service.GetFieldSchedule().Create(g, &params)
	if err != nil {
		response.HttpResponse(response.ParamHTTPRes{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  g,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPRes{
		Code: http.StatusCreated,
		Gin:  g,
	})
}

func (c *FieldScheduleController) GenerateScheduleForOneMonth(g *gin.Context) {
	var params dto.GenerateFieldScheduleForOneMonthRequest
	err := g.ShouldBindJSON(&params)
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

	err = c.service.GetFieldSchedule().GenerateScheduleForOneMonth(g, &params)
	if err != nil {
		response.HttpResponse(response.ParamHTTPRes{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  g,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPRes{
		Code: http.StatusCreated,
		Gin:  g,
	})
}

func (c *FieldScheduleController) Update(g *gin.Context) {
	var params dto.UpdateFieldScheduleRequest
	err := g.ShouldBindJSON(&params)
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

	result, err := c.service.GetFieldSchedule().Update(g, g.Param("uuid"), &params)
	if err != nil {
		response.HttpResponse(response.ParamHTTPRes{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  g,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPRes{
		Code: http.StatusCreated,
		Gin:  g,
		Data: result,
	})
}

func (c *FieldScheduleController) UpdateStatus(g *gin.Context) {
	var request dto.UpdateStatusFieldScheduleRequest
	err := g.ShouldBindJSON(&request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPRes{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  g,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
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

	err = c.service.GetFieldSchedule().UpdateStatus(g, &request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPRes{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  g,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPRes{
		Code: http.StatusCreated,
		Gin:  g,
	})
}

func (c *FieldScheduleController) Delete(g *gin.Context) {
	err := c.service.GetFieldSchedule().Delete(g, g.Param("uuid"))
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
		Gin:  g,
	})
}
