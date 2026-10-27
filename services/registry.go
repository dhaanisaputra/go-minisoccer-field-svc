package services

import (
	"field-service/common/gcs"
	"field-service/repositories"
	services "field-service/services/field"
	fieldschedule "field-service/services/fieldschedule"
	time "field-service/services/time"
)

type Registry struct {
	repository repositories.IRepositoryRegistry
	gcs        gcs.IGCSClient
}

type IServiceRegistry interface {
	GetField() services.IFieldService
	GetFieldSchedule() fieldschedule.IFieldScheduleService
	GetTime() time.ITimeService
}

func NewServiceRegistry(repository repositories.IRepositoryRegistry, gcs gcs.IGCSClient) IServiceRegistry {
	return &Registry{repository: repository, gcs: gcs}
}

func (r *Registry) GetField() services.IFieldService {
	return services.NewFieldService(r.repository, r.gcs)
}

func (r *Registry) GetFieldSchedule() fieldschedule.IFieldScheduleService {
	return fieldschedule.NewFieldScheduleService(r.repository)
}

func (r *Registry) GetTime() time.ITimeService {
	return time.NewTimeService(r.repository)
}
