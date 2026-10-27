package repositories

import (
	fieldRepo "field-service/repositories/field"
	fieldScheduleRepo "field-service/repositories/fieldschedule"
	timeRepo "field-service/repositories/time"

	"gorm.io/gorm"
)

type Regitry struct {
	db *gorm.DB
}

type IRepositoryRegistry interface {
	GetField() fieldRepo.IFieldRepository
	GetFieldSchedule() fieldScheduleRepo.IFieldScheduleRepository
	GetTime() timeRepo.ITimeRepository
}

func NewRepositoryRegistry(db *gorm.DB) IRepositoryRegistry {
	return &Regitry{db: db}
}

func (r *Regitry) GetField() fieldRepo.IFieldRepository {
	return fieldRepo.NewFieldRepository(r.db)
}

func (r *Regitry) GetFieldSchedule() fieldScheduleRepo.IFieldScheduleRepository {
	return fieldScheduleRepo.NewFieldScheduleRepository(r.db)
}

func (r *Regitry) GetTime() timeRepo.ITimeRepository {
	return timeRepo.NewTimeRepository(r.db)
}
