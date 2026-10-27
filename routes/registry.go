package routes

import (
	"field-service/clients"
	"field-service/controllers"
	fieldRoutes "field-service/routes/field"
	fieldScheduleroutes "field-service/routes/fieldschedule"
	timeRoutes "field-service/routes/time"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	client     clients.IClientRegistry
}

type IRegistry interface {
	Serve()
}

func NewRouteRegistry(controller controllers.IControllerRegistry, group *gin.RouterGroup, client clients.IClientRegistry) IRegistry {
	return &Registry{
		controller: controller,
		group:      group,
		client:     client,
	}
}

func (r *Registry) fieldRoute() fieldRoutes.IFieldRoute {
	return fieldRoutes.NewFieldRoute(r.controller, r.group, r.client)
}

func (r *Registry) fieldScheduleRoute() fieldScheduleroutes.IFieldScheduleRoute {
	return fieldScheduleroutes.NewFieldScheduleRoute(r.controller, r.group, r.client)
}

func (r *Registry) timeRoute() timeRoutes.ITimeRoute {
	return timeRoutes.NewTimeRoute(r.controller, r.group, r.client)
}

func (r *Registry) Serve() {
	r.fieldRoute().Run()
	r.fieldScheduleRoute().Run()
	r.timeRoute().Run()
}
