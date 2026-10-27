package routes

import (
	"field-service/clients"
	"field-service/constants"
	"field-service/controllers"
	"field-service/middlewares"

	"github.com/gin-gonic/gin"
)

type TimeRoute struct {
	controllers controllers.IControllerRegistry
	group       *gin.RouterGroup
	client      clients.IClientRegistry
}

type ITimeRoute interface {
	Run()
}

func NewTimeRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, client clients.IClientRegistry) ITimeRoute {
	return &TimeRoute{
		controllers: controller,
		group:       group,
		client:      client,
	}
}

func (r *TimeRoute) Run() {
	group := r.group.Group("/time")
	group.Use(middlewares.Authenticated())
	group.GET("", middlewares.CheckRole([]string{
		constants.Admin,
	}, r.client), r.controllers.GetTime().GetAll)
	group.GET("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
	}, r.client), r.controllers.GetTime().GetByUUID)
	group.POST("", middlewares.CheckRole([]string{
		constants.Admin,
	}, r.client), r.controllers.GetTime().Create)
}
