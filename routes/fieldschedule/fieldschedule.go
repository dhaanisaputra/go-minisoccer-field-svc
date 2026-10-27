package routes

import (
	"field-service/clients"
	"field-service/constants"
	"field-service/controllers"
	"field-service/middlewares"

	"github.com/gin-gonic/gin"
)

type FieldScheduleRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	client     clients.IClientRegistry
}

type IFieldScheduleRoute interface {
	Run()
}

func NewFieldScheduleRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, client clients.IClientRegistry) IFieldScheduleRoute {
	return &FieldScheduleRoute{
		controller: controller,
		group:      group,
		client:     client,
	}
}

func (r *FieldScheduleRoute) Run() {
	group := r.group.Group("/field/schedule")
	group.GET("", middlewares.AuthenticatedWithoutToken(), r.controller.GetFieldSchedule().GetAllByFieldIDAndDate)
	group.PATCH("", middlewares.AuthenticatedWithoutToken(), r.controller.GetFieldSchedule().UpdateStatus)
	group.Use(middlewares.Authenticated())
	group.GET("/pagination", middlewares.CheckRole([]string{
		constants.Admin,
		constants.Customer,
	}, r.client),
		r.controller.GetFieldSchedule().GetAllWithPagination)
	group.GET("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
		constants.Customer,
	}, r.client),
		r.controller.GetFieldSchedule().GetByUUID)
	group.POST("", middlewares.CheckRole([]string{
		constants.Admin,
	}, r.client),
		r.controller.GetFieldSchedule().Create)
	group.POST("/one-month", middlewares.CheckRole([]string{
		constants.Admin,
	}, r.client),
		r.controller.GetFieldSchedule().GenerateScheduleForOneMonth)
	group.PUT("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
	}, r.client),
		r.controller.GetFieldSchedule().Update)
	group.DELETE("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
	}, r.client),
		r.controller.GetFieldSchedule().Delete)
}
