package services

import (
	"field-service/common/storage"
	"field-service/repositories"
	services "field-service/services/field"
	fieldschedule "field-service/services/fieldschedule"
	time "field-service/services/time"
)

type Registry struct {
	repository repositories.IRepositoryRegistry
	storage    storage.IStorageClient
}

type IServiceRegistry interface {
	GetField() services.IFieldService
	GetFieldSchedule() fieldschedule.IFieldScheduleService
	GetTime() time.ITimeService
}

func NewServiceRegistry(repository repositories.IRepositoryRegistry, store storage.IStorageClient) IServiceRegistry {
	return &Registry{repository: repository, storage: store}
}

func (r *Registry) GetField() services.IFieldService {
	return services.NewFieldService(r.repository, r.storage)
}

func (r *Registry) GetFieldSchedule() fieldschedule.IFieldScheduleService {
	return fieldschedule.NewFieldScheduleService(r.repository)
}

func (r *Registry) GetTime() time.ITimeService {
	return time.NewTimeService(r.repository)
}
