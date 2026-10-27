package controllers

import (
	fieldControllers "field-service/controllers/field"
	fieldSchedulerController "field-service/controllers/fieldschedule"
	timeControllers3 "field-service/controllers/time"
	"field-service/services"
)

type Registry struct {
	service services.IServiceRegistry
}

type IControllerRegistry interface {
	GetField() fieldControllers.IFieldController
	GetFieldSchedule() fieldSchedulerController.IFieldScheduleController
	GetTime() timeControllers3.ITimeController
}

func NewControllerRegistry(service services.IServiceRegistry) IControllerRegistry {
	return &Registry{service: service}
}

func (r *Registry) GetField() fieldControllers.IFieldController {
	return fieldControllers.NewFieldController(r.service)
}

func (r *Registry) GetFieldSchedule() fieldSchedulerController.IFieldScheduleController {
	return fieldSchedulerController.NewFieldScheduleController(r.service)
}

func (r *Registry) GetTime() timeControllers3.ITimeController {
	return timeControllers3.NewTimeController(r.service)
}
