package routes

import (
	"field-service/clients"
	"field-service/constants"
	"field-service/controllers"
	"field-service/middlewares"

	"github.com/gin-gonic/gin"
)

type FieldRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	client     clients.IClientRegistry
}

type IFieldRoute interface {
	Run()
}

func NewFieldRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, client clients.IClientRegistry) IFieldRoute {
	return &FieldRoute{
		controller: controller,
		group:      group,
		client:     client,
	}
}

func (r *FieldRoute) Run() {
	group := r.group.Group("/field")
	group.GET("", middlewares.AuthenticatedWithoutToken(), r.controller.GetField().GetAllWithoutPagination)
	group.GET("/:uuid", middlewares.AuthenticatedWithoutToken(), r.controller.GetField().GetByUUID)
	group.Use(middlewares.Authenticated())
	group.GET("/pagination", middlewares.CheckRole([]string{
		constants.Admin,
		constants.Customer,
	}, r.client),
		r.controller.GetField().GetAllWithPagination)
	group.POST("", middlewares.CheckRole([]string{
		constants.Admin,
	}, r.client),
		r.controller.GetField().Create)
	group.PUT("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
	}, r.client),
		r.controller.GetField().Update)
	group.DELETE("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
	}, r.client),
		r.controller.GetField().Delete)
}
